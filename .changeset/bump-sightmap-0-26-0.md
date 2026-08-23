---
"@subtextdev/subtext-cli": patch
---

Update the bundled `github.com/sightmap/sightmap/go` library from v0.17.0 to v0.26.0. `sightmap upload` now parses `.sightmap/` corpora with the current loader and matcher — `$ref` expansion, view-scoped components, requests/messages/properties, and the newer route-matching semantics. No change to the CLI's own surface or flags.

Also remove dead "For workflow guidance" skill links from `subtext tunnel --help` and the root help text — the referenced skill pages don't exist. The live tool listing shown in each namespace's `--help` remains the guidance.
