

all: deps
	go build

deps:
	make -C extern/filecoin-ffi all


clean:
	make -C extern/filecoin-ffi clean
	go clean -cache -testcache .
