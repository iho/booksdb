check:
	gofumpt -l -w .
	gci -w .
	golangci-lint run

s:
	go test ./... -short

test:
	go clean -testcache
	go test ./... -coverprofile cover.out
	go tool cover -html=cover.out
