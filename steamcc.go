package steamcc

import (
	"encoding/hex"
	"fmt"
	"hash/crc32"

	"github.com/twonull/steamcc/crypto"
	"github.com/twonull/steamcc/path"
	"github.com/twonull/steamcc/types"
)

type Client struct {
	SteamDir string // Defaults to %ProgramFiles(X86)%/Steam/ on Windows
}

// Returns slice of Users in loginusers.vdf
func (c *Client) GetUsers() ([]types.User, error) {
	users, err := path.GetLoginUsers(c.SteamDir)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Takes in AccountName, returns AutoLogin Refresh Token if available
func (c *Client) GetTokenForUser(user string) (string, error) {
	cc, err := path.GetConnectCache(c.SteamDir)
	if err != nil {
		return "", err
	}

	checksum := fmt.Sprintf("%08x", crc32.ChecksumIEEE([]byte(user)))

	for u, v := range cc {
		if u[0:8] == checksum {
			dataIn, err := hex.DecodeString(v.(string))
			if err != nil {
				return "", err
			}
			dataOut, err := crypto.Decrypt([]byte(user), []byte(dataIn))
			if err != nil {
				return "", err
			}
			return string(dataOut), nil
		}
	}

	return "", fmt.Errorf("no connectcache entry for user")
}
