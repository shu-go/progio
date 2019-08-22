package progio_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/shu-go/gotwant"
	"github.com/shu-go/progio"
)

type buffer struct {
	//*os.File

	*bytes.Buffer
	size int64
}

func (b buffer) Close() {
	//b.File.Close()
}

func (b buffer) Size() int64 {
	/*
		stat, _ := b.File.Stat()
		return stat.Size()
	*/

	return b.size
}

func createBuffer() buffer {
	/*
		f, err := os.Open("./progio_test.go")
		if err != nil {
			panic(err)
		}
		return buffer{File: f}
	*/
	return buffer{
		Buffer: bytes.NewBufferString(strings.Repeat("hoge", 100)),
		size:   4 * 100,
	}
}

func TestReader(t *testing.T) {
	t.Run("Null", func(t *testing.T) {
		buf := createBuffer()
		defer buf.Close()

		r := progio.NewReader(buf, func(p int64) {
			gotwant.TestExpr(t, p, p <= buf.Size(), gotwant.Desc(fmt.Sprintf("size:%v", buf.Size())))
		})

		_, err := ioutil.ReadAll(r)
		gotwant.TestError(t, err, nil)
	})

	t.Run("Percent", func(t *testing.T) {
		buf := createBuffer()
		defer buf.Close()

		count := 0
		r := progio.NewReader(
			buf,
			func(p int64) {
				count++
				gotwant.TestExpr(t, p, 0 <= p)
				gotwant.TestExpr(t, p, p <= 100)
			},
			progio.Percent(buf.Size(), 1),
		)

		_, err := ioutil.ReadAll(r)
		gotwant.TestError(t, err, nil)

		gotwant.TestExpr(t, count, 0 < count)
	})

	t.Run("Time", func(t *testing.T) {
		buf := createBuffer()
		defer buf.Close()

		r := progio.NewReader(
			buf,
			func(p int64) {
				gotwant.TestExpr(t, p, p <= buf.Size())
			},
			progio.Time(5),
		)

		_, err := ioutil.ReadAll(r)
		gotwant.TestError(t, err, nil)
	})

	t.Run("TimeOver", func(t *testing.T) {
		buf := createBuffer()
		defer buf.Close()

		count := 0
		r := progio.NewReader(
			buf,
			func(p int64) {
				count++
				if count > 1 {
					t.Fail()
				}
			},
			progio.Time(5*time.Second),
		)

		_, err := ioutil.ReadAll(r)
		gotwant.TestError(t, err, nil)

		gotwant.TestExpr(t, count, count <= 1)
	})

	t.Run("ListenDuration", func(t *testing.T) {
		t.Parallel()

		buf := createBuffer()
		defer buf.Close()

		count := 0
		r := progio.NewReader(
			buf,
			func(p int64, d time.Duration) {
				time.Sleep(1 * time.Second)
				gotwant.TestExpr(t, d, d >= time.Duration(count)*time.Second)
				count++
			},
		)

		_, err := ioutil.ReadAll(r)
		gotwant.TestError(t, err, nil)
	})

	t.Run("ListenDurationPercent", func(t *testing.T) {
		t.Parallel()

		buf := createBuffer()
		defer buf.Close()

		count := 0
		r := progio.NewReader(
			buf,
			func(p int64, d time.Duration) {
				time.Sleep(1 * time.Second)
				gotwant.TestExpr(t, d, d >= time.Duration(count)*time.Second)
				count++
			},
			progio.Percent(buf.Size(), 1),
		)

		_, err := ioutil.ReadAll(r)
		gotwant.TestError(t, err, nil)
	})

	t.Run("ListenDurationTime", func(t *testing.T) {
		t.Parallel()

		buf := createBuffer()
		defer buf.Close()

		count := 0
		r := progio.NewReader(
			buf,
			func(p int64, d time.Duration) {
				time.Sleep(1 * time.Second)
				gotwant.TestExpr(t, d, d >= time.Duration(count)*time.Second)
				count++
			},
			progio.Time(5),
		)

		_, err := ioutil.ReadAll(r)
		gotwant.TestError(t, err, nil)
	})
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
