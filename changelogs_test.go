package readme

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangeLogs_List(t *testing.T) {
	t.Run("can list changelogs", func(t *testing.T) {
		client, err := NewClient(nil)

		assert.Nil(t, err)

		list, err := client.Changelogs.List(context.Background(), ChangelogsListOptions{
			PerPage: 2,
		})

		assert.Nil(t, err)
		assert.Greater(t, len(list.Items), 1, "expected nonempty list of items")

		assert.Equal(t, "", list.Pagination.Prev)
		assert.Equal(t, "/api/v1/changelogs?perPage=2&page=2", list.Pagination.Next)
		assert.Equal(t, "/api/v1/changelogs?perPage=2&page=4", list.Pagination.Last)
	})
}

func TestChangeLogs_CreateUpdateDelete(t *testing.T) {
	t.Run("can create changelogs", func(t *testing.T) {
		client, err := NewClient(nil)

		assert.Nil(t, err)

		opt := ChangelogCreateOptions{
			Title:  "Testing changelogs",
			Body:   "<p>I must be valid markdown</p>",
			Type:   ChangelogTypeAdded,
			Hidden: newBool(true),
			Metadata: Metadata{
				Description: "help I'm a useless description",
				Title:       "what is a meta title?",
			},
		}

		cl, err := client.Changelogs.Create(context.Background(), opt)
		assert.Nil(t, err)
		assert.NotEmpty(t, cl.Title)
		assert.NotEmpty(t, cl.Slug)
		assert.Equal(t, "help I'm a useless description", cl.Metadata.Description)
		assert.Equal(t, "what is a meta title?", cl.Metadata.Title)

		optu := ChangelogUpdateOptions{
			Title:  "Updated Testing changelogs",
			Type:   ChangelogTypeDeprecated,
			Hidden: newBool(false),
		}

		cu, err := client.Changelogs.Update(context.Background(), cl.Slug, optu)

		assert.Nil(t, err)
		assert.Equal(t, "Updated Testing changelogs", cu.Title)
		assert.Equal(t, ChangelogTypeDeprecated, cu.Type)
		assert.Equal(t, false, cu.Hidden)

		assert.Nil(t, client.Changelogs.Delete(context.Background(), cu.Slug))
	})

	t.Run("can parse error messages", func(t *testing.T) {
		client, err := NewClient(nil)

		assert.Nil(t, err)

		opt := ChangelogCreateOptions{
			Body:   "<p>I must be valid markdown</p>",
			Type:   ChangelogTypeAdded,
			Hidden: newBool(true),
		}

		_, err = client.Changelogs.Create(context.Background(), opt)
		assert.NotNil(t, err)
		assert.True(t, strings.HasPrefix(err.Error(), "CHANGELOG_INVALID: We couldn't save this changelog (Changelog title cannot be blank). (See "))
	})
}
