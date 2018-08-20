## Optimizing geth for the slicer (for fun and profit)

### Overview

We need to understand all we can about the way go-ethereum (geth from now)
optimizes its memory to manage its data. We should leverage from that,
and build our slice memory management overlay, so we can avoid as much
as possible interactions with the disk (which is way slower) at the time of
retrieving a slice (and the subsequent _deltas_ we want to get).

### Current State of the Art

(WIP)

### Cheatsheet

#### pprof command

	go tool pprof http://localhost:6060/debug/pprof/heap

#### geth command line option to add pprof

	--pprof --pprofaddr=0.0.0.0

#### memsize

	curl -L -X POST http://localhost:6060/memsize/scan?root=node

#### bonus

	curl http://localhost:6060/debug/pprof/goroutine?debug=1