//go:build windows && amd64

package windows_amd64

/*
#cgo LDFLAGS: -L${SRCDIR}/lib -lavdevice -lavformat -lavcodec -lswscale -lswresample -lavutil -lm -lsecur32 -lbcrypt -lws2_32 -lole32 -loleaut32 -lstrmiids -luuid -lgdi32 -lavicap32 -lvfw32 -luser32 -lshlwapi -lcrypt32 -lncrypt -lmfplat -lmfuuid
*/
import "C"
