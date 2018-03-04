Paginator [![Build Status](https://travis-ci.org/finejian/paginator.svg?branch=master)](https://travis-ci.org/finejian/paginator)
=========

Paginator是一个Golang分页工具，参考了原作者的项目结构，重新编写了项目核心分页处理逻辑，新增获取页码 URL 和生成 HTML 代码功能。

## 安装

	go get github.com/finejian/paginator

## 开始使用

获取分页页码:

```go
package main

import "github.com/finejian/paginator"

func main() {
	// 参数：传入数据总行数
	p := paginator.New(43)
	
	// 将 p 当作 template 对象 命名为 page 传到 "simple.html"
	// ...

	// 如果使用的是 gin 可以参考代码：
	router := gin.Default()
	router.LoadHTMLFiles("simple.html", "simple.html")
	router.GET("/simple", func(c *gin.Context) {
		p := paginator.New(43)

		c.HTML(http.StatusOK, "simple.html", gin.H{
			"page": p,
		})
	})
	router.Run(":8080")
}
```

`simple.html`

```html
{{if not .page.IsFirst}}[First](1){{end}}
{{if .page.HasPrevious}}[Previous]({{.page.Previous}}){{end}}

{{range .page.Pages}}
	{{if eq .Num -1}}
	...
	{{else}}
	{{.Num}}{{if .IsCurrent}}(current){{end}}
	{{end}}
{{end}}

{{if .page.HasNext}}[Next]({{.page.Next}}){{end}}
{{if not .page.IsLast}}[Last]({{.page.TotalPages}}){{end}}
```

输出结果:

```
[First](1) [Previous](2) ... 2 3(current) 4 ... [Next](4) [Last](5)
```


如果你在 html 代码中直接获取到相应页码的 url，必须在初始话 paginator 时候调用 Request 方法

```go
// 调用 Request 方法传入 http request 对象：
// c 为 网络请求的 Context 上下文对象
p := paginator.New(43).Request(c.Request)
```

在确保调用 Request 方法传入正确内容的情况下，可在 html 代码中得到如下使用
paginator 会获取 request 请求的 query 参数，并会在每个页面 url 中还原这些 query 请求参数
```html
{{if not .page.IsFirst}}<a href="{{.page.FristURL}}">首页</a>{{end}}
{{if .page.HasPrevious}}<a href="{{.page.PreviousURL}}">上一页</a>{{end}}

{{range .page.PageURLs}}
	{{if eq .Num -1}}
	...
	{{else}}
	<a href="{{.Path}}">{{.Num}}{{if .IsCurrent}}(current){{end}}</a>
	{{end}}
{{end}}

{{if .page.HasNext}}<a href="{{.page.NextURL}}">下一页</a>{{end}}
{{if not .page.IsLast}}<a href="{{.page.LastURL}}">尾页</a>{{end}}
```

也可以在 html 代码中直接使用如下方法，一次获取整段分页功能 html 代码
```html
{{.page.PageTemp}}
```

## 获取帮助

- [API Documentation](https://gowalker.org/github.com/finejian/paginator)
- [File An Issue](https://github.com/finejian/paginator/issues/new)

## License

This project is under Apache v2 License. See the [LICENSE](LICENSE) file for the full license text.