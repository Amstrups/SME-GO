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
				return
			}

			lowExec.P_IN <- v
			highExec.P_IN <- v
		}
	}()

	lowExec.ForwardTo(outs[t.SECRET], &wg)
	highExec.ForwardTo(outs[t.PUBLIC], &wg)

	wg.Wait()
}

// Main
func main() {
	var wg sync.WaitGroup
	cc := t.CreateCluster()

	ex := e.FOO()

	wg.Add(4)

	go func() {
		defer wg.Done()
		SME(ex.Function, cc.INS, cc.OUTS)
		cc.CloseWrites()
	}()

	go func() {
		defer wg.Done()
		ex.Inputs(cc.S_IN, cc.P_IN)
		cc.CloseReads()
	}()

	go func() {
		defer wg.Done()
		for {
			v, ok := <-cc.P_OUT
			if !ok {
				return
			}

			fmt.Println("Public:", v)
		}
	}()

	go func() {
		defer wg.Done()
		for {
			v, ok := <-cc.S_OUT
			if !ok {
				return
			}

			fmt.Println("Secret:", v)
		}
	}()

	wg.Wait()
}
