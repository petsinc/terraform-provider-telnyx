---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "telnyx_fqdn Resource - telnyx"
subcategory: ""
description: |-
  Resource for managing Telnyx FQDNs
---

# telnyx_fqdn (Resource)

Resource for managing Telnyx FQDNs



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `connection_id` (Number) ID of the connection associated with the FQDN
- `dns_record_type` (String) DNS record type
- `fqdn` (String) Fully Qualified Domain Name
- `port` (Number) Port associated with the FQDN

### Read-Only

- `created_at` (String) ISO 8601 formatted date indicating when the resource was created
- `id` (String) Unique identifier of the FQDN
- `updated_at` (String) ISO 8601 formatted date indicating when the resource was updated