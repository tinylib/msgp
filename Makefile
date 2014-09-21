
install:
	@go install

test: install
	@go test -v

bench: test
	@go test -v -bench . ./_generated