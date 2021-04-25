# Copyright (C) 2014,2020 Oscar Acena <oscaracena@gmail.com>
# This software is under the terms of Apache License v2 or later.

PLATFORM := $(shell uname -m)
TARGETS   = libgattlib.so
OBJECTS   = att.o crypto.o uuid.o gatt.o gattrib.o btio.o log.o utils.o

CFLAGS  += -fPIC $$(pkg-config --cflags glib-2.0) -I./bluez 
CFLAGS  += -DVERSION='"5.53"'

LDFLAGS  = -lbluetooth  $$(pkg-config --libs glib-2.0)

vpath %.c bluez/attrib
vpath %.c bluez/src
vpath %.c bluez/src/shared
vpath %.c bluez/lib
vpath %.c bluez/btio

all: $(TARGETS)
	go build

libgattlib.so: $(OBJECTS)
	$(CC) $(CFLAGS) -shared -o $@ $^ $(LDFLAGS)

test:
	LD_LIBRARY_PATH=$$(pwd) go test 

.PHONY: clean
clean:
	rm -f *.o *.so* *~

