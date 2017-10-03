package _generated

//go:generate msgp

// The leaves of interest in this crazy structs are strings. The test case
// looks for strings in the serialised msgpack and makes them unreadable.

type ErrorCtxMapChild struct {
	Val string
}

type ErrorCtxAsMap struct {
	Val      string
	Child    *ErrorCtxMapChild
	Children []*ErrorCtxMapChild
	Map      map[string]string

	Nest struct {
		Val      string
		Child    *ErrorCtxMapChild
		Children []*ErrorCtxMapChild
		Map      map[string]string

		Nest struct {
			Val      string
			Child    *ErrorCtxMapChild
			Children []*ErrorCtxMapChild
			Map      map[string]string
		}
	}
}

//msgp:tuple ErrorCtxTupleChild

type ErrorCtxTupleChild struct {
	Val string
}

//msgp:tuple ErrorCtxAsTuple

type ErrorCtxAsTuple struct {
	Val      string
	Child    *ErrorCtxTupleChild
	Children []*ErrorCtxTupleChild
	Map      map[string]string

	Nest struct {
		Val      string
		Child    *ErrorCtxTupleChild
		Children []*ErrorCtxTupleChild
		Map      map[string]string

		Nest struct {
			Val      string
			Child    *ErrorCtxTupleChild
			Children []*ErrorCtxTupleChild
			Map      map[string]string
		}
	}
}
