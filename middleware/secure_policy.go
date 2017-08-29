// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

type (
	// Secure is a middleware that helps setup a few basic security features. A single secure.Options struct can be
	// provided to configure which features should be enabled, and the ability to override a few of the default values.
	securePolicy struct {
		// Customize Secure with an Options struct.
		config       *SecureConfig
		fixedHeaders []header
	}

	header struct {
		key   string
		value []string
	}
)

// Constructs a new Policy instance with supplied options.
func newPolicy(config *SecureConfig) *securePolicy {
	policy := &securePolicy{}
	policy.loadConfig(config)
	return policy
}

func (p *securePolicy) loadConfig(config *SecureConfig) {
	p.config = config
	p.fixedHeaders = make([]header, 0, 5)

	// Frame Options header.
	if len(config.CustomFrameOptionsValue) > 0 {
		p.addHeader("X-Frame-Options", config.CustomFrameOptionsValue)
	} else if config.FrameDeny {
		p.addHeader("X-Frame-Options", "DENY")
	}

	// Content Type Options header.
	if config.ContentTypeNosniff {
		p.addHeader("X-Content-Type-Options", "nosniff")
	}

	// XSS Protection header.
	if config.BrowserXssFilter {
		p.addHeader("X-Xss-Protection", "1; mode=block")
	}

	// Content Security Policy header.
	if len(config.ContentSecurityPolicy) > 0 {
		p.addHeader("Content-Security-Policy", config.ContentSecurityPolicy)
	}

	// Strict Transport Security header.
	if config.STSSeconds != 0 {
		stsSub := ""
		if config.STSIncludeSubdomains {
			stsSub = "; includeSubdomains"
		}

		// TODO
		// "max-age=%d%s" refactor
		p.addHeader(
			"Strict-Transport-Security",
			fmt.Sprintf("max-age=%d%s", config.STSSeconds, stsSub))
	}
}

func (p *securePolicy) addHeader(key string, value string) {
	p.fixedHeaders = append(p.fixedHeaders, header{
		key:   key,
		value: []string{value},
	})
}

func (p *securePolicy) applyToContext(w http.ResponseWriter, r *http.Request) bool {
	if !p.config.IsDevelopment {
		p.writeSecureHeaders(w, r)

		if !p.checkAllowHosts(w, r) {
			return false
		}
		if !p.checkSSL(w, r) {
			return false
		}
	}
	return true
}

func (p *securePolicy) writeSecureHeaders(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	for _, pair := range p.fixedHeaders {
		header[pair.key] = pair.value
	}
}

func (p *securePolicy) checkAllowHosts(w http.ResponseWriter, r *http.Request) bool {
	if len(p.config.AllowedHosts) == 0 {
		return true
	}

	host := r.Host
	if len(host) == 0 {
		host = r.URL.Host
	}

	for _, allowedHost := range p.config.AllowedHosts {
		if strings.EqualFold(allowedHost, host) {
			return true
		}
	}

	return false
}

func (p *securePolicy) checkSSL(w http.ResponseWriter, r *http.Request) bool {
	if !p.config.SSLRedirect {
		return true
	}

	isSSLRequest := strings.EqualFold(r.URL.Scheme, "https") || r.TLS != nil
	if isSSLRequest {
		return true
	}

	// TODO
	// req.Host vs req.URL.Host
	url := r.URL
	url.Scheme = "https"
	url.Host = r.Host

	if len(p.config.SSLHost) > 0 {
		url.Host = p.config.SSLHost
	}

	status := http.StatusMovedPermanently
	if p.config.SSLTemporaryRedirect {
		status = http.StatusTemporaryRedirect
	}

	http.Redirect(w, r, url.String(), status)
	return false
}
