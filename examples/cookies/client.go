// +build ignore

package main

import (
	"fmt"

	"github.com/erikdubbelboer/fasthttp"
	"github.com/themester/fcookiejar"
)

func main() {
	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	req.SetRequestURI("http://localhost:8080")

	err := fasthttp.Do(req, res)
	if err != nil {
		panic(err)
	}

	// We can achieve cookies by different ways.
	// 1. Using 'key' parameter. Fasthttp will search using key in cookie field.
	cookie := fasthttp.AcquireCookie()
	// do not forget to release
	defer fasthttp.ReleaseCookie(cookie)
	cookie.SetKey("key")

	res.Header.Cookie(cookie)
	fmt.Printf("You cookie is '%s' found using '%s' key\n", cookie.Cookie(), cookie.Key())

	cookie.Reset()

	// 2. You can parse request header 'Set-Cookie' field.
	// But with this method you just can get the first cookie value.
	value := res.Header.Peek("Set-Cookie")
	if value != nil {
		cookie.ParseBytes(value)
		fmt.Printf("Your cookie is '%s' parsing bytes\n", cookie.Cookie())
	}

	// 3. You can use 'PeekCookie' from ResponseHeader using key.
	// This is the long way to get cookie instead of using 'Cookie' structure.
	value = res.Header.PeekCookie("key")
	if value != nil {
		fmt.Printf("Your cookie is '%s' parsing by key\n", value)
	}

	// 4. Or finally you can use CookieJar object.
	cookies := cookiejar.AcquireCookieJar()
	// Do not forget to release
	defer cookiejar.ReleaseCookieJar(cookies)

	// With this object you can collect cookies just executing:
	cookies.ResponseCookies(res)

	fmt.Printf("\nYour cookies:\n")
	for {
		cookie := cookies.Get()
		if cookie == nil {
			break
		}

		fmt.Printf("%s\n", cookie.Value())
	}
}
