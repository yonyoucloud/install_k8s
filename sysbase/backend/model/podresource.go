package model

// 设置PodResource表字段
type PodResource struct {
	ID         uint     `gorm:"primaryKey;autoIncrement;comment:自增ID"`
	PodID      uint     `gorm:"not null;default:0;index:idx_pod_resource,unique;comment:PodID"`
	ResourceID uint     `gorm:"not null;default:0;index:idx_pod_resource,unique;comment:ResourceID，该资源可以被多个Pod使用，比如mysql资源是可以被不同的Pod共用的，PodID和ResourceID唯一"`
	UpdatedAt  int64    `gorm:"autoUpdateTime:milli"`
	CreatedAt  int64    `gorm:"autoCreateTime:milli"`
	Pod        Pod      `gorm:"foreignKey:PodID;references:ID;constraint:OnDelete:NO ACTION,OnUpdate:NO ACTION"`
	Resource   Resource `gorm:"foreignKey:ResourceID;references:ID;constraint:OnDelete:NO ACTION,OnUpdate:NO ACTION"`
}

func (pr PodResource) InitTable() error {
	// 设置表信息
	//migrator := DBConn[DBName].Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='PodResources表'").Migrator()
	migrator := DBConn[DBName].Set("gorm:table_options", "").Migrator()

	// 判断表是否存在
	if migrator.HasTable(&PodResource{}) {
		// 存在就自动适配表，也就说原先没字段的就增加字段
		return migrator.AutoMigrate(&PodResource{})
	}

	// 不存在就创建新表
	return migrator.CreateTable(&PodResource{})
}

func (pr PodResource) Insert() (PodResource, error) {
	tx := DBConn[DBName].Create(&pr)
	return pr, tx.Error
}

func (pr PodResource) List() ([]PodResource, error) {
	var prs []PodResource
	result := DBConn[DBName].Order("id desc").Find(&prs)
	return prs, result.Error
}

func (pr PodResource) Delete() error {
	tx := DBConn[DBName].Delete(&pr)
	return tx.Error
}

func (pr PodResource) Edit(prData PodResource) error {
	tx := DBConn[DBName].Model(&pr)
	if tx.Error != nil {
		return tx.Error
	}

	return tx.Updates(prData).Error
}

func (pr PodResource) ListResource(id uint) ([]Resource, error) {
	var rs []Resource

	tx := DBConn[DBName].Model(&pr).
		Select("resource.*").
		Joins("left join resource on pod_resource.resource_id = resource.id").
		Where("resource.category != ? and pod_resource.pod_id = ?", "vps", id).
		Order("resource.id desc").
		Scan(&rs)
	return rs, tx.Error
}
