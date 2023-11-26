NAME=./out/ez-ex
FEATURES="sqlite_foreign_keys"

build-cli:
	go build -o ${NAME} -tags ${FEATURES} ./cmd/ez-ex-cli/ez-ex.go
run-cli:
	@go build -o ${NAME} -tags ${FEATURES} ./cmd/ez-ex-cli/ez-ex.go
	@${NAME}
clean:
	go clean
	go mod tidy
	go fmt ./...
test:
	go test -v -cover ./...