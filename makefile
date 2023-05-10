
all:
	go run /main.go

build:
	cd mdbook && ./../../mdBook/target/debug/mdbook serve