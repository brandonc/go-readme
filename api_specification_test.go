package readme

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApiSpecification_List(t *testing.T) {
	t.Run("list api specifications", func(t *testing.T) {
		client, err := NewClient(nil)

		assert.Nil(t, err)

		list, err := client.ApiSpecifications.List(context.Background(), "1.0", ApiSpecificationListOptions{})

		assert.Nil(t, err)
		assert.Greater(t, len(list.Items), 0)
	})
}

func TestApiSpecification_CreateUpdateDelete(t *testing.T) {
	t.Run("can upload, update, delete", func(t *testing.T) {
		client, err := NewClient(nil)

		assert.Nil(t, err)

		uploaded, err := client.ApiSpecifications.Upload(context.Background(), ApiSpecificationUploadOptions{
			SpecPath: "./fixtures/petstore.json",
			Version:  "1.0",
		})

		assert.Nil(t, err)

		assert.Equal(t, "Swagger Petstore", uploaded.Title)
		assert.NotEmpty(t, uploaded.ID)

		updated, err := client.ApiSpecifications.Update(context.Background(), uploaded.ID, ApiSpecificationUpdateOptions{
			SpecPath: "./fixtures/petstore.json",
		})

		assert.Nil(t, err)

		assert.Equal(t, "Swagger Petstore", updated.Title)
		assert.NotEmpty(t, updated.ID)

		err = client.ApiSpecifications.Delete(context.Background(), updated.ID)
		assert.Nil(t, err)
	})
}
