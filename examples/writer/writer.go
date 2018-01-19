package main

import (
	"os"

	"github.com/erikdubbelboer/fasthttp"
)

func main() {
	client := &fasthttp.Client{
		Name:   "Google, thanks for Golang, but we all hate Android.",
		Writer: os.Stdout,
	}
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	req.SetRequestURI(`https://google.es`)

	client.Do(req, res)
}
