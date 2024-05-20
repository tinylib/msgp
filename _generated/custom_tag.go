package _generated

//go:generate msgp
//msgp:tag mytag
//msgp:tag anothertag CustomTag2

type CustomTag struct {
	Foo string `mytag:"foo_custom_name"`
	Bar int    `mytag:"bar1234"`
}

type CustomTag2 struct {
	Foo string `anothertag:"foo_msgp_name" mytag:"foo_legacy_name"`
	Bar int    `anothertag:"bar5678" mytag:"not_used"`
}
