// +build ignore

package main

import (
	"fmt"
	"time"

	"github.com/erikdubbelboer/fasthttp"
)

func main() {
	fmt.Println(
		fasthttp.ListenAndServe(":8080", cookieHandler),
	)
}

func cookieHandler(ctx *fasthttp.RequestCtx) {
	cookie := fasthttp.AcquireCookie()
	cookie1 := fasthttp.AcquireCookie()
	cookie2 := fasthttp.AcquireCookie()
	// do not forget to release
	defer fasthttp.ReleaseCookie(cookie)
	defer fasthttp.ReleaseCookie(cookie1)
	defer fasthttp.ReleaseCookie(cookie2)

	cookie.SetDomain("make.fasthttp.great.again")
	cookie.SetKey("key")
	cookie.SetExpire(time.Now().Add(time.Hour))
	cookie.SetPath("/")
	cookie.SetValue("value")

	cookie1.SetDomain("make.fasthttp.great.again")
	cookie1.SetKey("use")
	cookie1.SetPath("/path")
	cookie1.SetValue("fasthttp")

	cookie2.SetKey("value")
	cookie2.SetValue("key")

	ctx.Response.Header.SetCookie(cookie)
	ctx.Response.Header.SetCookie(cookie1)
	ctx.Response.Header.SetCookie(cookie2)
}
