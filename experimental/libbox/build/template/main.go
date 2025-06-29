package main

import "C"
import (
	"github.com/sagernet/sing-box/experimental/libbox"
	"github.com/sagernet/sing-box/experimental/libbox/share"
)

func main() {}

//export CGoGetFreePorts
func CGoGetFreePorts(count int) *C.char {
	return C.CString(libbox.GetFreePorts(count))
}

//export CGoConvertShareLinksToXrayJson
func CGoConvertShareLinksToXrayJson(base64Text *C.char) *C.char {
	text := C.GoString(base64Text)
	return C.CString(libbox.ConvertShareLinksToXrayJson(text))
}

//export CGOConvertXrayJsonToShareLinks
func CGOConvertXrayJsonToShareLinks(base64Text *C.char) *C.char {
	text := C.GoString(base64Text)
	return C.CString(share.ConvertXrayJsonToShareLinks(text))
}

//export CGoLoadGeoData
func CGoLoadGeoData(base64Text *C.char) *C.char {
	text := C.GoString(base64Text)
	return C.CString(libbox.LoadGeoData(text))
}

//export CGoPing
func CGoPing(base64Text *C.char) *C.char {
	text := C.GoString(base64Text)
	return C.CString(libbox.Ping(text))
}

//export CGoQueryStats
func CGoQueryStats(base64Text *C.char) *C.char {
	text := C.GoString(base64Text)
	return C.CString(libbox.QueryStats(text))
}

//export CGoCustomUUID
func CGoCustomUUID(base64Text *C.char) *C.char {
	text := C.GoString(base64Text)
	return C.CString(libbox.CustomUUID(text))
}

//export CGoTestXray
func CGoTestXray(base64Text *C.char) *C.char {
	text := C.GoString(base64Text)
	return C.CString(libbox.TestXray(text))
}

//export CGoRunXray
func CGoRunXray(base64Text *C.char) *C.char {
	text := C.GoString(base64Text)
	return C.CString(libbox.RunXray(text))
}

//export CGoStopXray
func CGoStopXray() *C.char {
	return C.CString(libbox.StopXray())
}

//export CGoXrayVersion
func CGoXrayVersion() *C.char {
	return C.CString(libbox.XrayVersion())
}
