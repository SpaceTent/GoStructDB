fmt:
	@go mod tidy
	@goimports -w .
	@gofmt -w -s .
	@go clean ./...

run: fmt
	go run main.go


test:
	go test -v -coverprofile=profile.cov ./...

commit: fmt
	@git add .
	@git commit -a -m "$(m)"
	@git pull
	@git push

