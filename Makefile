all: format    \
     build     \
	 test      \
	 benchmark \
     coverage

format: 
	gofmt -w=true *.go

build: 
	go build 

test:
	go test

benchmark:
	go test -bench .

coverage:
	go test -cover .

