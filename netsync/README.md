netsync
=======

[![Build Status](http://img.shields.io/travis/ifishnet/hdfd.svg)](https://travis-ci.org/ifishnet/hdfd)
[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/ifishnet/hdfd/netsync)

## Overview

This package implements a concurrency safe block syncing protocol. The
SyncManager communicates with connected peers to perform an initial block
download, keep the chain and unconfirmed transaction pool in sync, and announce
new blocks connected to the chain. Currently the sync manager selects a single
sync peer that it downloads all blocks from until it is up to date with the
longest chain the sync peer is aware of.

## Installation and Updating

```bash
$ go get -u github.com/ifishnet/hdfd/netsync
```

## License

Package netsync is licensed under the [copyfree](http://copyfree.org) ISC License.
