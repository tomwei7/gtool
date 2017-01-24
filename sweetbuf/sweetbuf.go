package sweetbuf

import (
	"sync"
)

type SweetBuf struct {
	mx           sync.Mutex
	expireIn     int64
	updateTime   int64
	fetchFunc    func() ([]byte, error)
	decorateFunc func() (interface{}, error)
	cacheByte    []byte
	cacheType    interface{}
}

func (p *SweetBuf) SetFetchFunc(f func() ([]byte, error)) *SweetBuf {
	p.fetchFunc = f
	return p
}

func (p *SweetBuf) SetDecorateFunc(f func() (interface{}, error)) *SweetBuf {
	p.decorateFunc = f
	return p
}

func (p *SweetBuf) refreshCache() error {
}

func (p *SweetBuf) Get() ([]byte, error) {
}

func (p *SweetBuf) GetType() (interface{}, error) {
}
