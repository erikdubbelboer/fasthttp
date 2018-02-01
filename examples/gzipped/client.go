package main

import (
	"bytes"
	"fmt"

	"github.com/erikdubbelboer/fasthttp"
)

func main() {
	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()

	req.Header.Add("Accept-Encoding", "gzip")
	req.SetRequestURI("http://localhost:1313")

	err := fasthttp.Do(req, res)
	if err != nil {
		panic(err)
	}
	body := res.Body()
	if b := res.Header.Peek("Content-Encoding"); len(b) > 0 {
		if bytes.Index(b, []byte("gzip")) >= 0 {
			body, err = res.BodyGunzip()
			if err != nil {
				panic(err)
			}
		}
	}
	fmt.Printf("%s\n", body)
}
