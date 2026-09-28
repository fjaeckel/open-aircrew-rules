# Fragments

When several people (or agents) add credentials at the same time, the files everyone would
otherwise edit are not edited directly. Each contributor writes a **fragment** instead, and
an integration step merges all fragments into the shared files in one go:

| Directory | Merged into | Format |
| --- | --- | --- |
| [coverage/](coverage/README.md) | `coverage/articles.yaml` | changes per article (YAML) |
| [keys/](keys/README.md) | `messages/keys.yaml` | new message keys (YAML) |
| [policies/](policies/README.md) | `policies.yaml` | new policies (YAML) |
| [changelog/](changelog/README.md) | `CHANGELOG.md` | lines for "Unreleased" (Markdown) |
| [vocab-requests/](vocab-requests/README.md) | `vocabulary.yaml` (and engine, compiler, schema) | requests for new words (Markdown) |

Name each fragment after the work it belongs to, one file per credential or per family,
e.g. `easa-ratings-mep-land.yaml` or `faa-instructors.md`, so two contributors never write
the same file. Credential files, examples and `sources/` texts are not fragments: they are
new files and go where they belong.

## The gate

`go run ./cmd/rulescheck -strict -fragments` merges `coverage/`, `keys/` and `policies/`
fragments in memory, checks the result as if they were merged, validates each fragment
against its schema and lists every fragment still pending. Without `-fragments` the gate
fails on any fragment file (this README and the ones in the directories excepted), so the
main branch never ships unmerged fragments.

## Integration

1. Run the gate with `-fragments`: it must pass.
2. Merge each fragment into its shared file as its README says, keeping the shared file's
   order and style, and delete the fragment.
3. Implement accepted vocabulary requests (vocabulary, engine or compiler, schema, a worked
   example), then run the gate without `-fragments`: it must pass with no fragment left.
