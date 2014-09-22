
install:
	@go install ./...

generate:
	@go generate ./_generated

test: install generate
	@go test -v ./_generated

bench: test
	@go test -v -bench . ./_generated