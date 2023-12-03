APP_NAME=./out/ez-ex
MOCK_SCRIPT_NAME=./out/create-mock-data
FEATURES="sqlite_foreign_keys"

build-cli:
	go build -o ${APP_NAME} -tags ${FEATURES} ./cmd/ez-ex-cli/
run-cli:
	clear
	@go build -o ${APP_NAME} -tags ${FEATURES} ./cmd/ez-ex-cli/
	@${APP_NAME}
build-mock:
	@go build -o ${MOCK_SCRIPT_NAME} ./internal/datamock/
	@echo 'Run `${MOCK_SCRIPT_NAME} -h` for more info'
t:
clean:
	go clean
	go mod tidy
	go fmt ./...
test:
	go test -v -cover ./...