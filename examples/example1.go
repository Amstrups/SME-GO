package examples

import t "SME-GO/types"

func FOO() Example {
	return Example{
		Function:      foo,
		Inputs: foo_inputs,
	}
}

func foo(ins t.ReadChannels, outs t.WriteChannels) {
	x := <-ins[t.SECRET]
	outs[t.PUBLIC] <- x
}

func foo_inputs(sIn, _ t.Chan) {
	sIn <- 12
}
