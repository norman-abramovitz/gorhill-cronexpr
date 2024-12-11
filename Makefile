
all: cronexpr cronexpr.exe
	
.PHONY: test clean

cronexpr: cronexpr.go cronexpr_parse.go cronexpr_next.go cmd/cronexpr/main.go
	go build -o cronexpr ./cmd/cronexpr

cronexpr.exe: cronexpr.go cronexpr_parse.go cronexpr_next.go cmd/cronexpr/main.go
	GOOS=windows GOARCH=amd64 go build -o cronexpr.exe ./cmd/cronexpr

test:
	go test -v ./....

clean:
	rm -f cronexpr cronexpr.exe
