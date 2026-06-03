package _generated

//go:generate msgp -d "tags primary,fallback"

type CustomTags struct {
	Field1 string `primary:"f1_primary"`
	Field2 string `fallback:"f2_fallback"`
	Field3 string `primary:"f3_primary" fallback:"f3_fallback_ignored"`
	Field4 string `msg:"f4_msg"`
	Field5 string `msgpack:"f5_msgpack"`
	Field6 string `fallback:"f6_fallback" msg:"f6_msg_ignored"`
	Field7 string `fallback:"f7_fallback" msgpack:"f7_msgpack_ignored"`
	Field8 string `primary:"f8_primary" msg:"f8_msg_ignored" msgpack:"f8_msgpack_ignored"`
}
