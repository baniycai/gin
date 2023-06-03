// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

//go:build !nomsgpack

package binding

import "net/http"

// 在 HTTP 协议中，Content-Type 头字段的语法规则是 "type/subtype" 的形式。其中，“type”表示媒体类型的类别，而“subtype”表示具体的媒体类型。
// 在上述常见的 Content-Type 中，“application” 表示媒体类型属于应用程序数据，也就是说，这些 Content-Type 用于表示传输应用程序特定的数据格式或协议。
// 除了 application 类型之外，还有 text、image、audio、video 等其他类型，每种类型都用于表示特定的媒体格式或数据类型。

// Content-Type MIME of the most common data formats.
const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
	MIMEPROTOBUF          = "application/x-protobuf"
	MIMEMSGPACK           = "application/x-msgpack"
	MIMEMSGPACK2          = "application/msgpack"
	MIMEYAML              = "application/x-yaml"
	MIMETOML              = "application/toml"
)

// Binding describes the interface which needs to be implemented for binding the
// data present in the request such as JSON request body, query parameters or
// the form POST.
// todo 跟render一样，都是贼牛逼的接口设计，可以收藏
/* NOTE Name用来返回该绑定器的名称，如json、xml、form
Bind则是一个比较关键的方法，通常是取http.Request的body，根据绑定器的类型，将body的内容进行不同的处理，比如json解析、xml解析或者是表单(form)解析
之后将其值绑定到any类型的参数上，最后使用校验器来绑定的结果进行校验，主要是通过反射取binding这个tag来校验，比如binding的值为required，则该字段必须有值，否则返回err
*/
type Binding interface {
	Name() string
	Bind(*http.Request, any) error
}

// BindingBody adds BindBody method to Binding. BindBody is similar with Bind,
// but it reads the body from supplied bytes instead of req.Body.
type BindingBody interface {
	Binding
	BindBody([]byte, any) error
}

// BindingUri adds BindUri method to Binding. BindUri is similar with Bind,
// but it reads the Params.
type BindingUri interface {
	Name() string
	BindUri(map[string][]string, any) error
}

// StructValidator is the minimal interface which needs to be implemented in
// order for it to be used as the validator engine for ensuring the correctness
// of the request. Gin provides a default implementation for this using
// https://github.com/go-playground/validator/tree/v10.6.1.
type StructValidator interface {
	// ValidateStruct can receive any kind of type and it should never panic, even if the configuration is not right.
	// If the received type is a slice|array, the validation should be performed travel on every element.
	// If the received type is not a struct or slice|array, any validation should be skipped and nil must be returned.
	// If the received type is a struct or pointer to a struct, the validation should be performed.
	// If the struct is not valid or the validation itself fails, a descriptive error should be returned.
	// Otherwise nil must be returned.
	ValidateStruct(any) error

	// Engine returns the underlying validator engine which powers the
	// StructValidator implementation.
	Engine() any
}

// Validator is the default validator which implements the StructValidator
// interface. It uses https://github.com/go-playground/validator/tree/v10.6.1
// under the hood.
var Validator StructValidator = &defaultValidator{}

// These implement the Binding interface and can be used to bind the data
// present in the request to struct instances.
var (
	JSON          = jsonBinding{}
	XML           = xmlBinding{}
	Form          = formBinding{}
	Query         = queryBinding{}
	FormPost      = formPostBinding{}
	FormMultipart = formMultipartBinding{}
	ProtoBuf      = protobufBinding{}
	MsgPack       = msgpackBinding{}
	YAML          = yamlBinding{}
	Uri           = uriBinding{}
	Header        = headerBinding{}
	TOML          = tomlBinding{}
)

// Default returns the appropriate Binding instance based on the HTTP method
// and the content type.
// 结合请求方法和参数类型来返回不同的绑定器
func Default(method, contentType string) Binding {
	if method == http.MethodGet {
		return Form
	}

	switch contentType {
	case MIMEJSON:
		return JSON
	case MIMEXML, MIMEXML2:
		return XML
	case MIMEPROTOBUF:
		return ProtoBuf
	case MIMEMSGPACK, MIMEMSGPACK2:
		return MsgPack
	case MIMEYAML:
		return YAML
	case MIMETOML:
		return TOML
	case MIMEMultipartPOSTForm:
		return FormMultipart
	default: // case MIMEPOSTForm:
		return Form
	}
}

func validate(obj any) error {
	if Validator == nil {
		return nil
	}
	return Validator.ValidateStruct(obj)
}
