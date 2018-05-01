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

	cookie.SetDomain("make.fasthttp.great.again")
	cookie.SetKey("key")
	cookie.SetExpire(time.Now().Add(time.Hour))
	cookie.SetPath("/")
	cookie.SetValue("value")

	ctx.Response.Header.SetCookie(cookie)
}
