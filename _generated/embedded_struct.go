package _generated

//go:generate msgp

type GetUserRequestWithEmbeddedStruct struct {
	Common `msg:",flatten"`
	UserID uint32 `msg:"user_id"`
}

type GetUserRequest struct {
	RequestID uint32 `msg:"request_id"`
	Token     string `msg:"token"`
	UserID    uint32 `msg:"user_id"`
}

type Common struct {
	RequestID uint32 `msg:"request_id"`
	Token     string `msg:"token"`
}
