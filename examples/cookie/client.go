package main

import (
	"fmt"

	"github.com/erikdubbelboer/fasthttp"
)

func main() {
	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	req.SetRequestURI("http://localhost:8080")

	err := fasthttp.Do(req, res)
	if err != nil {
		panic(err)
	}

	// We can get the cookie from diferents ways.
	// 1. Using 'key' parameter. Fasthttp will search using key in cookie field.
	cookie := fasthttp.AcquireCookie()
	cookie.SetKey("key")

	res.Header.Cookie(cookie)
	fmt.Printf("You cookie is '%s' found using '%s' key\n", cookie.Cookie(), cookie.Key())

	cookie.Reset()

	// 2. You can parse request header 'Set-Cookie' field.
	value := res.Header.Peek("Set-Cookie")
	if value != nil {
		cookie.ParseBytes(value)
		fmt.Printf("Your cookie is '%s' parsing bytes\n", cookie.Cookie())
	}

	// 3. Or you can use 'PeekCookie' from ResponseHeader using key.
	// This is the long way to get cookie instead of using 'Cookie' structure.
	value = res.Header.PeekCookie("key")
	if value != nil {
		fmt.Printf("Your cookie is '%s' parsing by key\n", value)
	}
}
