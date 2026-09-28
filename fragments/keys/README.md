# Message key fragments

One YAML file per piece of work (`<authority>-<kind dir>-<name>.yaml`) in the format of
`messages/keys.yaml` (schema [`schema/messages.schema.json`](../../schema/messages.schema.json)):

```yaml
keys:
  - { key: rating.mep_revalidation_current, kind: message, params: [] }
  - { key: requirement.single_engine_landings, kind: requirement_name, params: [] }
```

Add only keys that do not exist yet; reuse existing keys wherever the words fit. With
`-fragments` the gate loads these keys after `messages/keys.yaml`; a key defined twice
(in the shared file or in two fragments) is an error, so check the other fragments before
adding a generic key.

Integration moves each key into the section of its family in `messages/keys.yaml` and
deletes the fragment.
