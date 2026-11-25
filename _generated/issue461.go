package _generated

//go:generate msgp

//msgp:binmarshal Issue461Bin
//msgp:binappend Issue461App
//msgp:textmarshal Issue461Text as:string
//msgp:textappend Issue461TextApp

// Marshalers have pointer receivers, so method calls need addressable values.
// Appenders have value receivers, as encode passes them as interface values.

type Issue461Bin struct{ V string }

func (t *Issue461Bin) MarshalBinary() ([]byte, error) { return []byte(t.V), nil }
func (t *Issue461Bin) UnmarshalBinary(b []byte) error { t.V = string(b); return nil }

type Issue461App struct{ V string }

func (t Issue461App) AppendBinary(b []byte) ([]byte, error) { return append(b, t.V...), nil }
func (t *Issue461App) UnmarshalBinary(b []byte) error       { t.V = string(b); return nil }

type Issue461Text struct{ V string }

func (t *Issue461Text) MarshalText() ([]byte, error) { return []byte(t.V), nil }
func (t *Issue461Text) UnmarshalText(b []byte) error { t.V = string(b); return nil }

type Issue461TextApp struct{ V string }

func (t Issue461TextApp) AppendText(b []byte) ([]byte, error) { return append(b, t.V...), nil }
func (t *Issue461TextApp) UnmarshalText(b []byte) error       { t.V = string(b); return nil }

type Issue461 struct {
	BinPS     *[]Issue461Bin
	BinPSP    *[]*Issue461Bin
	BinPA     *[2]Issue461Bin
	BinPP     **Issue461Bin
	BinPM     *map[string]Issue461Bin
	App       Issue461App
	AppS      []Issue461App
	AppA      [2]Issue461App
	AppM      map[string]Issue461App
	AppPS     *[]Issue461App
	AppPP     **Issue461App
	TextPS    *[]Issue461Text
	TextPP    **Issue461Text
	TextApp   Issue461TextApp
	TextAppPS *[]Issue461TextApp
}

// Issue461List gets methods on a (*z) receiver.
type Issue461List []Issue461App
