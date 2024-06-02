package main

import (
	"fmt"
	"sync"
)

type T int

type lvChan chan T
type rlvChan <-chan T
type wlvChan chan<- T
type ReadChannels map[string]rlvChan
type WriteChannels map[string]wlvChan
type F func(in ReadChannels, out WriteChannels)

const (
	PUBLIC = "PUBLIC"
	SECRET = "SECRET"
)

// utils
func MakeReadOnly(ch chan T) <-chan T {
	return ch
}

func MakeWriteOnly(ch chan T) chan<- T {
	return ch
}

// ChannelCluster
type ChannelCluster struct {
	sIn, pIn   lvChan
	sOut, pOut lvChan
	ins        ReadChannels
	outs       WriteChannels
}

func CreateCluster() ChannelCluster {
	sIn := make(lvChan)
	pIn := make(lvChan)
	sOut := make(lvChan)
	pOut := make(lvChan)

	return ChannelCluster{
		sIn:  sIn,
		pIn:  pIn,
		sOut: sOut,
		pOut: pOut,
		ins:  ReadChannels{PUBLIC: MakeReadOnly(pIn), SECRET: MakeReadOnly(sIn)},
		outs: WriteChannels{PUBLIC: MakeWriteOnly(pOut), SECRET: MakeWriteOnly(sOut)},
	}
}

func (c *ChannelCluster) CloseWrites() {
	close(c.sOut)
	close(c.pOut)
}

func (c *ChannelCluster) CloseReads() {
	close(c.sIn)
	close(c.pIn)
}

func (c *ChannelCluster) ForwardTo(ch wlvChan, wg *sync.WaitGroup) {
	wg.Add(2)
	// read public outs
	go func() {
		defer wg.Done()
		for {
			v, ok := <-c.pOut
			if !ok {
				return
			}
			ch <- v
		}
	}()

	// read secret outs
	go func() {
		defer wg.Done()
		for {
			v, ok := <-c.sOut
			if !ok {
				return
			}
			ch <- v
		}
	}()
}

// Channel helpers
func BuildChannels(sIn, pIn, sOut, pOut lvChan) (ReadChannels, WriteChannels) {
	ins := ReadChannels{PUBLIC: MakeReadOnly(pIn), SECRET: MakeReadOnly(sIn)}
	outs := WriteChannels{PUBLIC: MakeWriteOnly(pOut), SECRET: MakeWriteOnly(sOut)}

	return ins, outs
}

// SME-GO Function
func SME(f F, ins ReadChannels, outs WriteChannels) {
	highCluster := CreateCluster()
	lowCluster := CreateCluster()

	var wg sync.WaitGroup

	wg.Add(2)
	//low execution
	go func() {
		defer wg.Done()
		f(lowCluster.ins, lowCluster.outs)
		lowCluster.CloseWrites()
	}()

	//high execution
	go func() {
		defer wg.Done()
		f(highCluster.ins, highCluster.outs)
		highCluster.CloseWrites()
	}()

	wg.Add(2)
	// forward (some) inputs in secret channel
	go func() {
		defer wg.Done()
		for {
			v, ok := <-ins[SECRET]
			if !ok {
				return
			}

			highCluster.sIn <- 0 // changed to constant
			lowCluster.sIn <- v
		}
	}()

	// forward inputs on public channel
	go func() {
		defer wg.Done()
		for {
			v, ok := <-ins[PUBLIC]
			if !ok {
				return
			}

			highCluster.pIn <- v
			lowCluster.pIn <- v
		}
	}()

	lowCluster.ForwardTo(outs[SECRET], &wg)
	highCluster.ForwardTo(outs[PUBLIC], &wg)

	wg.Wait()
}

// Main
func main() {
	var wg sync.WaitGroup
	cc := CreateCluster()

	wg.Add(4)

	go func() {
		defer wg.Done()
		SME(baa, cc.ins, cc.outs)
		cc.CloseWrites()
	}()

	go func() {
		defer wg.Done()
		baa_inputs(cc.sIn, cc.pIn)
		cc.CloseReads()
	}()

	go func() {
		defer wg.Done()
		for {
			v, ok := <-cc.pOut
			if !ok {
				return
			}

			fmt.Println("Public:", v)
		}
	}()

	go func() {
		defer wg.Done()
		for {
			v, ok := <-cc.sOut
			if !ok {
				return
			}

			fmt.Println("Secret:", v)
		}
	}()

	wg.Wait()
}

// Testable function(s)
func foo(ins ReadChannels, outs WriteChannels) {
	x := <-ins[SECRET]
	fmt.Println(x)
	outs[PUBLIC] <- x
}

func foo_inputs(sIn, _ lvChan) {
	fmt.Println("writing 12 to sec")
	sIn <- 12
}

func baa(ins ReadChannels, outs WriteChannels) {
	sin := <-ins[SECRET]
	pin := <-ins[PUBLIC]
	sout := sin + pin
	pout := sin * pin
	fmt.Println(sin, pin)

	outs[PUBLIC] <- pout
	outs[SECRET] <- sout
}

func baa_inputs(sIn, pIn lvChan) {
	fmt.Println("writing 12 to sec")
	pIn <- 20
	sIn <- 12
}
