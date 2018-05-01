# Cookie example

We can achieve cookies by different ways:

1. Using 'key' parameter. Fasthttp will search using key in cookie field.

2. You can parse request header 'Set-Cookie' field.

3. You can use 'PeekCookie' from ResponseHeader using key. This is the long way to get cookie instead of using 'Cookie' structure.

4. Or finally you can use CookieJar object.
  1. Iterating with Response.VisitAllCookie

  2. You can use CookieJar.ResponseCookies.

# How to build

```
make
```

# How to run

```
./server & # This command will start http server in background
./client
```
