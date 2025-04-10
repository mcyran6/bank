docker run --rm -v $(pwd):/src -w /src sqlc/sqlc generate

go mod tidy

go get 

go test