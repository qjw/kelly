package store

import (
	"net/http"
	"time"
)

type Options struct {
	Path   string
	Domain string
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	MaxAge   *int
	Secure   *bool
	HttpOnly *bool

	// 服务器 session可用，后端key的前缀
	KeyPrefix string
}

func (this *Options) SetMaxAge(maxAge int) {
	this.MaxAge = &maxAge
}
func (this *Options) SetSecure(secure bool) {
	this.Secure = &secure
}
func (this *Options) SetHttpOnly(httpOnly bool) {
	this.HttpOnly = &httpOnly
}

func conbineOptions(options *Options) *Options {
	newOptions := &Options{
		Path:     "/",
		MaxAge:   newInt(86400 * 30),
		Secure:   newBool(false),
		HttpOnly: newBool(false),
	}
	if options == nil {
		return newOptions
	}

	// 合并
	if len(options.Path) > 0 {
		newOptions.Path = options.Path
	}
	if len(options.Domain) > 0 {
		newOptions.Domain = options.Domain
	}
	if options.MaxAge != nil {
		newOptions.MaxAge = options.MaxAge
	}
	if options.Secure != nil {
		newOptions.Secure = options.Secure
	}
	if options.HttpOnly != nil {
		newOptions.HttpOnly = options.HttpOnly
	}
	if len(options.KeyPrefix) > 0 {
		newOptions.KeyPrefix = options.KeyPrefix
	}
	return newOptions
}

func newCookie(name, value string, options *Options) *http.Cookie {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   *options.MaxAge,
		Secure:   *options.Secure,
		HttpOnly: *options.HttpOnly,
	}
	if *options.MaxAge > 0 {
		d := time.Duration(*options.MaxAge) * time.Second
		cookie.Expires = time.Now().Add(d)
	} else if *options.MaxAge < 0 {
		// Set it to the past to expire now.
		cookie.Expires = time.Unix(1, 0)
	}
	return cookie
}
