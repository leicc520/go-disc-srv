module github.com/leicc520/go-disc-srv

go 1.16

require (
	github.com/dchest/captcha v1.0.0
	github.com/gin-gonic/gin v1.7.3
	github.com/gobuffalo/packr/v2 v2.8.3
	github.com/leicc520/go-gin-http v1.0.0
	github.com/leicc520/go-orm v1.0.1
	github.com/mattn/go-sqlite3 v1.14.6
	github.com/rogpeppe/go-internal v1.8.1 // indirect
	github.com/stretchr/objx v0.1.1 // indirect
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa // indirect
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b // indirect
	golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab // indirect
)

replace (
	github.com/leicc520/go-gin-http v1.0.0 => ../go-gin-http
	github.com/leicc520/go-orm v1.0.1 => ../go-orm
)
