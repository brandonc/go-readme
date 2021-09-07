package readme

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomPages_List(t *testing.T) {
	t.Run("can list custompages", func(t *testing.T) {
		client, err := NewClient(nil)

		assert.Nil(t, err)

		list, err := client.CustomPages.List(context.Background(), CustomPagesListOptions{
			PerPage: 2,
			Page:    1,
		})

		assert.Nil(t, err)
		assert.Greater(t, len(list.Items), 1, "expected nonempty list of items")

		assert.Equal(t, "", list.Pagination.Prev)
		assert.Equal(t, "/api/v1/custompages?perPage=2&page=2", list.Pagination.Next)
		assert.Equal(t, "/api/v1/custompages?perPage=2&page=4", list.Pagination.Last)
	})
}

func TestCustomPages_CreateUpdateDelete(t *testing.T) {
	t.Run("can create custompages", func(t *testing.T) {
		client, err := NewClient(nil)

		assert.Nil(t, err)

		opt := CustomPageCreateOptions{
			Title:    "Testing custompages",
			Body:     "# I must be valid markdown",
			HtmlMode: newBool(false),
			Hidden:   newBool(true),
			Metadata: Metadata{
				Description: "help I'm a useless description",
				Title:       "what is a meta title?",
			},
		}

		cl, err := client.CustomPages.Create(context.Background(), opt)
		assert.Nil(t, err)
		assert.NotEmpty(t, cl.Title)
		assert.NotEmpty(t, cl.Slug)
		assert.Equal(t, "help I'm a useless description", cl.Metadata.Description)
		assert.Equal(t, "what is a meta title?", cl.Metadata.Title)

		optu := CustomPageUpdateOptions{
			Title:  "Updated Testing custompages",
			Hidden: newBool(false),
		}

		cu, err := client.CustomPages.Update(context.Background(), cl.Slug, optu)

		assert.Nil(t, err)
		assert.Equal(t, "Updated Testing custompages", cu.Title)
		assert.Equal(t, false, cu.Hidden)

		assert.Nil(t, client.CustomPages.Delete(context.Background(), cl.Slug))
	})

	t.Run("can parse error messages", func(t *testing.T) {
		client, err := NewClient(nil)

		assert.Nil(t, err)

		opt := CustomPageCreateOptions{
			Body:   "# I must be valid markdown",
			Hidden: newBool(true),
		}

		_, err = client.CustomPages.Create(context.Background(), opt)
		assert.NotNil(t, err)
		assert.True(t, strings.HasPrefix(err.Error(), "CUSTOMPAGE_INVALID: We couldn't save this page (Custom page title cannot be blank). (See "))
	})
}
