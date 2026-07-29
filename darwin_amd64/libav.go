//go:build darwin && amd64 && !ios

package darwin_amd64

/*
#cgo LDFLAGS: -L${SRCDIR}/lib -lavdevice -lavformat -lavcodec -lswscale -lswresample -lavutil -lm
#cgo LDFLAGS: -framework AVFoundation -framework CoreMedia -framework CoreVideo -framework CoreFoundation -framework Foundation -framework CoreGraphics -framework VideoToolbox -framework Security
*/
import "C"
