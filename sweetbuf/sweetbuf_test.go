package sweetbuf

import "testing"

func TestSweetBuf(t *testing.T) {
	buf := NewSweetBuf().SetFetchFunc(FetchHttpGet("http://news-at.zhihu.com/api/4/start-image/1080*1776"))
	if data, err := buf.GetByte(); err != nil {
		t.Error(err.Error())
	} else {
		t.Log(string(data))
	}
}
