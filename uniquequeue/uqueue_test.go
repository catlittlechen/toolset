// Author: catlittlechen@gmail.com

package uqueue

import (
	"strconv"
	"testing"
)

type Str string

func (s Str) UniuqID() string {
	return string(s)
}

func TestUQueue(t *testing.T) {
	count := 10
	uq := New()
	for i := 0; i < count; i++ {
		s := strconv.Itoa(i)
		uq.Push(Str(s))
		uq.Print()
	}
	str := uq.Pop()
	for str != nil {
		uq.Print()
		str = uq.Pop()
	}

}
