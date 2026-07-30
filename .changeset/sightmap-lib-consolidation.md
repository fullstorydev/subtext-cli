---
"@subtextdev/subtext-cli": minor
---

`sightmap upload` now collects `.sightmap/` via the shared `github.com/sightmap/sightmap/go` library instead of a local, from-scratch YAML walker. Two visible changes:

- Uploaded components now carry `tags` (authored classification labels, e.g. `defect`) when a `.sightmap/` definition declares them — previously dropped entirely.
- A child component no longer inherits its parent's `source:` when it doesn't declare its own (matches the shared library's convention, which also doesn't inherit `memory`/`stability`/`properties`). Add an explicit `source:` to a child that needs one.
