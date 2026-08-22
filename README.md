# GridPulse Relay

电网继电保护跳闸 / 保护出口投递

Read `PROJECT.md` for the product scope, who uses it, and what the records mean. This is not a generic CMS.

## Run

```bash
set GOTOOLCHAIN=local
set CGO_ENABLED=0
go test ./... -count=1
go run ./cmd/server -addr 127.0.0.1:8080 -db ./data.sqlite
go run ./cmd/seed -db ./data.sqlite
```

Open http://127.0.0.1:8080/ for the console.

