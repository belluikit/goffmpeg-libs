//go:build darwin && arm64 && !ios

package darwin_arm64

/*
#cgo LDFLAGS: -L${SRCDIR}/lib -lavdevice -lavformat -lavcodec -lswscale -lswresample -lavutil -lm
#cgo LDFLAGS: -framework AVFoundation -framework CoreMedia -framework CoreVideo -framework CoreFoundation -framework Foundation -framework CoreGraphics -framework VideoToolbox -framework Security
*/
import "C"
