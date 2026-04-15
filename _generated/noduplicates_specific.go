package _generated

//go:generate msgp

//msgp:noduplicates NoDupTargeted NoDupTargetedMap

// NoDupTargeted has the directive applied - should reject duplicates.
type NoDupTargeted struct {
	Alpha string `msg:"alpha"`
	Beta  int    `msg:"beta"`
	Gamma bool   `msg:"gamma"`
}

// NoDupTargetedMap has the directive applied as a map type.
type NoDupTargetedMap map[string]int

// NoDupUntargeted does NOT have the directive - duplicates should be silently accepted.
type NoDupUntargeted struct {
	Alpha string `msg:"alpha"`
	Beta  int    `msg:"beta"`
}

// NoDupUntargetedMap does NOT have the directive - duplicates silently accepted.
type NoDupUntargetedMap map[string]string
