---
"@subtextdev/subtext-cli": minor
---

`sightmap upload` now sends the whole `.sightmap/` corpus — views with their routes, nested components and their `properties:`, requests, messages and memory — as the canonical JSON wire `sightmap export` writes, instead of a flat name/selectors list. Session review uses it to scope components by view and to put component properties on interaction signals (e.g. `[FlightCard flight="TP 799" price="398"]`); the flat list dropped all of that. The help and README now also name `review open`'s `sightmap_upload_url` as an upload target.

Also bumps the bundled `github.com/sightmap/sightmap/go` library from v0.26.0 to v0.36.1, which v0.26.0 needed to load corpora using object-form `extract:` (it failed with "cannot unmarshal !!map into string").
