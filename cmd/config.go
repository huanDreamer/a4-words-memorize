package main

import (
	"flag"
	"os"
	"time"
)

// Config 是服务运行配置。取值优先级：命令行参数 > 环境变量 > 默认值。
type Config struct {
	Addr         string        // 监听地址
	DataDir      string        // 本地文件存储目录
	Release      bool          // 是否使用 gin release 模式
	ReadTimeout  time.Duration // 读超时
	WriteTimeout time.Duration // 写超时
}

// LoadConfig 解析命令行参数与环境变量，返回运行配置。
func LoadConfig() Config {
	var cfg Config

	defAddr := envOr("WORDS_ADDR", ":8900")
	defData := envOr("WORDS_DATA_DIR", "data")
	defRelease := os.Getenv("GIN_MODE") == "release" || os.Getenv("WORDS_RELEASE") == "1"

	flag.StringVar(&cfg.Addr, "addr", defAddr, "HTTP 监听地址")
	flag.StringVar(&cfg.DataDir, "data", defData, "数据目录")
	flag.BoolVar(&cfg.Release, "release", defRelease, "是否使用 production(release) 模式")
	flag.DurationVar(&cfg.ReadTimeout, "read-timeout", 10*time.Second, "HTTP 读超时")
	flag.DurationVar(&cfg.WriteTimeout, "write-timeout", 10*time.Second, "HTTP 写超时")
	flag.Parse()

	return cfg
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
