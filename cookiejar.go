package fasthttp

import "io"

// CookieJar is container of cookies
//
// This object is used to handle multiple cookies
type CookieJar map[string]*Cookie

// Set sets cookie using key-value
//
// This function can replace an existent cookie
func (cj *CookieJar) Set(key, value string) {
	setCookie(cj, key, value)
}

// Get returns and delete a value from cookiejar.
func (cj *CookieJar) Get() *Cookie {
	for k, v := range *cj {
		delete(*cj, k)
		return v
	}
	return nil
}

// SetBytesK sets cookie using key=value
//
// This function can replace an existent cookie.
func (cj *CookieJar) SetBytesK(key []byte, value string) {
	setCookie(cj, b2s(key), value)
}

// SetBytesV sets cookie using key=value
//
// This function can replace an existent cookie.
func (cj *CookieJar) SetBytesV(key string, value []byte) {
	setCookie(cj, key, b2s(value))
}

// SetBytesKV sets cookie using key=value
//
// This function can replace an existent cookie.
func (cj *CookieJar) SetBytesKV(key, value []byte) {
	setCookie(cj, b2s(key), b2s(value))
}

func setCookie(cj *CookieJar, key, value string) {
	c, ok := (*cj)[key]
	if !ok {
		c = AcquireCookie()
	}
	c.SetKey(key)
	c.SetValue(value)
	(*cj)[key] = c
}

// SetCookie sets cookie using its key.
//
// After that you can use Peek function to get cookie value.
func (cj *CookieJar) SetCookie(cookie *Cookie) {
	(*cj)[b2s(cookie.Key())] = cookie
}

// Peek peeks cookie value using key.
//
// This function does not delete cookie
func (cj *CookieJar) Peek(key string) *Cookie {
	return (*cj)[key]
}

// Release releases all cookie values.
func (cj *CookieJar) Release() {
	for k, v := range *cj {
		ReleaseCookie(v)
		delete(*cj, k)
	}
}

// ReleaseCookie releases a cookie specified by parsed key.
func (cj *CookieJar) ReleaseCookie(key string) {
	c, ok := (*cj)[key]
	if ok {
		ReleaseCookie(c)
		delete(*cj, key)
	}
}

// PeekValue returns value of specified cookie-key.
func (cj *CookieJar) PeekValue(key string) []byte {
	c, ok := (*cj)[key]
	if ok {
		return c.Value()
	}
	return nil
}

// ResponseCookies gets all response cookies and stores it in cj.
func (cj *CookieJar) ResponseCookies(r *Response) {
	r.Header.VisitAllCookie(func(key, value []byte) {
		cj.SetBytesKV(key, value)
	})
}

// RequestCookies gets all request cookies and stores it in cj.
func (cj *CookieJar) RequestCookies(r *Request) {
	r.Header.VisitAllCookie(func(key, value []byte) {
		cj.SetBytesKV(key, value)
	})
}

// WriteTo writes all cookies representation to w.
func (cj *CookieJar) WriteTo(w io.Writer) (n int64, err error) {
	for _, c := range *cj {
		nn, err := c.WriteTo(w)
		n += nn
		if err != nil {
			break
		}
	}
	return
}
