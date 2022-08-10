module github.com/leicc520/go-disc-srv

go 1.16

require (
	github.com/dchest/captcha v1.0.0
	github.com/gin-gonic/gin v1.7.3
	github.com/leicc520/go-gin-http v1.0.0
	github.com/leicc520/go-orm v1.0.1
	github.com/mattn/go-sqlite3 v1.14.6
)

replace (
	github.com/leicc520/go-gin-http v1.0.0 => ../go-gin-http
	github.com/leicc520/go-orm v1.0.1 => ../go-orm
)
