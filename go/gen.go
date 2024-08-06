package gen

//go:generate go run ./cmd/tools/ternDotEnv/main.go
//go:generate sqlc generate -f ./internal/db/postgres/sqlc.yaml
