package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
	"words/application"
	"words/interfaces/response"
)

type StudyWord struct {
}

// index页面
func (b StudyWord) Index(ctx *gin.Context) {
	h := gin.H{}
	bookList := application.NewBookApplication(ctx.Request.Context()).BookList()

	h["bookList"] = bookList
	ctx.HTML(http.StatusOK, "index.html", h)
}

// 背单词页面
func (b StudyWord) Study(ctx *gin.Context) {
	h := gin.H{}
	bookId := ctx.Query("bookId")
	planId := ctx.Query("planId")

	var studyPlan response.StudyWordResp
	var err error
	if bookId != "" {
		var pid int64
		pid, err = application.NewStudyApplication(ctx.Request.Context()).GenerateStudyPlan(bookId)
		ctx.Redirect(http.StatusFound, fmt.Sprintf("study?planId=%d", pid))
		return
	} else if planId != "" {
		studyPlan, err = application.NewStudyApplication(ctx.Request.Context()).GetStudyPlan(cast.ToInt64(planId))
	}
	if err != nil {
		return
	}
	h["studyPlan"] = studyPlan

	ctx.HTML(http.StatusOK, "study.html", h)
}

// 标记为已完成
func (b StudyWord) MarkStudied(ctx *gin.Context) {
	h := gin.H{}
	planId := ctx.Query("planId")
	err := application.NewStudyApplication(ctx.Request.Context()).MarkStudied(cast.ToInt64(planId))
	if err != nil {
		h["err"] = err.Error()
	}
	ctx.JSON(http.StatusOK, h)
}
