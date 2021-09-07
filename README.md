# go-readme
---

A go HTTP client for interacting with the [readme.com](https://readme.com/) [API](https://docs.readme.com/reference/intro-to-the-readme-api)

### Example

```go
client, err := NewClient(&Config{
  ApiKey: "12345",
})

response, err := client.ApiSpecifications.Upload(context.Background(), ApiSpecificationUploadOptions{
  SpecPath: "./helloOpenAPI.json",
  Version: "1.0"
})
```
