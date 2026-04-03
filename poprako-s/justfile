set shell := ["pwsh.exe", "-c"]

default:
    just --list
    
swag:
    swag init

check paths="./...":
    golangci-lint-v2 run {{ paths }}

fmt:
    golangci-lint-v2 fmt ./...

cloc:
    cloc internal/

dev: fmt swag
    go run main.go
    
build:
    $env:GOOS="linux";$env:GOARCH="amd64";go build -o output/labelplus-next-server
    
mgr-add script-name:
    sqlx migrate add -r {{script-name}}
    
mgr-run:
    sqlx migrate run

mgr-rvt mode="step":
    {{if mode == "all" {
        "sqlx migrate revert --target-version 0"
    } else {
        "sqlx migrate revert"
    }}}
    
seed-base:
    bun run script/seed/seedBase.ts
