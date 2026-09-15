package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"words/domain/repository"
)

func main() {
	cfg := LoadConfig()

	// 本地文件存储目录，可通过 WORDS_DATA_DIR 或 -data 覆盖，默认 data
	repository.Default = repository.NewStore(cfg.DataDir)

	if cfg.Release {
		gin.SetMode(gin.ReleaseMode)
	}
	g := gin.Default()
	g.SetTrustedProxies(nil)
	SetRouters(g)

	server := &http.Server{
		Addr:           cfg.Addr,
		Handler:        g,
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		fmt.Printf("server start success pid:%d addr:%s data:%s\n", os.Getpid(), cfg.Addr, repository.Default.Dir())
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(os.Stderr, "服务启动失败: %v\n", err)
			os.Exit(1)
		}
	}()

	// 优雅退出：收到信号后停止接收新请求，并给在途请求一段处理时间。
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		_ = server.Close()
	}
	fmt.Println("程序关闭")
}
