package util

import (
	"fmt"
	"testing"
)

func TestUtil(t *testing.T) {
	pool := ThreadPool{}
	pool.HandleFunc(func(i int, s string) {
		fmt.Println(i)
		fmt.Println(s)
	}, 1, "hello")
}

func TestUtil2(t *testing.T) {
	type Key struct {
		Sequence string
		Quality  string
	}
	key1 := Key{
		Sequence: "A",
		Quality:  "1",
	}
	key2 := Key{
		Sequence: "A",
		Quality:  "2",
	}
	m := make(map[Key]int)
	m[key1] = 1
	m[key2] = 2
	fmt.Println(m)
	key1 = Key{
		Sequence: "A",
		Quality:  "3",
	}
	fmt.Println(m)
}
