# @subtextdev/subtext-cli

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
