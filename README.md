![TackDB](docs/media/tackdb.svg?raw=true "TackDB logo")
=======================================================

[summary]::
TackDB is an in-memory key-value store.

[![Build Status](https://travis-ci.org/mtso/tackdb.svg?branch=master)](https://travis-ci.org/mtso/tackdb)

- [Building](#building)
- [Configuration](#configuration)
- [Commands](#commands)

## Building
```
$ go get github.com/tackdb/tackdb
$ go install github.com/tackdb/tackdb/cmd/tackdb
```

## Configuration

There are two ways to configure the settings for a TackDB instance.
The first way is to pass the settings as options from the command line.
The second way is to edit the `tackdb.conf` file. Note that the
`tackdb.conf` file overrides command line options.

### Example of `tackdb.conf`

The only used attribute in the `tackdb.conf` file at the moment is the port string.

```json
{
  "port": "3750"
}
```

### Command Line Options

```
Usage of tackdb:
  --confname string
        Filename of TackDB runtime configuration file. (default "tackdb.conf")
  -d, --dir string
        Directory location of runtime configuration file (.tackrc). (default "/Users/matthewtso")
  -p, --port string
        TCP service port. (default "3750")
```

For example:
```
$ tackdb --port=3750
```

## Commands

`SET [key] [value]`  
Saves a reference to `value` by the given `key`.

`GET [key]`  
Retrieves the `value` for the given `key`.

`UNSET [key]`  
Removes the `value` referenced by `key`.

`NUMEQUALTO [value]`  
Retrieves the count of saved values match the given `value`.

`BEGIN`  
Begins a transaction block. while the current connection has an open transaction,
no other connections can execute commands.

`ROLLBACK`  
Returns the datastore to the state of the last `BEGIN` command.

`COMMIT`  
Saves all outstanding transactions.
