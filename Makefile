test:
	go test ./client


teste2e:
	go test ./internal/e2e


testall: test teste2e
