.PHONY: init test mock-repository

init: 
	go mod tidy
	go mod vendor

test:
	go test -short -coverprofile coverage.out -v ./...

mock-repository:
	$(shell go env GOPATH)/bin/mockgen -source src/repository/wallet_repository.go -destination src/mock/repository/wallet_repository.go