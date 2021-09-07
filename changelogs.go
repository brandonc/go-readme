package readme

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

type changelogs struct {
	client *Client
}

// ChangelogsListOptions is a url-encodable data structure for options you can provide to the List endpoint
type ChangelogsListOptions struct {
	PerPage int `url:"perPage,omitempty"`
	Page    int `url:"page,omitempty"`
}

// Changelogs describes the API methods available for the Changelogs API https://docs.readme.com/reference/getchangelogs
type Changelogs interface {
	List(ctx context.Context, options ChangelogsListOptions) (*ChangelogsList, error)
	Create(ctx context.Context, changelog ChangelogCreateOptions) (*Changelog, error)
	Update(ctx context.Context, slug string, changelog ChangelogUpdateOptions) (*Changelog, error)
	Delete(ctx context.Context, slug string) error
}

const (
	// ChangelogTypeAdded represents an "added" change type
	ChangelogTypeAdded = "added"

	// ChangelogTypeAdded represents a "fixed" change type
	ChangelogTypeFixed = "fixed"

	// ChangelogTypeAdded represents an "improved" change type
	ChangelogTypeImproved = "improved"

	// ChangelogTypeAdded represents a "deprecated" change type
	ChangelogTypeDeprecated = "deprecated"

	// ChangelogTypeAdded represents a "removed" change type
	ChangelogTypeRemoved = "removed"
)

// Changelog contains the details of each Changelog
type Changelog struct {
	Metadata              Metadata `json:"metadata"`
	Title                 string   `json:"title"`
	Slug                  string   `json:"slug"`
	Body                  string   `json:"body"`
	Type                  string   `json:"type"`
	Hidden                bool     `json:"hidden"`
	PendingAlgoliaPublish bool     `json:"pendingAlgoliaPublish"`
	CreatedAt             string   `json:"createdAt"`
	Html                  string   `json:"html"`
}

// ChangelogCreateOptions is the API request body when creating a Changelog
type ChangelogCreateOptions struct {
	Title    string   `json:"title"`
	Body     string   `json:"body"`
	Type     string   `json:"type,omitempty"`
	Hidden   *bool    `json:"hidden,omitempty"`
	Metadata Metadata `json:"metadata,omitempty"`
}

// ChangelogUpdateOptions is the API request body when updating a Changelog
type ChangelogUpdateOptions struct {
	Title  string `json:"title,omitempty"`
	Body   string `json:"body,omitempty"`
	Type   string `json:"type,omitempty"`
	Hidden *bool  `json:"hidden,omitempty"`
}

// ChangelogsList is the API response details of the List method
type ChangelogsList struct {
	Pagination *Pagination
	Items      []*Changelog
}

// Delete a changelog by slug
func (c *changelogs) Delete(ctx context.Context, slug string) error {
	_, err := c.client.delete(ctx, "changelogs/"+slug)
	return err
}

// Update an existing changelog by slug
func (c *changelogs) Update(ctx context.Context, slug string, changelog ChangelogUpdateOptions) (*Changelog, error) {
	bodyBytes, err := json.Marshal(changelog)

	if err != nil {
		return nil, fmt.Errorf("could not marshal request body: %w", err)
	}

	response, err := c.client.put(ctx, "changelogs/"+slug, bytes.NewBuffer(bodyBytes))

	if err != nil {
		return nil, err
	}

	result := Changelog{}
	return &result, c.client.decodeAndClose(response.Body, &result)
}

// Create a new changelog
func (c *changelogs) Create(ctx context.Context, changelog ChangelogCreateOptions) (*Changelog, error) {
	bodyBytes, err := json.Marshal(changelog)

	if err != nil {
		return nil, fmt.Errorf("could not marshal request body: %w", err)
	}

	response, err := c.client.post(ctx, "changelogs", bytes.NewBuffer(bodyBytes))

	if err != nil {
		return nil, err
	}

	result := Changelog{}
	return &result, c.client.decodeAndClose(response.Body, &result)
}

// List the changelogs according to some paging options
func (c *changelogs) List(ctx context.Context, options ChangelogsListOptions) (*ChangelogsList, error) {
	response, pagination, err := c.client.getPaged(ctx, "changelogs", options)

	if err != nil {
		return nil, err
	}

	result := ChangelogsList{
		Pagination: pagination,
	}

	return &result, c.client.decodeAndClose(response.Body, &result.Items)
}
