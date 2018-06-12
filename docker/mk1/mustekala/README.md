# Mustekala Box

## Overview

Designed to operate on a coreOS box with default configs.
This is a full functioning set of containers running a reverse proxy in `443`,
with forwarding from `80` to the `ws` port in the IPFS bridge, which has
a read only parity-ipfs datastore.

The generation of SSL certificates is made with `Let's encrypt` automatically.

## Vision

* People wanting to increase the density of the network by running one of these boxes in their laptops or servers.
* Different flavors (such as)
  * Running the box with a fast synced parity box.
  * Running the box "only the ipfs node, thank you".
  * Just run the diagnostics.

## Requirements

In a CoreOS machine, install `docker-compose`

```
mkdir -p /opt/bin
curl -L "https://github.com/docker/compose/releases/download/1.14.0/docker-compose-$(uname -s)-$(uname -m)" -o /opt/bin/docker-compose
chmod +x /opt/bin/docker-compose
```

## Execution

0. Use coreOS as the account `core`.
1. Download this repository into `/home/core`.
2. Add your peers (ethereum and ipfs) into a directory `/home/core/peers`. [Like this](https://github.com/MetaMask/eth-ipfs-browser-client/blob/master/peers).
2. `cd /home/core/ipfs-eth-bridge/mustekala`
3. `./run_box`

## Customization of domains

Just change the environmental variables in `mustekala/docker-compose.yml`
to your own domain, subdomains and admin email.

## Adding peers

The API port of the bridge (`5001`) is opened in docker. (Make sure is firewalled from the internet!)

Just do

```
curl -s http://localhost:5001/api/v0/swarm/connect?arg=<my peer multiaddress>
```

For example

```
curl -s http://localhost:5001/api/v0/swarm/connect?arg=/ip4/52.232.121.251/tcp/4001/ipfs/QmXFdPj3FuVpkgmNHNTFitkp4DSmVuF6HxNX6tCZr4LFz9
```

Will add you `tiger.musteka.la` as a peer.

### Bulk adding using `curl` and the `API`

Using this ~ugly hack~ script, we can sequentially add nodes from a file containing a list of peers

```
./add_ipfs_peer_list <path to the file>
```
