package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/i4n-co/driplimit"
	"github.com/i4n-co/driplimit/pkg/config"
)

// httpClient is a driplimit http client that implements the driplimit service
type httpClient struct {
	client *http.Client
	cfg    *config.Config
	logger *slog.Logger
}

// New creates a new driplimit http client
func New(cfg *config.Config) driplimit.ServiceWithToken {
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: cfg.UpstreamTimeout,
		}).Dial,
		TLSHandshakeTimeout: cfg.UpstreamTimeout,
	}

	return &httpClient{
		cfg: cfg,
		client: &http.Client{
			Timeout:   cfg.UpstreamTimeout,
			Transport: transport,
		},
		logger: cfg.Logger().With("component", "client"),
	}
}

// httpClientWithToken is a driplimit http client that implements the driplimit service with a token
type httpClientWithToken struct {
	*httpClient
	token string
}

// WithToken creates a new driplimit http client instance with a token
func (c *httpClient) WithToken(token string) driplimit.Service {
	return &httpClientWithToken{
		httpClient: c,
		token:      token,
	}
}

// do sends a request to the upstream server
func (c *httpClientWithToken) do(ctx context.Context, action string, payload any, result any) error {
	jsn, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", c.cfg.UpstreamURL+action, bytes.NewBuffer(jsn))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "driplimit")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return driplimit.ErrFromHTTPCode(resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

func (c *httpClientWithToken) KeyCheck(ctx context.Context, payload driplimit.KeysCheckPayload) (key *driplimit.Key, err error) {
	key = new(driplimit.Key)
	err = c.do(ctx, "/v1/keys.check", payload, key)
	return key, err
}

func (c *httpClientWithToken) KeyCreate(ctx context.Context, payload driplimit.KeyCreatePayload) (key *driplimit.Key, token *string, err error) {
	key = new(driplimit.Key)
	err = c.do(ctx, "/v1/keys.create", payload, key)
	return key, &key.Token, err
}

func (c *httpClientWithToken) KeyGet(ctx context.Context, payload driplimit.KeyGetPayload) (key *driplimit.Key, err error) {
	key = new(driplimit.Key)
	err = c.do(ctx, "/v1/keys.get", payload, key)
	return key, err
}

func (c *httpClientWithToken) KeyList(ctx context.Context, payload driplimit.KeyListPayload) (klist *driplimit.KeyList, err error) {
	klist = new(driplimit.KeyList)
	err = c.do(ctx, "/v1/keys.list", payload, klist)
	return klist, err
}

func (c *httpClientWithToken) KeyDelete(ctx context.Context, payload driplimit.KeyDeletePayload) (err error) {
	return c.do(ctx, "/v1/keys.delete", payload, nil)
}

func (c *httpClientWithToken) KeyspaceGet(ctx context.Context, payload driplimit.KeyspaceGetPayload) (keyspace *driplimit.Keyspace, err error) {
	keyspace = new(driplimit.Keyspace)
	err = c.do(ctx, "/v1/keyspaces.get", payload, keyspace)
	return keyspace, err
}

func (c *httpClientWithToken) KeyspaceCreate(ctx context.Context, payload driplimit.KeyspaceCreatePayload) (keyspace *driplimit.Keyspace, err error) {
	keyspace = new(driplimit.Keyspace)
	err = c.do(ctx, "/v1/keyspaces.create", payload, keyspace)
	return keyspace, err
}

func (c *httpClientWithToken) KeyspaceList(ctx context.Context, payload driplimit.KeyspaceListPayload) (kslist *driplimit.KeyspaceList, err error) {
	kslist = new(driplimit.KeyspaceList)
	err = c.do(ctx, "/v1/keyspaces.list", payload, kslist)
	return kslist, err
}

func (c *httpClientWithToken) KeyspaceDelete(ctx context.Context, payload driplimit.KeyspaceDeletePayload) (err error) {
	return c.do(ctx, "/v1/keyspaces.delete", payload, nil)
}

// ServiceKeyGet returns a service key based on the given payload
func (c *httpClientWithToken) ServiceKeyGet(ctx context.Context, payload driplimit.ServiceKeyGetPayload) (sk *driplimit.ServiceKey, err error) {
	sk = new(driplimit.ServiceKey)
	err = c.do(ctx, "/v1/serviceKeys.get", payload, sk)
	return sk, err
}

func (c *httpClientWithToken) ServiceKeyCreate(ctx context.Context, payload driplimit.ServiceKeyCreatePayload) (sk *driplimit.ServiceKey, err error) {
	sk = new(driplimit.ServiceKey)
	err = c.do(ctx, "/v1/serviceKeys.create", payload, sk)
	return sk, err
}

func (c *httpClientWithToken) ServiceKeyList(ctx context.Context, payload driplimit.ServiceKeyListPayload) (sklist *driplimit.ServiceKeyList, err error) {
	sklist = new(driplimit.ServiceKeyList)
	err = c.do(ctx, "/v1/serviceKeys.list", payload, sklist)
	return sklist, err
}

func (c *httpClientWithToken) ServiceKeyDelete(ctx context.Context, payload driplimit.ServiceKeyDeletePayload) (err error) {
	return c.do(ctx, "/v1/serviceKeys.delete", payload, nil)
}
