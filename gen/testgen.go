package gen

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
)

func WriteTestNBench(w io.Writer, s *Struct) error {
	var buf bytes.Buffer

	err := benTemplate.Execute(&buf, s)
	if err != nil {
		return fmt.Errorf("error executing template: %s", err)
	}

	bts, err := format.Source(buf.Bytes())
	if err != nil {
		w.Write(bts)
		return fmt.Errorf("go/format: %s", err)
	}

	_, err = w.Write(bts)
	return err
}
