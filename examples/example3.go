package examples

import (
	"SME-GO/types"
)

func ZUU() Example {
	return Example{
		Function: zuu,
		Inputs:   zuu_inputs,
	}
}

func zuu(ins types.ReadChannels, outs types.WriteChannels) {
	var acc types.T

reads:
	for {
		v, ok := <-ins[types.PUBLIC]
		if !ok {
			break reads
		}
		acc += v * v
	}

	for i := 0; i < 100; i++ {
		outs[types.PUBLIC] <- acc
		acc += 1

	}
	outs[types.SECRET] <- acc

}

func zuu_inputs(S_IN, P_IN types.WriteChan) {
	for j := 0; j < 100; j++ {
		i := types.T(j)
		P_IN <- i * i
	}
}
