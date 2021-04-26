package gattlib

import (
	"fmt"
	"unsafe"
)

// go:generate make

/*
#cgo pkg-config: glib-2.0 bluez
#cgo CFLAGS: -I./bluez
#cgo LDFLAGS: -L./ -lgattlib

#include <stdbool.h>
#include <glib.h>
#include <bluetooth/bluetooth.h>

#include "lib/uuid.h"
#include "btio/btio.h"
#include "attrib/att.h"
#include "attrib/utils.h"
#include "attrib/gattrib.h"

struct connection_result {
    bool done;
    GError *err;
    GAttrib *attrib;

    // go-data
    void *goCtx;
};
typedef struct connection_result connection_result_t;

void my_connect_cb(GIOChannel *io, GError *err, gpointer user_data);
*/
import "C"

type Connection struct {
	attrib *C.GAttrib
}

type goCtx struct {
	doneCh chan bool
}

//export connectionDone
func connectionDone(res *C.connection_result_t) {
	ctx := (*goCtx)(unsafe.Pointer(res.goCtx))
	ctx.doneCh <- true
}

func Connect(dest string) (*Connection, error) {
	optDest := C.CString(dest)
	// XXX: make configurable
	optSrc := C.CString("hci0")
	optDstType := C.CString("public")
	optSecLevel := C.CString("low")
	optPsm := C.int(0)
	optMtu := C.int(0)
	var gerr *C.GError

	ctx := goCtx{
		doneCh: make(chan bool),
	}
	res := C.connection_result_t{goCtx: unsafe.Pointer(&ctx)}
	chann := C.gatt_connect(optSrc, optDest, optDstType, optSecLevel, optPsm, optMtu, (C.BtIOConnect)(C.my_connect_cb), &gerr, (C.gpointer)(unsafe.Pointer(&res)))
	if chann == nil {
		return nil, fmt.Errorf("cannot try to connect: %v", gerr.message)
	}
	<-ctx.doneCh
	if res.err != nil {
		return nil, fmt.Errorf("cannot connect: %v", res.err.message)
	}

	return &Connection{
		attrib: res.attrib,
	}, nil
}
