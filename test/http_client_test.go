package test

import (
	"ginskeleton/app/global/variable"
	_ "ginskeleton/bootstrap"
	"github.com/qifengzhang007/goCurl"
	"testing"
)

func TestHttpClient(t *testing.T) {
	cli := goCurl.CreateHttpClient()
	if resp, err := cli.Get("http://hq.sinajs.cn/list=sh601360"); err == nil {
		content, err := resp.GetContents()
		if err != nil {
			t.Errorf("单元测试未通过,返回值不符合要求：%s\n", content)
		}
		t.Log(content)
	}
}

func TestPprof(t *testing.T) {
	cli := goCurl.CreateHttpClient()
	for i := 1; i <= 500; i++ {
		resp, err := cli.Get("http://127.0.0.1:20191/api/v1/home/news", goCurl.Options{
			FormParams: map[string]interface{}{
				"newsType": "portal",
				"page":     "2",
				"limit":    "52",
			},
		})
		if err == nil {
			if txt, err := resp.GetContents(); err == nil {
				if i == 500 {

					variable.ZapLog.Info(txt)
				}
			}
		} else {
			t.Log(err.Error())
		}
	}
}
