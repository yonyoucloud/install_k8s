package model

// 设置Resources表字段
type Resource struct {
	ID   uint   `gorm:"primaryKey;autoIncrement;comment:自增ID"`
	Name string `gorm:"type:varchar(127);not null;default:'';comment:必要描述名称，便于查看选择"`
	//Category  string `gorm:"type:enum('vps','mysql','redis','mongodb','rabbitmq','elasticsearch','kafka');default:'vps';comment:资源类别"`
	//Scope     string `gorm:"type:enum('default','master','sources','replicas','publish','node','etcd','etcdlb','masterlb','lvs','registry','pridns','newnode','newetcd','newmaster');default:'default';comment:特定资源的特定描述"`
	Category  string `gorm:"type:varchar(15);default:'vps';comment:资源类别"`
	Scope     string `gorm:"type:varchar(15);default:'default';comment:特定资源的特定描述"`
	Host      string `gorm:"type:varchar(255);not null;default:'';comment:主机地址，也可以是IP"`
	Port      uint32 `gorm:"not null;default:0;comment:端口号"`
	User      string `gorm:"type:varchar(127);not null;default:'';comment:用户名"`
	Password  string `gorm:"type:varchar(127);not null;default:'';comment:密码"`
	UpdatedAt int64  `gorm:"autoUpdateTime:milli"`
	CreatedAt int64  `gorm:"autoCreateTime:milli"`
}

func (r Resource) InitTable() error {
	// 设置表信息
	//migrator := DBConn[DBName].Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Resources表'").Migrator()
	migrator := DBConn[DBName].Set("gorm:table_options", "").Migrator()

	// 判断表是否存在
	if migrator.HasTable(&Resource{}) {
		// 存在就自动适配表，也就说原先没字段的就增加字段
		return migrator.AutoMigrate(&Resource{})
	}

	// 不存在就创建新表
	return migrator.CreateTable(&Resource{})
}

func (r Resource) Insert() (Resource, error) {
	tx := DBConn[DBName].Create(&r)
	return r, tx.Error
}

func (r Resource) List() ([]Resource, error) {
	var rs []Resource
	result := DBConn[DBName].Order("id desc").Find(&rs)
	return rs, result.Error
}

func (r Resource) Delete() error {
	tx := DBConn[DBName].Delete(&r)
	return tx.Error
}

func (r Resource) Edit(rData Resource) error {
	tx := DBConn[DBName].Model(&r)
	if tx.Error != nil {
		return tx.Error
	}

	return tx.Updates(rData).Error
}

func (r Resource) ListK8sCluster() ([]Resource, error) {
	var rs []Resource

	tx := DBConn[DBName].Model(&r).
		Select("resource.*").
		Joins("left join k8s_cluster_resource on resource.id = k8s_cluster_resource.resource_id").
		Where("resource.category = ? AND k8s_cluster_resource.id is null", "vps").
		Order("resource.id desc").
		Scan(&rs)
	return rs, tx.Error
}

func (r Resource) ListPod(podID uint) ([]Resource, error) {
	var rs []Resource

	if podID == 0 {
		result := DBConn[DBName].
			Where("category != ?", "vps").
			Order("id desc").
			Find(&rs)
		return rs, result.Error
	}

	tx := DBConn[DBName].Model(&r).
		Select("resource.*").
		Joins("left join pod_resource on resource.id = pod_resource.resource_id").
		Where("resource.category != ?", "vps").
		Where(DBConn[DBName].
			Where("pod_resource.resource_id not in (?)", DBConn[DBName].
				Table("pod_resource").
				Select("pod_resource.resource_id").
				Where("pod_resource.pod_id = ?", podID).
				Find(&r)).
			Or("pod_resource.pod_id is null")).
		Order("resource.id desc").
		Scan(&rs)
	return rs, tx.Error
}
