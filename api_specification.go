package readme

import (
	"context"
	"net/http"
)

type api_specification struct {
	client *Client
}

// ApiSpecifications describes the API methods available for the api-specifications API https://docs.readme.com/reference/getapispecification
type ApiSpecifications interface {
	List(ctx context.Context, version string, opt ApiSpecificationListOptions) (*ApiSpecificationList, error)
	Upload(ctx context.Context, opt ApiSpecificationUploadOptions) (*ApiSpecificationStub, error)
	Update(ctx context.Context, id string, opt ApiSpecificationUpdateOptions) (*ApiSpecificationStub, error)
	Delete(ctx context.Context, id string) error
}

// ApiSpecificationListOptions is the options available for the list endpoint of the api-specifications
type ApiSpecificationListOptions struct {
	PerPage int `url:"perPage,omitempty"`
	Page    int `url:"page,omitempty"`
}

// ApiSpecificationList is the result from the api-specifications list endpoint
type ApiSpecificationList struct {
	Pagination *Pagination
	Items      []*ApiSpecification
}

// ApiSpecificationCategory is the category created for uploaded api specifications
type ApiSpecificationCategory struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
	Order int    `json:"order"`
	ID    string `json:"_id"`
}

// ApiSpecificationUploadOptions are the options available when uploading a new api specification
type ApiSpecificationUploadOptions struct {
	SpecPath string
	Version  string
}

// ApiSpecificationUpdateOptions are the options available when updating an existing api specification
type ApiSpecificationUpdateOptions struct {
	SpecPath string
}

// ApiSpecificationStub is the result of the upload and update endpoints when interacting with api specifications
type ApiSpecificationStub struct {
	Title string `json:"title"`
	ID    string `json:"_id"`
}

// ApiSpecification represents the metadata details from an api specification
type ApiSpecification struct {
	Title      string                    `json:"title"`
	Source     string                    `json:"source"`
	Version    string                    `json:"version"`
	LastSynced string                    `json:"lastSynced"`
	Category   *ApiSpecificationCategory `json:"category"`
	Type       string                    `json:"oas"`
	ID         string                    `json:"id"`
}

// List the api specifications according to some paging options
func (a *api_specification) List(ctx context.Context, version string, opt ApiSpecificationListOptions) (*ApiSpecificationList, error) {
	url, err := addOptions("api-specification", opt)

	if err != nil {
		return nil, err
	}

	versionHeader := make(http.Header)
	versionHeader.Add("x-readme-version", version)

	response, err := a.client.do(ctx, "GET", url, nil, versionHeader)

	if err != nil {
		return nil, err
	}

	pagination, err := a.client.parsePaginationResponse(&response.Header)

	if err != nil {
		return nil, err
	}

	result := ApiSpecificationList{
		Pagination: pagination,
	}

	return &result, a.client.decodeAndClose(response.Body, &result.Items)
}

// Upload the api specification at the specified path for the specified version
func (a *api_specification) Upload(ctx context.Context, opt ApiSpecificationUploadOptions) (*ApiSpecificationStub, error) {
	body, header, err := a.client.newFileUploadBody(opt.SpecPath)

	if err != nil {
		return nil, err
	}

	response, err := a.client.do(ctx, "POST", "api-specification", body, header)

	if err != nil {
		return nil, err
	}

	result := ApiSpecificationStub{}
	return &result, a.client.decodeAndClose(response.Body, &result)
}

// Update an existing api specification using the specified path for the specified version
func (a *api_specification) Update(ctx context.Context, id string, opt ApiSpecificationUpdateOptions) (*ApiSpecificationStub, error) {
	body, header, err := a.client.newFileUploadBody(opt.SpecPath)

	if err != nil {
		return nil, err
	}

	response, err := a.client.do(ctx, "PUT", "api-specification/"+id, body, header)

	if err != nil {
		return nil, err
	}

	result := ApiSpecificationStub{}
	return &result, a.client.decodeAndClose(response.Body, &result)
}

// Delete an existing api specification
func (a *api_specification) Delete(ctx context.Context, id string) error {
	_, err := a.client.delete(ctx, "api-specification/"+id)
	return err
}
