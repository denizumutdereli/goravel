test:
	@go test -v ./...

cover:
	@go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

coverage:
	@go test -cover ./...

build_cli:
	@go build -o ../goravel/gortisan.exe ./cmd/cli