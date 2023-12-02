NAME=./out/ez-ex
FEATURES="sqlite_foreign_keys"

build-cli:
	go build -o ${NAME} -tags ${FEATURES} ./cmd/ez-ex-cli/
run-cli:
	clear
	@go build -o ${NAME} -tags ${FEATURES} ./cmd/ez-ex-cli/
	@${NAME}
clean:
	go clean
	go mod tidy
	go fmt ./...
test:
	go test -v -cover ./...