package main

/*
#cgo LDFLAGS: -lace_napi.z
#include <stdint.h>
#include <node_api.h>

static void finish(napi_threadsafe_function fn, uint32_t id, uint32_t result) {
    uintptr_t packed = ((uintptr_t)id << 32) | result;
    if (napi_call_threadsafe_function(fn, (void *)packed, napi_tsfn_nonblocking) != napi_closing) {
        napi_release_threadsafe_function(fn, napi_tsfn_release);
    }
}
*/
import "C"

import "time"

//export StartProbe
func StartProbe(fn C.napi_threadsafe_function, id, delay C.uint32_t) {
	go func() {
		time.Sleep(time.Duration(delay) * time.Millisecond)
		C.finish(fn, id, id*id)
	}()
}

func main() {}
