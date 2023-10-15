package conditions

import (
	"math/rand"
	"testing"
)

func TestRandomSession(t *testing.T) {
	r1 := SessionRandom()
	for i := 0; i < 500; i++ {
		r := SessionRandom()
		if r != r1 {
			t.Fatal("RandomSession not consistent")
		}
	}
}

func TestRandomKey(t *testing.T) {
	rs := RandomForKey("key1", 1)
	if rs != 292785326893130985 {
		t.Fatal("randForKey not stable")
	}
	rs = RandomForKey("key2", 2)
	if rs != 1378833688500478092 {
		t.Fatal("randForKey not stable")
	}
	rs = RandomForKey("key3", 3)
	if rs != 4152708728521114743 {
		t.Fatal("randForKey not stable")
	}
	for i := 0; i < 1_000_000; i++ {
		rs = RandomForKey("keyx", rand.Int63())
		if rs < 0 || rs > 1<<62 {
			t.Fatal("RandomForKey out of range")
		}
	}
}
