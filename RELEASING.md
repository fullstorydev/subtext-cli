# Releasing the Subtext CLI

## Prerequisites

- Write access to this repo (`fullstorydev/subtext-cli`).
- (First release only) `@subtextdev` npm scope access and a Trusted Publisher configured on npmjs.com (see below).

## Release process

### 1. Decide on the version

We use standard `vX.Y.Z` tags. The CLI version lives only in the git tag — there is no version file to update.

### 2. Update the npm wrapper version

The npm package version must be bumped manually before tagging so that `npm publish` picks up the right version and `install.js` points at the correct release download URL. The release workflow will fail the validation step if the tag and `package.json` version disagree.

```bash
# In npm/package.json, bump "version" to match your tag (without the "v" prefix)
# e.g. if you're tagging v0.2.0, set "version": "0.2.0"
```

Commit:

```bash
git add npm/package.json
git commit -m "bump npm version to 0.2.0"
git push
```

### 3. Tag and push

```bash
git tag v0.2.0
git push origin v0.2.0
```

This triggers the `release-cli` GitHub Actions workflow.

### 4. Watch the workflow

Go to **Actions → Release CLI** in the repo. It will:

1. Run `go test ./...` to confirm tests pass.
2. Run GoReleaser, which:
   - Builds binaries for darwin/linux/windows × amd64/arm64.
   - Creates tar.gz / zip archives and a `checksums.txt`.
   - Creates a GitHub Release named `vX.Y.Z`.
3. Publishes the npm package to `@subtextdev/subtext-cli` via OIDC (requires a Trusted Publisher configured on npmjs.com — no `NPM_TOKEN` secret needed).

### 5. Smoke test

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
- [ ] Configure a Trusted Publisher on npmjs.com: GitHub Actions, org `fullstorydev`, repo `subtext-cli`, workflow `release-cli.yml`, allow `npm publish`.
- [ ] Test `npx @subtextdev/subtext-cli auth whoami` after publish.
