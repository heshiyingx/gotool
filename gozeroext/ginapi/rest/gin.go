package rest

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

type RunOption func(*gin.Engine)

func MustNewServer(c rest.RestConf, opts ...RunOption) *gin.Engine {
	server, err := NewServer(c, opts...)
	if err != nil {
		logx.Must(err)
	}
	//engine := gin.Default()
	return server
}
func NewServer(c rest.RestConf, opts ...RunOption) (*gin.Engine, error) {

	engine := gin.Default()
	for _, opt := range opts {
		opt(engine)
	}
	return engine, nil
}
func StartServer(s *gin.Engine, c rest.RestConf) {
	if err := c.SetUp(); err != nil {
		logx.Must(err)
	}
	srv := &http.Server{
		Addr:    c.Host + ":" + strconv.Itoa(c.Port),
		Handler: s,
		//ReadTimeout: time.Millisecond * time.Duration(c.Timeout),

	}
	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
