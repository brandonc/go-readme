package readme

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupEnvVars(token string) func() {
	origToken := os.Getenv("README_API_KEY")

	os.Setenv("README_API_KEY", token)

	return func() {
		os.Setenv("README_API_KEY", origToken)
	}
}

func newBool(v bool) *bool {
	b := v
	return &b
}

func TestClient_NewClient(t *testing.T) {
	t.Run("has expected default config", func(t *testing.T) {
		defer setupEnvVars("testKey")()

		client, err := NewClient(nil)

		assert.Nil(t, err, "error occurred while creating client")
		assert.Equal(t, client.baseUrl.String(), DefaultAddress+DefaultBasePath+"/", "default client does not have default address")
		assert.Equal(t, client.apiKey, "testKey", "default client does not get api key from env")
	})

	t.Run("adds specified config to default config", func(t *testing.T) {
		cfg := Config{
			Headers: make(http.Header),
		}
		cfg.Headers.Add("x-custom-test", "hello, world")

		client, err := NewClient(&cfg)

		assert.Nil(t, err, "error occurred while creating client")
		assert.Equal(t, "go-readme", client.headers.Get("user-agent"), "user-agent header not set")
		assert.Equal(t, "hello, world", client.headers.Get("x-custom-test"), "custom request header not present")
		assert.NotNil(t, client.Changelogs, "Changelogs is not set")
	})

	t.Run("returns an error if url can't be parsed", func(t *testing.T) {
		cfg := Config{
			Address: "../dir/",
		}

		client, err := NewClient(&cfg)

		assert.NotNil(t, err, "expected NewClient to return an error")
		assert.Nil(t, client, "expected nil client")
	})
}
