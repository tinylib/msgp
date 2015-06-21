
# NOTE: This Makefile is only necessary if you 
# plan on developing the msgp tool and library.
# Installation can still be performed with a
# normal `go install`.

# generated integration test files
GGEN = ./_generated/generated.go ./_generated/generated_test.go
# generated unit test files
MGEN = ./msgp/defgen_test.go

install:
	go install ./...

$(GGEN): ./_generated/def.go
	go generate ./_generated

$(MGEN): ./msgp/defs_test.go
	go generate ./msgp

test: install $(GGEN) $(MGEN)
	go test -v ./msgp
	go test -v ./_generated

bench: $(GGEN) $(MGEN) install
	go test -bench . ./msgp
	go test -bench . ./_generated

clean:
	$(RM) $(GGEN) $(MGEN)

get-deps:
	go get -d -t ./...

# travis CI enters here
travis: get-deps test
