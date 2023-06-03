// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import "net/http"

// Render interface is to be implemented by JSON, XML, HTML, YAML and so on.
// todo 这个包和接口的设计都贼牛逼，可以收藏起来参考！！！
/* NOTE 该接口用来渲染http响应内容，支持JSON, XML, HTML, YAML等多种响应体
具体实现中，WriteContentType()接口用来写入响应头的Content-Type值，JSON实现则写入application/json，XML实现则写入application/xml...当然，这里的写入是只是写到w的header map中，还没刷到io中
Render就比较关键了，在Render中完成对Data any这个结构的处理，进行json化，或xml化...,之后将前面调用WriteContentType()写入的map的header(Content-Type)写入w，再将json/xml化后的data写入w，返回，完美！！！
*/
type Render interface {
	// Render writes data with custom ContentType.
	Render(http.ResponseWriter) error
	// WriteContentType writes custom ContentType.
	WriteContentType(w http.ResponseWriter)
}

var (
	_ Render     = JSON{}
	_ Render     = IndentedJSON{}
	_ Render     = SecureJSON{}
	_ Render     = JsonpJSON{}
	_ Render     = XML{}
	_ Render     = String{}
	_ Render     = Redirect{}
	_ Render     = Data{}
	_ Render     = HTML{}
	_ HTMLRender = HTMLDebug{}
	_ HTMLRender = HTMLProduction{}
	_ Render     = YAML{}
	_ Render     = Reader{}
	_ Render     = AsciiJSON{}
	_ Render     = ProtoBuf{}
	_ Render     = TOML{}
)

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
