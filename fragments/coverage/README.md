# Coverage fragments

One YAML file per piece of work (`<authority>-<kind dir>-<name>.yaml`), schema
[`schema/coverage-fragment.schema.json`](../../schema/coverage-fragment.schema.json):

```yaml
articles:
  - article: "easa:FCL.740.A"
    evaluations: [easa.rating.mep-land#revalidation]   # added to the article's evaluations
    done_pending: [easa/ratings/mep-land]              # removed from its pending files
  - article: "easa:FCL.060"
    evaluations: [easa.rating.mep-land#passengers_day, easa.rating.mep-land#passengers_night]
    done_pending: [easa/ratings/mep-land]
    note: "Still to convert beyond the evaluations above: passenger recency for SET land, helicopter types and the LAPL(A)."
```

| Key | Merge |
| --- | --- |
| `article` | the article (`coverage/articles.yaml` entry); an article not listed yet is added |
| `evaluations` | appended to the article's `evaluations` |
| `pending` | appended to `pending` (a credential file still to write) |
| `done_pending` | removed from `pending`; it must be pending there |
| `note` | replaces the note; the note is dropped when no pending file is left |
| `not_evaluated` | replaces the reason |

Every compiled evaluation must end up under the article of its `source`, and a pending
file that now exists fails the gate, exactly as for `coverage/articles.yaml`. Rewrite the
`note` when you remove a pending file and others remain, so it only names what is still to
convert.

Integration applies the same rules to `coverage/articles.yaml` by hand, in the file's
order, then deletes the fragment.
