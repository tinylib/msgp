package msgp

// NonStrMapKey must be implemented to allow non-string Go types to be used as MessagePack map keys.
// Msgp maps must have string keys for JSON interop.
// NonStrMapKey enables conversion from type to string and vice versa.
type NonStrMapKey interface {
	MsgpStrMapKey() string
	MsgpFromStrMapKey(s string) NonStrMapKey
	MsgpStrMapKeySize() int
}
