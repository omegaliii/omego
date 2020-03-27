package main

import (
	"net/http"
	"omego"
	"fmt"
	"time"
	"html/template"
)

type student struct {
	Name string
	Age  int8
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := omego.New()
	r.Use(omego.Logger())
	r.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")

	stu1 := &student{Name: "omegoktutu", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	
	r.GET("/", func(c *omego.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})

	r.GET("/students", func(c *omego.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", omego.H{
			"title":  "omego",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	r.GET("/date", func(c *omego.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", omego.H{
			"title": "omego",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	r.Run(":9999")
}