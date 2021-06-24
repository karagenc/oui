# OUI

This is a fork of the fork of the fork the fork of [go-ouitools](https://github.com/dutchcoders/go-ouitools) package to work with MAC addresses and OUI. This package includes an OUI database.

## Example

```go
package main

import (
	"fmt"
	"os"

	"github.com/tomruk/oui"
)

func main() {
	// You may also use NewDB and NewDBFromReader functions.
	// Or use the ouidata.NewDB function to load the embedded DB.
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

```bash
go test
```

## Update Embedded OUI DB

```bash
wget -O oui.txt https://gitlab.com/wireshark/wireshark/-/raw/master/manuf
```

## References

* Wireshark OUI database (aka Wireshark manufacturer database)
  * See https://www.wireshark.org/tools/oui-lookup.html
  * See https://gitlab.com/wireshark/wireshark/-/raw/master/manuf

## Contributors

[See here](https://github.com/dutchcoders/go-ouitools#contributors)
