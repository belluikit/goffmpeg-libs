//go:build linux && amd64 && !android

package linux_amd64

/*
#cgo LDFLAGS: -L${SRCDIR}/lib -lavdevice -lavformat -lavcodec -lswscale -lswresample -lavutil -lm -lssl -lcrypto -lpthread
*/
import "C"
