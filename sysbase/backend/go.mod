module sysbase

go 1.15

replace sysbase => ../../sysbase

require (
	github.com/fvbock/endless v0.0.0-20170109170031-447134032cb6
	github.com/gin-gonic/gin v1.9.0
	github.com/hnakamur/go-scp v1.0.2
	github.com/pkg/sftp v1.13.5
	github.com/spf13/pflag v1.0.5
	golang.org/x/crypto v0.8.0
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/sqlite v1.5.7
	gorm.io/gorm v1.25.7-0.20240204074919-46816ad31dde
)
