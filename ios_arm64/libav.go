//go:build ios && arm64 && !iossimulator

package ios_arm64

/*
#cgo LDFLAGS: -L${SRCDIR}/lib -lavformat -lavcodec -lswscale -lswresample -lavutil -lm
#cgo LDFLAGS: -framework VideoToolbox -framework CoreMedia -framework CoreVideo -framework CoreFoundation -framework Foundation -framework Security
*/
import "C"
