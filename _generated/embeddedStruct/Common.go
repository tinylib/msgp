package embeddedStruct

//go:generate msgp -file ./

type Common struct {
	RequestID uint32 `msg:"request_id"`
	Token     string `msg:"token"`
}
