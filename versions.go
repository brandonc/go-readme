package readme

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

type versions struct {
	client *Client
}

// Changelogs describes the API methods available for the Changelogs API https://docs.readme.com/reference/getcustompages
type Versions interface {
	List(ctx context.Context) (*VersionsList, error)
	Get(ctx context.Context, versionId string) (*Version, error)
	Create(ctx context.Context, version VersionCreateOptions) (*Version, error)
	Update(ctx context.Context, versionId string, version VersionUpdateOptions) (*Version, error)
	Delete(ctx context.Context, versionId string) error
}

// VersionListItem contains the details of each Version from the List endpoint
type VersionListItem struct {
	Version      string `json:"version"`
	CodeName     string `json:"codename"`
	IsStable     bool   `json:"is_stable"`
	IsBeta       bool   `json:"is_beta"`
	IsHidden     bool   `json:"is_hidden"`
	IsDeprecated bool   `json:"is_deprecated"`
	ID           string `json:"_id"`
	CreatedAt    string `json:"createdAt"`
}

// Version contains the details of each Version from the List endpoint
type Version struct {
	Version      string   `json:"version"`
	VersionClean string   `json:"version_clean"`
	Categories   []string `json:"categories"`
	CodeName     string   `json:"codename"`
	IsStable     bool     `json:"is_stable"`
	IsBeta       bool     `json:"is_beta"`
	IsHidden     bool     `json:"is_hidden"`
	IsDeprecated bool     `json:"is_deprecated"`
	ID           string   `json:"_id"`
	CreatedAt    string   `json:"createdAt"`
	ForkedFrom   string   `json:"forked_from"`
	ReleaseDate  string   `json:"releaseDate"`
	Project      string   `json:"project"`
}

// VersionCreateOptions is the API request body when creating a Version
type VersionCreateOptions struct {
	Version      string `json:"version"`
	CodeName     string `json:"codename,omitempty"`
	From         string `json:"from"`
	IsStable     *bool  `json:"is_stable,omitempty"`
	IsBeta       *bool  `json:"is_beta,omitempty"`
	IsHidden     *bool  `json:"is_hidden,omitempty"`
	IsDeprecated *bool  `json:"is_deprecated,omitempty"`
}

// VersionUpdateOptions is the API request body when updating a Version
type VersionUpdateOptions struct {
	Version      string `json:"version"`
	CodeName     string `json:"codename,omitempty"`
	From         string `json:"from"`
	IsStable     *bool  `json:"is_stable,omitempty"`
	IsBeta       *bool  `json:"is_beta,omitempty"`
	IsHidden     *bool  `json:"is_hidden,omitempty"`
	IsDeprecated *bool  `json:"is_deprecated,omitempty"`
}

// VersionsList is the API response details of the List method
type VersionsList struct {
	Items []*VersionListItem
}

// Delete a custompage by slug
func (c *versions) Delete(ctx context.Context, versionId string) error {
	_, err := c.client.delete(ctx, "version/"+versionId)
	return err
}

// Update an existing custompage by versionId (semver, ex "1.0")
func (c *versions) Update(ctx context.Context, versionId string, version VersionUpdateOptions) (*Version, error) {
	bodyBytes, err := json.Marshal(version)

	if err != nil {
		return nil, fmt.Errorf("could not marshal request body: %w", err)
	}

	response, err := c.client.put(ctx, "version/"+versionId, bytes.NewBuffer(bodyBytes))

	if err != nil {
		return nil, err
	}

	result := Version{}
	return &result, c.client.decodeAndClose(response.Body, &result)
}

// Create a new changelog
func (c *versions) Create(ctx context.Context, version VersionCreateOptions) (*Version, error) {
	bodyBytes, err := json.Marshal(version)

	if err != nil {
		return nil, fmt.Errorf("could not marshal request body: %w", err)
	}

	response, err := c.client.post(ctx, "version", bytes.NewBuffer(bodyBytes))

	if err != nil {
		return nil, err
	}

	result := Version{}
	return &result, c.client.decodeAndClose(response.Body, &result)
}

// Get the Version specified by the versionId (semver, ex. "1.0")
func (c *versions) Get(ctx context.Context, versionId string) (*Version, error) {
	response, err := c.client.get(ctx, "version/"+versionId, nil)

	if err != nil {
		return nil, err
	}

	result := Version{}

	return &result, c.client.decodeAndClose(response.Body, &result)
}

// List the Versions
func (c *versions) List(ctx context.Context) (*VersionsList, error) {
	response, err := c.client.get(ctx, "version", nil)

	if err != nil {
		return nil, err
	}

	result := VersionsList{}

	return &result, c.client.decodeAndClose(response.Body, &result.Items)
}
