//go:build windows

package crypto

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	crypt32                = syscall.NewLazyDLL("crypt32.dll")
	kernel32               = syscall.NewLazyDLL("kernel32.dll")
	procCryptUnprotectData = crypt32.NewProc("CryptUnprotectData")
	procLocalFree          = kernel32.NewProc("LocalFree")
)

type DATA_BLOB struct {
	cbData uint32
	pbData uintptr
}

// Decryption via DPAPI CryptUnprotectData with AccountName as entropy
func Decrypt(entropy []byte, dataIn []byte) ([]byte, error) {
	var entropyBlob DATA_BLOB
	var dataInBlob DATA_BLOB
	var dataOutBlob DATA_BLOB
	var pDescrOut uintptr

	if len(entropy) > 0 {
		entropyBlob.cbData = uint32(len(entropy))
		entropyBlob.pbData = uintptr(unsafe.Pointer(&entropy[0]))
	}

	dataInBlob.cbData = uint32(len(dataIn))
	dataInBlob.pbData = uintptr(unsafe.Pointer(&dataIn[0]))

	ret, _, err := procCryptUnprotectData.Call(
		uintptr(unsafe.Pointer(&dataInBlob)),
		uintptr(unsafe.Pointer(&pDescrOut)),
		uintptr(unsafe.Pointer(&entropyBlob)),
		0,   // NULL
		0,   // NULL
		0x1, // CRYPTPROTECT_UI_FORBIDDEN
		uintptr(unsafe.Pointer(&dataOutBlob)),
	)

	if ret == 0 {
		return nil, fmt.Errorf("CryptUnprotectData failed: %v", err)
	}

	// Convert to Go slice
	resultSize := int(dataOutBlob.cbData)
	result := make([]byte, resultSize)
	for i := 0; i < resultSize; i++ {
		result[i] = *(*byte)(unsafe.Pointer(dataOutBlob.pbData + uintptr(i)))
	}

	if dataOutBlob.pbData != 0 {
		procLocalFree.Call(dataOutBlob.pbData)
	}
	if pDescrOut != 0 {
		procLocalFree.Call(pDescrOut)
	}

	return result, nil
}
