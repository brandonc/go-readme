package readme

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocs_Search(t *testing.T) {
	t.Run("can search by query", func(t *testing.T) {
		client, err := NewClient(nil)

		assert.Nil(t, err)

		results, err := client.Docs.Search(context.Background(), "readme")

		assert.Nil(t, err)
		assert.Greater(t, len(results.Results), 0)
	})

	t.Run("get doc by slug", func(t *testing.T) {
		client, err := NewClient(nil)

		assert.Nil(t, err)

		result, err := client.Docs.Get(context.Background(), "getting-started")

		assert.Nil(t, err)
		assert.Equal(t, "Getting Started with go-readme-int-test", result.Title)
		assert.Equal(t, "getting-started", result.Slug)
		assert.NotEmpty(t, result.Body)
		assert.NotEmpty(t, result.BodyHTML)
		assert.NotEmpty(t, result.Category)
		assert.NotEmpty(t, result.CreatedAt)
		assert.NotEmpty(t, result.UpdatedAt)
		assert.Equal(t, "basic", result.Type)
	})
}
