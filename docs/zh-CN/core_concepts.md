# Seago 核心概念
经典 Seago

为了更快速的启用 Seago，seago.Classic 提供了一些默认的组件以方便 Web 开发:

	s := seago.Classic()
  	// ... 可以在这里使用中间件和注册路由
  	s.Run()
下面是 seago.Classic 已经包含的功能：

请求/响应日志 - seago.Logger
容错恢复 - seago.Recovery
静态文件服务 - seago.Static
Seago 实例

任何类型为 seago.Seago 的对象都可以被认为是 Macaron 的实例，您可以在单个应用中使用任意数量的 Seago 实例。

处理器

处理器是 Seago 的灵魂和核心所在. 一个处理器基本上可以是任何的函数:

	s.Get("/", func() {
    	println("hello world")
	})
返回值

当一个处理器返回结果的时候, Seago 将会把返回值作为字符串写入到当前的 http.ResponseWriter 里面：

	s.Get("/", func() string {
    	return "hello world" // HTTP 200 : "hello world"
	})
另外你也可以选择性的返回状态码:

	s.Get("/", func() (int, string) {
    	return 418, "i'm a teapot" // HTTP 418 : "i'm a teapot"
	})
服务注入

处理器是通过反射来调用的，Seago 通过 依赖注入 来为处理器注入参数列表。 这样使得 Seago 与 Go 语言的 http.HandlerFunc 接口完全兼容。

如果你加入一个参数到你的处理器, Seago 将会搜索它参数列表中的服务，并且通过类型判断来解决依赖关系：

	m.Get("/", func(resp http.ResponseWriter, req *http.Request) { 
	    // resp 和 req 是由 Macaron 默认注入的服务
	    resp.WriteHeader(200) // HTTP 200
	})
下面的这些服务已经被包含在经典 Macaron 中（seago.Classic）：

*seago.Context - HTTP 请求上下文
*log.Logger - Seago 全局日志器
http.ResponseWriter - HTTP 响应流
*http.Request - HTTP 请求对象
中间件机制

中间件处理器是工作于请求和路由之间的。本质上来说和 Macaron 其他的处理器没有分别. 您可以使用如下方法来添加一个中间件处理器到队列中:

	s.Use(func() {
	  // 处理中间件事物
	})
你可以通过 Handlers 函数对中间件队列实现完全的控制. 它将会替换掉之前的任何设置过的处理器:

	s.Handlers(
	    Middleware1,
	    Middleware2,
	    Middleware3,
	)
中间件处理器可以非常好处理一些功能，包括日志记录、授权认证、会话（sessions）处理、错误反馈等其他任何需要在发生在 HTTP 请求之前或者之后的操作:

	// 验证一个 API 密钥
	s.Use(func(ctx *seago.Context) {
	    if ctx.Requset.Header.Get("X-API-KEY") != "secret123" {
	        ctx.Response.WriteHeader(http.StatusUnauthorized)
	    }
	})
Seago 环境变量

一些 Seago 处理器依赖 seago.Env 全局变量为开发模式和部署模式表现出不同的行为，不过更建议使用环境变量 SEAGO_ENV=production 来指示当前的模式为部署模式。

