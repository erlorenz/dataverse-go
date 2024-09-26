

test:
	go test ./client/...

teste2e:
	go test -v -race ./internal/e2e/...


testall: test teste2e
