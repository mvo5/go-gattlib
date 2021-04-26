#include "_cgo_export.h"

#include <stdbool.h>
#include <glib.h>
#include <bluetooth/bluetooth.h>

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

   connectionDone(res);
}
