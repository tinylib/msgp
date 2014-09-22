
install:
	@go install ./...

test: install
	@go test -v

bench: test
	@go test -bench . ./_generated