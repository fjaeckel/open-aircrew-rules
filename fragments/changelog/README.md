# Changelog fragments

One Markdown file per piece of work (`<authority>-<kind dir>-<name>.md`) holding the lines
to add under "Unreleased" in `CHANGELOG.md`, grouped by the headings used there:

```markdown
### Added

#### EASA

- `easa.rating.mep-land` MEP (land) class rating: revalidation (proficiency check), passenger
  recency by day and by night, and its licence.
```

Write for pilots and integrators; name each credential or evaluation
(`<credential id>#<evaluation id>`). The gate does not read these files; with `-fragments`
it lists them.

Integration merges the lines into `CHANGELOG.md` under the matching headings and deletes
the fragment.
