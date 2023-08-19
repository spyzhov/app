package app

import (
	"unsafe"
)

// isNil is a safe and fast method to check is current interface value is nil.
//
// Reason is incident, that nil != nil in some cases.
//
// Full description of the problem: https://dev.to/pauljlucas/go-tcha-when-nil--nil-hic
func isNil(i interface{}) bool {
	return (*[2]uintptr)(unsafe.Pointer(&i))[1] == 0
}
