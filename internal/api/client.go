package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/mathiiiiiis/hitori-tui/internal/mono"
)

const defaultBase = "https://api.playhitori.de"

// Client talks to backend
type Client struct {
	base  string
	token string
	http  *http.Client
}

func New(token string) *Client {
	base := os.Getenv("HITORI_API")
	if base == "" {
		base = defaultBase
	}
	return &Client{
		base:  base,
		token: token,
		http:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) do(method, path string, body any) (*http.Response, error) {
	var req *http.Request
	var err error
	if body != nil {
		b, mErr := json.Marshal(body)
		if mErr != nil {
			return nil, mErr
		}
		req, err = http.NewRequest(method, c.base+path, bytes.NewBuffer(b))
	} else {
		req, err = http.NewRequest(method, c.base+path, nil)
	}
	if err != nil {
		return nil, err
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return c.http.Do(req)
}

// ==== auth ====

type CLIInitResponse struct {
	SessionID string `json:"session_id"`
	AuthURL   string `json:"auth_url"`
}

func (c *Client) CLIAuthInit(provider string) (*CLIInitResponse, error) {
	resp, err := c.do("GET", "/auth/cli/init?provider="+provider, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("init failed: %s", resp.Status)
	}
	var r CLIInitResponse
	return &r, json.NewDecoder(resp.Body).Decode(&r)
}

// CLIAuthPoll returns ("", nil) while pending, (token, nil) when ready
func (c *Client) CLIAuthPoll(sessionID string) (string, error) {
	resp, err := c.do("GET", "/auth/cli/poll?session_id="+sessionID, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("session expired")
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("poll failed: %s", resp.Status)
	}
	var r struct {
		Status string `json:"status"`
		Token  string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", err
	}
	if r.Status == "ready" {
		return r.Token, nil
	}
	return "", nil
}

// ==== save ====

type saveResponse struct {
	Data      json.RawMessage `json:"data"`
	UpdatedAt *time.Time      `json:"updated_at"`
}

// LoadState fetches saved
// Returns (nil, nil) when theres no save yet
func (c *Client) LoadState() (*mono.State, error) {
	resp, err := c.do("GET", "/save", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("load failed: %s", resp.Status)
	}
	var sr saveResponse
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return nil, err
	}
	// empty save object => no Mono created yet
	if len(sr.Data) == 0 || string(sr.Data) == "{}" || string(sr.Data) == "null" {
		return nil, nil
	}
	var s mono.State
	if err := json.Unmarshal(sr.Data, &s); err != nil {
		return nil, err
	}
	if s.Name == "" {
		return nil, nil // unfinished/empty save
	}
	return &s, nil
}

func (c *Client) SaveState(s *mono.State) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	resp, err := c.do("PUT", "/save", json.RawMessage(data))
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("save failed: %s", resp.Status)
	}
	return nil
}
