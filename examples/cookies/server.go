// +build ignore

package main

import (
	"fmt"
	"time"

	"github.com/erikdubbelboer/fasthttp"
	"github.com/themester/fcookiejar"
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
	// IMPORTANT: If you are going to use cookiejar DO NOT release manually.
	// CookieJar package releases cookies automatically.
	//
	//defer fasthttp.ReleaseCookie(cookie)
	//defer fasthttp.ReleaseCookie(cookie1)
	//defer fasthttp.ReleaseCookie(cookie2)

	cookie.SetDomain("make.fasthttp.great.again")
	cookie.SetKey("key")
	cookie.SetExpire(time.Now().Add(time.Hour))
	cookie.SetPath("/")
	cookie.SetValue("value")

	cookie1.SetDomain("make.fasthttp.great.again")
	cookie1.SetKey("use")
	cookie1.SetPath("/path")
	cookie1.SetValue("fasthttp")

	cookie2.SetKey("hello")
	cookie2.SetValue("world")

	// You can use SetCookie calls
	/*
		ctx.Response.Header.SetCookie(cookie)
		ctx.Response.Header.SetCookie(cookie1)
		ctx.Response.Header.SetCookie(cookie2)
	*/

	// or if you want to enhance your
	// cookie control over multiple users with multiple cookies.
	cookies := cookiejar.AcquireCookieJar()
	defer cookiejar.ReleaseCookieJar(cookies)

	cookies.Put(cookie)
	cookies.Put(cookie1)
	cookies.Put(cookie2)

	// writting to the response header
	cookies.WriteToResponse(&ctx.Response)
}
