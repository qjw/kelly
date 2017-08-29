// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package middleware

import (
	"github.com/qjw/kelly"
	"net/http"
)

// Options is a struct for specifying configuration options for the secure.
type SecureConfig struct {
	// AllowedHosts is a list of fully qualified domain names that are allowed.
	//Default is empty list, which allows any and all host names.
	AllowedHosts []string
	// If SSLRedirect is set to true, then only allow https requests.
	// Default is false.
	SSLRedirect bool
	// If SSLTemporaryRedirect is true, the a 302 will be used while redirecting.
	// Default is false (301).
	SSLTemporaryRedirect bool
	// SSLHost is the host name that is used to redirect http requests to https.
	// Default is "", which indicates to use the same host.
	SSLHost string
	// STSSeconds is the max-age of the Strict-Transport-Security header.
	// Default is 0, which would NOT include the header.
	STSSeconds int64
	// If STSIncludeSubdomains is set to true, the `includeSubdomains` will
	// be appended to the Strict-Transport-Security header. Default is false.
	STSIncludeSubdomains bool
	// If FrameDeny is set to true, adds the X-Frame-Options header with
	// the value of `DENY`. Default is false.
	FrameDeny bool
	// CustomFrameOptionsValue allows the X-Frame-Options header value
	// to be set with a custom value. This overrides the FrameDeny option.
	CustomFrameOptionsValue string
	// If ContentTypeNosniff is true, adds the X-Content-Type-Options header
	// with the value `nosniff`. Default is false.
	ContentTypeNosniff bool
	// If BrowserXssFilter is true, adds the X-XSS-Protection header with
	// the value `1; mode=block`. Default is false.
	BrowserXssFilter bool
	// ContentSecurityPolicy allows the Content-Security-Policy header value
	// to be set with a custom value. Default is "".
	// http://www.ruanyifeng.com/blog/2016/09/csp.html  XSS攻击
	ContentSecurityPolicy string
	// When true, the whole secury policy applied by the middleware is disable
	// completely.
	IsDevelopment bool
	//// Handlers for when an error occurs (ie bad host).
	BadHostHandler kelly.HandlerFunc
}

func DefaultSecureConfig() *SecureConfig {
	return &SecureConfig{
		SSLRedirect:           true,
		IsDevelopment:         false,
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
	}
}

func Secure(config *SecureConfig) kelly.HandlerFunc {
	policy := newPolicy(config)
	if config.BadHostHandler == nil {
		config.BadHostHandler = func(c *kelly.Context) {
			c.WriteString(http.StatusForbidden, http.StatusText(http.StatusForbidden))
		}
	}

	return func(c *kelly.Context) {
		if policy.applyToContext(c, c.Request()) {
			c.InvokeNext()
		} else {
			config.BadHostHandler(c)
		}
	}
}
