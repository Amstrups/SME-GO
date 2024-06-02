package main

import (
	e "SME-GO/examples"
	t "SME-GO/types"
	"fmt"
	"sync"
)

// SME-GO Function
func SME(f e.F, ins t.ReadChannels, outs t.WriteChannels) {
	highExec := t.CreateCluster()
	lowExec := t.CreateCluster()

	var wg sync.WaitGroup

	wg.Add(2)

	//low execution
	go func() {
		defer wg.Done()
		f(lowExec.GetChannels())
		lowExec.CloseWrites()
	}()

	//high execution
	go func() {
		defer wg.Done()
		f(highExec.GetChannels())
		highExec.CloseWrites()
	}()

	wg.Add(2)

	// forward (some) inputs in secret channel
	go func() {
		defer wg.Done()
		for {
			v, ok := <-ins[t.SECRET]

			if !ok {
				close(lowExec.S_IN)
				close(highExec.S_IN)
				return
			}

			lowExec.S_IN <- 0 // change to constant for low exec
			highExec.S_IN <- v
		}
	}()

	// forward inputs on public channel
	go func() {
		defer wg.Done()
		for {
			v, ok := <-ins[t.PUBLIC]
			if !ok {
				close(lowExec.P_IN)
				close(highExec.P_IN)
				return
			}

			lowExec.P_IN <- v
			highExec.P_IN <- v
		}
	}()

	// handle outputs from f
	go lowExec.Forward(outs[t.PUBLIC], lowExec.P_OUT, lowExec.S_OUT, &wg)
	go highExec.Forward(outs[t.SECRET], highExec.S_OUT, highExec.P_OUT, &wg)

	wg.Wait()
}

func Run(ex e.Example) {
	var wg sync.WaitGroup
	cc := t.CreateCluster()

	wg.Add(4)

	go func() {
		defer wg.Done()
		SME(ex.Function, cc.INS, cc.OUTS)
		cc.CloseWrites()
	}()

	go func() {
		defer wg.Done()
		ex.Inputs(t.MakeWriteOnly(cc.S_IN), t.MakeWriteOnly(cc.P_IN))
		cc.CloseReads()
	}()

	go func() {
		defer wg.Done()
		public_ints := []t.T{}
		for {
			v, ok := <-cc.P_OUT
			if !ok {
				fmt.Println("Publics:", public_ints)
				return
			}
			public_ints = append(public_ints, v)

		}
	}()

	go func() {
		secret_ints := []t.T{}
		defer wg.Done()

		for {
			v, ok := <-cc.S_OUT
			if !ok {
				fmt.Println("Secrets:", secret_ints)
				return
			}
			secret_ints = append(secret_ints, v)

		}
	}()

	wg.Wait()

}

func main() {
	Run(e.FOO())
}
