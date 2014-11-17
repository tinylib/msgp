
install:
	@go install ./...

generate:
	@go generate ./_generated

test: install generate
	@go test -v ./_generated

test-pkg: install
	@export GOFILE=./_generated/ && msgp -o ./_generated/generated.go
	@go test -v ./_generated

bench: install generate
	@go test -bench . ./_generated

clean:
	rm ./_generated/generated.go && rm ./_generated/generated_test.go