.PHONY: build build-rust build-go clean run test

build: build-rust build-go

build-rust:
	cd rust && cargo build --release

build-go:
	go build -o systui ./cmd/systui

clean:
	cd rust && cargo clean
	rm -f systui
	rm -f *.so

run: build
	./systui

test:
	go test ./...
	cd rust && cargo test
