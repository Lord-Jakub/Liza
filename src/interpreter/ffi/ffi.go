//go:build linux

package ffi

import (
	"fmt"
	"lizalang/ast"
	"lizalang/interpreter/object"
	"unsafe"
)

// https://github.com/libffi/libffi
// https://www.youtube.com/watch?v=0o8Ex8mXigU
// https://stackoverflow.com/questions/27506579/calling-functions-in-an-so-file-from-go

//#cgo LDFLAGS: -ldl -lffi
//#include <dlfcn.h>
//#include <ffi.h>
//#include <stdlib.h>
import "C"

func LoadLib(external string) unsafe.Pointer {
	handle := C.dlopen(C.CString(external), C.RTLD_LAZY)
	fmt.Println(C.GoString((*C.char)(C.dlerror())))
	return handle
}

func GetSymbol(symbol string, handle unsafe.Pointer) unsafe.Pointer {
	s := C.dlsym(handle, C.CString(symbol))
	return s
}

func LoadFn(functionDeclaration ast.FunctionDeclarationStatement) unsafe.Pointer {
	cif := (*C.ffi_cif)(C.malloc(C.sizeof_ffi_cif))
	nargs := C.uint(len(functionDeclaration.Args))
	rtype := GetType(functionDeclaration.Type.T())
	argtypes := make([]*C.ffi_type, len(functionDeclaration.Args))
	for i, arg := range functionDeclaration.Args {
		argtypes[i] = GetType(arg.Type.T())
	}
	var argtypes_ptr **C.ffi_type
	if len(argtypes) > 0 {
		argtypes_ptr = &argtypes[0]
	}

	C.ffi_prep_cif(cif, C.FFI_DEFAULT_ABI, nargs, rtype, argtypes_ptr) // ffi_status ffi_prep_cif(ffi_cif *cif, ffi_abi abi, unsigned int nargs, ffi_type *rtype, ffi_type **argtypes)
	fmt.Println(C.GoString((*C.char)(C.dlerror())))

	return unsafe.Pointer(&cif)
}

func Call(cif unsafe.Pointer, fn unsafe.Pointer, rtype string, args []object.Object) object.Object {
	// allocate argument array in C
	// otherwise get panic: runtime error: cgo argument has Go pointer to unpinned Go pointer
	// C cannot acces Go heap
	cArgs := C.malloc(C.size_t(len(args)) * C.size_t(unsafe.Sizeof(uintptr(0))))
	argPtrs := (*[1 << 30]unsafe.Pointer)(cArgs)[:len(args):len(args)]

	for i, arg := range args {
		switch v := arg.GetValue().(type) {
		case int64:
			ptr := C.malloc(C.size_t(unsafe.Sizeof(C.int(0))))
			*(*C.int)(ptr) = C.int(int32(v))
			argPtrs[i] = ptr
		case float64:
			ptr := C.malloc(C.size_t(unsafe.Sizeof(C.float(0))))
			*(*C.float)(ptr) = C.float(v)
			argPtrs[i] = ptr
		case string:
			cstr := C.CString(v)
			argPtrs[i] = unsafe.Pointer(cstr)
		default:
			panic("unsupported type")
		}
	}
	cArgs = unsafe.Pointer(&argPtrs[0])
	rettype := GetType(rtype)
	ret := C.malloc(C.size_t(rettype.size))
	C.ffi_call((*C.ffi_cif)(cif), (*[0]byte)(fn), ret, (*unsafe.Pointer)(cArgs)) // void ffi_call (ffi_cif *cif, void *fn, void *rvalue, void **avalues)
	var result object.Object = &object.VoidObject{Value: *(*any)(ret)}
	switch rtype {
	case "int":
		result = &object.IntObject{Value: int64(*(*int)(ret))}
	case "float":
		result = &object.FloatObject{Value: float64(*(*float32)(ret))}
	case "string":
		result = &object.StringObject{Value: C.GoString(*(**C.char)(ret))}
	}

	C.free(cArgs)
	C.free(ret)
	return result
}

func GetType(t string) *C.ffi_type {
	switch t {
	case "int":
		return &C.ffi_type_sint32
	case "string":
		return &C.ffi_type_pointer
	case "void":
		return &C.ffi_type_void
	case "float":
		return &C.ffi_type_float
	}
	return &C.ffi_type{}
}
