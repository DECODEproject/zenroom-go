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

	genKeysScript := []byte(`
		octet = require 'octet'
		ecdh = require 'ecdh'
		json = require 'json'
		keyring = ecdh.new('ec25519')
		keyring:keygen()
		
		output = json.encode({
			public = keyring:public():base64(),
			private = keyring:private():base64()
		})
		print(output)
	`)
	keys, err := zenroom.Exec(genKeysScript, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Keys: %s", keys)

 ```

 * Zenroom documentation https://zenroom.dyne.org/