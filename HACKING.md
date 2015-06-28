## Hacking

The `msgp` git repo contains a top-level `Makefile` that makes it simple to test changes to the code generation libraries (`/parse` and `/gen`) without having to keep track of the build state of the tool or the generated test code.

Running `make test` runs unit tests in `/msgp` and functional tests in `/_generated`. `make clean` removes those files, and `make wipe` removes generated files *and* the `msgp` tool binary from `$GOBIN`. When you switch branches, you should probably run `make wipe && make all`.

All changes to the library that could have an impact on the performance of the generated code should confirm that there are no performance regressions by running `make bench` in both the master and feature branches, followed by `make benchcmp`. Running `make bench` generates two files: `<branch>-unit-bench.txt` and `<branch>-gen-bench.txt`. Unit benchmarks reflect performances changes in the `/msgp` library, while generated benchmarks reflect changes in the quality of the generated code.

