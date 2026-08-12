module github.com/fastsprint/project-service

go 1.21

require (
	github.com/fastsprint/common v0.0.0
	github.com/gin-gonic/gin v1.9.1
	github.com/spf13/viper v1.17.0
	gorm.io/gorm v1.25.5
)

replace github.com/fastsprint/common => ../common
