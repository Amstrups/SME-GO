package types

import (
	"sync"
)

// ChannelCluster
type ChannelCluster struct {
	S_IN, P_IN   Chan
	S_OUT, P_OUT Chan
	INS          ReadChannels
	OUTS         WriteChannels
}

func CreateCluster() ChannelCluster {
	sIn := make(Chan)
	pIn := make(Chan)
	sOut := make(Chan)
	pOut := make(Chan)

	return ChannelCluster{
		S_IN:  sIn,
		P_IN:  pIn,
		S_OUT: sOut,
		P_OUT: pOut,
		INS:   ReadChannels{PUBLIC: MakeReadOnly(pIn), SECRET: MakeReadOnly(sIn)},
		OUTS:  WriteChannels{PUBLIC: MakeWriteOnly(pOut), SECRET: MakeWriteOnly(sOut)},
	}
}

func (c *ChannelCluster) GetChannels() (ReadChannels, WriteChannels) {
	return c.INS, c.OUTS

}

func (c *ChannelCluster) CloseWrites() {
	close(c.S_OUT)
	close(c.P_OUT)
}

func (c *ChannelCluster) CloseReads() {
	close(c.S_IN)
	close(c.P_IN)
}

func (c *ChannelCluster) ForwardTo(ch WriteChan, wg *sync.WaitGroup) {
	wg.Add(2)
	// read public outs
	go func() {
		defer wg.Done()
		for {
			v, ok := <-c.P_OUT
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
			v, ok := <-c.S_OUT
			if !ok {
				return
			}
			ch <- v
		}
	}()
}
