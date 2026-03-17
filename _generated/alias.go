package _generated

//go:generate msgp

type AliasedFieldsV3 struct {
	CurrentName string `msg:"newname,alias=name"`
	MultiAlias  int    `msg:"key3,alias=key2;key1"`
	NoAlias     bool   `msg:"no_alias"`
}

type AliasedFieldsV2 struct {
	CurrentName string `msg:"newname,alias=name"`
	MultiAlias  int    `msg:"key2,alias=key1"`
	NoAlias     bool   `msg:"no_alias"`
}

type AliasedFieldsV1 struct {
	CurrentName string `msg:"name"`
	MultiAlias  int    `msg:"key1"`
	NoAlias     bool   `msg:"no_alias"`
}
