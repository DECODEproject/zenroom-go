package zenroom

/*
#cgo CFLAGS: -I./zenroom
#cgo LDFLAGS: -L./zenroom -lzenroom
#include <stdio.h>
#include <stdlib.h>
#include "zenroom.h"

extern int zenroom_exec(char *script, char *conf, char *keys,
                 char *data, int verbosity);

extern int zenroom_exec_tobuf(char *script, char *conf, char *keys,
                       char *data, int verbosity,
                       char *stdout_buf, size_t stdout_len,
                       char *stderr_buf, size_t stderr_len);

const char *zenroom(char *script, char *keys, char *data) {
  if (freopen("/dev/null", "a", stderr) == NULL)
    return NULL;

  char *outbuffer = (char *)calloc(sizeof(char) * MAX_STRING, 0);
  if (outbuffer == NULL) {
    free(outbuffer);
    return NULL;
  }

  fflush(stdout);
  setvbuf(stdout, outbuffer, _IOLBF, MAX_STRING);

  if (zenroom_exec(script, NULL, keys, data, 1) != 0) {
    free(outbuffer);
    return NULL;
  }

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

	//defer C.free(unsafe.Pointer(pt))

	if pt == nil {
		return "", fmt.Errorf("error calling zenroom lib")
	}
	res := strings.TrimSpace(C.GoString(pt))
	return res, nil
}

// ExecToBuf ...
func ExecToBuf(script, keys, data string) (string, error) {
	if len(script) == 0 {
		return "", fmt.Errorf("no lua script to process")
	}
	stdout := emptyString(1024)
	stderr := emptyString(1024)
	defer C.free(unsafe.Pointer(stdout))
	defer C.free(unsafe.Pointer(stderr))

	res := C.zenroom_exec_tobuf(C.CString(script), nil, C.CString(keys), C.CString(data), 1,
		(*C.char)(stdout), 1024,
		(*C.char)(stderr), 1024)

	if res != 0 {
		return "", fmt.Errorf("error calling zenroom: %s ", C.GoString(stderr))
	}

	return C.GoString(stdout), nil
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
