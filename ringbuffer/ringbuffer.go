/*
RingBuffer for sort
*/

package ringbuffer

import (
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	// MinSeqID start of ring
	MinSeqID = 1024
	// MaxSeqID end of ring
	MaxSeqID = 9223372036854775
)

// Message of ring
type Message interface {
	GetSeqID() uint64
}

// RingBuffer .
type RingBuffer struct {
	length int
	sleep  time.Duration

	buffer []Message

	now    int
	seqID  uint64 //min seqID
	batter uint64

	channel chan Message

	closed     bool
	runChannel chan bool
}

// NewRingBuffer new()
func NewRingBuffer(length int, seqID uint64, sleep time.Duration) *RingBuffer {
	if seqID < MinSeqID {
		seqID = MinSeqID
	}
	r := &RingBuffer{
		length:     length,
		sleep:      sleep,
		buffer:     make([]Message, length),
		now:        0,
		seqID:      seqID,
		batter:     seqID + uint64(length),
		channel:    make(chan Message, 10),
		runChannel: make(chan bool, 1),
	}
	if r.batter > MaxSeqID {
		r.batter = r.batter - MaxSeqID + MinSeqID - 1
	}
	go r.Run()
	return r
}

// Run start roll the ring and get the msg
func (r *RingBuffer) Run() {
	defer func() {
		r.runChannel <- true
	}()
	for {
		for {
			if r.closed {
				return
			}
			if r.buffer[r.now] != nil {
				if r.buffer[r.now].GetSeqID() == r.seqID {
					log.Debugf("get data index:%d expectSeqID:%d", r.now, r.seqID)
					break
				}
				log.Errorf("ring bug! index:%d seqID:%d expectSeqID:%d", r.now, r.buffer[r.now].GetSeqID(), r.seqID)
			}
			log.Debugf("waiting for data... index:%d expectSeqID:%d", r.now, r.seqID)
			time.Sleep(r.sleep)
		}

		r.channel <- r.buffer[r.now]
		r.buffer[r.now] = nil

		r.now++
		if r.now == r.length {
			r.now = 0
			r.batter = r.batter + uint64(r.length)
			if r.batter > MaxSeqID {
				r.batter = r.batter - MaxSeqID + MinSeqID - 1
			}
		}

		r.seqID++
		if r.seqID > MaxSeqID {
			r.seqID = MinSeqID
		}
	}
}

// Set set the msg into ring
func (r *RingBuffer) Set(cm Message) {
	log.Debugf("batter %d now %d binlog:%d ringLen:%d min: %d max: %d", r.batter, r.now, cm.GetSeqID(), r.length, MinSeqID, MaxSeqID)

	var (
		seqID  uint64
		batter uint64
	)
	for {
		if r.closed {
			return
		}
		// drop invalid ID
		seqID = cm.GetSeqID()
		if seqID < r.seqID && MaxSeqID-MinSeqID+1 > 2*(r.seqID-seqID) {
			log.Warnf("duplicate %d seqID:%d min:%d max:%d", seqID, r.seqID, MinSeqID, MaxSeqID)
			break
		}

		// compress
		batter = r.batter
		if batter <= r.seqID {
			batter = batter + MaxSeqID - MinSeqID + 1
		}
		if seqID < r.seqID {
			seqID = seqID + MaxSeqID - MinSeqID + 1
		}

		if seqID < batter {
			r.buffer[r.length-int(batter-seqID)] = cm
			log.Debugf("binlog msg: seqID:%d pos:%d", cm.GetSeqID(), r.length-int(batter-seqID))
			break
		}
		if seqID < r.seqID+uint64(r.length) {
			r.buffer[seqID-batter] = cm
			log.Debugf("binlog msg: seqID:%d pos:%d", cm.GetSeqID(), seqID-batter)
			break
		}
		time.Sleep(r.sleep)
	}
}

// Close .
func (r *RingBuffer) Close() {
	if r.closed {
		return
	}
	r.closed = true
	<-r.runChannel
	close(r.channel)
}
