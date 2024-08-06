Driplimit works with a JSON RPC API. All calls are made via HTTP using the POST method. 
The API returns 200 OK or 204 CREATED on success.

See [Errors](https://github.com/i4n-co/driplimit/wiki/RPC-API-V1#errors) section in case of failure.

{{ range $namespace, $docs := .RPCs }}
## {{ title $namespace }}
{{ range $doc := $docs }}
### `POST /v1{{ $doc.RPCDocumentation.Path }}`
{{ $doc.RPCDocumentation.Description }}

**Headers**

* `Content-Type: application/json` - tells the service you wish to communicate with json
* `Authorization: Bearer <token>`  - the service key token

**Parameters**

{{ range $field := $doc.ParamFields }}
* {{if ne $field.Name ""}}`{{ $field.Name }}` {{end}}<span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">{{ $field.Type }}</span> - {{ $field.Description }}{{ if $field.Required }} (required){{ end }}
{{ range $subfield := $field.SubFields }}
  * {{if ne $subfield.Name ""}}`{{ $subfield.Name }}` {{end}}<span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">{{ $subfield.Type }}</span> - {{ $subfield.Description }}{{ if $subfield.Required }} (required){{ end }}
{{ range $subsubfield := $subfield.SubFields }}
    * {{if ne $subsubfield.Name ""}}`{{ $subsubfield.Name }}` {{end}}<span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">{{ $subsubfield.Type }}</span> - {{ $subsubfield.Description }}{{ if $subsubfield.Required }} (required){{ end }}
{{ end }}
{{ end }}
{{ end }}

<details>
<summary> <b>cURL example</b> </summary>

```bash
$ curl -X POST \
       -H "Content-Type: application/json" \
       -H "Authorization: Bearer <token>" \{{ if ne $doc.RPCDocumentation.Parameters nil }}
       --data '
{{ toPrettyJson $doc.RPCDocumentation.Parameters | indent 8 }}'{{end}} https://demo.driplim.it/v1{{ $doc.RPCDocumentation.Path }}
```

```json
{{ toPrettyJson $doc.RPCDocumentation.Response }}
```
</details>
{{ end }}
{{ end }}

## Errors

If HTTP response code is greater than or equal to 400, the api returns a json object indicating the reason of the failure:

```json
{
  "error": "the reason of the failure",
  "invalid_fields": [
    "field1",
    "field2"
  ]
}
```

`invalid_fields` can also be integrated in the error response if one or more input parameters are invalids.

### HTTP response code

* `200` ok
* `204` created
* `400` invalid payload
* `401` unauthorized
* `403` cannot delete itself
* `404` not found
* `409` already exists
* `419` key expired
* `429` rate limit exceeded
