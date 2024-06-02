package examples

import t "SME-GO/types"

type F func(in t.ReadChannels, out t.WriteChannels)
type Inputs func(in t.Chan, out t.Chan)

type Example struct {
	Function     F
	Inputs Inputs
}
