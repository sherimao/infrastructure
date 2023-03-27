package network

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	"github.com/valyala/fasthttp"
)

// 默认请求超时时间
const defaultRequestTimeout = 5 * time.Second

// Logger 日志接口
type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Error(v ...interface{})
}

var _ Logger = (*defaultLogger)(nil)

// 实现默认日志
type defaultLogger struct {
	Ctx context.Context
}

// Debug Debug 级别日志
func (p defaultLogger) Debug(v ...interface{}) {
}

// Info Info 级别日志
func (p defaultLogger) Info(v ...interface{}) {
}

// Error Error 级别日志
func (p defaultLogger) Error(v ...interface{}) {
}

// 选项
type options struct {
	logger Logger
}

// Option 选项
type Option interface {
	apply(*options)
}

type loggerOption struct {
	Log Logger
}

func (l loggerOption) apply(opts *options) {
	opts.logger = l.Log
}

// WithLogger 带日志
func WithLogger(log Logger) Option {
	return loggerOption{Log: log}
}

// HTTPRequest ...
func HTTPRequest(content Content, opts ...Option) (status int, res []byte, err error) {
	localOption := options{
		logger: defaultLogger{},
	}
	for _, o := range opts {
		o.apply(&localOption)
	}
	t := Http{
		logger: localOption.logger,
	}
	status, res, err = t.HTTPRequest(content)
	return
}

// Http ...
type Http struct {
	logger Logger
}

// Content ...
type Content struct {
	URL         string
	Method      string
	Headers     map[string]string
	ContentType string
	Body        []byte
	Timeout     time.Duration
}

func getJSON(data interface{}) string {
	j, _ := json.Marshal(data)
	return string(j)
}

// HTTPRequest 发起 HTTP 请求
func (p Http) HTTPRequest(content Content) (status int, res []byte, err error) {
	p.logger.Info("URL: ", content.URL)
	p.logger.Info("Method: ", content.Method)
	p.logger.Info("Headers: ", getJSON(content.Headers))
	p.logger.Info("ContentType: ", content.ContentType)
	p.logger.Info("Body: ", string(content.Body))
	p.logger.Info("Timeout: ", content.Timeout)

	start := time.Now()
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
		p.logger.Info("HTTPRequest cost(seconds): ", time.Since(start).Seconds())
	}()

	u, err := url.Parse(content.URL)
	if err != nil {
		return
	}
	req.SetRequestURI(u.String())
	req.Header.SetMethod(content.Method)
	for k, v := range content.Headers {
		req.Header.Set(k, v)
	}
	contentType := content.ContentType
	if contentType == "" {
		contentType = "application/json; charset=utf-8"
	}
	req.Header.SetContentType(contentType)
	req.SetBody(content.Body)
	timeout := content.Timeout
	if timeout == 0 {
		timeout = defaultRequestTimeout
	}

	err = fasthttp.DoTimeout(req, resp, timeout)
	if err != nil {
		if err == fasthttp.ErrTimeout {
			resp.SetStatusCode(504)
		}
		return
	}

	res = resp.Body()
	p.logger.Info("rsp: ", string(res))
	p.logger.Info("StatusCode: ", resp.StatusCode())
	status = resp.StatusCode()
	// 在释放了resp之后，resp.Body()的内容可能不再是正确的，这里先解出来再序列化，相当于复制一个副本
	var tmp map[string]interface{}
	_ = json.Unmarshal(res, &tmp)
	res, _ = json.Marshal(tmp)
	return
}
