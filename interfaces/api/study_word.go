package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"words/application"
)

type StudyWord struct {
}

// Index 首页：单词本列表 + 最近学习记录 + 顶部统计
func (b StudyWord) Index(ctx *gin.Context) {
	data := application.NewBookApplication(ctx.Request.Context()).IndexData()
	ctx.HTML(http.StatusOK, "index.html", data)
}

// Study 背单词页面：
//   - 带 bookId：生成一份新计划并重定向到具体计划
//   - 带 planId：展示已有计划
func (b StudyWord) Study(ctx *gin.Context) {
	bookId := ctx.Query("bookId")
	planId := ctx.Query("planId")

	switch {
	case bookId != "":
		pid, err := application.NewStudyApplication(ctx.Request.Context()).GenerateStudyPlan(bookId)
		if err != nil {
			b.renderError(ctx, fmt.Errorf("生成学习计划失败: %w", err))
			return
		}
		ctx.Redirect(http.StatusFound, fmt.Sprintf("study?planId=%d", pid))
	case planId != "":
		studyPlan, err := application.NewStudyApplication(ctx.Request.Context()).GetStudyPlan(cast.ToInt64(planId))
		if err != nil {
			b.renderError(ctx, fmt.Errorf("获取学习计划失败: %w", err))
			return
		}
		ctx.HTML(http.StatusOK, "study.html", gin.H{"studyPlan": studyPlan})
	default:
		b.renderError(ctx, fmt.Errorf("缺少 bookId 或 planId 参数"))
	}
}

// MarkStudied 标记计划为已完成
func (b StudyWord) MarkStudied(ctx *gin.Context) {
	planId := cast.ToInt64(ctx.Query("planId"))
	if err := application.NewStudyApplication(ctx.Request.Context()).MarkStudied(planId); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}

// renderError 记录错误并返回一个可读的错误页，避免用户看到空白页。
func (b StudyWord) renderError(ctx *gin.Context, err error) {
	log.Printf("[%s %s] %v", ctx.Request.Method, ctx.Request.URL.Path, err)
	ctx.String(http.StatusInternalServerError, "出错了：%s", err.Error())
}
