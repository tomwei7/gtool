package bigqueue

import (
	"fmt"
	"testing"
)

var bq *BigQueue

func init() {
	bq = NewBigQueue(2)
	bq.TempDir = "/Users/weicheng"
}

func TestBigQueue(t *testing.T) {
	for i := 0; i < 10; i++ {
		if err := bq.Put([]byte(fmt.Sprintf("I'am %d", i))); err != nil {
			t.Error(err.Error())
		}
	}
	for {
		v, err := bq.Get()
		if v == nil && err == nil {
			t.Log("End Of Queue")
			break
		} else if err != nil {
			t.Error(err.Error())
		} else {
			t.Log(string(v))
		}
	}
}

func BenchmarkBigQueue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := bq.Put([]byte(fmt.Sprintf("I'am %d", i))); err != nil {
			b.Error(err.Error())
		}
	}
}
