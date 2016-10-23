package cfg

import (
	"flag"
)

type MsgpConfig struct {
	Out        string
	GoFile     string
	Encode     bool
	Marshal    bool
	Tests      bool
	Unexported bool
}

// call DefineFlags before myflags.Parse()
func (c *MsgpConfig) DefineFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.Out, "o", "", "output file")
	fs.StringVar(&c.GoFile, "file", "", "input file")
	fs.BoolVar(&c.Encode, "io", true, "create Encode and Decode methods")
	fs.BoolVar(&c.Marshal, "marshal", true, "create Marshal and Unmarshal methods")
	fs.BoolVar(&c.Tests, "tests", true, "create tests and benchmarks")
	fs.BoolVar(&c.Unexported, "unexported", false, "also process unexported types")
}

// call c.ValidateConfig() after myflags.Parse()
func (c *MsgpConfig) ValidateConfig() error {
	return nil
}
