package gen

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"text/template"
)

var (
	decTemplate *template.Template
	encTemplate *template.Template
	marTemplate *template.Template
	unmTemplate *template.Template
)

func init() {
	gopath := os.Getenv("GOPATH")
	prefix := gopath + "/src/github.com/philhofer/msgp/gen/"

	decTemplate = template.Must(template.ParseFiles(prefix+"decode.tmpl", prefix+"elem_dec.tmpl"))
	encTemplate = template.Must(template.ParseFiles(prefix+"encode.tmpl", prefix+"elem_enc.tmpl"))
	marTemplate = template.Must(template.ParseFiles(prefix + "marshal.tmpl"))
	unmTemplate = template.Must(template.ParseFiles(prefix + "unmarshal.tmpl"))
}

// WriteDecoderMethod writes the DecodeMsg(io.Reader) method.
// Ptr.Value should be a *Struct. If 'unmarshal' is true, the Unmarshal
// method will also be written.
func WriteDecoderMethod(w io.Writer, p *Ptr) error {
	var buf bytes.Buffer
	err := decTemplate.Execute(&buf, p)
	if err != nil {
		return fmt.Errorf("template: %s", err)
	}

	bts, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("gofmt: %s", err)
	}

	_, err = w.Write(bts)
	if err != nil {
		return err
	}

	return nil
}

// WriteEncoderMethod writes the EncodeMsg(io.Writer) method.
// Ptr.Value should be a *Struct. If 'marshal' is true, the Marshal
// method will also be written.
func WriteEncoderMethod(w io.Writer, p *Ptr) error {
	var buf bytes.Buffer
	err := encTemplate.Execute(&buf, p)
	if err != nil {
		return err
	}

	bts, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	_, err = w.Write(bts)
	if err != nil {
		return err
	}
	return nil
}
