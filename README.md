# OUI

This is a fork of the fork of the fork the fork of [go-ouitools](https://github.com/dutchcoders/go-ouitools) package to work with MAC addresses and OUI. This package includes an OUI database.

## Example

```
package main

import (
	"fmt"
	"os"

	"github.com/karagenc/oui"
)

func main() {
	// You may also use NewDB and NewDBFromReader functions.
	db, err := oui.NewDBFromFile("oui.txt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
  
	mac := "00:16:e0:3d:f4:4c"
	vendor, err := db.Lookup(mac)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("%s => %s\n", mac, vendor)
}
```

## Testing

```
go test
```

## References

* Wireshark OUI database

## Contributors

[See here](https://github.com/dutchcoders/go-ouitools#contributors)
