//go:build android && arm

package android_arm

/*
#cgo LDFLAGS: -L${SRCDIR}/lib -lavformat -lavcodec -lswscale -lswresample -lavutil -lm -llog -landroid -lmediandk
*/
import "C"
