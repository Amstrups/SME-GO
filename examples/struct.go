package examples

import t "SME-GO/types"

type F func(in t.ReadChannels, out t.WriteChannels)
type Inputs func(SECRET_INPUT t.WriteChan, PUBLIC_INPUT t.WriteChan)

type Example struct {
	Function F
	Inputs   Inputs
}
