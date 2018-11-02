# Zenroom-Go

[![wercker status](https://app.wercker.com/status/87881dcb0b0ab25390300f91b96a9bf3/s/master "wercker status")](https://app.wercker.com/project/byKey/87881dcb0b0ab25390300f91b96a9bf3)
[![GoDoc](https://godoc.org/github.com/thingful/zenroom-go?status.svg)](https://godoc.org/github.com/thingful/zenroom-go)

Zenroom Binding for Go

## Introduction

Zenroom is a brand new virtual machine for fast cryptographic operations on Elliptic Curves. The Zenroom VM has no external dependencies, includes a cutting edge selection of C99 libraries and builds a small executable ready to run on: desktop, embedded, mobile, cloud and browsers (webassembly). This library adds a CGO wrapper for Zenroom, which aims to make the Zenroom VM easy to use from Go.

## Installation

Currently the bindings are only available for Linux machines, but if this is your current environment you should be able to just do:

```bash
$ go get github.com/thingful/zenroom-go
```

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/thingful/zenroom-go"
)

	genKeysScript := []byte(`
		keyring = ecdh.new('ec25519')
		keyring:keygen()
		
		output = JSON.encode({
			public = keyring:public():base64(),
			private = keyring:private():base64()
		})
		print(output)
	`)
	
	keys, err := zenroom.Exec(genKeysScript)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Keys: %s", keys)

 ```

## More Documentation

 * Zenroom documentation https://zenroom.dyne.org/
