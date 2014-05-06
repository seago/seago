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

func (c *Context) Url() string {
	return c.Request.URL.String()
}

func (c *Context) Uri() string {
	return c.Request.RequestURI
}

func (c *Context) Scheme() string {
	if c.Request.URL.Scheme != "" {
		return c.Request.URL.Scheme
	} else if c.Request.TLS != nil {
		return "https"
	} else {
		return "http"
	}

}

func (c *Context) GetHeader(key string) string {
	return c.Request.Header.Get(key)
}

func (c *Context) UserAgent() string {
	return c.GetHeader("USER-AGENT")
}

func (c *Context) Referer() string {
	return c.GetHeader("REFERER")
}

func (c *Context) Site() string {
	return c.Scheme() + "://" + c.Host()
}

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

func (c *Context) GetCookie(key string) string {
	cookie, err := c.Request.Cookie(key)
	if err != nil {
		return ""
	}
	return cookie.Value
}

/**
* others[0] is cookie path
* others[1] is cookie domain
* others[2] is cookie Secure
* others[3] is cookie httponly
 */
func (c *Context) SetCookie(key, value string, maxAge int, others ...interface{}) {
	cookie := &http.Cookie{
		Name:     key,
		Value:    value,
		MaxAge:   maxAge,
		Path:     "/",
		Secure:   false,
		HttpOnly: false,
	}
	if len(others) > 0 {
		cookie.Path = others[0].(string)
	}
	if len(others) > 1 {
		cookie.Domain = others[1].(string)
	}
	if len(others) > 2 {
		cookie.Secure = others[2].(bool)
	}
	if len(others) > 3 {
		cookie.HttpOnly = others[3].(bool)
	}
	c.Response.SetHeader("Set-Cookie", cookie.String())
}

func (c *Context) Proxy() []string {
	if ip_proxy := c.GetHeader("HTTP_X_FORWARDED_FOR"); ip_proxy != "" {
		return strings.Split(ip_proxy, ",")
	}
	return []string{}
}

func (c *Context) Protocol() string {
	return c.Request.Proto
}

func (c *Context) GetParam(key string) Value {
	return c.params[key]
}

func (c *Context) SetParam(key, value string) {
	c.params[key] = Value(value)
}

func (c *Context) CkeckParam(key string) bool {
	_, ok := c.params[key]

	return ok
}
