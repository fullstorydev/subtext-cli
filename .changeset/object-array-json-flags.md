---
"@subtextdev/subtext-cli": patch
---

Passthrough flags for `object` parameters are now parsed as JSON, so the documented `subtext review zoom --resolution '{"network":"detail"}'` works (it was sent as a string and rejected by the server). `array` flags also accept a whole JSON array (`--rule_ids '["a","b"]'`) alongside repeated flags.
