package readme

import "context"

type Category struct {
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	Order     int    `json:"order"`
	Reference bool   `json:"reference"`
	IsAPI     bool   `json:"isAPI"`
	Version   string `json:"version"`
	Project   string `json:"project"`
	CreatedAt string `json:"createdAt"`
	ID        string `json:"_id"`
}

type categories struct {
	client *Client
}

// CategoriesListOptions is a url-encodable data structure for options you can provide to the List endpoint
type CategoriesListOptions struct {
	PerPage int `url:"perPage,omitempty"`
	Page    int `url:"page,omitempty"`
}

// CategoriesList is the API response details of the List method
type CategoriesList struct {
	Pagination *Pagination
	Items      []*Category
}

// Categories describes the API methods available for the Categories API https://docs.readme.com/reference/getcategories
type Categories interface {
	List(ctx context.Context, options CategoriesListOptions) (*CategoriesList, error)
	Get(ctx context.Context, slug string) (*Category, error)
}

// List the categories according to some paging options
func (c *categories) List(ctx context.Context, opt CategoriesListOptions) (*CategoriesList, error) {
	response, pagination, err := c.client.getPaged(ctx, "categories", opt)

	if err != nil {
		return nil, err
	}

	result := CategoriesList{
		Pagination: pagination,
	}

	return &result, c.client.decodeAndClose(response.Body, &result.Items)
}

// Get a category using the specified slug
func (c *categories) Get(ctx context.Context, slug string) (*Category, error) {
	response, err := c.client.get(ctx, "categories/"+slug, nil)

	if err != nil {
		return nil, err
	}

	result := Category{}
	return &result, c.client.decodeAndClose(response.Body, &result)
}
