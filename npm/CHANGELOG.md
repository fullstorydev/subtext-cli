# @subtextdev/subtext-cli

## 1.1.7

### Patch Changes

- cc4a963: Fix the release job's npm upgrade step failing with EBADENGINE by running on Node 24 (npm@latest now requires Node ^22.22.2 || ^24.15.0 || >=26.0.0), restoring the Node version used before the changesets rewrite.

## 1.1.6

### Patch Changes

- 93df32e: Fix npm publish failing with ENEEDAUTH by upgrading npm to a version that supports OIDC trusted publishing before the publish step.

## 1.1.5

### Patch Changes

- 0d5b845: Verify the changesets release automation end-to-end: no functional change.
