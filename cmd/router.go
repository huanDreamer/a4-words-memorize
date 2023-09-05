package main

import (
	"github.com/gin-gonic/gin"
	"words/interfaces/api"
)

func SetRouters(r *gin.Engine) {

	r.LoadHTMLFiles("./web/index.html", "./web/study.html")

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
