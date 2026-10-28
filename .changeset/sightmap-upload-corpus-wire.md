---
"@subtextdev/subtext-cli": minor
---

`sightmap upload` now sends the whole `.sightmap/` corpus (the same wire as `sightmap export --url`) instead of a flat component list. Views keep their routes, `requests:` and `messages:` are included, and components keep their `properties:`, so review snapshots and signals get view scoping, request names, and extracted property values. Previously all of that was dropped on upload.

`sightmap upload --help` now also names `review open` as a source of the upload URL, and notes that the URL is single-use and short-lived.
