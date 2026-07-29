//go:build android && amd64

package android_amd64

/*
#cgo LDFLAGS: -L${SRCDIR}/lib -lavformat -lavcodec -lswscale -lswresample -lavutil -lm -llog -landroid -lmediandk
*/
import "C"
