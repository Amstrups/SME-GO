package examples

import (
	t "SME-GO/types"
)

func BAA() Example {
	return Example{
		Function: baa,
		Inputs:   baa_inputs,
	}
}

func baa(ins t.ReadChannels, outs t.WriteChannels) {
	sin := <-ins[t.SECRET]
	pin := <-ins[t.PUBLIC]
	sout := sin + pin
	pout := sin * pin

	outs[t.PUBLIC] <- pout
	outs[t.SECRET] <- sout
}

func baa_inputs(sIn, pIn t.WriteChan) {
	sIn <- 12
	pIn <- 20
}
