package readme

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

type custompages struct {
	client *Client
}

// ChangelogsListOptions is a url-encodable data structure for options you can provide to the List endpoint
type CustomPagesListOptions struct {
	PerPage int `url:"perPage,omitempty"`
	Page    int `url:"page,omitempty"`
}

// Changelogs describes the API methods available for the Changelogs API https://docs.readme.com/reference/getcustompages
type CustomPages interface {
	List(ctx context.Context, options CustomPagesListOptions) (*CustomPagesList, error)
	Get(ctx context.Context, slug string) (*CustomPage, error)
	Create(ctx context.Context, changelog CustomPageCreateOptions) (*CustomPage, error)
	Update(ctx context.Context, slug string, changelog CustomPageUpdateOptions) (*CustomPage, error)
	Delete(ctx context.Context, slug string) error
}

// Changelog contains the details of each Changelog
type CustomPage struct {
	Metadata              Metadata `json:"metadata"`
	Title                 string   `json:"title"`
	Slug                  string   `json:"slug"`
	Body                  string   `json:"body"`
	Hidden                bool     `json:"hidden"`
	Fullscreen            bool     `json:"fullscreen"`
	Html                  string   `json:"html"`
	HtmlMode              bool     `json:"htmlmode"`
	CreatedAt             string   `json:"createdAt"`
	PendingAlgoliaPublish bool     `json:"pendingAlgoliaPublish"`
}

// ChangelogCreateOptions is the API request body when creating a Changelog
type CustomPageCreateOptions struct {
	Title    string   `json:"title"`
	Body     string   `json:"body"`
	Html     string   `json:"html"`
	HtmlMode *bool    `json:"htmlmode,omitempty"`
	Hidden   *bool    `json:"hidden,omitempty"`
	Metadata Metadata `json:"metadata,omitempty"`
}

// ChangelogUpdateOptions is the API request body when updating a Changelog
type CustomPageUpdateOptions struct {
	Title    string   `json:"title,omitempty"`
	Body     string   `json:"body,omitempty"`
	Html     string   `json:"html"`
	HtmlMode *bool    `json:"htmlmode,omitempty"`
	Hidden   *bool    `json:"hidden,omitempty"`
	Metadata Metadata `json:"metadata,omitempty"`
}

// ChangelogsList is the API response details of the List method
type CustomPagesList struct {
	Pagination *Pagination
	Items      []*CustomPage
}

// Delete a custompage by slug
func (c *custompages) Delete(ctx context.Context, slug string) error {
	_, err := c.client.delete(ctx, "custompages/"+slug)
	return err
}

// Update an existing custompage by slug
func (c *custompages) Update(ctx context.Context, slug string, custompage CustomPageUpdateOptions) (*CustomPage, error) {
	bodyBytes, err := json.Marshal(custompage)

	if err != nil {
		return nil, fmt.Errorf("could not marshal request body: %w", err)
	}

	response, err := c.client.put(ctx, "custompages/"+slug, bytes.NewBuffer(bodyBytes))

	if err != nil {
		return nil, err
	}

	result := CustomPage{}
	return &result, c.client.decodeAndClose(response.Body, &result)
}

// Create a new changelog
func (c *custompages) Create(ctx context.Context, custompage CustomPageCreateOptions) (*CustomPage, error) {
	bodyBytes, err := json.Marshal(custompage)

	if err != nil {
		return nil, fmt.Errorf("could not marshal request body: %w", err)
	}

	response, err := c.client.post(ctx, "custompages", bytes.NewBuffer(bodyBytes))

	if err != nil {
		return nil, err
	}

	result := CustomPage{}
	return &result, c.client.decodeAndClose(response.Body, &result)
}

// Get the custompage specified by the slug
func (c *custompages) Get(ctx context.Context, slug string) (*CustomPage, error) {
	response, err := c.client.get(ctx, "custompages/"+slug, nil)

	if err != nil {
		return nil, err
	}

	result := CustomPage{}

	return &result, c.client.decodeAndClose(response.Body, &result)
}

// List the custompages according to some paging options
func (c *custompages) List(ctx context.Context, options CustomPagesListOptions) (*CustomPagesList, error) {
	response, pagination, err := c.client.getPaged(ctx, "custompages", options)

	if err != nil {
		return nil, err
	}

	result := CustomPagesList{
		Pagination: pagination,
	}

	return &result, c.client.decodeAndClose(response.Body, &result.Items)
}
