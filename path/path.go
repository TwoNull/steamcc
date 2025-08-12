//go:build !windows

package path

import (
	"fmt"

	"github.com/twonull/steamcc/types"
)

// Unimplmented
func GetLoginUsers(steamDir string) ([]types.User, error) {
	return nil, fmt.Errorf("unimplemented")
}

// Unimplmented
func GetConnectCache(steamDir string) (map[string]string, error) {
	return nil, fmt.Errorf("unimplemented")
}
