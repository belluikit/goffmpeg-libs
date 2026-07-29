//go:build windows && arm64

package windows_arm64

/*
#cgo LDFLAGS: -L${SRCDIR}/lib -lavdevice -lavformat -lavcodec -lswscale -lswresample -lavutil -lm -lsecur32 -lbcrypt -lws2_32 -lole32 -loleaut32 -lstrmiids -luuid -lgdi32 -lavicap32 -lvfw32 -luser32 -lshlwapi -lcrypt32 -lncrypt -lmfplat -lmfuuid
*/
import "C"
