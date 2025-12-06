//go:build linux

package ffi

import (
	"lizalang/ast"
	"unsafe"
)

// https://github.com/libffi/libffi
// https://www.youtube.com/watch?v=0o8Ex8mXigU
// https://stackoverflow.com/questions/27506579/calling-functions-in-an-so-file-from-go

//#cgo LDFLAGS: -ldl
//#include <dlfcn.h>
//#include <ffi.h>
import "C"

func LoadLib(external string) unsafe.Pointer {
	handle := C.dlopen(C.CString(external), C.RTLD_LAZY)
	return handle
}

func GetSymbol(symbol string, handle unsafe.Pointer) unsafe.Pointer {
	s := C.dlsym(handle, C.CString(symbol))
	return s
}

func LoadFn(functionDeclaration ast.FunctionDeclarationStatement) unsafe.Pointer {
	var cif C.ffi_cif
	C.ffi_prep_cif(&cif, C.FFI_DEFAULT_ABI, 0, nil, nil) //ffi_status ffi_prep_cif(ffi_cif *cif, ffi_abi abi, unsigned int nargs, ffi_type *rtype, ffi_type **argtypes)
	return unsafe.Pointer(&cif)
}
