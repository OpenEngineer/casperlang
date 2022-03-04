package main

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"

	billy "gopkg.in/src-d/go-billy.v4"
	billymemfs "gopkg.in/src-d/go-billy.v4/memfs"
	gitcore "gopkg.in/src-d/go-git.v4"
	gitconfig "gopkg.in/src-d/go-git.v4/config"
	gitplumbing "gopkg.in/src-d/go-git.v4/plumbing"
	gittransport "gopkg.in/src-d/go-git.v4/plumbing/transport"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	gitmemory "gopkg.in/src-d/go-git.v4/storage/memory"
)

func correctGitUrl(url_ *String, sshKey string) *String {
	url := url_.Value()

	if !strings.HasSuffix(url, ".git") {
		url += ".git"
	}

	if sshKey == "" {
		if !strings.HasPrefix(url, "https://") {
			url = "https://" + url
		}
	} else {
		if strings.HasPrefix(url, "https://") {
			url = strings.TrimPrefix(url, "https://")
		} else if strings.HasPrefix(url, "http://") {
			url = strings.TrimPrefix(url, "http://")
		}

		url = "git@" + strings.Replace(url, "/", ":", 1)
	}

	return NewString(url, url_.Context())
}

func newGitAuthMethod(sshKey string) gittransport.AuthMethod {
	if sshKey == "" {
		return nil
	}

	// pem decoding is done inside gitssh
	authMethod, err := gitssh.NewPublicKeys("git", []byte(sshKey), "")
	if err != nil {
		panic(err)
	}

	return authMethod
}

// returns nil if not found
func loopGitReferenceNames(url *String, sshKey string, cond func(gitplumbing.ReferenceName) error) error {
	storer := gitmemory.NewStorage()

	remoteConfig := &gitconfig.RemoteConfig{
		Name: "origin",
		URLs: []string{url.Value()},
	}

	if err := remoteConfig.Validate(); err != nil {
		panic(err)
	}

	remote := gitcore.NewRemote(storer, remoteConfig)

	lstOptions := &gitcore.ListOptions{
		Auth: newGitAuthMethod(sshKey),
	}

	lst, err := remote.List(lstOptions)
	if err != nil {
		return url.Context().Error("fetch error (" + err.Error() + ")")
	}

	for _, ref := range lst {
		if err := cond(ref.Name()); err != nil {
			return err
		}
	}

	return nil
}

func parseGitTag(ref_ gitplumbing.ReferenceName) (string, bool) {
	ref := string(ref_)

	parts := strings.Split(ref, "/")

	if parts[1] != "tags" {
		return "", false
	}

	tag := parts[2]

	return tag, true
}

func selectGitTag(url *String, tag *String, sshKey string) (gitplumbing.ReferenceName, error) {
	found := false
	var ref gitplumbing.ReferenceName

	if err := loopGitReferenceNames(url, sshKey, func(ref_ gitplumbing.ReferenceName) error {
		if !found {
			if ref_.IsTag() {
				if refTag, ok := parseGitTag(ref_); ok {
					if refTag == tag.Value() {
						ref = ref_
						found = true
					}
				}
			}
		}

		return nil
	}); err != nil {
		return ref, err
	}

	if !found {
		return ref, tag.Context().Error("tag \"" + tag.Value() + "\" not found")
	}

	return ref, nil
}

func writeFile(fs billy.Filesystem, src string, dst string) error {
	// XXX: only write the files that pass parser tests?
	fIn, err := fs.Open(src)
	if err != nil {
		return err
	}

	defer fIn.Close()

	fOut, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer fOut.Close()

	wOut := bufio.NewWriter(fOut)

	if _, err := io.Copy(wOut, fIn); err != nil {
		return err
	}

	wOut.Flush()

	return nil
}

func writeDir(fs billy.Filesystem, dirSrc string, dirDst string) error {
	if err := os.MkdirAll(dirDst, 0755); err != nil {
		return err
	}

	files, err := fs.ReadDir(dirSrc)
	if err != nil {
		return err
	}

	for _, file := range files {
		src := fs.Join(dirSrc, file.Name())
		dst := filepath.Join(dirDst, file.Name())

		if file.IsDir() {
			if err := writeDir(fs, src, dst); err != nil {
				return err
			}
		} else {
			ext := filepath.Ext(file.Name())
			if ext == ".cas" || file.Name() == "package.json" {
				if err := writeFile(fs, src, dst); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func hasPackageConf(fs billy.Filesystem) bool {
	_, err := fs.Stat("/package.json")
	if err == nil {
		return true
	} else {
		return false
	}
}

func writeWorktree(fs billy.Filesystem, dst string) error {
	return writeDir(fs, "/", dst)
}

func cloneGitRef(url *String, ref gitplumbing.ReferenceName, sshKey string, dst *String) error {
	wt := billymemfs.New()

	storer := gitmemory.NewStorage()

	cloneOptions := &gitcore.CloneOptions{
		URL:               url.Value(),
		Auth:              newGitAuthMethod(sshKey),
		ReferenceName:     ref,
		SingleBranch:      true,
		NoCheckout:        true, // checkout follows further down
		RecurseSubmodules: gitcore.NoRecurseSubmodules,
		Progress:          nil,
	}

	if err := cloneOptions.Validate(); err != nil {
		return url.Context().Error("git fetch error (" + err.Error() + ")")
	}

	repo, err := gitcore.Clone(storer, wt, cloneOptions)
	if err != nil {
		return url.Context().Error("git fetch error (" + err.Error() + ")")
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return url.Context().Error("git fetch error (" + err.Error() + ")")
	}

	if err := worktree.Checkout(&gitcore.CheckoutOptions{
		Branch: ref,
	}); err != nil {
		return url.Context().Error("git fetch error (" + err.Error() + ")")
	}

	if !hasPackageConf(wt) {
		return url.Context().Error("git repo doesn't contain have package.json")
	}

	if err := writeWorktree(wt, dst.Value()); err != nil {
		return url.Context().Error("git writer error (" + err.Error() + ")")
	}

	return nil
}

func FetchGitRepo(url *String, version *String, sshKey string, dst *String) error {
	url = correctGitUrl(url, sshKey)

	tagRef, err := selectGitTag(url, version, sshKey)
	if err != nil {
		return err
	}

	if isFile(dst.Value()) {
		return dst.Context().Error("git clone dst \"" + dst.Value() + "\" is a file")
	}

	if !isDir(dst.Value()) {
		if err := cloneGitRef(url, tagRef, sshKey, dst); err != nil {
			return err
		}
	}

	return nil
}
