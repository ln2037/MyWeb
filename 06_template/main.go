package main

/*
(1) render array
$ curl http://localhost:9999/date
<html>
<body>
    <p>hello, myweb</p>
    <p>Date: 2019-08-17</p>
</body>
</html>
*/

/*
(2) custom render function
$ curl http://localhost:9999/students
<html>
<body>
    <p>hello, myweb</p>
    <p>0: mywebktutu is 20 years old</p>
    <p>1: Jack is 22 years old</p>
</body>
</html>
*/

/*
(3) serve static files
$ curl http://localhost:9999/assets/css/mywebktutu.css
p {
    color: orange;
    font-weight: 700;
    font-size: 20px;
}
*/

import (
	"MyWeb/06_template/myweb"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := myweb.New()
	r.Use(myweb.Logger())
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("06_template/templates/*")
	r.Static("/assets", "./static")
	stu1 := &student{Name: "shiwen", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	r.GET("/", func(context *myweb.Context) {
		context.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.GET("/students", func(context *myweb.Context) {
		context.HTML(http.StatusOK, "arr.tmpl", myweb.H{
			"title":  "gee",
			"stuArr": [2]*student{stu1, stu2},
		})
	})
	r.GET("/date", func(context *myweb.Context) {
		context.HTML(http.StatusOK, "custom_func.tmpl", myweb.H{
			"title": "gee",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})
	r.Run(":9999")
}
