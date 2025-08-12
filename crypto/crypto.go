//go:build !windows

package crypto

import "fmt"

type DATA_BLOB struct {
	cbData uint32
	pbData uintptr
}

// Unimplemented
func Decrypt(entropy []byte, dataIn []byte) ([]byte, error) {
	return nil, fmt.Errorf("unimplemented")
}
