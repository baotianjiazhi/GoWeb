package main

import (
	"bluebell/controller"
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/logger"
	"bluebell/pkg/snowflake"
	"bluebell/routers"
	"bluebell/settings"
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"go.uber.org/zap"
)


func init() {


	var filename string
	flag.StringVar(&filename, "filename", "conf/app.ini", "文件名")
	if len(filename) <= 0 {
		panic("need config file.eg: webapp config.ini")
		return
	}
	// 1. 加载配置--视频中使用的是viper但是我用的是goini这个库
	settings.Setup(filename)

	if settings.AppSetting.Mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	// 2. 初始化日志
	if err := logger.Init(); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
	}

	zap.L().Debug("logger init success")
	// 3. 初始化MySQL连接
	if err := mysql.InitDB(); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
	}
	// 4. 初始化Redis连接
	if err := redis.Init(); err != nil {
		fmt.Printf("init redist failed, err:%v\n", err)
	}

	// 5. 初始化雪花算法
	if err := snowflake.Init(settings.SnowFlakeSetting.StartTime, settings.SnowFlakeSetting.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}

	// 6. 国际化
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init trans failed, err:%v\n", err)
		return
	}

}
func main() {

	// 5. 注册路由
	r := routers.SetUp()

	defer mysql.Close()
	defer redis.Close()
	defer zap.L().Sync()
	// 启动服务（优雅关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", settings.AppSetting.Port),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen error", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // 阻塞在此， 当接收到上述两种信号时才会往下进行
	zap.L().Info("ShutDown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（处理未处理完的请求再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server shutdown:", zap.Error(err))
	}

	zap.L().Info("server exiting")
}
