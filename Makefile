all: cronexpr

cronexpr: cronexpr.go cronexpr_parse.go cronexpr_next.go cmd/cronexpr/main.go
	go build -o cronexpr ./cmd/cronexpr

windows: cronexpr.go cronexpr_parse.go cronexpr_next.go cmd/cronexpr/main.go
	GOOS=windows GOARCH=amd64 go build -o cronexpr.exe ./cmd/cronexpr
