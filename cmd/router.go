package main

import (
	"html/template"
	"strings"

	"github.com/gin-gonic/gin"

	"words/interfaces/api"
)

func SetRouters(r *gin.Engine) {

	r.SetFuncMap(template.FuncMap{
		// firstRune 取首个字符，用于生成单词本图标文字（按 rune 取，中文安全）
		"firstRune": func(s string) string {
			for _, c := range strings.TrimSpace(s) {
				return string(c)
			}
			return "#"
		},
	})

	r.LoadHTMLFiles("./web/index.html", "./web/study.html")

	// 静态资源（样式表等）。只暴露 web/static，避免把模板文件当静态资源暴露出去。
	r.Static("/static", "./web/static")

	bookApi := new(api.StudyWord)

	html := r.Group("")
	{
		html.GET("", bookApi.Index)
		html.GET("study", bookApi.Study)
		html.PUT("markStudied", bookApi.MarkStudied)
	}

	appApi := r.Group("/app")
	{
		appApi.GET("", bookApi.Index)
	}
}
