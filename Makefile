testunit:
	go test . && go test ./internal/auth


teste2e:
	go test ./internal/e2e


test: testunit
test: teste2e


