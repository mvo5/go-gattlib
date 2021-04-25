package gattlib

import (
	"fmt"
	"time"
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
};
typedef struct connection_result connection_result_t;

void my_connect_cb(GIOChannel *io, GError *err, gpointer user_data) {
    printf("my_connect_cb %x\n", err);
    struct connection_result *res = (struct connection_result*)(user_data);
    res->done = true;
    res->err = err;
    if (err) {
             return;
    }
    // XXX: detect mtu
    uint16_t mtu = ATT_DEFAULT_LE_MTU;
    res->attrib = g_attrib_new(io, mtu);
    // XXX: support notify, cf gatttool.c connect_cb()
}
*/
import "C"

type Connection struct {
	attrib *C.GAttrib
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

	res := C.connection_result_t{}
	chann := C.gatt_connect(optSrc, optDest, optDstType, optSecLevel, optPsm, optMtu, (C.BtIOConnect)(C.my_connect_cb), &gerr, (C.gpointer)(unsafe.Pointer(&res)))
	if chann == nil {
		return nil,fmt.Errorf("cannot try to connect: %v", gerr.message)
	}
	// XXX: horrible
	for {
		if res.done {
			break
		}
		time.Sleep(1*time.Second)
	}
	if res.err != nil {
		return nil,fmt.Errorf("cannot connect: %v", res.err.message)
	}
	
	return &Connection{
		attrib: res.attrib,
	}, nil
}
