testunit:
	go test ./client


teste2e:
	go test ./internal/e2e


test: testunit
test: teste2e


