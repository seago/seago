## 模板引擎
目前 Seago 应用有两款官方模板引擎中间件可供选择，即 `seago.Renderer` 和 pongo2.Pongoer。

您可以自由选择使用哪一款模板引擎，并且您只能为一个 Seago 实例 注册一款模板引擎。

共有特性：

* 均支持 XML、JSON 和原始数据格式的响应，它们之间的不同只体现在 HTML 渲染上。
* 均使用 templates 作为默认模板文件目录。
* 均使用 .tmpl 和 .html 作为默认模板文件后缀。
* 均支持通过 Seago 环境变量 来判断是否缓存模板文件（当 seago.Env == seago.PROD 时）。
* 渲染 HTML

## Go 模板引擎

该服务可以通过函数 `seago.Renderer` 来注入，并通过类型 `seago.Render` 来体现。该服务为可选，一般情况下可直接使用 `*seago.Context.Render`。该服务使用 Go 语言内置的模板引擎来渲染 HTML。如果想要了解更多有关使用方面的信息，请参见 官方文档。

使用示例

假设您的应用拥有以下目录结构：

    main/
        |__ main.go
        |__ templates/
                |__ hello.tmpl
    hello.tmpl：

    <h1>Hello {{.Name}}</h1>
    main.go：

    package main

    import "github.com/seago/seago"

    func main() {
        s := seago.Classic()
        s.Use(`seago.Renderer`())
        
        s.Get("/", func(ctx *seago.Context) {
            ctx.Data["Name"] = "jeremy"
            ctx.HTML(200, "hello") // 200 为响应码
        })
        
        s.Run()
    }

## 自定义选项

该服务允许接受一个参数来进行自定义选项（`seago.RenderOptions`）：

    package main

    import "github.com/seago/seago"

    func main() {
        s := seago.Classic()
        s.Use(`seago.Renderer`(`seago.RenderOptions`{
            // 模板文件目录，默认为 "templates"
            Directory: "templates",
            // 模板文件后缀，默认为 [".tmpl", ".html"]
            Extensions: []string{".tmpl", ".html"},
            // 模板函数，默认为 []
            Funcs: []template.FuncMap{map[string]interface{}{
                "AppName": func() string {
                    return "Seago"
                },
                "AppVer": func() string {
                    return "1.0.0"
                },
            }},
            // 模板语法分隔符，默认为 ["{{", "}}"]
            Delims: seago.Delims{"{{", "}}"},
            // 追加的 Content-Type 头信息，默认为 "UTF-8"
            Charset: "UTF-8",
            // 渲染具有缩进格式的 JSON，默认为不缩进
            IndentJSON: true,
            // 渲染具有缩进格式的 XML，默认为不缩进
            IndentXML: true,
            // 渲染具有前缀的 JSON，默认为无前缀
            PrefixJSON: []byte("Seago"),
            // 渲染具有前缀的 XML，默认为无前缀
            PrefixXML: []byte("Seago"),
            // 允许输出格式为 XHTML 而不是 HTML，默认为 "text/html"
            HTMLContentType: "text/html",
        }))        
        // ...
    }

## Pongo2 模板引擎

该服务可以通过函数 `pongo2.Pongoer` 来注入，并通过类型 `seago.Render`来体现。该服务为可选，一般情况下可直接使用 `*seago.Context.Render`。该服务使用 Pongo2 v3 模板引擎来渲染 HTML。如果想要了解更多有关使用方面的信息，请参见 官方文档。

## 使用示例

假设您的应用拥有以下目录结构：

    main/
        |__ main.go
        |__ templates/
                |__ hello.tmpl
    hello.tmpl：

    <h1>Hello {{Name}}</h1>
    main.go：

    package main

    import (
        "github.com/seago/seago"
        "github.com/seago-contrib/pongo2"
    )

    func main() {
        s := seago.Classic()
        s.Use(pongo2.Pongoer())
        
        s.Get("/", func(ctx *seago.Context) {
            ctx.Data["Name"] = "jeremy"
            ctx.HTML(200, "hello") // 200 is the response code.
        })
        
        s.Run()
    }

## 自定义选项

该服务允许接受一个参数来进行自定义选项（pongo2.Options）：

    package main

    import (
        "github.com/seago/seago"
        "github.com/seago-contrib/pongo2"
    )

    func main() {
        s := seago.Classic()
        s.Use(pongo2.Pongoer(pongo2.Options{
            // 模板文件目录，默认为 "templates"
            Directory: "templates",
            // 模板文件后缀，默认为 [".tmpl", ".html"]
            Extensions: []string{".tmpl", ".html"},
            // 追加的 Content-Type 头信息，默认为 "UTF-8"
            Charset: "UTF-8",
            // 渲染具有缩进格式的 JSON，默认为不缩进
            IndentJSON: true,
            // 渲染具有缩进格式的 XML，默认为不缩进
            IndentXML: true,
            // 允许输出格式为 XHTML 而不是 HTML，默认为 "text/html"
            HTMLContentType: "text/html",
        }))        
        // ...
    }

## 模板集

当您的应用存在多套模板时，就需要使用模板集来实现运行时动态设置需要渲染的模板。

* Go 模板引擎的使用方法：

        // ...
        s.Use(`seago.Renderer`(`seago.RenderOptions`{
            Directory: "templates/default",
        }, "theme1:templates/theme1", "theme2:templates/theme2"))

        s.Get("/foobar", func(ctx *seago.Context) {
            ctx.HTML(200, "hello", "jeremy")
        })

        s.Get("/foobar1", func(ctx *seago.Context) {
            ctx.HTMLSet(200, "theme1", "hello", "jeremy")
        })

        s.Get("/foobar2", func(ctx *seago.Context) {
            ctx.HTMLSet(200, "theme2", "hello", "jeremy")
        })
        // ...

* Pongo2 模板引擎的使用方法：

        // ...
        s.Use(`pongo2.Pongoer`(`pongo2.Options`{
            Directory: "templates/default",
        }, "theme1:templates/theme1", "theme2:templates/theme2"))

        s.Get("/foobar", func(ctx *seago.Context) {
            ctx.HTML(200, "hello", "jeremy")
        })

        s.Get("/foobar1", func(ctx *seago.Context) {
            ctx.HTMLSet(200, "theme1", "hello", "jeremy")
        })

        s.Get("/foobar2", func(ctx *seago.Context) {
            ctx.HTMLSet(200, "theme2", "hello", "jeremy")
        })
        // ...
正如您所看到的那样，其实就是 2 个方法的不同：`seago.Renderers` 和 `pongo2.Pongoers`。

第一个配置参数用于指定默认的模板集和配置选项，之后则是一个模板集名称和目录（通过 : 分隔）的列表。

如果您的模板集名称和模板集路径的最后一部分相同，则可以省略名称：

        // ...
        s.Use(`seago.Renderer`(`seago.RenderOptions`{
            Directory: "templates/default",
        }, "templates/theme1", "templates/theme2"))

        s.Get("/foobar", func(ctx *seago.Context) {
            ctx.HTML(200, "hello", "jeremy")
        })

        s.Get("/foobar1", func(ctx *seago.Context) {
            ctx.HTMLSet(200, "theme1", "hello", "jeremy")
        })

        s.Get("/foobar2", func(ctx *seago.Context) {
            ctx.HTMLSet(200, "theme2", "hello", "jeremy")
        })
        // ...

## 模板集辅助方法

检查某个模板集是否存在：

        // ...
        s.Get("/foobar", func(ctx *seago.Context) {
            ok := ctx.HasTemplateSet("theme2")
            // ...
        })
        // ...
修改模板集的目录：

        // ...
        s.Get("/foobar", func(ctx *seago.Context) {
            ctx.SetTemplatePath("theme2", "templates/new/theme2")
            // ...
        })
        // ...

## 小结

也许您已经发现，除了在 HTML 语法上的不同之外，两款引擎在代码层面的用法是完全一样的。

如果您只是想要得到 HTML 渲染后的结果，则可以调用方法 `*seago.Context.Render.HTMLString`：

    package main

    import "github.com/seago/seago"

    func main() {
        s := seago.Classic()
        s.Use(`seago.Renderer`())
        
        s.Get("/", func(ctx *seago.Context) {
            ctx.Data["Name"] = "jeremy"
            output, err := ctx.HTMLString("hello")
            // 进行其它操作
        })
        
        s.Run()
    }

* 渲染 XML、JSON 和原始数据

相对于渲染 HTML 而言，渲染 XML、JSON 和原始数据的工作要简单的多。

        package main

        import (
            "github.com/seago/seago"
        )

        type Person struct {
            Name string
            Age  int
            Sex  string
        }

        func main() {
            s := seago.Classic()
            s.Use(`seago.Renderer`())

            s.Get("/xml", func(ctx *seago.Context) {
                p := Person{"Seago", 21, "male"}
                ctx.XML(200, &p)
            })
            s.Get("/json", func(ctx *seago.Context) {
                p := Person{"Seago", 21, "male"}
                ctx.JSON(200, &p)
            })
            s.Get("/raw", func(ctx *seago.Context) {
                ctx.RawData(200, []byte("raw data goes here"))
            })

            s.Run()
        }

* 响应状态码、错误和重定向

如果您希望响应指定状态码、错误和重定向操作，则可以参照以下代码：

        package main

        import (
            "github.com/seago/seago"
        )

        func main() {
            s := seago.Classic()
            s.Use(`seago.Renderer`())

            s.Get("/status", func(ctx *seago.Context) {
                ctx.Status(403)
            })
            s.Get("/error", func(ctx *seago.Context) {
                ctx.Error(500, "Internal Server Error")
            })
            s.Get("/redirect", func(ctx *seago.Context) {
                ctx.Redirect("/") // 第二个参数为响应码，默认为 302
            })

            s.Run()
        }

* 运行时修改模板路径

如果您希望在运行时修改应用的模板路径，则可以调用方法 `*seago.Context.SetTemplatePath`。需要注意的是，修改操作是全局生效的，而不只是针对当前请求。

使用示例

假设您的应用拥有以下目录结构：

    main/
        |__ main.go
        |__ templates/
                |__ hello.tmpl
        |__ templates2/
                |__ hello.tmpl
    templates/hello.tmpl：

    <h1>Hello {{.Name}}</h1>
    templates2/hello.tmpl：

    <h1>What's up, {{.Name}}</h1>
    main.go：

    package main

    import (
        "github.com/seago/seago"
    )

    func main() {
        s := seago.Classic()
        s.Use(`seago.Renderer`())

        s.Get("/old", func(ctx *seago.Context) {
            ctx.Data["Name"] = "Seago"
            ctx.HTML(200, "hello")
            ctx.SetTemplatePath("templates2")
        })
        s.Get("/new", func(ctx *seago.Context) {
            ctx.Data["Name"] = "Seago"
            ctx.HTML(200, "hello")
        })

        s.Run()
    }
当您首次请求 /old 页面时，响应结果为 `<h1>Hello Seago</h1>`，然后便执行了修改模板路径为 template2。此时，当您请求 /new 页面时，响应结果会变成 `<h1>What's up, Seago</h1>`。