main: ./build/casper 

./build/casper: ./src/*.go | ./build
	go build -o ./build/casper  $^

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
