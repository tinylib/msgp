package _generated

//go:generate msgp -d "tag mytag"

type CustomTag struct {
	Foo string `mytag:"foo_custom_name"`
	Bar int    `mytag:"bar1234"`
}
