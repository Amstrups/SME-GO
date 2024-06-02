package examples

import (
	t "SME-GO/types"
)

func FOO() Example {
	return Example{
		Function: foo,
		Inputs:   foo_inputs,
	}
}

func foo(ins t.ReadChannels, outs t.WriteChannels) {
	x := <-ins[t.SECRET]
	y := <-ins[t.PUBLIC]
	outs[t.SECRET] <- y
	outs[t.PUBLIC] <- x
}

func foo_inputs(sIn, pIn t.WriteChan) {
	sIn <- 12
	pIn <- 11
}
