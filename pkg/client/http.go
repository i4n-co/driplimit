package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/i4n-co/driplimit"
)

// HTTP is a driplimit http client that implements the driplimit service
type HTTP struct {
	upstreamURL     string
	sendRequestFunc func(req *http.Request) (*http.Response, error)
	serviceToken    string
}

// New creates a new driplimit http client
func New(upstreamURL string, timeout ...time.Duration) *HTTP {
	timeoutDuration := 5 * time.Second
	if len(timeout) > 0 && timeout[0] > 0 {
		timeoutDuration = timeout[0]
	}
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: timeoutDuration,
		}).DialContext,
		TLSHandshakeTimeout: timeoutDuration,
	}
	client := &http.Client{
		Timeout:   timeoutDuration,
		Transport: transport,
	}
	return &HTTP{
		upstreamURL:     upstreamURL,
		sendRequestFunc: client.Do,
	}
}

// WithSendRequestFunc replace the default client http Do func by a custom one.
// This method is mainly used for tests
func (h *HTTP) WithSendRequestFunc(f func(req *http.Request) (*http.Response, error)) *HTTP {
	h.sendRequestFunc = f
	return h
}

// WithServiceToken sets a default service token used in all requests
func (h *HTTP) WithServiceToken(token string) *HTTP {
	h.serviceToken = token
	return h
}

// do sends a request to the upstream server
func do[K any](ctx context.Context, c *HTTP, action string, payload driplimit.Payload, target ...K) (err error) {
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(payload)
	if err != nil {
		return fmt.Errorf("failed to encode payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.upstreamURL+action, buf)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "driplimit")
	if c.serviceToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.serviceToken)
	}

	if payload != nil && payload.ServiceToken() != "" {
		req.Header.Set("Authorization", "Bearer "+payload.ServiceToken())
	}

	resp, err := c.sendRequestFunc(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return driplimit.ErrFromHTTPCode(resp.StatusCode)
	}

	bodybuf := new(bytes.Buffer)
	io.Copy(bodybuf, resp.Body)

	r := bytes.NewReader(bodybuf.Bytes())
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	fmt.Println(string(b))

	r.Seek(0, 0)
	if len(target) > 0 {
		err = json.NewDecoder(r).Decode(target[0])
		if err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

func (c *HTTP) KeyCheck(ctx context.Context, payload driplimit.KeysCheckPayload) (key *driplimit.Key, err error) {
	key = new(driplimit.Key)
	err = do(ctx, c, "/v1/keys.check", payload, key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (c *HTTP) KeyCreate(ctx context.Context, payload driplimit.KeyCreatePayload) (key *driplimit.Key, token *string, err error) {
	key = new(driplimit.Key)
	err = do(ctx, c, "/v1/keys.create", payload, key)
	if err != nil {
		return nil, nil, err
	}
	return key, &key.Token, nil
}

func (c *HTTP) KeyGet(ctx context.Context, payload driplimit.KeyGetPayload) (key *driplimit.Key, err error) {
	key = new(driplimit.Key)
	err = do(ctx, c, "/v1/keys.get", payload, key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (c *HTTP) KeyList(ctx context.Context, payload driplimit.KeyListPayload) (klist *driplimit.KeyList, err error) {
	klist = new(driplimit.KeyList)
	err = do(ctx, c, "/v1/keys.list", payload, klist)
	if err != nil {
		return nil, err
	}
	return klist, nil
}

func (c *HTTP) KeyDelete(ctx context.Context, payload driplimit.KeyDeletePayload) (err error) {
	err = do[driplimit.Key](ctx, c, "/v1/keys.delete", payload)
	if err != nil {
		return err
	}
	return nil
}

func (c *HTTP) KeyspaceGet(ctx context.Context, payload driplimit.KeyspaceGetPayload) (keyspace *driplimit.Keyspace, err error) {
	keyspace = new(driplimit.Keyspace)
	err = do(ctx, c, "/v1/keyspaces.get", payload, keyspace)
	if err != nil {
		return nil, err
	}
	return keyspace, nil
}

func (c *HTTP) KeyspaceCreate(ctx context.Context, payload driplimit.KeyspaceCreatePayload) (keyspace *driplimit.Keyspace, err error) {
	keyspace = new(driplimit.Keyspace)
	err = do(ctx, c, "/v1/keyspaces.create", payload, keyspace)
	if err != nil {
		return nil, err
	}
	return keyspace, nil
}

func (c *HTTP) KeyspaceList(ctx context.Context, payload driplimit.KeyspaceListPayload) (kslist *driplimit.KeyspaceList, err error) {
	kslist = new(driplimit.KeyspaceList)
	err = do(ctx, c, "/v1/keyspaces.list", payload, kslist)
	if err != nil {
		return nil, err
	}
	return kslist, nil
}

func (c *HTTP) KeyspaceDelete(ctx context.Context, payload driplimit.KeyspaceDeletePayload) (err error) {
	err = do(ctx, c, "/v1/keyspaces.delete", payload, make(map[any]any))
	if err != nil {
		return err
	}
	return nil
}

// ServiceKeyGet returns a service key based on the given payload
func (c *HTTP) ServiceKeyCurrent(ctx context.Context) (sk *driplimit.ServiceKey, err error) {
	sk = new(driplimit.ServiceKey)
	err = do(ctx, c, "/v1/serviceKeys.current", nil, sk)
	if err != nil {
		return nil, err
	}
	return sk, nil
}

// ServiceKeyGet returns a service key based on the given payload
func (c *HTTP) ServiceKeyGet(ctx context.Context, payload driplimit.ServiceKeyGetPayload) (sk *driplimit.ServiceKey, err error) {
	sk = new(driplimit.ServiceKey)
	err = do(ctx, c, "/v1/serviceKeys.get", payload, sk)
	if err != nil {
		return nil, err
	}
	return sk, nil
}

func (c *HTTP) ServiceKeyCreate(ctx context.Context, payload driplimit.ServiceKeyCreatePayload) (sk *driplimit.ServiceKey, err error) {
	sk = new(driplimit.ServiceKey)
	err = do(ctx, c, "/v1/serviceKeys.create", payload, sk)
	if err != nil {
		return nil, err
	}
	return sk, nil
}

func (c *HTTP) ServiceKeyList(ctx context.Context, payload driplimit.ServiceKeyListPayload) (sklist *driplimit.ServiceKeyList, err error) {
	sklist = new(driplimit.ServiceKeyList)
	err = do(ctx, c, "/v1/serviceKeys.list", payload, sklist)
	if err != nil {
		return nil, err
	}
	return sklist, nil
}

func (c *HTTP) ServiceKeyDelete(ctx context.Context, payload driplimit.ServiceKeyDeletePayload) (err error) {
	err = do(ctx, c, "/v1/serviceKeys.delete", payload, make(map[any]any))
	if err != nil {
		return err
	}
	return nil
}

func (c *HTTP) ServiceKeySetToken(ctx context.Context, payload driplimit.ServiceKeySetTokenPayload) (err error) {
	err = do(ctx, c, "/v1/serviceKeys.set_token", payload, make(map[any]any))
	if err != nil {
		return err
	}
	return nil
}
