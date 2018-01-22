## The DevP2P / LibP2P Hot Importer

### Why "_hot importer_"?

Because we are just plugging into the devp2p network, as another peer, trying
to get the data as hot as it comes.

### Very Quick Start

```
./run --devp2p-bootnodes ./config/bootnodes-devp2p
```

**Moar** Tracing with `bridge-debug`

```
./run --devp2p-bootnodes ./config/bootnodes-devp2p --bridge-debug
```

**EVEN MOAR** Tracing. Add `--devp2p-lib-debug`. Will show you what the `go-ethereum/p2p` is doing.

```
./run --devp2p-bootnodes ./config/bootnodes-devp2p --bridge-debug --devp2p-lib-debug
```

### DevP2P Network Status File

If you want to have a `.csv` file with all the nodes you have dialed, and their status,
just use the flag `--devp2p-network-status` with some chosen filename.

```
./run --devp2p-bootnodes ./config/bootnodes-devp2p --devp2p-network-status /tmp/devp2p-stats
```

### Documentation

At this stage, the hot importer is a very experimental package.

Documentation is in code, starting from `main.go`. If you want to contribute
to this application, don't be shy to be verbose: Reading another person's code
is where our billed hours go the most.
