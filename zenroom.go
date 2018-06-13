package zenroom

/*
#cgo CFLAGS: -IC:${SRCDIR}
#cgo LDFLAGS: -L${SRCDIR}/lib -Wl,-rpath=${SRCDIR}/lib -lzenroom
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "zenroom.h"

extern int zenroom_exec(char *script, char *conf, char *keys,
                 char *data, int verbosity);

extern int zenroom_exec_tobuf(char *script, char *conf, char *keys,
                       char *data, int verbosity,
                       char *stdout_buf, size_t stdout_len,
                       char *stderr_buf, size_t stderr_len);
*/
import (
	"C"
)

import (
	"fmt"
	"unsafe"

	_ "github.com/thingful/zenroom-go/lib"
)

// maxString is zenroom defined buffer MAX_STRING size
const maxString = 4096

// Exec calls zenroom_exec_tobuf function with the next params.
// script: Lua script to execute.
// keys: Optional field mapped to KEYS zenroom global var.
// data: Optional field mapped to DATA zenroom global var.
// Returns: a string with zenroom output and error which can be a zenroom stderr
func Exec(script, keys, data []byte) ([]byte, error) {

	var cScript, optKeys, optData *C.char

	if script == nil {
		return nil, fmt.Errorf("no lua script to process")
	}
	cScript = (*C.char)(unsafe.Pointer(&script[0]))

	if keys != nil {
		optKeys = (*C.char)(unsafe.Pointer(&keys[0]))
	}
	if data != nil {
		optData = (*C.char)(unsafe.Pointer(&data[0]))
	}

	stdout := emptyString(maxString)
	stderr := emptyString(maxString)
	defer C.free(unsafe.Pointer(stdout))
	defer C.free(unsafe.Pointer(stderr))

	res := C.zenroom_exec_tobuf(cScript, nil, optKeys, optData, 1,
		(*C.char)(stdout), maxString,
		(*C.char)(stderr), maxString)

	if res != 0 {
		return nil, fmt.Errorf("error calling zenroom: %s ", C.GoString(stderr))
	}

	return C.GoBytes(unsafe.Pointer(stdout), C.int(C.strlen(stdout))), nil
}

// reimplementation of https://golang.org/src/strings/strings.go?s=13172:13211#L522
func emptyString(size int) *C.char {
	p := C.malloc(C.size_t(size + 1))
	// largest array size that can be used on all architectures
	pp := (*[1 << 30]byte)(p)
	bp := copy(pp[:], " ")
	for bp < size {
		copy(pp[bp:], pp[:bp])
		bp *= 2
	}
	pp[size] = 0
	return (*C.char)(p)
}
