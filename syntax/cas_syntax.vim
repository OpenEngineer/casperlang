" Vim syntax file
" Language: casperlang
" Filenames: *.cas

if exists("b:current_syntax")
  finish
endif

syn keyword Keyword import panic

syn match Special "\\\d\d\d\|\\." contained
syn region String start=+"+  skip=+\\\\\|\\"+  end=+"\|$+	contains=Special
syn match Special "'\\.'" contained
syn match Macro "$\(\d\+\)\?"

syn keyword Type String Int Any IO Float Bool Path

syn match Constant "-\=\<\d\+L\=\>\|0[b][0-1]\+\>"
syn match Constant "-\=\<\d\+L\=\>\|0[o][0-7]\+\>"
syn match Constant "-\=\<\d\+L\=\>\|0[x][0-9a-fA-F]\+\>"
syn match Constant '\<\zs\d\+\(\.\d\+\([e][-]\?\d\+\)\?\)\?\ze'

syn keyword Todo contained TODO XXX NOTE
syn match	Comment	"#.*" contains=Todo

let b:current_syntax = "casper"
