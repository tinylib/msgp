package _generated

import "testing"

func TestIssue461RoundTrip(t *testing.T) {
	bin := &Issue461Bin{V: "b"}
	app := &Issue461App{V: "a"}
	text := &Issue461Text{V: "t"}
	in := Issue461{
		BinPS:     &[]Issue461Bin{{V: "1"}, {V: "2"}},
		BinPSP:    &[]*Issue461Bin{{V: "3"}},
		BinPA:     &[2]Issue461Bin{{V: "4"}, {V: "5"}},
		BinPP:     &bin,
		BinPM:     &map[string]Issue461Bin{"k": {V: "6"}},
		App:       Issue461App{V: "7"},
		AppS:      []Issue461App{{V: "8"}},
		AppA:      [2]Issue461App{{V: "9"}, {V: "10"}},
		AppM:      map[string]Issue461App{"k": {V: "11"}},
		AppPS:     &[]Issue461App{{V: "12"}},
		AppPP:     &app,
		TextPS:    &[]Issue461Text{{V: "13"}},
		TextPP:    &text,
		TextApp:   Issue461TextApp{V: "14"},
		TextAppPS: &[]Issue461TextApp{{V: "15"}},
	}
	issue459RoundTrip(t, in, in)
	issue459RoundTrip(t, Issue461List{{V: "x"}, {V: "y"}}, Issue461List{{V: "x"}, {V: "y"}})
}
