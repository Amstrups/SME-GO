package main

import (
	"SME-GO/examples"
	t "SME-GO/types"
	"sync"
	"testing"
)

func runRegular(ex examples.Example, N int) {

	for i := 0; i < N; i++ {
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
			ex.Function(ins, outs)
			close(sout)
			close(pout)
		}()

		go func() {
			defer wg.Done()
			ex.Inputs(sin, pin)
			close(sin)
			close(pin)
		}()

		go func() {
			defer wg.Done()
			for {
				_, ok := <-pout
				if !ok {
					return
				}
			}
		}()

		go func() {
			defer wg.Done()
			for {
				_, ok := <-sout
				if !ok {
					return
				}
			}
		}()

		wg.Wait()
	}
}

func runSME(ex examples.Example, N int) {

	for i := 0; i < N; i++ {
		wg := sync.WaitGroup{}
		c := t.CreateCluster()

		wg.Add(4)

		go func() {
			defer wg.Done()
			SME(ex.Function, c.INS, c.OUTS)
			c.CloseWrites()
		}()

		go func() {
			defer wg.Done()
			ex.Inputs(t.MakeWriteOnly(c.S_IN), t.MakeWriteOnly(c.P_IN))
			c.CloseReads()
		}()

		go func() {
			defer wg.Done()
			for {
				_, ok := <-c.P_OUT
				if !ok {
					return
				}
			}
		}()

		go func() {
			defer wg.Done()
			for {
				_, ok := <-c.S_OUT
				if !ok {
					return
				}
			}
		}()

		wg.Wait()
	}
}

// =============== Example 1: Foo ===============
func BenchmarkRegular_EX1(b *testing.B) {
	foo := examples.FOO()
	runRegular(foo, b.N)

}

func BenchmarkSME_EX1(b *testing.B) {
	foo := examples.FOO()
	runSME(foo, b.N)
}

// =============== Example 2: Baa ===============

func BenchmarkRegular_EX2(b *testing.B) {
	baa := examples.BAA()
	runRegular(baa, b.N)

}

func BenchmarkSME_EX2(b *testing.B) {
	baa := examples.BAA()
	runSME(baa, b.N)
}

// =============== Example 3: Zuu ===============
func BenchmarkRegular_EX3(b *testing.B) {
	zuu := examples.ZUU()

	runRegular(zuu, b.N)
}

func BenchmarkSME_EX3(b *testing.B) {
	zuu := examples.ZUU()
	runSME(zuu, b.N)
}

// =============== Native foo-function ===============

func foo(a, b int) (int, int) {
	x := a
	y := b

	return x, y

}

func BenchmarkSimpleFoo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		foo(12, 10)
	}
}
