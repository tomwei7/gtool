// 这其实并不是一个队列
package bigqueue

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	//"runtime"
	"sync"
)

type BigQueue struct {
	mcap       int
	mlen       int
	m          sync.Mutex
	c          chan []byte
	fqueue     *os.File
	fqueueHead int64
	fqueueTail int64
	TempDir    string
}

func (p *BigQueue) Get() (v []byte, err error) {
	if len(p.c) > 0 {
		return <-p.c, nil
	}
	//for {
	//	v, err = p.fGet()
	//	if v == nil && err == nil {
	//		runtime.Gosched()
	//		continue
	//	}
	//	break
	//}
	return p.fGet()
}

func (p *BigQueue) fGet() (v []byte, err error) {
	p.m.Lock()
	defer p.m.Unlock()
	if p.fqueueHead >= p.fqueueTail {
		if err = p.fqueue.Close(); err != nil {
			return
		}
		if err = os.Remove(p.fqueue.Name()); err == nil {
			p.fqueueHead = 0
			p.fqueueTail = 0
			p.fqueue = nil
		}
		return
	}
	if _, err = p.fqueue.Seek(p.fqueueHead, 0); err != nil {
		return
	}
	var dlen int64
	if err = binary.Read(p.fqueue, binary.LittleEndian, &dlen); err != nil {
		return
	}
	v = make([]byte, int(dlen))
	if n, err := p.fqueue.Read(v); err != nil {
		return nil, err
	} else if n != int(dlen) {
		return nil, errors.New(fmt.Sprintf("Data Length Error Expect %d Read %d", dlen, n))
	}
	p.fqueueHead = p.fqueueHead + dlen + 8
	return
}

func (p *BigQueue) Put(v []byte) (err error) {
	if p.mlen < p.mcap {
		p.c <- v
		p.mlen++
		return
	}
	return p.fPut(v)
}

func (p *BigQueue) fPut(v []byte) (err error) {
	p.m.Lock()
	defer p.m.Unlock()
	if p.fqueue == nil {
		if p.fqueue, err = ioutil.TempFile(p.TempDir, "bigqueue_"); err != nil {
			return
		}
	}
	if _, err = p.fqueue.Seek(p.fqueueTail, 0); err != nil {
		return
	}
	dlen := int64(len(v))
	if err = binary.Write(p.fqueue, binary.LittleEndian, dlen); err != nil {
		return
	}
	if _, err = p.fqueue.Write(v); err != nil {
		return
	}
	// int64 为8个字节
	p.fqueueTail = p.fqueueTail + dlen + 8
	return
}

func (p *BigQueue) init() {
	p.c = make(chan []byte, p.mcap)
}

// mlen memory queue len
func NewBigQueue(mcap int) *BigQueue {
	bq := &BigQueue{mcap: mcap}
	bq.init()
	return bq
}
