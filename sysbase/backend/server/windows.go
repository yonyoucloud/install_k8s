//go:build windows
// +build windows

package server

import (
	"fmt"

	//"gopkg.in/gin-contrib/pprof.v1"

	"sysbase/args"
	"sysbase/config"
	"sysbase/model"
	"sysbase/router"
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

	router.Run(fmt.Sprintf("%s:%d", args.Holder.GetBindAddress(), args.Holder.GetPort()))

	return nil
}
