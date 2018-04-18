# Zenroom-Go

[![wercker status](https://app.wercker.com/status/87881dcb0b0ab25390300f91b96a9bf3/s/master "wercker status")](https://app.wercker.com/project/byKey/87881dcb0b0ab25390300f91b96a9bf3)

Zenroom Binding for go

## How to use

* Install zenroom bindings
``` go get github.com/thingful/zenroom-go```

* Ensure to have installed the contents of zenroom folder in /usr/local/lib or /usr/lib

* Run ` sudo ldconfig` 

* Have Fun!
```
package main

import (
	"fmt"
	"log"

	"github.com/thingful/zenroom-go"
)

func main() {
	script := `print("holass")`
    keys := ""
    data := ""
	res, err := zenroom.Exec(script, keys, data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}
 ```

 * Zenroom documentation https://zenroom.dyne.org/