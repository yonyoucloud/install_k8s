module git.yonyou.com/sysbase/backend

go 1.15

replace git.yonyou.com/sysbase => ../../sysbase

require (
	github.com/fvbock/endless v0.0.0-20170109170031-447134032cb6
	github.com/gin-gonic/gin v1.6.3
	github.com/hnakamur/go-scp v1.0.1
	github.com/pkg/sftp v1.12.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/mysql v1.0.4
	gorm.io/gorm v1.20.12
	gorm.io/plugin/dbresolver v1.1.0
)
