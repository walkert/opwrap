# opwrap

The opwrap module can be used to return an [op](https://github.com/walkert/op).Op object while leveraging [stash](https://github.com/walkert/stash) if it has been configured.

## Installation

To install the module, simply run:

`$ go get github.com/walkert/opwrap`


## Usage

The primary function of `opwrap` is to return an op.Op object that has been configured with the `op.WithPassword` option using a password read using the `stash` client. If `stash` hasn't been configured or cannot be contacted, a standard `op.Op` object will be returned.

### Leveraging `stash`

If you have a `stash` server running, you can configure `opwrap` to use it by setting the following environment variables:

```shell
OPWRAP_CERT - The location of your stash cert file
OPWRAP_CERT - The location of your stash config file
OPWRAP_PORT - The port number your stash server is listening on
```

Below is a simple example:

```go
package main

import (
    "fmt"
    "log"

    "github.com/walkert/opwrap"
)

func main() {
    o, err := opwrap.GetOP()
    if err != nil {
        log.Fatal(err)
    }
    user, pass, err := o.GetUserPassword("vault item")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(user, pass)
}
```
