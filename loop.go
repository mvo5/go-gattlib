package gattlib

/*
#cgo pkg-config: glib-2.0
#include "glib.h"
*/
import "C"

import (
	"runtime"
)

type GMainLoop struct {
	C *C.GMainLoop
}

func GMainLoopNew() (loop *GMainLoop) {
	CLoop := C.g_main_loop_new(nil, C.gboolean(0))
	loop = &GMainLoop{C: CLoop}
	runtime.SetFinalizer(loop, func(loop *GMainLoop) {
		C.g_main_loop_unref(loop.C)
	})
	return loop
}

func (l *GMainLoop) Run() {
	C.g_main_loop_run(l.C)
}

func (l *GMainLoop) Quit() {
	C.g_main_loop_quit(l.C)
}
