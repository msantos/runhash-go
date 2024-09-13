# SYNOPSIS

runhash *options* *subcommand* [*key* [*...*]]

# DESCRIPTION

runhash: command line interface for distributed node selection

`runhash` uses [Rendezvous or highest random weight
hashing](http://www.eecs.umich.edu/techreports/cse/96/CSE-TR-316-96.pdf)
to deterministically choose one or more nodes from a pool of nodes:

> Rendezvous or highest random weight (HRW) hashing is
> an algorithm that allows clients to achieve distributed
> agreement on a set of `k` options out of a possible set of
> `n` options. A typical application is when clients need
> to agree on which sites (or proxies) objects are assigned
> to. ([Wikipedia](https://en.wikipedia.org/wiki/Rendezvous_hashing))

`runhash` deterministically chooses:

* a node from a group of nodes to run a command

* a node or nodes to use for a service

## SUBCOMMANDS

### sort

Select one or more nodes from a list of nodes. Nodes are read:

* as command line arguments
* space delimited nodes from the RUNHASH_NODES environment variable
* from stdin if '-' is provided as the argument

#### Example

```
runhash sort mykey 127.0.0.1 192.168.1.1 10.0.0.1

RUNHASH_NODES="127.0.0.1 192.168.1.1 10.0.0.1" runhash sort mykey

echo -e "127.0.0.1\n192.168.1.1\n10.0.0.1" | runhash sort mykey -
```

### exec

Execute a command on a subset of nodes chosen from the list of nodes
based on the key. The number of nodes chosen is set by the `-n` option.

The command will always run on this node if either:
* RUNHASH_NODES is empty
* the node is included in `RUNHASH_NODES` and `-n` is 0 (all nodes)

#### Example

```
RUNHASH_NODES="$(uname -n) foo bar" runhash -n 1 exec mykey ls -al
```

### xargs

Failover a command between a sorted list of nodes. The command is executed
sequentially on the sorted list of nodes. If the command exits non-0,
the command is run again with the next node in the list.

'{}' in the command arguments is replaced with the selected node.

#### Example

```
RUNHASH_NODES="127.0.0.1 127.1.1.1" runhash xargs mykey nc "{}" 443

# set an environment variable for a command
RUNHASH_NODES="127.0.0.1 127.1.1.1" runhash xargs mykey \
  env TEST_VARIABLE="{}" env
```

# Build

```
go install codeberg.org/msantos/runhash-go/cmd/runhash@latest
```

To build a reproducible executable from the git repository:

```
CGO_ENABLED=0 go build -trimpath -ldflags "-w" ./cmd/runhash

# to include the version number
make
```

# OPTIONS

n *int*
: number of nodes to return. Set to `0` to return all nodes.

node *string*
: overrides `RUNHASH_NODE` environment variable

nodes *string*
: overrides `RUNHASH_NODES` environment variable

sorted *true|false*
: use the existing sort order for nodes, provided in the environment or
command line (default false)

## xargs

okexit *true|false*
: The "--okexit" option is the opposite of bash's "set -o errexit": `xargs`
terminates if the command exits with status 0.

If "--okexit=false", xargs will run the command on all nodes in
the list.

replace *string*
: template string (default "{}")

# ENVIRONMENT VARIABLES

RUNHASH_NODE="*hostname*"
: Tag used to identify this node in the node list. Usually this will be
set to the hostname or an IP address. Defaults to `uname -n`.

RUNHASH_NODES=""
: Pool of nodes to choose from.
