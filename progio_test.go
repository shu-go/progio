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

func BenchmarkRead(b *testing.B) {
	org := make([]byte, 0, 10*10000)
	for i := 0; i < 10000; i++ {
		org = append(org, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}...)
	}
	b.ResetTimer()

	b.Run("std", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r := bytes.NewBuffer(org)
			ioutil.ReadAll(r)
		}
	})

	b.Run("NullThrottling", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r := progio.NewReader(bytes.NewBuffer(org), func(p int64) {
			})
			ioutil.ReadAll(r)
		}
	})

	b.Run("Percent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r := progio.NewReader(bytes.NewBuffer(org), func(p int64) {
				//rog.Print(p)
			}, progio.Percent(10*10000, 5))
			ioutil.ReadAll(r)
		}
	})

	b.Run("Time", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r := progio.NewReader(bytes.NewBuffer(org), func(p int64) {
			}, progio.Time(10))
			ioutil.ReadAll(r)
		}
	})
}
