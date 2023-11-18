NAME=./out/ez-ex
FEATURES="sqlite_foreign_keys"

build:
	go build -o ${NAME} -tags ${FEATURES} ex-ez.go
run:
	@go build -o ${NAME} -tags ${FEATURES} ex-ez.go
	@${NAME}