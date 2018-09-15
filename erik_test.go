package fasthttp

import (
	"bytes"
	"net"
	//"os"
	//"runtime"
	"testing"

	//"github.com/erikdubbelboer/bpprof"
	"github.com/valyala/fasthttp/fasthttputil"
)

func BenchmarkErikHttpSimple(b *testing.B) {
	//runtime.MemProfileRate = 1

	body := []byte("test")
	ln := fasthttputil.NewInmemoryListener()

	go Serve(ln, func(ctx *RequestCtx) {
		ctx.SetStatusCode(200)
		ctx.SetBody(body)
	})

	c := Client{
		Dial: func(addr string) (net.Conn, error) {
			return ln.Dial()
		},
	}

	var reuseBody []byte
	for i := 0; i < b.N; i++ {
		_, responseBody, err := c.Get(reuseBody, "http://example.com")
		if err != nil {
			panic(err)
		}
		reuseBody = responseBody
	}

	//bpprof.Heap(os.Stderr, "allocobjects")
}

func BenchmarkErikHttpBodyReader(b *testing.B) {
	//runtime.MemProfileRate = 1

	requestBody := []byte("request")
	responseBody := []byte("response")
	ln := fasthttputil.NewInmemoryListener()

	go Serve(ln, func(ctx *RequestCtx) {
		if !bytes.Equal(requestBody, ctx.PostBody()) {
			panic("request body is wrong")
		}
		ctx.SetStatusCode(200)
		ctx.SetBody(responseBody)
	})

	c := Client{
		Dial: func(addr string) (net.Conn, error) {
			return ln.Dial()
		},
	}

	var body bytes.Buffer

	requestURI := []byte("http://example.com")
	requestMethod := []byte("POST")

	for i := 0; i < b.N; i++ {
		req := AcquireRequest()
		res := AcquireResponse()

		req.SetRequestURIBytes(requestURI)
		req.SetBody(requestBody)
		req.Header.SetMethodBytes(requestMethod)

		res.ReturnBodyReader = true

		if err := c.Do(req, res); err != nil {
			panic(err)
		}

		body.Reset()
		if _, err := body.ReadFrom(res.BodyReader); err != nil {
			panic(err)
		}

		if !bytes.Equal(responseBody, body.Bytes()) {
			panic("response body is wrong")
		}

		res.BodyReader.Close()
		ReleaseBodyReader(res.BodyReader)
		ReleaseRequest(req)
		ReleaseResponse(res)
	}

	//bpprof.Heap(os.Stderr, "allocobjects")
}
