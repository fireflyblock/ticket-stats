module github.com/fireflyblock/ticket-stats

go 1.15

require (
	github.com/filecoin-project/go-address v0.0.5
	github.com/filecoin-project/go-state-types v0.1.0
	github.com/filecoin-project/lotus v1.8.0
	github.com/ipfs/go-log/v2 v2.1.3
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
)
replace github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi
