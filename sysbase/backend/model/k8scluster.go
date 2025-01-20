package model

import (
	"strconv"
	"strings"
)

// 设置K8sClusters表字段
type K8sCluster struct {
	ID        uint   `gorm:"primaryKey;autoIncrement;comment:自增ID"`
	Name      string `gorm:"type:varchar(255);not null;default:'';comment:集群名称"`
	UpdatedAt int64  `gorm:"autoUpdateTime:milli"`
	CreatedAt int64  `gorm:"autoCreateTime:milli"`
}

func (kc K8sCluster) InitTable() error {
	// 设置表信息
	//migrator := DBConn[DBName].Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='K8sClusters表'").Migrator()
	migrator := DBConn[DBName].Set("gorm:table_options", "").Migrator()

	// 判断表是否存在
	if migrator.HasTable(&K8sCluster{}) {
		// 存在就自动适配表，也就说原先没字段的就增加字段
		return migrator.AutoMigrate(&K8sCluster{})
	}

	// 不存在就创建新表
	return migrator.CreateTable(&K8sCluster{})
}

func (kc K8sCluster) Insert(resourceID string) (K8sCluster, error) {
	tx := DBConn[DBName].Create(&kc)
	err := tx.Error
	if err == nil {
		kc.K8sClusterResourceInsert(resourceID)
	}
	return kc, err
}

func (kc K8sCluster) List() ([]K8sCluster, error) {
	var kcs []K8sCluster
	result := DBConn[DBName].Order("id desc").Find(&kcs)
	return kcs, result.Error
}

func (kc K8sCluster) Delete() error {
	tx := DBConn[DBName].Delete(&kc)
	return tx.Error
}

func (kc K8sCluster) Edit(kcData K8sCluster, resourceID string) error {
	tx := DBConn[DBName].Model(&kc)
	if tx.Error != nil {
		return tx.Error
	}

	err := tx.Updates(kcData).Error
	if err == nil {
		kc.K8sClusterResourceInsert(resourceID)
	}
	return err
}

func (kc K8sCluster) K8sClusterResourceInsert(resourceID string) {
	resourceIDArr := strings.Split(resourceID, ",")
	for _, rID := range resourceIDArr {
		rIDInt, _ := strconv.Atoi(rID)
		if rIDInt <= 0 {
			continue
		}
		kcr := K8sClusterResource{
			K8sClusterID: kc.ID,
			ResourceID:   uint(rIDInt),
		}
		kcr.Insert()
	}
}

func (kc K8sCluster) Get() (K8sCluster, error) {
	var kcData K8sCluster
	err := DBConn[DBName].Where(&kc).Find(&kcData).Error
	return kcData, err
}
