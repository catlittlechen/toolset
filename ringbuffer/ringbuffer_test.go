package ringbuffer

import (
	"math/rand"
	"testing"
	"time"
)

type ConsumerMessage struct {
	SeqID uint64
}

func (cm *ConsumerMessage) GetSeqID() uint64 {
	return cm.SeqID
}

func TestRing(t *testing.T) {
	r := NewRingBuffer(1000, MinSeqID, time.Second)

	go func() {
		for {
			data := <-r.channel
			t.Log(data.GetSeqID())
		}
	}()

	count := 3
	chans := make([]chan *ConsumerMessage, count)
	real := make([]chan *ConsumerMessage, count)
	for i := 0; i < count; i++ {
		chans[i] = make(chan *ConsumerMessage, 1024)
		real[i] = make(chan *ConsumerMessage, 1024)
		go func(j int) {
			for cm := range real[j] {
				r.Set(cm)
			}
		}(i)
	}

	go func() {
		for i := MinSeqID; i < 1000000; i++ {
			time.Sleep(time.Microsecond)
			if i > MaxSeqID {
				i = MinSeqID
			}
			random := rand.Int() % count
			chans[random] <- &ConsumerMessage{
				SeqID: uint64(i),
			}
		}
	}()

	go func() {
		for {
			random := rand.Int() % count
			select {
			case i := <-chans[0]:
				real[random] <- i
			case i := <-chans[1]:
				real[random] <- i
			case i := <-chans[2]:
				real[random] <- i
			}
		}
	}()

	time.Sleep(120 * time.Second)
}
