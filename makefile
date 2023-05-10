build:
	cd mdbook && ./../../mdBook/target/debug/mdbook serve
	
run:
	go run cmd/main.go