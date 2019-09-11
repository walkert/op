# op

A module that wraps calls to the 1Password `op` command to retrieve credentials from an active session.

## Installation

To install the module, simply run the below or just import it into your code if using modules:

`$ go get github.com/walkert/op`


## Details

The op module expects that you have `op` installed and have signed-in to an account before (thus generating ~/.op/config). If you already have an OP\_SESSION environment variable set, the module will attempt to use that to communicate with `op`. If it's unset, it will call `op` and prompt you to sign-in (this will be taken care of by `op` - not the module) unless you pass an explicit password to `New`.

If an account is not specified, the module will attempt to use the only account in the config file. If there is more than one, you'll need to specify its name explicitly.

## Getting Started

### Get the username and password from a specific 1password account

```go
package main

import (
    "fmt"
    "log"

    "github.com/walkert/op"
)

func main() {
    o, err := op.New(op.WithAccount("myaccount"))
    if err != nil {
        log.Fatal(err)
    }
    user, pass, err := o.GetUserPass("item name")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Got user and pass:", user, pass)
}
```

### Get the totp for an item from the default account

```go
package main

import (
    "fmt"
    "log"

    "github.com/walkert/op"
)

func main() {
    totp, err := op.GetTotp("my item")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Got totp:", totp)
}
```
