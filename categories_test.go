package readme

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategories_List(t *testing.T) {
	t.Run("can list all categories", func(t *testing.T) {
		client, err := NewClient(nil)
		assert.Nil(t, err)

		opt := CategoriesListOptions{
			PerPage: 1,
		}

		list, err := client.Categories.List(context.Background(), opt)

		assert.Nil(t, err)

		assert.Greater(t, len(list.Items), 0, "expected nonempty list of items")

		assert.Equal(t, "", list.Pagination.Prev)
		assert.Equal(t, "/api/v1/categories?perPage=1&page=2", list.Pagination.Next)
		assert.NotEmpty(t, list.Pagination.Last)
	})

	t.Run("can get a category", func(t *testing.T) {
		client, err := NewClient(nil)
		assert.Nil(t, err)

		category, err := client.Categories.Get(context.Background(), "documentation")

		assert.Nil(t, err)
		assert.Equal(t, "Documentation", category.Title)
	})

	t.Run("returns error when get nonexisting category", func(t *testing.T) {
		client, err := NewClient(nil)
		assert.Nil(t, err)

		category, err := client.Categories.Get(context.Background(), "snazzy")

		assert.Nil(t, category)
		assert.True(t, strings.HasPrefix(err.Error(), "CATEGORY_NOTFOUND: The category with the slug 'snazzy' couldn't be found. (See"))
	})
}
