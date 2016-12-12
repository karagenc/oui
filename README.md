go-ouitools
===========

Golang tools to work with Mac addresses and oui. Includes oui database to resolve to vendor. 

## Sample
```
package main

import (
	"fmt"
	"os"

	"github.com/davidbb/ouidb"
)

var (
    db *ouidb.OuiDB
)

func main() {
	db = ouidb.New("oui.txt")
	if db == nil {
		fmt.Println("database not initialized")
		os.Exit(1)
	}
  
	mac := "00:16:e0:3d:f4:4c"
	v, err := db.Lookup(mac)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	
	fmt.Printf("%s => %s\n", mac, v)
}

```

## Testing
```
go test
```

## References
* wireshark oui database

## Contributors
* Remco Verhoef (Dutchcoders) @remco_verhoef
* Claudio Matsuoka
* David Barrera

