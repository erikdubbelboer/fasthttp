# Cookie example

We can get the cookie from diferents ways.
1. Using 'key' parameter. Fasthttp will search using key in cookie field.
2. You can parse request header 'Set-Cookie' field.
3. Or you can use 'PeekCookie' from ResponseHeader using key. This is the long way to get cookie instead of using 'Cookie' structure.

# How to build

```
make
```

# How to run

```
./server & # This command will start http server in background
./client
```
