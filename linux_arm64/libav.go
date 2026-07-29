//go:build linux && arm64 && !android

package linux_arm64

/*
#cgo LDFLAGS: -L${SRCDIR}/lib -lavdevice -lavformat -lavcodec -lswscale -lswresample -lavutil -lm -lssl -lcrypto -lpthread
*/
import "C"
