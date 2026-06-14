package auth

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	keyring "github.com/zalando/go-keyring"
)

const (
	service = "hitori"
	tokenKey = "token"
)

// LoadToken returns stored JWT or "" if not found
func LoadToken() string {
	t, err := keyring.Get(service, tokenKey)
	if err != nil {
		return ""
	}
	return t
}

// StoreToken saves JWT to system keyring
func StoreToken(token string) error {
	return keyring.Set(service, tokenKey, token)
}

// DeleteToken clears stored token
func DeleteToken() error {
	return keyring.Delete(service, tokenKey)
}

// OpenBrowser opens URL in browser
func OpenBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux":
		// prefer xdg-open; fall back to sensible-browser
		if _, err := exec.LookPath("xdg-open"); err == nil {
			cmd, args = "xdg-open", []string{url}
		} else {
			cmd, args = "sensible-browser", []string{url}
		}
	case "darwin":
		cmd, args = "open", []string{url}
	case "windows":
		cmd, args = "cmd", []string{"/c", "start", url}
	default:
		fmt.Fprintf(os.Stderr, "open this URL to authenticate:\n%s\n", url)
		return nil
	}

	return exec.Command(cmd, args...).Start()
}
