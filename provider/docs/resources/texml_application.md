---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "telnyx_texml_application Resource - telnyx"
subcategory: ""
description: |-
  Resource for managing Telnyx TeXML applications
---

# telnyx_texml_application (Resource)

Resource for managing Telnyx TeXML applications



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `friendly_name` (String) User-assigned name for the application
- `inbound` (Attributes) Inbound settings for the TeXML application (see [below for nested schema](#nestedatt--inbound))
- `outbound` (Attributes) Outbound settings for the TeXML application (see [below for nested schema](#nestedatt--outbound))
- `voice_url` (String) URL to deliver XML Translator webhooks

### Optional

- `active` (Boolean) Specifies whether the connection can be used
- `anchorsite_override` (String) Anchorsite Override
- `dtmf_type` (String) DTMF Type
- `first_command_timeout` (Boolean) Specifies whether calls should hangup after timing out
- `first_command_timeout_secs` (Number) Specifies how many seconds to wait before timing out a dial command
- `status_callback` (String) URL for status callback
- `status_callback_method` (String) HTTP request method for status callback
- `voice_fallback_url` (String) Fallback URL to deliver XML Translator webhooks if the primary URL fails
- `voice_method` (String) HTTP request method for voice webhooks

### Read-Only

- `created_at` (String) Creation time of the TeXML application
- `id` (String) Unique identifier of the TeXML application
- `updated_at` (String) Last update time of the TeXML application

<a id="nestedatt--inbound"></a>
### Nested Schema for `inbound`

Optional:

- `channel_limit` (Number) Limits the total number of inbound calls
- `shaken_stir_enabled` (Boolean) Enables Shaken/Stir data for inbound calls
- `sip_subdomain` (String) Subdomain for receiving inbound calls
- `sip_subdomain_receive_settings` (String) Receive calls from specified endpoints


<a id="nestedatt--outbound"></a>
### Nested Schema for `outbound`

Optional:

- `channel_limit` (Number) Limits the total number of outbound calls
- `outbound_voice_profile_id` (String) Associated outbound voice profile ID