package sweetbuf

import (
	"errors"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type SweetBuf struct {
	mx           sync.Mutex
	expireIn     int64 //有效期
	expired      int64 //过期时间
	updateTime   int64 //更新时间
	fetchFunc    func() ([]byte, error)
	decorateFunc func([]byte) (interface{}, error)
	cacheByte    []byte
	cacheValue   interface{}
}

// 设置刷新数据函数
func (p *SweetBuf) SetFetchFunc(f func() ([]byte, error)) *SweetBuf {
	p.fetchFunc = f
	return p
}

// 设置数据decode函数
func (p *SweetBuf) SetDecorateFunc(f func([]byte) (interface{}, error)) *SweetBuf {
	p.decorateFunc = f
	return p
}

func (p *SweetBuf) refreshCache() error {
	if p.fetchFunc == nil {
		return errors.New("FetchFunc Not Set")
	}
	if data, err := p.fetchFunc(); err != nil {
		return err
	} else {
		p.cacheByte = data
		p.updateTime = time.Now().Unix()
		p.expired = p.updateTime + p.expireIn
		if p.decorateFunc != nil {
			// ignore decorate error
			p.cacheValue, _ = p.decorateFunc(p.cacheByte)
		}
		return nil
	}
}

func (p *SweetBuf) GetByte() ([]byte, error) {
	if p.expired > time.Now().Unix() {
		return p.cacheByte, nil
	}
	// Expired
	p.mx.Lock()
	defer p.mx.Unlock()
	if p.expired > time.Now().Unix() {
		return p.cacheByte, nil
	}
	if err := p.refreshCache(); err != nil {
		return nil, err
	} else {
		return p.cacheByte, nil
	}
}

func (p *SweetBuf) GetValue() (interface{}, error) {
	if p.decorateFunc == nil {
		return nil, errors.New("DecorateFunc Not Set")
	}
	if p.cacheValue != nil {
		return p.cacheValue, nil
	}
	if v, err := p.decorateFunc(p.cacheByte); err != nil {
		return nil, err
	} else {
		p.cacheValue = v
		return v, nil
	}
}
func NewSweetBuf() *SweetBuf {
	return &SweetBuf{}
}

func FetchHttpGet(url string) func() ([]byte, error) {
	return func() ([]byte, error) {
		if r, err := http.Get(url); err != nil {
			return nil, err
		} else {
			defer r.Body.Close()
			return ioutil.ReadAll(r.Body)
		}
	}
}
