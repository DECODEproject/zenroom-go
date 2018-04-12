package zenroom

/*
#cgo CFLAGS: -I./zenroom
#cgo LDFLAGS: -L./zenroom -lzenroom
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "zenroom.h"

char* zenroom(char* script, char* keys, char* data){
	if (freopen("/dev/null", "a", stderr) == NULL)
		return NULL;

	char *sbuffer = (char *) malloc(sizeof(char) * 1024);
	char *outbuffer = (char *) malloc(sizeof(char) * 1024);

	if (sbuffer == NULL || outbuffer == NULL)
		return NULL;

	fflush(stdout);
	setvbuf(stdout, outbuffer, _IOFBF, 1024);

	if (zenroom_exec(script, NULL, keys, data, 1) != 0)
		return NULL;

	strncpy(sbuffer,outbuffer,1024);
	setbuf(stdout, NULL);

	fflush(stdout);
	return outbuffer;
}
*/
import (
	"C"
)

import (
	"fmt"
	"strings"
	"unsafe"
)

// Exec ...
func Exec(script, keys, data string) (string, error) {
	if len(script) == 0 {
		return "", fmt.Errorf("no lua script to process")
	}
	pt := C.zenroom(C.CString(script), C.CString(keys), C.CString(data))

	defer C.free(unsafe.Pointer(pt))

	if pt == nil {
		return "", fmt.Errorf("error calling zenroom lib")
	}
	res := strings.TrimSpace(C.GoString(pt))
	return res, nil
}
