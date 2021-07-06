

all: deps
	go build

deps:
    git submodule update --init --recursive
	make -C extern/filecoin-ffi all


clean:
	make -C extern/filecoin-ffi clean
	go clean -cache -testcache .
