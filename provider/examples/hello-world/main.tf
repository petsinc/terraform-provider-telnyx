provider "telnyx" {
  endpoint = "http://httpbin.org/post"
}

resource "telnyx_request" "test" {
  message = "Hello, World!"
}
