package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"words/domain/repository"
)

func main() {
	// 初始化mongo
	repository.InitMongo()

	gin.SetMode(gin.ReleaseMode)
	g := gin.Default()
	SetRouters(g)
	server := &http.Server{
		Addr:           ":8900",
		Handler:        g,
		ReadTimeout:    time.Duration(3) * time.Second,
		WriteTimeout:   time.Duration(3) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	QuitSignal(func() {
		repository.CloseMongo()
		_ = server.Close()
		fmt.Println("程序关闭")
	})
}

func QuitSignal(quitFunc func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	fmt.Printf("server start success pid:%d\n", os.Getpid())
	for s := range c {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			quitFunc()
			return
		default:
			return
		}
	}
}
