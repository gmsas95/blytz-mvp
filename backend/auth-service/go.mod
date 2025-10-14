module github.com/blytz/auth-service

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/golang-jwt/jwt/v5 v5.2.0
	firebase.google.com/go/v4 v4.13.0
	github.com/go-redis/redis/v8 v8.11.5
	github.com/prometheus/client_golang v1.17.0
	gorm.io/gorm v1.25.5
	gorm.io/driver/postgres v1.5.4
	go.uber.org/zap v1.26.0
	golang.org/x/crypto v0.17.0
)

replace github.com/blytz/shared => ../shared