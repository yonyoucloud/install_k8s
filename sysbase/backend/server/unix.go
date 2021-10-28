// +build darwin freebsd linux netbsd openbsd

package server

import (
	"fmt"
	"log"
	"syscall"

	"github.com/fvbock/endless"

	//"gopkg.in/gin-contrib/pprof.v1"

	"git.yonyou.com/sysbase/backend/args"
	"git.yonyou.com/sysbase/backend/config"
	"git.yonyou.com/sysbase/backend/model"
	"git.yonyou.com/sysbase/backend/router"
)

type Server struct {
	config *config.Config
}

func NewServer(c *config.Config) *Server {
	s := &Server{
		config: c,
	}
	return s
}

func (s *Server) Run() error {
	// 连接数据库
	model.InitDB(s.config.Mysql)

	// 初始化路由
	router := router.InitRouter(s.config)

	/*
	   pprof.Register(router, &pprof.Options{
	       // default is "debug/pprof"
	       RoutePrefix: "debug/pprof",
	   })
	*/

	// 平滑重启
	/*
	   kill -SIGHUP 2591 will trigger a fork/restart
	   kill -SIGINT[SIGTERM] 2591 will trigger a shutdown of the server (it will finish running requests)
	*/
	server := endless.NewServer(fmt.Sprintf("%s:%d", args.Holder.GetBindAddress(), args.Holder.GetPort()), router)
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d\n", syscall.Getpid())
		// save it somehow
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}

	return nil
}
