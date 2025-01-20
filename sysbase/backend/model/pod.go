package model

import (
	"strconv"
	"strings"
)

// 设置Pods表字段
type Pod struct {
	ID           uint   `gorm:"primaryKey;autoIncrement;comment:自增ID"`
	K8sClusterID uint   `gorm:"not null;default:0;comment:对应的K8S集群ID，一个Pod只能对应一个K8S集群，一个K8S集群可以对应多个Pod"`
	Name         string `gorm:"type:char(50);not null;default:'';comment:Pod名，规范：中国大陆-阿里云-100"`
	Code         string `gorm:"type:char(10);not null;default:'';index:idx_host,unique;comment:Pod的代码号，作为二级域名，和domain字段作为Pod的唯一访问入口"`
	Domain       string `gorm:"type:char(20);not null;default:'';index:idx_host,unique;comment:Pod的根域名，和code字段作为Pod的唯一访问入口"`
	Cap          uint32 `gorm:"not null;default:0;comment:Pod可以容纳的租户个数"`
	//Iaas         string     `gorm:"type:enum('aliyun','huaweiyun','tencent', 'amazon');default:'aliyun';comment:Pod所属IaaS厂商标识"`
	Iaas       string     `gorm:"type:varchar(15);default:'aliyun';comment:Pod所属IaaS厂商标识"`
	UpdatedAt  int64      `gorm:"autoUpdateTime:milli"`
	CreatedAt  int64      `gorm:"autoCreateTime:milli"`
	K8sCluster K8sCluster `gorm:"foreignKey:K8sClusterID;references:ID;constraint:OnDelete:NO ACTION,OnUpdate:NO ACTION"`
}

func (p Pod) InitTable() error {
	// 设置表信息
	//migrator := DBConn[DBName].Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Pods表'").Migrator()
	migrator := DBConn[DBName].Set("gorm:table_options", "").Migrator()

	// 判断表是否存在
	if migrator.HasTable(&Pod{}) {
		// 存在就自动适配表，也就说原先没字段的就增加字段
		return migrator.AutoMigrate(&Pod{})
	}

	// 不存在就创建新表
	return migrator.CreateTable(&Pod{})
}

func (p Pod) Insert(resourceID string) (Pod, error) {
	tx := DBConn[DBName].Create(&p)
	if tx.Error == nil {
		p.PodResourceInsert(resourceID)
	}

	return p, tx.Error
}

func (p Pod) List() ([]Pod, error) {
	var ps []Pod
	result := DBConn[DBName].Order("id desc").Find(&ps)
	return ps, result.Error
}

func (p Pod) Delete() error {
	tx := DBConn[DBName].Delete(&p)
	return tx.Error
}

func (p Pod) Edit(pData Pod, resourceID string) error {
	tx := DBConn[DBName].Model(&p)
	if tx.Error != nil {
		return tx.Error
	}

	err := tx.Updates(pData).Error
	if err == nil {
		p.PodResourceInsert(resourceID)
	}
	return err
}

func (p Pod) PodResourceInsert(resourceID string) {
	resourceIDArr := strings.Split(resourceID, ",")
	for _, rID := range resourceIDArr {
		rIDInt, _ := strconv.Atoi(rID)
		if rIDInt <= 0 {
			continue
		}
		pr := PodResource{
			PodID:      p.ID,
			ResourceID: uint(rIDInt),
		}
		pr.Insert()
	}
}

func (p Pod) Get() (Pod, error) {
	var pData Pod
	err := DBConn[DBName].Where(&p).Find(&pData).Error
	return pData, err
}
