package readme

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersions_List(t *testing.T) {
	t.Run("can list versions", func(t *testing.T) {
		client, err := NewClient(nil)

		assert.Nil(t, err)

		list, err := client.Versions.List(context.Background())

		assert.Nil(t, err)
		assert.Greater(t, len(list.Items), 0, "expected nonempty list of items")
	})
}

func TestVersions_CreateUpdateDelete(t *testing.T) {
	t.Run("can create versions", func(t *testing.T) {
		client, err := NewClient(nil)

		assert.Nil(t, err)

		opt := VersionCreateOptions{
			Version:      "2.0",
			From:         "1.0",
			IsStable:     newBool(false),
			IsBeta:       newBool(false),
			IsDeprecated: newBool(false),
			IsHidden:     newBool(true),
		}

		one, err := client.Versions.Get(context.Background(), "1.0")

		assert.Nil(t, err)

		ver, err := client.Versions.Create(context.Background(), opt)
		assert.Nil(t, err)
		assert.Equal(t, "2.0", ver.Version)
		assert.Equal(t, one.ID, ver.ForkedFrom)
		assert.False(t, ver.IsStable)
		assert.True(t, ver.IsHidden)
		assert.False(t, ver.IsDeprecated)
		assert.False(t, ver.IsBeta)

		optu := VersionUpdateOptions{
			IsHidden: newBool(false),
		}

		cu, err := client.Versions.Update(context.Background(), ver.Version, optu)

		assert.Nil(t, err)
		assert.True(t, cu.IsHidden)
		assert.Nil(t, client.Versions.Delete(context.Background(), cu.Version))
	})
}
