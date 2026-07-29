# Releasing the Subtext CLI

Releases are automated via [changesets](https://github.com/changesets/changesets).
There is no manual version bump and no manual `git tag`.

## Prerequisites

- Write access to this repo (`fullstorydev/subtext-cli`).
- (First release only) `@subtextdev` npm scope access and a Trusted Publisher
  configured on npmjs.com (see below). Publishing uses OIDC provenance; no
  `NPM_TOKEN` secret.

## Release process

### 1. Add a changeset with your PR

If your change is user-facing (should bump the version / appear in the
changelog), run:

```sh
npm run changeset
```

Pick `patch` / `minor` / `major`, write a short summary, and commit the
resulting `.changeset/*.md` file with your PR. See
[`.changeset/README.md`](.changeset/README.md) for the full workflow.

### 2. Merge to main

The `release` workflow (`.github/workflows/release.yml`) notices pending
changesets and opens/updates a **"Version Packages" PR** that bumps
`npm/package.json` and writes `npm/CHANGELOG.md`.

### 3. Merge the Version Packages PR

This is the release trigger. With no changesets left to consume, the same
workflow:

1. Runs `go test ./...`.
2. Tags the commit (`vX.Y.Z`, matching the version the PR just set) and pushes
   the tag.
3. Runs GoReleaser, which builds binaries for darwin/linux/windows ×
   amd64/arm64, creates archives + `checksums.txt`, and creates a GitHub
   Release named `vX.Y.Z`.
4. Publishes `@subtextdev/subtext-cli` to npm via OIDC (no `NPM_TOKEN` needed).

Watch it under **Actions → Release** in the repo.

### 4. Smoke test

```bash
# via npm
npx @subtextdev/subtext-cli@0.2.0 --version

# via go install
go install github.com/fullstorydev/subtext-cli/cmd/subtext@v0.2.0
subtext --version
```

## Snapshot builds (local testing)

```bash
goreleaser release --snapshot --clean --skip=publish
./dist/subtext_darwin_arm64_v8.0/subtext version
```

No tag or GitHub credentials needed. Binaries land in `dist/`.

## First release checklist

- [ ] Confirm `@subtextdev` npm scope exists on npmjs.com.
- [ ] Configure a Trusted Publisher on npmjs.com: GitHub Actions, org
      `fullstorydev`, repo `subtext-cli`, workflow `release.yml`, allow
      `npm publish`.
- [ ] Test `npx @subtextdev/subtext-cli auth whoami` after publish.
