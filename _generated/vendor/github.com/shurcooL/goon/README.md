goon
====

Go Object Notation

Objectives
----------

Goon aims to be an open standard for human-readable data interchange, primarily between an existing Go program and a future Go program.

Example
-------

The following Go program:

```go
package main

import "github.com/shurcooL/go-goon"

func main() {
	type Inner struct {
		Field1 string
		Field2 int
	}
	type Lang struct {
		Name  string
		Year  int
		URL   string
		Inner *Inner
	}

	x := Lang{
		Name: "Go",
		Year: 2009,
		URL:  "http",
		Inner: &Inner{
			Field1: "Secret!",
		},
	}

	goon.Dump(x)
}
```

Will produce this goon to stdout:

```go
(Lang)(Lang{
	Name: (string)("Go"),
	Year: (int)(2009),
	URL:  (string)("http"),
	Inner: (*Inner)(&Inner{
		Field1: (string)("Secret!"),
		Field2: (int)(0),
	}),
})

```

Spec
----

* goon is valid gofmt-ed Go code.
* goon is a [Go literal](http://golang.org/ref/spec#Literal).

Relevant things
---------------

- https://github.com/davecgh/go-spew/issues/6
- https://github.com/davecgh/go-spew
- http://golang.org/pkg/encoding/gob/
- http://en.wikipedia.org/wiki/JSON
- http://golang.org/ref/spec
- gofmt
- http://semver.org/

Seriously?
----------

Yeah. I went there.

But why?
--------

Because I want programming to be as easy as possible. As long as Go programs cannot easily spit out at least parts of Go code (i.e. variable values) that I can easily look at, cut relevant parts out of, modify and paste into new Go programs, there will always be things that are harder than they have to be.

More concretely, I hack around with Go ASTs quite often. `parser.ParseFile(...)`, `for _, d := range file.Decls {`, `if f, ok := d.(*ast.FuncDecl); ok {`, `spew.Dump(x);`, that kind of stuff. Right now my iterations are slowed down because I can't easily extract various AST segments that I want to further manipulate.

Did you just steal this template from Tom's TOML?
-------------------------------------------------

Stealing implies I prevented someone from [owning](https://twitter.com/shurcooL/status/266294572949327872) something. When it comes to digital things, one typically copies. However, it's possible to steal [credit](https://twitter.com/shurcooL/status/266294965053820928) for having created something. I don't think that's the case here since I'm acknowledging Tom's original efforts. Thank you Tom, you rock! :)

Implementations
---------------

If you have an implementation, please send a pull request adding it to this list.

- Go (@shurcooL) - https://github.com/shurcooL/go-goon - not completely finished, but mostly works

Contributing
------------

Contributions are highly welcome. Think of the objective/motivation as the key, and everything else as the value that is open to any improvements.

License
-------

- [MIT License](http://opensource.org/licenses/mit-license.php)
