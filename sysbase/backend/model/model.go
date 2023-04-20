package model

import (
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	"sysbase/config"
)

type Model interface {
	InitTable()
}

var (
	_ Model = Pod{}
	_ Model = Resource{}
	_ Model = K8sCluster{}
	_ Model = K8sClusterResource{}
	_ Model = PodResource{}
	_ Model = TenantPod{}

	db *gorm.DB
)

func InitDB(c config.Mysql) {
	var err error
	db, err = gorm.Open(mysql.Open(c.MasterDsn), &gorm.Config{})

	var sources []gorm.Dialector
	for _, dsn := range c.SourcesDsn {
		sources = append(sources, mysql.Open(dsn))
	}

	var replicas []gorm.Dialector
	for _, dsn := range c.ReplicasDsn {
		replicas = append(replicas, mysql.Open(dsn))
	}

	db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  sources,
		Replicas: replicas,
		// sources/replicas 负载均衡策略
		Policy: dbresolver.RandomPolicy{},
	}).
		//设置了连接可复用的最大时间
		SetConnMaxIdleTime(time.Hour).
		//设置了连接可复用的最大时间
		SetConnMaxLifetime(c.ConnMaxLifetime).
		// 设置空闲连接池中连接的最大数量
		SetMaxIdleConns(c.MaxIdleConns).
		// 设置打开数据库连接的最大数量
		SetMaxOpenConns(c.MaxOpenConns))

	if err != nil {
		log.Panicln(err)
	}

	Resource{}.InitTable()
	K8sCluster{}.InitTable()
	Pod{}.InitTable()
	K8sClusterResource{}.InitTable()
	PodResource{}.InitTable()
	TenantPod{}.InitTable()
}
