package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuthClient_GetProfile_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/dev-profiles/profile", r.URL.Path)

		cookie, err := r.Cookie("jwt")
		assert.NoError(t, err)
		assert.Equal(t, "valid-token", cookie.Value)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id": 123}`))
	}))
	defer server.Close()

	client := NewAuthClient(server.URL)

	// act
	id, err := client.GetProfile("valid-token")

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, id)
	assert.Equal(t, int64(123), *id)
}

func TestAuthClient_GetProfile_UnexpectedStatusCode(t *testing.T) {
	// arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewAuthClient(server.URL)

	// act
	id, err := client.GetProfile("invalid-token")

	// assert
	assert.Nil(t, id)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status code")
}

func TestAuthClient_GetProfile_InvalidJSON(t *testing.T) {
	// arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`invalid-json`))
	}))
	defer server.Close()

	client := NewAuthClient(server.URL)

	// act
	id, err := client.GetProfile("valid-token")

	// assert
	assert.Nil(t, id)
	assert.Error(t, err)
}

type failingRoundTripper struct{}

func (f *failingRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("network error")
}

func TestAuthClient_GetProfile_DoError(t *testing.T) {
	// arrange
	httpClient := &http.Client{
		Transport: &failingRoundTripper{},
		Timeout:   1 * time.Second,
	}

	client := &AuthClient{
		baseURL:    "http://fake-url",
		httpClient: httpClient,
	}

	// act
	id, err := client.GetProfile("any-token")

	// assert
	assert.Nil(t, id)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "network error")
}
