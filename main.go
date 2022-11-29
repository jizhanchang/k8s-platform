package main

import (
	"context"
	"flag"
	"fmt"
	"k8s-platfrom/logger"
	routes "k8s-platfrom/routers"
	"k8s-platfrom/settings"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {

	//1、加载配置
	cfg := flag.String("c", "cfg.yml", "configuration file")
	// 参数解析
	flag.Parse()

	if err := settings.Init(cfg); err != nil {
		log.Fatalf("init setttings failed, err:[%s]\n", err.Error())
	}
	//2.初始化日志库
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		log.Fatalf("init logger failed, err:[%s]\n", err.Error())
	}
	// 退出之前把内存中的日志dump到文件中
	defer zap.L().Sync()
	zap.L().Info("init zap logger success")
	r := routes.SetUp(settings.Conf.Mode)

	zap.L().Debug("register route success...")

	//6、启动服务（优雅重启关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.AppConfig.Port),
		Handler: r,
	}
	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: %s\n", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}
