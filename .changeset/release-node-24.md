---
"@subtextdev/subtext-cli": patch
---

Fix the release job's npm upgrade step failing with EBADENGINE by running on Node 24 (npm@latest now requires Node ^22.22.2 || ^24.15.0 || >=26.0.0), restoring the Node version used before the changesets rewrite.
