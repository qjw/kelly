// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package middleware

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func generateNormalCorsHeaders(c *CorsConfig) http.Header {
	headers := make(http.Header)
	if c.AllowCredentials {
		headers.Set("Access-Control-Allow-Credentials", "true")
	}
	if len(c.ExposeHeaders) > 0 {
		exposeHeaders := convert(normalize(c.ExposeHeaders), http.CanonicalHeaderKey)
		headers.Set("Access-Control-Expose-Headers", strings.Join(exposeHeaders, ","))
	}
	if c.AllowAllOrigins {
		headers.Set("Access-Control-Allow-Origin", "*")
	} else {
		headers.Set("Vary", "Origin")
	}
	return headers
}

func generatePreflightCorsHeaders(c *CorsConfig) http.Header {
	headers := make(http.Header)
	if c.AllowCredentials {
		headers.Set("Access-Control-Allow-Credentials", "true")
	}
	if len(c.AllowMethods) > 0 {
		allowMethods := convert(normalize(c.AllowMethods), strings.ToUpper)
		value := strings.Join(allowMethods, ",")
		headers.Set("Access-Control-Allow-Methods", value)
	}
	if len(c.AllowHeaders) > 0 {
		allowHeaders := convert(normalize(c.AllowHeaders), http.CanonicalHeaderKey)
		value := strings.Join(allowHeaders, ",")
		headers.Set("Access-Control-Allow-Headers", value)
	}
	if c.MaxAge > time.Duration(0) {
		value := strconv.FormatInt(int64(c.MaxAge/time.Second), 10)
		headers.Set("Access-Control-Max-Age", value)
	}
	if c.AllowAllOrigins {
		headers.Set("Access-Control-Allow-Origin", "*")
	} else {
		// Always set Vary headers
		// see https://github.com/rs/cors/issues/10,
		// https://github.com/rs/cors/commit/dbdca4d95feaa7511a46e6f1efb3b3aa505bc43f#commitcomment-12352001

		headers.Add("Vary", "Origin")
		headers.Add("Vary", "Access-Control-Request-Method")
		headers.Add("Vary", "Access-Control-Request-Headers")
	}
	return headers
}

type corsConfig struct {
	allowAllOrigins  bool
	allowCredentials bool
	allowOriginFunc  func(string) bool
	allowOrigins     []string
	exposeHeaders    []string
	normalHeaders    http.Header
	preflightHeaders http.Header
}

func newCors(config *CorsConfig) *corsConfig {
	if err := config.Validate(); err != nil {
		panic(err.Error())
	}
	return &corsConfig{
		allowOriginFunc:  config.AllowOriginFunc,
		allowAllOrigins:  config.AllowAllOrigins,
		allowCredentials: config.AllowCredentials,
		allowOrigins:     normalize(config.AllowOrigins),
		normalHeaders:    generateNormalCorsHeaders(config),
		preflightHeaders: generatePreflightCorsHeaders(config),
	}
}

func (cors *corsConfig) applyCors(w http.ResponseWriter, r *http.Request) (res bool) {
	res = true
	origin := r.Header.Get("Origin")
	if len(origin) == 0 {
		// request is not a CORS request
		return
	}
	if !cors.validateOrigin(origin) {
		// c.AbortWithStatus(http.StatusForbidden)
		setHttpCode(w, http.StatusForbidden)
		res = false
		return
	}

	headers := w.Header()
	if r.Method == "OPTIONS" {
		cors.handlePreflight(w)
		// 必须在最后才设置http code
		defer setHttpCode(w, http.StatusOK)
		res = false
		// 不能直接退出，因为后面要加上Access-Control-Allow-Origin
	} else {
		cors.handleNormal(headers)
	}

	if !cors.allowAllOrigins && !cors.allowCredentials {
		headers.Set("Access-Control-Allow-Origin", origin)
	}
	return
}

func (cors *corsConfig) validateOrigin(origin string) bool {
	if cors.allowAllOrigins {
		return true
	}
	for _, value := range cors.allowOrigins {
		if value == origin {
			return true
		}
	}
	if cors.allowOriginFunc != nil {
		return cors.allowOriginFunc(origin)
	}
	return false
}

func (cors *corsConfig) handlePreflight(w http.ResponseWriter) {
	for key, value := range cors.preflightHeaders {
		log.Printf("key %s value %s", key, value)
		w.Header()[key] = value
	}
}

func (cors *corsConfig) handleNormal(headers http.Header) {
	for key, value := range cors.normalHeaders {
		headers[key] = value
	}
}
