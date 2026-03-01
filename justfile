set shell := ["pwsh.exe", "-c"]

default:
    just --list

dev:
    go run main.go
    
build: 
    $env:GOOS="linux";$env:GOARCH="amd64";go build -o output/labelplus-next-server
    
mgr-add script-name:
    sqlx migrate add -r {{script-name}}

mgr-rvt mode="step":
    {{if mode == "all" {
        "sqlx migrate revert --target-version 0"
    } else {
        "sqlx migrate revert"
    }}}
