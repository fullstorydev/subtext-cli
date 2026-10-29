# @subtextdev/subtext-cli

## 1.3.0

### Minor Changes

- 750f663: `sightmap upload` now sends the whole `.sightmap/` corpus — views with their routes, nested components and their `properties:`, requests, messages and memory — as the canonical JSON wire `sightmap export` writes, instead of a flat name/selectors list. Session review uses it to scope components by view and to put component properties on interaction signals (e.g. `[FlightCard flight="TP 799" price="398"]`); the flat list dropped all of that. The help and README now also name `review open`'s `sightmap_upload_url` as an upload target.

  Also bumps the bundled `github.com/sightmap/sightmap/go` library from v0.26.0 to v0.36.1, which v0.26.0 needed to load corpora using object-form `extract:` (it failed with "cannot unmarshal !!map into string").

### Patch Changes

- ca3c3f5: Update the bundled `github.com/sightmap/sightmap/go` library from v0.17.0 to v0.26.0. `sightmap upload` now parses `.sightmap/` corpora with the current loader and matcher — `$ref` expansion, view-scoped components, requests/messages/properties, and the newer route-matching semantics. No change to the CLI's own surface or flags.

  Also remove dead "For workflow guidance" skill links from `subtext tunnel --help` and the root help text — the referenced skill pages don't exist. The live tool listing shown in each namespace's `--help` remains the guidance.

## 1.2.0

### Minor Changes

- bd48d45: `sightmap upload` now collects `.sightmap/` via the shared `github.com/sightmap/sightmap/go` library instead of a local, from-scratch YAML walker. Two visible changes:

  - Uploaded components now carry `tags` (authored classification labels, e.g. `defect`) when a `.sightmap/` definition declares them — previously dropped entirely.
  - A child component no longer inherits its parent's `source:` when it doesn't declare its own (matches the shared library's convention, which also doesn't inherit `memory`/`stability`/`properties`). Add an explicit `source:` to a child that needs one.

## 1.1.7

### Patch Changes

- cc4a963: Fix the release job's npm upgrade step failing with EBADENGINE by running on Node 24 (npm@latest now requires Node ^22.22.2 || ^24.15.0 || >=26.0.0), restoring the Node version used before the changesets rewrite.

## 1.1.6

### Patch Changes

- 93df32e: Fix npm publish failing with ENEEDAUTH by upgrading npm to a version that supports OIDC trusted publishing before the publish step.

## 1.1.5

### Patch Changes

- 0d5b845: Verify the changesets release automation end-to-end: no functional change.
