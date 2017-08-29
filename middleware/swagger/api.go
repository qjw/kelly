// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package swagger

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"strings"
	"github.com/qjw/kelly"
	"github.com/qjw/kelly/binding"
	"regexp"
)

type DocLoader func(key string) ([]byte, error)

type options struct {
	// 是否调试模式
	debugFlag bool
	// doc文档定义路径
	docPath string
	// url前缀
	baseUrl string
	// 文件和里面内容的缓存
	swaggerData map[string]SwaggerEntry
	// 文档对象缓存
	docData map[string]*SwaggerDocFile
	// 当不从文件中加载doc时，获取swagger数据的loader
	docLoader DocLoader
}

func newOptions(config *Config) *options {
	opt := &options{
		debugFlag:   config.Debug,
		docPath:     config.DocFilePath,
		baseUrl:     config.BasePath,
		swaggerData: make(map[string]SwaggerEntry),
		docData:     make(map[string]*SwaggerDocFile),
	}
	// 非调试模式，不允许从外部加载doc文件
	if !opt.debugFlag {
		opt.docPath = ""
	}
	if len(opt.docPath) > 0 {
		if path, err := filepath.Abs(opt.docPath); err != nil {
			panic(err)
			return nil
		} else {
			if strings.HasSuffix(path, "/") {
				opt.docPath = path
			} else {
				opt.docPath = path + "/"
			}
		}
	}
	return opt
}

var (
	gOption *options = nil
)

func swaggerImp(docFile *SwaggerDocFile, path string, method string, entry string) {
	if data, ok := (*docFile)[entry]; ok {
		swaggerFinish(path, method, &data)
	} else {
		panic(errors.New("file'" + entry + "' have no entry '" + entry + "'"))
	}
}

func parseFileNode(entry string) (filepath string, node string, err error) {
	entrys := strings.Split(entry, ":")
	if len(entrys) != 2 {
		err = errors.New("invalid swagger entry '" + entry + "'")
		return
	}

	filepath = entrys[0]
	node = entrys[1]
	if len(filepath) == 0 {
		err = errors.New("invalid swagger entry '" + entry + "',invalid filepath")
		return
	}
	if len(node) == 0 {
		err = errors.New("invalid swagger entry '" + entry + "',invalid node")
		return
	}
	return
}

func realPath(r kelly.Router, path string) string {
	var burl = r.Path()
	if len(gOption.baseUrl) > 0 && strings.HasPrefix(burl, gOption.baseUrl) {
		burl = strings.TrimPrefix(burl, gOption.baseUrl)
	}
	if strings.HasPrefix(path, "/") {
		path = burl + path
	} else {
		path = burl + "/" + path
	}

	// /aaa/bb/:cc/:dd/ee =>  /aaa/bb/{cc}/{dd}/ee
	re := regexp.MustCompile(`:([0-9a-zA-Z]+)`)
	return re.ReplaceAllString(path, "{$1}")
}

const SwaggerGlobalParam = "_swaggerGlobalParamKey"

func SetGlobalParam(extra *StructParam) kelly.AnnotationHandlerFunc {
	return func(c *kelly.AnnotationContext) {
		c.R().Set(SwaggerGlobalParam,extra)
	}
}

func Swagger(extra *StructParam) kelly.AnnotationHandlerFunc {
	return func(c *kelly.AnnotationContext) {
		// 检查全局信息
		gParam := c.R().Get(SwaggerGlobalParam)
		if gParam != nil{
			gParam2 := gParam.(*StructParam)
			if extra.Tags == nil && gParam2.Tags != nil{
				extra.Tags = gParam2.Tags
			}
		}

		path := realPath(c.R(), c.Path())
		swaggerFinish(path, c.Method(), NewSwaggerMethodEntry(extra))
	}
}

func SwaggerFile(entry string) kelly.AnnotationHandlerFunc {
	return func(c *kelly.AnnotationContext) {
		path := realPath(c.R(), c.Path())
		// 解析文件路径和内部路径
		var err error
		rfilepath, node, err := parseFileNode(entry)
		if err != nil {
			panic(err)
		}

		// 是否有缓存
		if docFile, ok := gOption.docData[rfilepath]; ok {
			swaggerImp(docFile, path, c.Method(), node)
			return
		}

		var yamlFile []byte
		if len(gOption.docPath) > 0 {
			rfilepath = gOption.docPath + rfilepath

			// 加载文件
			yamlFile, err = ioutil.ReadFile(rfilepath)
			if err != nil {
				panic(err)
			}
		} else {
			yamlFile, err = gOption.docLoader(rfilepath)
			if err != nil {
				panic(err)
			}
		}

		var docFile SwaggerDocFile
		err = yaml.Unmarshal(yamlFile, &docFile)
		if err != nil {
			panic(err)
		}

		// 写入缓存
		gOption.docData[rfilepath] = &docFile
		swaggerImp(&docFile, path, c.Method(), node)
	}
}

func swaggerFinish(path string, method string, entry *SwaggerMethodEntry) {
	if err := binding.Validate(entry); err != nil {
		panic(err)
		return
	}

	var sentry SwaggerEntry
	if v, ok := gOption.swaggerData[path]; ok {
		sentry = v
	} else {
		sentry = SwaggerEntry{}
	}
	sentry.SetMethod(method, *entry)
	gOption.swaggerData[path] = sentry
}

