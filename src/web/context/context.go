package context

import (
	"errors"
	"mime/multipart"
	"net/http"
	"strings"
)

type Context struct {
	Request  *http.Request
	Response *Response
	params   map[string]Value
	File     map[string]*multipart.FileHeader
}

func NewContext(rw http.ResponseWriter, r *http.Request) *Context {
	return &Context{r, NewResponse(rw), make(map[string]Value), make(map[string]*multipart.FileHeader)}
}

/**
 * 获取URL
 */
func (c *Context) Url() string {
	return c.Request.URL.String()
}

/**
 * 获取请求地址
 */
func (c *Context) Uri() string {
	return c.Request.RequestURI
}

/**
 * 获取scheme
 */
func (c *Context) Scheme() string {
	if c.Request.URL.Scheme != "" {
		return c.Request.URL.Scheme
	} else if c.Request.TLS != nil {
		return "https"
	} else {
		return "http"
	}

}

/**
 * 获取request头内容
 */
func (c *Context) GetHeader(key string) string {
	return c.Request.Header.Get(key)
}

/**
 * 获取request的user agent
 */
func (c *Context) UserAgent() string {
	return c.GetHeader("USER-AGENT")
}

/**
 * 获取请求referer
 */
func (c *Context) Referer() string {
	return c.GetHeader("REFERER")
}

/**
 * 获取应用网址
 */
func (c *Context) Site() string {
	return c.Scheme() + "://" + c.Host()
}

/**
 * 获取Host
 */
func (c *Context) Host() string {
	if c.Request.Host != "" {
		hosts := strings.Split(c.Request.Host, ":")
		if len(hosts) > 0 {
			return hosts[0]
		}
		return c.Request.Host
	}
	return "127.0.0.1"
}

/**
 * 获取客户端IP
 */
func (c *Context) Ip() string {
	ip_proxy := c.Proxy()
	if len(ip_proxy) > 0 && ip_proxy[0] != "" {
		return ip_proxy[0]
	}
	index := strings.Index(c.Request.RemoteAddr, ":")
	if index != -1 {
		return c.Request.RemoteAddr[:index]
	}
	return "127.0.0.1"
}

func (c *Context) ParseForm(maxMemory int64) error {
	if c.Request.Method == "POST" || c.Request.Method == "PUT" {
		if strings.Contains(c.GetHeader("content-type"), "multipart/form-data") {
			if err := c.Request.ParseMultipartForm(maxMemory); err != nil {
				return errors.New("ParseMultipartForm err:" + err.Error())
			}
		} else {
			if err := c.Request.ParseForm(); err != nil {
				return errors.New("ParseForm err:" + err.Error())
			}
		}
	}
	return nil
}

/**
 * 获取代理IP
 */
func (c *Context) Proxy() []string {
	if ip_proxy := c.GetHeader("HTTP_X_FORWARDED_FOR"); ip_proxy != "" {
		return strings.Split(ip_proxy, ",")
	}
	return []string{}
}

/**
 * 获取Proto HTTP/1.1
 */
func (c *Context) Protocol() string {
	return c.Request.Proto
}

/**
 * 获取param值
 */
func (c *Context) GetParam(key string) Value {
	return c.params[key]
}

func (c *Context) SetParam(key, value string) {
	c.params[key] = Value(value)
}

/**
 * 检查param值
 */
func (c *Context) CkeckParam(key string) bool {
	_, ok := c.params[key]

	return ok
}
