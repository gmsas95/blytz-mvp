module github.com/gmsas95/blytz-mvp/services/order-service

go 1.23.0

toolchain go1.24.9

replace github.com/gmsas95/blytz-mvp/shared => ../../shared

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/gmsas95/blytz-mvp/shared v0.0.0-00010101000000-000000000000
	github.com/prometheus/client_golang v1.17.0
	go.uber.org/zap v1.26.0
	gorm.io/driver/postgres v1.5.2
	gorm.io/gorm v1.25.5
)
