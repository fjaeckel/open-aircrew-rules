# Policy fragments

One YAML file per piece of work in the format of `policies.yaml` (schema
[`schema/policies.schema.json`](../../schema/policies.schema.json)):

```yaml
policies:
  - id: example-convention
    statement: What the catalogue does, in one sentence.
    rationale: Why, and why no regulation says it.
```

A policy is a convention with no legal text behind it, cited as `ref: policy:<id>`. First
look for an existing policy in `policies.yaml` that says the same; a new one is rare and a
reviewer will ask why it is not law. With `-fragments` the gate accepts `policy:` refs to
these ids; every policy must be cited at least once, and an id declared twice is an error.

Integration appends each policy to `policies.yaml` and deletes the fragment.
