package msgp

// EndlessReader is an io.Reader
// that loops over the same data
// endlessly. It is primarily useful
// for testing and benchmarking.
type EndlessReader struct {
	data   []byte
	offset int
}

// NewEndlessReader returns a new endless reader
func NewEndlessReader(b []byte) *EndlessReader {
	return &EndlessReader{data: b, offset: 0}
}

// Read implements io.Reader. In practice, it
// always returns (len(p), nil).
func (c *EndlessReader) Read(p []byte) (int, error) {
	var n int
	l := len(p)
	m := len(c.data)
	for n < l {
		nn := copy(p[n:], c.data[c.offset:])
		n += nn
		c.offset += nn
		c.offset %= m
	}
	return n, nil
}
