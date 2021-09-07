package readme

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/brandonc/go-weblinks"
	"github.com/google/go-querystring/query"
)

const (
	userAgent = "go-readme"

	// DefaultAddress is the default host of the readme API
	DefaultAddress = "https://dash.readme.com"

	// DefaultAddress is the default base path for all reqeusts
	DefaultBasePath = "/api/v1"
)

// Pagination holds pagination response data
type Pagination struct {
	// TotalCount is the total number of records available
	TotalCount int

	// Next is the API path of the next page
	Next string

	// Prev is the API path of the previous page
	Prev string

	// Last is the API path of the last page
	Last string
}

type Config struct {
	// Address is the host of the readme API
	Address string

	// ApiKey to authenticate requests
	ApiKey string

	// Headers are the request headers sent with all requests
	Headers http.Header

	// HttpClient is a default pooled http client
	HttpClient *http.Client
}

// Client is the primary object used to interact with the readme API
type Client struct {
	baseUrl *url.URL
	apiKey  string
	headers http.Header
	http    *http.Client

	// Changelogs allows interactions with Changelog API resources
	Changelogs Changelogs

	// CustomPages allows interactions with CustomPages API resources
	CustomPages CustomPages

	// Docs allows interactions with Docs API resources
	Docs Docs

	// Categories allows interactions with Categories API resources
	Categories Categories

	// Project allows fetching metadata about the current project
	Project Project

	// Versions allows interactions with the Versions API resources
	Versions Versions

	// ApiSpecification allows interactions with ApiSpecification API resources
	ApiSpecifications ApiSpecifications
}

// ErrorResponse is the standard json error details given for server errors
type errorResponse struct {
	ErrorCode  string `json:"error"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion"`
	DocsUrl    string `json:"docs"`
	Help       string `json:"help"`
}

// Metadata contains the metadata details of each page
type Metadata struct {
	Image       []string `json:"image"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
}

// DefaultConfig returns a default config using environment settings if available
func DefaultConfig() *Config {
	config := &Config{
		Address:    DefaultAddress,
		ApiKey:     os.Getenv("README_API_KEY"),
		Headers:    make(http.Header),
		HttpClient: http.DefaultClient,
	}

	config.Headers.Add("user-agent", userAgent)
	config.Headers.Add("content-type", "application/json")

	return config
}

func handleErrorResponse(response *http.Response) error {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("could not read error body: %w", err)
	}

	if strings.HasPrefix(response.Header.Get("content-type"), "application/json") {
		serverError := errorResponse{}
		if err = json.Unmarshal(body, &serverError); err != nil {
			return err
		}

		return fmt.Errorf("%s: %s (See %s)", serverError.ErrorCode, serverError.Message, serverError.DocsUrl)
	} else {
		return fmt.Errorf("server error (%v): %s", response.StatusCode, body)
	}
}

func (c *Client) do(ctx context.Context, method string, path string, reader io.Reader, header http.Header) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, method, c.baseUrl.String()+path, reader)

	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	request.Header = c.headers

	for name, value := range header {
		request.Header.Set(name, strings.Join(value, ", "))
	}

	response, err := c.http.Do(request)
	if err != nil {
		return nil, fmt.Errorf("could not perform %s request: %w", method, err)
	}

	if response.StatusCode >= 400 {
		return nil, handleErrorResponse(response)
	}

	return response, nil
}

func (c *Client) post(ctx context.Context, path string, reader io.Reader) (*http.Response, error) {
	return c.do(ctx, "POST", path, reader, nil)
}

func (c *Client) put(ctx context.Context, path string, reader io.Reader) (*http.Response, error) {
	return c.do(ctx, "PUT", path, reader, nil)
}

func (c *Client) get(ctx context.Context, path string, opts interface{}) (*http.Response, error) {
	url, err := addOptions(path, opts)

	if err != nil {
		return nil, fmt.Errorf("could not encode url options: %w", err)
	}

	return c.do(ctx, "GET", url, nil, nil)
}

func (c *Client) delete(ctx context.Context, path string) (*http.Response, error) {
	return c.do(ctx, "DELETE", path, nil, nil)
}

func (c *Client) getPaged(ctx context.Context, path string, opts interface{}) (*http.Response, *Pagination, error) {
	response, err := c.get(ctx, path, opts)

	if err != nil {
		return nil, nil, err
	}

	pagination, err := c.parsePaginationResponse(&response.Header)
	return response, pagination, err
}

func (c *Client) parsePaginationResponse(responseHeaders *http.Header) (*Pagination, error) {
	totalCount, err := strconv.Atoi(responseHeaders.Get("x-total-count"))

	if err != nil {
		return nil, fmt.Errorf("could not get total count from paginated response: %w", err)
	}

	log.Printf("[DEBUG] Total Count = %v, Link header = \"%v\"", responseHeaders.Get("x-total-count"), responseHeaders.Get("Link"))

	links, err := weblinks.Parse(responseHeaders.Get("Link"))

	if err != nil {
		return nil, fmt.Errorf("could not parse Link header: %w", err)
	}

	return &Pagination{
		Next:       links["next"].URI.String(),
		Prev:       links["prev"].URI.String(),
		Last:       links["last"].URI.String(),
		TotalCount: totalCount,
	}, nil
}

func (c *Client) decodeAndClose(body io.ReadCloser, v interface{}) error {
	defer body.Close()
	if v != nil {
		return json.NewDecoder(body).Decode(v)
	}
	return nil
}

// NewClient creates a new Client that can be used to interact with all readme.com APIs
func NewClient(cfg *Config) (*Client, error) {
	config := DefaultConfig()

	if cfg != nil {
		if cfg.Address != "" {
			config.Address = cfg.Address
		}

		if cfg.ApiKey != "" {
			config.ApiKey = cfg.ApiKey
		}

		for k, v := range cfg.Headers {
			config.Headers[k] = v
		}
	}

	baseUrl, err := url.ParseRequestURI(config.Address)
	if err != nil {
		return nil, fmt.Errorf("invalid address: %w", err)
	}

	baseUrl.Path = DefaultBasePath
	if !strings.HasSuffix(baseUrl.Path, "/") {
		baseUrl.Path += "/"
	}

	config.Headers.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(config.ApiKey+":")))

	client := &Client{
		baseUrl: baseUrl,
		apiKey:  config.ApiKey,
		headers: config.Headers,
		http:    config.HttpClient,
	}

	client.Changelogs = &changelogs{client: client}
	client.CustomPages = &custompages{client: client}
	client.Docs = &docs{client: client}
	client.Categories = &categories{client: client}
	client.Project = &project{client: client}
	client.Versions = &versions{client: client}
	client.ApiSpecifications = &api_specification{client: client}

	return client, nil
}

func (c *Client) newFileUploadBody(path string) (io.Reader, http.Header, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open file specified: %w", err)
	}

	if err != nil {
		return nil, nil, fmt.Errorf("could not read file specified: %w", err)
	}

	fi, err := file.Stat()
	if err != nil {
		return nil, nil, fmt.Errorf("could not stat the file specified: %w", err)
	}
	defer file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("spec", fi.Name())
	if err != nil {
		return nil, nil, err
	}
	io.Copy(part, file)

	header := make(http.Header)
	header.Add("Content-Type", writer.FormDataContentType())

	err = writer.Close()

	if err != nil {
		return nil, nil, err
	}

	return body, header, nil
}

func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
