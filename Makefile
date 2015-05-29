
install:
	@go install ./...

generate:
	@go generate ./_generated

test: install generate
	@go test -v ./_generated

test-pkg: install
	@msgp -o ./_generated/generated.go -file ./_generated
	@go test -v ./_generated

bench: install generate
	@go test -bench . ./_generated

clean:
	$(RM) ./_generated/generated.go && $(RM) ./_generated/generated_test.go