package progio_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"bitbucket.org/shu_go/gotwant"
	"bitbucket.org/shu_go/progio"
)

func TestReader(t *testing.T) {
	buf := bytes.Buffer{}
	buf.WriteString("hoge")
	r := progio.NewReader(&buf, func(p int64) {
		gotwant.TestExpr(t, p, p <= 4)
	})

	content, err := ioutil.ReadAll(r)
	gotwant.TestError(t, err, nil)
	gotwant.Test(t, string(content), "hoge")
}

func TestWriter(t *testing.T) {
	buf := bytes.Buffer{}
	r := progio.NewWriter(&buf, func(p int64) {
		gotwant.TestExpr(t, p, p <= 4)
	})
	r.Write([]byte("hoge"))

	gotwant.Test(t, buf.String(), "hoge")
}
