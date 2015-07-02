//+build !noasm,amd64

package msgp

//go:noescape
func putBinHdr(p *byte, sz int) int

//go:noescape
func putStrHdr(p *byte, sz int) int

//go:noescape
func putArrayHdr(p *byte, sz int) int

//go:noescape
func putMapHdr(p *byte, sz int) int

//go:noescape
func getUnix(b []byte) (sec int64, nsec int32)

//go:noescape
func putUnix(b []byte, sec int64, nsec int32)

//go:noescape
func putMint64(b []byte, i int64)

//go:noescape
func getMint64(b []byte) int64

//go:noescape
func putMint32(b []byte, i int32)

//go:noescape
func getMint32(b []byte) int32

//go:noescape
func putMuint64(b []byte, u uint64)

//go:noescape
func getMuint64(b []byte) uint64

//go:noescape
func prefixu64(b []byte, sz uint64, pre byte)

//go:noescape
func prefixu32(b []byte, sz uint32, pre byte)

//go:noescape
func put232(p []byte, a float32, b float32)

//go:noescape
func put264(p []byte, a float64, b float64)
