package model

import (
	"fmt"
	"strings"
)

// 设置K8sClusterResource表字段
type K8sClusterResource struct {
	ID           uint       `gorm:"primaryKey;autoIncrement;comment:自增ID"`
	K8sClusterID uint       `gorm:"not null;default:0;comment:K8sClusterID"`
	ResourceID   uint       `gorm:"not null;default:0;unique;comment:ResourceID，该字段唯一，确保该资源被某一个K8S集群唯一使用"`
	UpdatedAt    int64      `gorm:"autoUpdateTime:milli"`
	CreatedAt    int64      `gorm:"autoCreateTime:milli"`
	K8sCluster   K8sCluster `gorm:"foreignKey:K8sClusterID;references:ID;constraint:OnDelete:NO ACTION,OnUpdate:NO ACTION"`
	Resource     Resource   `gorm:"foreignKey:ResourceID;references:ID;constraint:OnDelete:NO ACTION,OnUpdate:NO ACTION"`
}

func (kcr K8sClusterResource) InitTable() error {
	// 设置表信息
	//migrator := DBConn[DBName].Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='K8sClusterResources表'").Migrator()
	migrator := DBConn[DBName].Set("gorm:table_options", "").Migrator()

	// 判断表是否存在
	if migrator.HasTable(&K8sClusterResource{}) {
		// 存在就自动适配表，也就说原先没字段的就增加字段
		return migrator.AutoMigrate(&K8sClusterResource{})
	}

	// 不存在就创建新表
	return migrator.CreateTable(&K8sClusterResource{})
}

func (kcr K8sClusterResource) Insert() (K8sClusterResource, error) {
	tx := DBConn[DBName].Create(&kcr)
	return kcr, tx.Error
}

func (kcr K8sClusterResource) List() ([]K8sClusterResource, error) {
	var kcrs []K8sClusterResource
	result := DBConn[DBName].Order("id desc").Find(&kcrs)
	return kcrs, result.Error
}

func (kcr K8sClusterResource) Delete() error {
	tx := DBConn[DBName].Delete(&kcr)
	return tx.Error
}

func (kcr K8sClusterResource) Edit(kcrData K8sClusterResource) error {
	tx := DBConn[DBName].Model(&kcr)
	if tx.Error != nil {
		return tx.Error
	}

	return tx.Updates(kcrData).Error
}

func (kcr K8sClusterResource) ListResource(id uint, scopes []string) ([]Resource, error) {
	var rs []Resource

	tx := DBConn[DBName].Model(&kcr).
		Select("resource.*").
		Joins("left join resource on k8s_cluster_resource.resource_id = resource.id").
		Where("resource.category = ? and k8s_cluster_resource.k8s_cluster_id = ?", "vps", id)

	if len(scopes) > 0 {
		var scopesInterface []interface{}
		var whereParam []string
		for _, scope := range scopes {
			whereParam = append(whereParam, "resource.scope")
			scopesInterface = append(scopesInterface, scope)
		}
		tx.Where(fmt.Sprintf("%s = ?", strings.Join(whereParam, " = ? or ")), scopesInterface...)
	}

	tx.Order("resource.id desc").Scan(&rs)
	return rs, tx.Error
}
