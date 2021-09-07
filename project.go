package readme

import "context"

// ProjectMetadata is information about the project associated with the API Key/
type ProjectMetadata struct {
	Name      string `json:"name"`
	Subdomain string `json:"subdomain"`
	JwtSecret string `json:"jwtSecret"`
	BaseURL   string `json:"baseUrl"`
	Plan      string `json:"plan"`
}

// Project allows interaction with the Project API resource
type Project interface {
	Get(ctx context.Context) (*ProjectMetadata, error)
}

type project struct {
	client *Client
}

func (p *project) Get(ctx context.Context) (*ProjectMetadata, error) {
	response, err := p.client.get(ctx, "", nil)

	if err != nil {
		return nil, err
	}

	result := ProjectMetadata{}
	return &result, p.client.decodeAndClose(response.Body, &result)
}
