package model

import (
	"errors"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"sysbase/config"
)

type Model interface {
	InitTable() error
}

var (
	_ Model = Pod{}
	_ Model = Resource{}
	_ Model = K8sCluster{}
	_ Model = K8sClusterResource{}
	_ Model = PodResource{}
	_ Model = TenantPod{}

	DBConn map[string]*gorm.DB
)

const DBName = "sysbase"

func InitDB(dbs config.Db) error {

	DBConn = make(map[string]*gorm.DB)
	for _, db := range dbs {
		switch db.Type {
		case "sqlite":
			if len(db.Dsn) == 0 {
				return errors.New(fmt.Sprintf("%s 数据库配置为空", db.Name))
			}
			var err error
			conn, err := gorm.Open(sqlite.Open(db.Dsn), &gorm.Config{
				NamingStrategy: schema.NamingStrategy{
					TablePrefix:   db.TablePrefix,
					SingularTable: true,
				},
			})
			if err != nil {
				return err
			}

			sqlDB, err := conn.DB()
			if err != nil {
				return err
			}

			// 设置打开数据库连接的最大数量
			sqlDB.SetMaxOpenConns(db.MaxOpenConns)

			DBConn[db.Name] = conn
			break
		}
	}

	_ = Resource{}.InitTable()
	_ = K8sCluster{}.InitTable()
	_ = Pod{}.InitTable()
	_ = K8sClusterResource{}.InitTable()
	_ = PodResource{}.InitTable()
	_ = TenantPod{}.InitTable()

	return nil
}
