package readme

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProject(t *testing.T) {
	t.Run("can fetch project info", func(t *testing.T) {
		client, err := NewClient(nil)
		assert.Nil(t, err)

		project, err := client.Project.Get(context.Background())
		assert.Nil(t, err)

		assert.Equal(t, "go-readme-int-test", project.Name)
	})
}
