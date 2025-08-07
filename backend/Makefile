coverage-html:
	@go test ./... -coverprofile=coverage.out -coverpkg=./... -timeout 30s
	@go tool cover -html=./coverage.out -o ./coverage.html
	@rm ./coverage.out