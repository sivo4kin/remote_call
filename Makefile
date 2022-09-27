PHONY: test

test: dep gen
	go test -v -count 1 ./test

gen:
	cd src;go generate

dep:
	go mod tidy
