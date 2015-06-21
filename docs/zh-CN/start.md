#开始使用 Seago

在我们开始之前，必须明确的一点就是，文档不会教授您任何有关 Go 语言的基础知识。所有对 Seago 使用的讲解均是基于您已有的知识基础上展开的。

通过执行以下命令来安装 Seago：

	go get github.com/seago/seago

##最简示例

创建一个名为 main.go 的文件，然后输入以下代码：

	package main

	import "github.com/seago/seago"

	func main() {
    	s := seago.Classic()
    	s.Get("/", func() string {
        	return "Hello world!"
    	})
    	s.Run()
	}
函数 seago.Classic 创建并返回一个 经典 Seago 实例。

方法 s.Get 是用于注册针对 HTTP GET 请求的路由。在本例中，我们注册了针对根路径 / 的路由，并提供了一个 处理器 函数来进行简单的处理操作，即返回内容为 Hello world! 的字符串作为响应。

您可能会问，为什么处理器函数可以返回一个字符串作为响应？这是由于 返回值 所带来的特性。换句话说，我们在本例中使用了 Seago 中处理器的一个特殊语法来将返回值作为响应内容。

最后，我们调用 s.Run 方法来让服务器启动。在默认情况下，Seago 实例 会监听 0.0.0.0:4000。

接下来，就可以执行命令 go run main.go 运行程序。您应该在程序启动后看到一条日志信息：

[Seago] listening on 0.0.0.0:4000 (development)
现在，打开您的浏览器然后访问 localhost:4000。您会发现，一切是如此的美好！

##扩展示例

现在，让我们对 main.go 做出一些修改，以便进行更多的练习。

	package main

	import (
    	"log"
    	"net/http"

    	"github.com/seago/seago"
	)

	func main() {
    	s := seago.Classic()
    	s.Get("/", myHandler)

    	log.Println("Server is running...")
    	log.Println(http.ListenAndServe("0.0.0.0:4000", s))
	}

	func myHandler(ctx *seago.Context) string {
    	return "the request path is: " + ctx.Req.RequestURI
	}
当您再次执行命令 go run main.go 运行程序的时候，您会看到屏幕上显示的内容为 the request path is: /。

首先，我们依旧使用了 经典 Seago 来为根路径 / 注册针对 HTTP GET 请求的路由。但我们不再使用匿名函数，而是改用名为 myHandler 的函数作为处理器。需要注意的是，注册路由时，不需要在函数名称后面加上括号，因为我们不需要在此时调用这个函数。

函数 myHandler 接受一个类型为 *seago.Context 的参数，并返回一个字符串。您可能已经发现我们并没有告诉 Web 需要传递什么参数给处理器，而且当您查看 s.Get 方法的声明时会发现，Seago 实际上将所有的处理器（seago.Handler）都当作类型 interface{} 来处理。那么，Web 又是怎么知道需要传递什么参数来调用处理器并执行逻辑的呢？

这就涉及到 服务注入 的概念了， *seago.Context 就是默认注入的服务之一，所以您可以直接使用它作为参数。如果您不明白怎么注入您自己的服务，没关系，反正还不是时候知道这些。

和之前的例子一样，我们需要让服务器监听在某个地址上。这一次，我们使用 Go 标准库的函数 http.ListenAndServe 来完成这项操作。如此一来，您便可以发现，任一 Seago 实例 都是和标准库完全兼容的。

##了解更多

您现在已经知道怎么基于 Seago 来书写简单的代码，请尝试修改上文中的两个示例，并确保您已经完全理解上文中的所有内容。

当您觉得自己已经原地满血复活后，就可以继续学习之后的内容了。
