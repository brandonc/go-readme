package readme

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

type docs struct {
	client *Client
}

const (
	DocTypeBasic = "basic"
	DocTypeError = "error"
	DocTypeLink  = "link"
)

type Doc struct {
	Metadata              Metadata `json:"metadata"`
	Title                 string   `json:"title"`
	Type                  string   `json:"type"`
	Slug                  string   `json:"slug"`
	Excerpt               string   `json:"excerpt"`
	Body                  string   `json:"body"`
	Order                 int      `json:"order"`
	IsReference           bool     `json:"isReference"`
	Hidden                bool     `json:"hidden"`
	LinkURL               string   `json:"link_url"`
	LinkExternal          bool     `json:"link_external"`
	PendingAlgoliaPublish bool     `json:"pendingAlgoliaPublish"`
	PreviousSlug          string   `json:"previousSlug"`
	SlugUpdatedAt         string   `json:"slugUpdatedAt"`
	User                  string   `json:"user"`
	Project               string   `json:"project"`
	Category              string   `json:"category"`
	CreatedAt             string   `json:"createdAt"`
	UpdatedAt             string   `json:"updatedAt"`
	Version               string   `json:"version"`
	IsAPI                 bool     `json:"isApi"`
	ID                    string   `json:"id"`
	BodyHTML              string   `json:"body_html"`
}

type DocError struct {
	Code string `json:"code"`
}

type DocCreateOptions struct {
	Title     string    `json:"title"`
	Category  string    `json:"category"`
	Type      string    `json:"type,omitempty"`
	Hidden    *bool     `json:"hidden,omitempty"`
	Order     *int      `json:"order,omitempty"`
	ParentDoc string    `json:"parentDoc,omitempty"`
	Error     *DocError `json:"error,omitempty"`
}

type DocUpdateOptions struct {
	Title     string    `json:"title"`
	Category  string    `json:"category"`
	Type      string    `json:"type,omitempty"`
	Hidden    *bool     `json:"hidden,omitempty"`
	Order     *int      `json:"order,omitempty"`
	ParentDoc string    `json:"parentDoc,omitempty"`
	Error     *DocError `json:"error,omitempty"`
}

type DocSearchResults struct {
	Results []*DocSearchResult `json:"results"`
}

type DocSearchResult struct {
	IndexName    string `json:"indexName"`
	Title        string `json:"title"`
	Slug         string `json:"slug"`
	Project      string `json:"project"`
	ReferenceID  string `json:"referenceId"`
	Subdomain    string `json:"subdomain"`
	InternalLink string `json:"internalLink"`
	ObjectID     string `json:"objectID"`
	URL          string `json:"url"`
}

// Changelogs describes the API methods available for the Changelogs API https://docs.readme.com/reference/getcustompages
type Docs interface {
	Get(ctx context.Context, slug string) (*Doc, error)
	Create(ctx context.Context, doc DocCreateOptions) (*Doc, error)
	Update(ctx context.Context, slug string, doc DocUpdateOptions) (*Doc, error)
	Delete(ctx context.Context, slug string) error
	Search(ctx context.Context, search string) (*DocSearchResults, error)
}

func (d *docs) Get(ctx context.Context, slug string) (*Doc, error) {
	response, err := d.client.get(ctx, "docs/"+slug, nil)

	if err != nil {
		return nil, err
	}

	result := Doc{}
	return &result, d.client.decodeAndClose(response.Body, &result)
}

func (d *docs) Create(ctx context.Context, doc DocCreateOptions) (*Doc, error) {
	bodyBytes, err := json.Marshal(doc)

	if err != nil {
		return nil, fmt.Errorf("could not marshal request body: %w", err)
	}

	response, err := d.client.post(ctx, "docs", bytes.NewBuffer(bodyBytes))

	if err != nil {
		return nil, err
	}

	result := Doc{}
	return &result, d.client.decodeAndClose(response.Body, &result)
}

func (d *docs) Update(ctx context.Context, slug string, doc DocUpdateOptions) (*Doc, error) {
	bodyBytes, err := json.Marshal(doc)

	if err != nil {
		return nil, fmt.Errorf("could not marshal request body: %w", err)
	}

	response, err := d.client.post(ctx, "docs/"+slug, bytes.NewBuffer(bodyBytes))

	if err != nil {
		return nil, err
	}

	result := Doc{}
	return &result, d.client.decodeAndClose(response.Body, &result)
}

func (d *docs) Delete(ctx context.Context, slug string) error {
	_, err := d.client.delete(ctx, "docs/"+slug)
	return err
}

func (d *docs) Search(ctx context.Context, search string) (*DocSearchResults, error) {
	response, err := d.client.post(ctx, "docs/search?search="+url.QueryEscape(search), nil)

	if err != nil {
		return nil, err
	}

	result := DocSearchResults{}
	return &result, d.client.decodeAndClose(response.Body, &result)
}
