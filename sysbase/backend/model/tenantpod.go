package model

// 设置TenantPod表字段
type TenantPod struct {
	ID         uint   `gorm:"tpimaryKey;autoIncrement;comment:自增ID"`
	TenantID   uint   `gorm:"not null;default:0;index:idx_tenant_pod,unique;comment:租户ID"`
	PodID      uint   `gorm:"not null;default:0;index:idx_tenant_pod,unique;comment:PodID，和TenantID一起作为唯一索引"`
	TenantName string `gorm:"type:varchar(127);not null;default:'';comment:租户名称"`
	UpdatedAt  int64  `gorm:"autoUpdateTime:milli"`
	CreatedAt  int64  `gorm:"autoCreateTime:milli"`
	Pod        Pod    `gorm:"foreignKey:PodID;references:ID;constraint:OnDelete:NO ACTION,OnUpdate:NO ACTION"`
}

func (tp TenantPod) InitTable() error {
	// 设置表信息
	//migrator := DBConn[DBName].Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='TenantPods表'").Migrator()
	migrator := DBConn[DBName].Set("gorm:table_options", "").Migrator()

	// 判断表是否存在
	if migrator.HasTable(&TenantPod{}) {
		// 存在就自动适配表，也就说原先没字段的就增加字段
		return migrator.AutoMigrate(&TenantPod{})
	}

	// 不存在就创建新表
	return migrator.CreateTable(&TenantPod{})
}

func (tp TenantPod) Insert() (TenantPod, error) {
	tx := DBConn[DBName].Create(&tp)
	return tp, tx.Error
}

func (tp TenantPod) List() ([]TenantPod, error) {
	var tps []TenantPod
	result := DBConn[DBName].Order("id desc").Find(&tps)
	return tps, result.Error
}

func (tp TenantPod) Delete() error {
	tx := DBConn[DBName].Delete(&tp)
	return tx.Error
}

func (tp TenantPod) Edit(tpData TenantPod) error {
	tx := DBConn[DBName].Model(&tp)
	if tx.Error != nil {
		return tx.Error
	}

	return tx.Updates(tpData).Error
}

func (tp TenantPod) Get() (TenantPod, error) {
	var tpData TenantPod
	err := DBConn[DBName].Where(&tp).Find(&tpData).Error
	return tpData, err
}
