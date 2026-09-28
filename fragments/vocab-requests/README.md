# Vocabulary requests

`vocabulary.yaml` is not edited during parallel work: a new word needs engine or compiler
code, a schema entry and a worked example, and two contributors adding words at once
conflict. When a credential cannot be written with the existing words, write one Markdown
file per request (`<authority>-<kind dir>-<name>-<word>.md`):

```markdown
# Request: count word `single_engine_landings`

- Section: credential_vocabulary.counts
- Needed by: easa.rating.mep-land#revalidation (easa:FCL.740.A(b)(3))
- Compiles to: metric `landings.total` with filter `engines: 1`
- Why the existing words do not fit: ...
- Worked example: the record and the expected outcome that will show it
```

Until the request is implemented the credential cannot use the word: leave the evaluation
out (keep the file in `pending` in your coverage fragment) or express it with existing
words and an interpretation that says what is approximated. The gate does not read these
files; with `-fragments` it lists them.

Integration implements accepted requests (vocabulary, engine or compiler, schema, a worked
example), then the requesting contributor completes the evaluation; the request file is
deleted.
