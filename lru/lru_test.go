// Author: catlittlechen@gmail.com

package lru

import (
	"strconv"
	"testing"
)

func TestLRU(t *testing.T) {
	count := 10
	lru, err := New(count, nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	count = count * count
	for i := 0; i < count; i++ {
		str := strconv.Itoa(i)
		lru.Set(str, str)
		lru.Print()
	}

	list := []string{"99", "98", "90"}
	for i := 0; i < len(list); i++ {
		lru.Set(list[i], list[i])
		lru.Print()
	}
	list = []string{"99", "98", "90", "33"}
	for i := 0; i < len(list); i++ {
		lru.Del(list[i])
		lru.Print()
	}
	if lru.left != 3 {
		t.Fatalf("wrong left is %s", lru.left)
	}
	for i := 0; i < len(list); i++ {
		lru.Set(list[i], list[i])
		lru.Print()
	}

	if lru.left != 0 {
		t.Fatalf("wrong left is %s", lru.left)
	}

}
