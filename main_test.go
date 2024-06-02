package main

import (
	"SME-GO/examples"
	t "SME-GO/types"
	"sync"
	"testing"
)

func f(ins t.ReadChannels, outs t.WriteChannels) {
	sin := <-ins[t.SECRET]
	pin := <-ins[t.PUBLIC]

	sout := sin + 10
	pout := pin * 10

	outs[t.SECRET] <- sout
	outs[t.PUBLIC] <- pout
}

func input(sIn, pIn t.Chan) {
	sIn <- 12
	pIn <- 20
}

func BenchmarkRegularExec(b *testing.B) {

	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		sin := make(chan t.T)
		pin := make(chan t.T)
		sout := make(chan t.T)
		pout := make(chan t.T)
		ins := map[string]t.ReadChan{
			t.SECRET: sin,
			t.PUBLIC: pin,
		}
		outs := map[string]t.WriteChan{
			t.SECRET: sout,
			t.PUBLIC: pout,
		}

		wg.Add(4)
		go func() {
			defer wg.Done()
			f(ins, outs)
			close(sout)
			close(pout)
		}()

		go func() {
			defer wg.Done()
			input(sin, pin)
			close(sin)
			close(pin)
		}()

		go func() {
			defer wg.Done()
			for range sout {
			}
		}()

		go func() {
			defer wg.Done()
			for range pout {
			}
		}()

		wg.Wait()
	}
}

func BenchmarkSME(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		c := t.CreateCluster()

		wg.Add(4)

		go func() {
			defer wg.Done()
			SME(f, c.INS, c.OUTS)
			c.CloseWrites()
		}()

		go func() {
			defer wg.Done()
			input(c.S_IN, c.P_IN)
			c.CloseReads()
		}()

		go func() {
			defer wg.Done()
			for range c.P_OUT {
			}
		}()

		go func() {
			defer wg.Done()
			for range c.S_OUT {
			}
		}()

		wg.Wait()
	}
}

// =============== Longer test
func BenchmarkRegularExec_HEAVY(b *testing.B) {
	zuu := examples.ZUU()

	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		sin := make(chan t.T)
		pin := make(chan t.T)
		sout := make(chan t.T)
		pout := make(chan t.T)
		ins := map[string]t.ReadChan{
			t.SECRET: sin,
			t.PUBLIC: pin,
		}
		outs := map[string]t.WriteChan{
			t.SECRET: sout,
			t.PUBLIC: pout,
		}

		wg.Add(4)
		go func() {
			defer wg.Done()
			zuu.Function(ins, outs)
			close(sout)
			close(pout)
		}()

		go func() {
			defer wg.Done()
			zuu.Inputs(sin, pin)
			close(sin)
			close(pin)
		}()

		go func() {
			defer wg.Done()
			for range sout {
			}
		}()

		go func() {
			defer wg.Done()
			for range pout {
			}
		}()

		wg.Wait()
	}
}

func BenchmarkSME_HEAVY(b *testing.B) {
	zuu := examples.ZUU()

	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		c := t.CreateCluster()

		wg.Add(4)

		go func() {
			defer wg.Done()
			SME(zuu.Function, c.INS, c.OUTS)
			c.CloseWrites()
		}()

		go func() {
			defer wg.Done()
			zuu.Inputs(t.MakeWriteOnly(c.S_IN), t.MakeWriteOnly(c.P_IN))
			c.CloseReads()
		}()

		go func() {
			defer wg.Done()
			for {

				_, ok := <-c.P_OUT
				if !ok {
					break
				}

			}
		}()

		go func() {
			defer wg.Done()
			for range c.S_OUT {
			}
		}()

		wg.Wait()
	}
}
