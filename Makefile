main: ./build/casper 

build_windows_amd64=$(abspath ./build/windows_amd64)
dst_windows_amd64= $(addprefix $(build_windows_amd64)/,casper)

windows: $(dst_windows_amd64)

./build/casper: ./src/*.go | ./build
	go build -o ./build/casper  $^

$(dst_windows_amd64): ./src/*.go | ./build
	export GOOS=windows; \
	export GOARCH=amd64; \
	cd $(dir $<); \
	go build -o $(abspath $@)

./build:
	mkdir -p ./build

install: ./build/casper
	sudo cp ./build/casper /usr/local/bin/casper

test: ./build/casper
	# run test suite
	cd ./examples; \
	../build/casper ./test_all.cas

install-vim: install-vim-syntax install-vim-indent

install-vim-%:
ifneq ($(wildcard $(HOME)/.vim/$*),)
	cp ./syntax/cas_$*.vim $(HOME)/.vim/$*/cas.vim
else ifneq ($(wildcard $(HOME)/.config/nvim/$*),)
	cp ./syntax/cas_$*.vim $(HOME)/.config/nvim/$*/cas.vim
else
	$(warning no vim $* directory found)
endif

package: ./build/casper
	tar -czf ./build/casper-0.1.0-linux_x86_64.tar.gz ./build/casper

package-windows: ./build/windows_amd64/casper
	zip ./build/casper-0.1.0-windows_amd64.zip ./build/windows_amd64/casper
