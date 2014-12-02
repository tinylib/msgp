package gen

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"text/template"
)

var (
	decTemplate         *template.Template
	encTemplate         *template.Template
	marTemplate         *template.Template
	unmTemplate         *template.Template
	benTemplate         *template.Template
	sizTemplate         *template.Template
	marshalTestTemplate *template.Template
	encodeTestTemplate  *template.Template
)

func init() {
	_, prefix, _, _ := runtime.Caller(0)
	prefix = filepath.Dir(prefix) + "/"

	decTemplate = template.Must(template.ParseFiles(prefix+"decode.tmpl", prefix+"elem_dec.tmpl"))
	encTemplate = template.Must(template.ParseFiles(prefix+"encode.tmpl", prefix+"elem_enc.tmpl"))
	marTemplate = template.Must(template.ParseFiles(prefix+"marshal.tmpl", prefix+"marshal_enc.tmpl"))
	unmTemplate = template.Must(template.ParseFiles(prefix+"unmarshal.tmpl", prefix+"elem_unm.tmpl"))
	sizTemplate = template.Must(template.ParseFiles(prefix+"size.tmpl", prefix+"size_enc.tmpl"))

	marshalTestTemplate = template.Must(template.ParseFiles(prefix + "testMarshal.tmpl"))
	encodeTestTemplate = template.Must(template.ParseFiles(prefix + "testEncode.tmpl"))
}

// execAndFormat executes a template and formats the output, using buf as temporary storage
func execAndFormat(t *template.Template, w io.Writer, i interface{}, buf *bytes.Buffer) error {
	if buf == nil {
		buf = bytes.NewBuffer(nil)
	}
	buf.Reset()
	err := t.Execute(buf, i)
	if err != nil {
		return fmt.Errorf("template: %s", err)
	}
	bts, err := format.Source(buf.Bytes())
	if err != nil {
		w.Write(buf.Bytes())
		return fmt.Errorf("gofmt: %s", err)
	}
	_, err = w.Write(bts)
	return err
}

// WriteEncodeDecode writes the EncodeMsg, EncodeTo, DecodeMsg, and DecodeFrom methods,
// using buf as scratch space.
func WriteEncodeDecode(w io.Writer, p *Ptr, buf *bytes.Buffer) error {
	err := execAndFormat(decTemplate, w, p, buf)
	if err != nil {
		return err
	}
	return execAndFormat(encTemplate, w, p, buf)
}

// WriteMarhsalUnmarshal writes the MarshalMsg, UnmarshalMsg, AppendMsg, and Maxsize
// methods using buf as scratch space.
func WriteMarshalUnmarshal(w io.Writer, p *Ptr, buf *bytes.Buffer) error {
	err := execAndFormat(marTemplate, w, p, buf)
	if err != nil {
		return err
	}
	err = execAndFormat(unmTemplate, w, p, buf)
	if err != nil {
		return err
	}
	return execAndFormat(sizTemplate, w, p, buf)
}
