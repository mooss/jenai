## New OpenRouter field

When OpenRouter adds a new field, a message like the following can appear:

```bash
jenai capa glm
```
> 2026/10/20 19:18:45 json: unable to unmarshal JSON string into Go openrouter.ModelEndpoint: unknown object member name "service_tier" within "/data/endpoints/0"
> [...]

### Download and resolve

```bash
curl -fSL https://openrouter.ai/openapi.json > openrouter.openapi.json
go run ./scripts/resolve-openapi-refs.go -ref-siblings merge openrouter.openapi.json > openrouter-resolved.openapi.json
```

### Identify missing path

First identify the initial path that was queried by following the code breadcrumbs (here `models/{author}/{slug}/endpoints`).

Select this path using `jq`:
```bash
jq '.paths["/models/{author}/{slug}/endpoints"]' openrouter-resolved.openapi.json > models-endpoints.json
du -sh models-endpoints.json openrouter-resolved.openapi.json
```
> 128K    models-endpoints.json
> 47M     openrouter-resolved.openapi.json


Restricting the path makes it easier to search and keeping the file at hand can help when the issue is not as straightforward as in this guide.
```
jq -c 'paths | select(.[-1] | tostring | endswith("service_tier"))' models-endpoints.json
```
```json
["get","responses","200","content","application/json","schema","properties","data","properties","endpoints","items","example","service_tier"]
["get","responses","200","content","application/json","schema","properties","data","properties","endpoints","items","properties","service_tier"]
```

Here the first path is clearly an example and the second one is pertinent.
```bash
jq --argjson path\
  '["get","responses","200","content","application/json","schema","properties","data","properties","endpoints","items","properties","service_tier"]'\
  'getpath($path)' models-endpoints.json
```
```json
{
  "description": "The service tier this endpoint sells, or null for a standard endpoint. The same value whether the provider selects the tier with a request parameter or with a dedicated model id. Always the canonical tier: priority covers the fast request alias and the turbo, highspeed, and prime provider model-id spellings, none of which this field reports.",
  "enum": [
    "flex",
    "priority",
    "ultrafast",
    null
  ],
  "example": "priority",
  "type": [
    "string",
    "null"
  ],
  "x-speakeasy-unknown-values": "allow"
}
```

### Generate missing field

Piping the previous command to `jenai` with some added context should be enough to generate the missing content:
```bash
jq --argjson path\
  '["get","responses","200","content","application/json","schema","properties","data","properties","endpoints","items","properties","service_tier"]'\
  'getpath($path)' models-endpoints.json |\
  jenai o Explain how existing Go structs can be updated to add support for this OpenAPI spec for the path\
  '["get","responses","200","content","application/json","schema","properties","data","properties","endpoints","items","properties","service_tier"]'.\
  Also include YAML tags. Make it in the same style as the provided files. The generated structs will live in the same package so pay attention to name collision and reuse the types already here everywhere possible.\
  Make sure to specify the paths to the files that must be updated.\
  --tt missing-service_tier.md --model luna --file go/llm/openrouter/*-types.go
```

`jenai` will output instructions on how to update the code.

## Updating model zoo

To update stale information from the model zoo or to populate the fields from a newly-added model, use:
```bash
jenai model-zoo regen go/models/models.yaml
```
