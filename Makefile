compile:
	gox -osarch "linux/amd64 darwin/amd64 windows/amd64" -output "build/{{.Dir}}_{{.OS}}_{{.Arch}}"

fmt:
	find ! -path "./vendor/*" -name "*.go" -exec gofmt -s -w {} \;

lint:
	golangci-lint run

run-tests:
	./test.sh

test-all: run-tests lint
