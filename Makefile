
# NOTE: This Makefile is only necessary if you 
# plan on developing the msgp tool and library.
# Installation can still be performed with a
# normal `go install`.

# For more information, please see HACKING.md

GGEN=./_generated/generated.go ./_generated/generated_test.go
MGEN=./msgp/defgen_test.go
SHELL=/bin/bash
BIN=$(GOBIN)/msgp
BRANCH=$(shell git symbolic-ref --short HEAD)

.PHONY: clean wipe install get-deps bench all lint travis

$(BIN): *.go
	go install ./...

$(GGEN): ./_generated/def.go
	go generate ./_generated

$(MGEN): ./msgp/defs_test.go
	go generate ./msgp

$(BRANCH)-gen-bench.txt: all
	go test ./_generated -run=NONE -bench . | tee $(BRANCH)-gen-bench.txt

$(BRANCH)-unit-bench.txt: all
	go test ./msgp -run=NONE -bench . | tee $(BRANCH)-unit-bench.txt

$(GOBIN)/benchcmp:
	go get golang.org/x/tools/cmd/benchcmp

bench: $(BRANCH)-gen-bench.txt $(BRANCH)-unit-bench.txt

install: $(BIN)

test: all
	go test -v ./msgp
	go test - v./_generated

benchcmp: $(BRANCH)-gen-bench.txt master-gen-bench.txt $(GOBIN)/benchcmp
	benchcmp master-gen-bench.txt $(BRANCH)-gen-bench.txt

clean:
	$(RM) $(GGEN) $(MGEN)

wipe: clean
	$(RM) $(BIN)

get-deps:
	go get -d -t ./...

all: install $(GGEN) $(MGEN)

travis: get-deps test
