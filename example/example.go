package main

import (
	"fmt"
	"io"
	. "middleware"
	"os"
	"strconv"
	"web"
	"web/context"
)

func main() {
	DefaultMiddleware.Add("test", "test middleware")
	seago := web.NewSeago(":8080")

	seago.Get("/test", func(ctx *context.Context) string {

		test := ctx.GetParam("test_get").String()
		i := ctx.GetParam("id_get").Int()
		return "hello " + test + " " + strconv.Itoa(i) + " " + DefaultMiddleware.Get("test").(string)
	})
	seago.Post("/test", func(ctx *context.Context) string {
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
	seago.Run()

}
