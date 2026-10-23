## New OpenRouter field

When OpenRouter adds a new field, a message like the following can appear:

```bash
jenai capa glm
```
> 2026/10/20 19:18:45 json: unable to unmarshal JSON string into Go openrouter.ModelEndpoint: unknown object member name "service_tier" within "/data/endpoints/0"
> [...]

This means the generated OpenRouter structs needs to be updated by running `./scripts/regen-oapi-structs.bash` from the root of the repo.

## Updating model zoo

To update stale information from the model zoo or to populate the fields from a newly-added model, use:
```bash
jenai model-zoo regen go/models/models.yaml
```
