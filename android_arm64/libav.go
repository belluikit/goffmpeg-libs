//go:build android && arm64

package android_arm64

/*
#cgo LDFLAGS: -L${SRCDIR}/lib -lavformat -lavcodec -lswscale -lswresample -lavutil -lm -llog -landroid -lmediandk
*/
import "C"
