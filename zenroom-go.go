package zenroom

/*
#cgo CFLAGS: -I./zenroom
#cgo LDFLAGS: -L./zenroom -lzenroom
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "zenroom.h"

const char *zenroom(char *script, char *keys, char *data) {
  if (freopen("/dev/null", "a", stderr) == NULL)
    return NULL;

  char *outbuffer = (char *)malloc(sizeof(char) * 1024);
  if (outbuffer == NULL) {
    free(outbuffer);
    return NULL;
  }

  fflush(stdout);
  setvbuf(stdout, outbuffer, _IOLBF, 1024);

  if (zenroom_exec(script, NULL, keys, data, 1) != 0) {
    free(outbuffer);
    return NULL;
  }

  setbuf(stdout, NULL);

  fflush(stdout);
  //printf("len outbuffer=%zu\n", strlen(outbuffer));
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
