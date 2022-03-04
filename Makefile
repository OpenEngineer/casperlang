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
