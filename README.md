## Seago

[![Build Status](https://drone.io/github.com/seago/seago/status.png)](https://drone.io/github.com/seago/seago/latest)

Seago is a Go Framework.

Seagois a simple web framework. The router package is inspired by web.go and martinit.

## Features
* RESTful support
* Session support
* Cache support
* Middleware support
##Getting Started

After installing Go and setting up your GOPATH, create your first .go file. We'll call it server.go.
~~~ go
package main

import (
	"fmt"
	"github.com/seago/seago"
	"github.com/seago/seago/context"
	. "github.com/seago/seago/middleware"
	"io"
	"os"
	"strconv"
)

func main() {
	DefaultMiddleware.Add("test", "test middleware")
	app := seago.NewSeago(":8080")

	//http://localhost:8080/test/Miller
	app.Get("/test/:test_get", func(ctx *context.Context, test string) string {
		i := ctx.GetParam("id_get").Int()
		return "hello " + test + " " + strconv.Itoa(i) + " " + DefaultMiddleware.Get("test").(string)
	})
	app.Post("/test", func(ctx *context.Context) string {
		test := ctx.GetParam("test_post").String()
		i := ctx.GetParam("id_post").Int()
		file := ctx.File["file"]
		f, err := file.Open()
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()
		fi, err := os.Create(file.Filename)
		if err != nil {
			fmt.Println(err)
		}
		_, err = io.Copy(fi, f)
		if err != nil {
			fmt.Println(err)
		}
		defer fi.Close()

		return "hello " + test + " " + strconv.Itoa(i)
	})
	app.Get("/", func() string {
		return `<html>
					<head>
					<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
					<title>test file upload</title>
					</head>
					<body>
						<form action="http://localhost:8080/test" method="post" enctype="multipart/form-data">
							test:<input type="text" name="test_post" value="" /><br />
							id:<input type="text" name="id_post" value="" /><br />
							<input type="file" name="file" /><br />
							<input type="submit" name="submit" value="upload" />
						</form>
					</body>
				</html>`
	})
	app.Profile()
	app.Server.SetMaxMemory(200 << 20)
	app.Run()

}
~~~
Then install the Seago package (go 1.1 and greater is required):

go get github.com/seago/seago
Then run your server:

go run server.go
You will now have a Martini webserver running on localhost:8080.