test:
	go test ./...

examples:
	go run example/full/full.go
	go run example/incr/incr.go
	go run example/multi/multi.go
	go run example/simple/simple.go

.PHONY: test examples
