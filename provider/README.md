# Terraform Telnyx Provider

This is a Terraform provider for Telnyx. Currently, it includes a placeholder "Hello World" resource for testing.

## Example Usage

```hcl
provider "telnyx" {
  endpoint = "http://httpbin.org/post"
}

resource "telnyx_request" "test" {
  message = "Hello, World!"
}
```
