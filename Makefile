compile:
	gox -output "build/{{.Dir}}_{{.OS}}_{{.Arch}}"

fmt:
	find ! -path "./vendor/*" -name "*.go" -exec gofmt -s -w {} \;

lint:
	golangci-lint run

run-tests:
	./test.sh

test-all: run-tests lint
