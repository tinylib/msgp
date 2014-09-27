
install:
	@go install ./...

generate:
	@go generate ./_generated

test: install generate
	@go test -v ./_generated

bench: install generate
	@go test -v -bench . ./_generated