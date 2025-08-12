//go:build windows

package path

import (
	"os"
	"path/filepath"

	"github.com/twonull/steamcc/types"

	"github.com/andygrunwald/vdf"
)

func GetLoginUsers(steamDir string) ([]types.User, error) {
	var users []types.User

	if steamDir == "" {
		steamDir = filepath.Join(os.Getenv("ProgramFiles(X86)"), "Steam")
	}

	file, err := os.Open(filepath.Join(steamDir, "config", "loginusers.vdf"))
	if err != nil {
		return nil, err
	}

	p := vdf.NewParser(file)
	m, err := p.Parse()
	if err != nil {
		return nil, err
	}

	for u, v := range m["users"].(map[string]map[string]string) {
		users = append(users, types.User{
			AccountName: v["AccountName"],
			PersonaName: v["PersonaName"],
			Steam64:     u,
			AutoLogin:   v["AllowAutoLogin"] == "1",
		})
	}

	return users, nil
}

// Attempts to obtain ConnectCache from local.vdf, falls back to config.vdf (legacy) if not present
func GetConnectCache(steamDir string) (map[string]string, error) {
	if steamDir == "" {
		steamDir = filepath.Join(os.Getenv("ProgramFiles(X86)"), "Steam")
	}

	cDir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	_, err = os.Stat(filepath.Join(cDir, "Steam", "local.vdf"))
	if err != nil {
		file, err := os.Open(filepath.Join(steamDir, "config", "config.vdf"))
		if err != nil {
			return nil, err
		}

		p := vdf.NewParser(file)
		m, err := p.Parse()
		if err != nil {
			return nil, err
		}

		return m["InstallConfigStore"].(map[string]interface{})["Software"].(map[string]interface{})["Valve"].(map[string]interface{})["Steam"].(map[string]interface{})["ConnectCache"].(map[string]string), nil
	}

	file, err := os.Open(filepath.Join(cDir, "Steam", "local.vdf"))
	if err != nil {
		return nil, err
	}

	p := vdf.NewParser(file)
	m, err := p.Parse()
	if err != nil {
		return nil, err
	}

	return m["MachineUserConfigStore"].(map[string]interface{})["Software"].(map[string]interface{})["Valve"].(map[string]interface{})["Steam"].(map[string]interface{})["ConnectCache"].(map[string]string), nil
}
