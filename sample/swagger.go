// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package main

import (
	"bytes"
	"github.com/qjw/kelly"
	"github.com/qjw/kelly/middleware"
	"github.com/qjw/kelly/middleware/swagger"
	"github.com/qjw/kelly/toolkits"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"net/http"
	"strconv"
)

type swaggerParam struct {
	NonceStr  string `json:"nonceStr"`
	Timestamp string `json:"timestamp"`
	Url       string `json:"url"`
}

func InitSwagger(r kelly.Router) {
	// 增加中间件处理跨域问题
	router := r.Group("/swagger", middleware.Cors(&middleware.CorsConfig{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type"},
	})).GlobalAnnotation(swagger.SetGlobalParam(&swagger.StructParam{
		Tags: []string{"API接口"},
	})).OPTIONS("/*path", func(c *kelly.Context) {
		c.ResponseStatusOK()
	})

	router.Annotation(swagger.Swagger(&swagger.StructParam{
		ResponseData: &swagger.SuccessResp{},
		QueryData:    &swaggerParam{},
		Summary:      "api1",
	})).GET("/api1", func(c *kelly.Context) {
		c.ResponseStatusOK()
	})

	router.Annotation(swagger.Swagger(&swagger.StructParam{
		ResponseData: &swagger.SuccessResp{},
		FormData:     &swaggerParam{},
		Summary:      "api1",
	})).POST("/api1", func(c *kelly.Context) {
		c.ResponseStatusOK()
	})

	router.Annotation(swagger.Swagger(&swagger.StructParam{
		ResponseData: &swagger.SuccessResp{},
		JsonData:     &swaggerParam{},
		Summary:      "api1",
	})).PUT("/api1", func(c *kelly.Context) {
		c.ResponseStatusOK()
	})

	router.Annotation(swagger.Swagger(&swagger.StructParam{
		ResponseData: &swagger.SuccessResp{},
		PathData:     &swaggerParam{},
		Summary:      "api1",
	})).DELETE("/api1/:nonceStr/:timestamp/:url", func(c *kelly.Context) {
		c.ResponseStatusOK()
	})

	router.Annotation(swagger.Swagger(&swagger.StructParam{
		ResponseData: &swagger.SuccessResp{},
		FormData:     &swaggerParam{},
		Summary:      "api1",
	})).PATCH("/api1", func(c *kelly.Context) {
		c.ResponseStatusOK()
	})

	router.Annotation(
		swagger.SwaggerFile("swagger.yaml:upload_material"),
	).POST("/upload_material", func(c *kelly.Context) {
		c.ResponseStatusOK()
	})

	router.Annotation(
		swagger.SwaggerFile("swagger.yaml:get_material"),
	).GET("/get_material", func(c *kelly.Context) {
		m := image.NewRGBA(image.Rect(0, 0, 240, 240))
		blue := color.RGBA{0, 0, 255, 255}
		draw.Draw(m, m.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)

		var img image.Image = m
		buffer := new(bytes.Buffer)
		if err := jpeg.Encode(buffer, img, nil); err != nil {
			c.WriteString(http.StatusInternalServerError, err.Error())
			return
		}

		c.SetHeader("Content-Type", "image/jpeg")
		c.SetHeader("Content-Length", strconv.Itoa(len(buffer.Bytes())))
		if _, err := c.Write(buffer.Bytes()); err != nil {
			c.WriteString(http.StatusInternalServerError, err.Error())
			return
		}
	})

	router.Annotation(
		swagger.SwaggerFile("swagger.yaml:qrcode"),
	).GET("/qrcode", func(c *kelly.Context) {
		qrcode, _ := toolkits.NewQRCode(c.MustGetQueryVarible("content"), toolkits.QrcodeMedium)
		qrcode.WriteKelly(400, c)
	})
}
