package sightmap

import (
	"fmt"
	"os"
	"path/filepath"

	sm "github.com/sightmap/sightmap/go/sightmap"
)

// Collect loads root/.sightmap/ via the shared sightmap library and returns the corpus
// to upload. Upload sends it as the library's canonical wire (json.Marshal of
// sightmap.Corpus: globals, views with their routes and view-scoped components,
// requests, messages, signals, environments and memory), the same bytes `sightmap
// export` produces, so component properties, nesting and URL-aware view scoping
// reach the server intact. (This used to send a flat name/selectors/source/memory/tags
// list, which dropped all of that.)
//
// Returns nil (not an error) when no .sightmap/ directory exists — sightmap.Load has no
// such tolerance (it walks the directory and errors if it's missing), so that contract
// is preserved here explicitly.
func Collect(root string) (*sm.Corpus, error) {
	sdir := filepath.Join(root, ".sightmap")
	if _, err := os.Stat(sdir); os.IsNotExist(err) {
		return nil, nil
	}

	corpus, err := sm.Load(sdir)
	if err != nil {
		return nil, fmt.Errorf("sightmap: load %s: %w", sdir, err)
	}
	return corpus, nil
}

// IsEmpty reports whether a collected corpus has nothing to upload.
func IsEmpty(c *sm.Corpus) bool {
	return c == nil || (len(c.AllComponents()) == 0 && len(c.Memory) == 0 && len(c.Requests) == 0 && len(c.Messages) == 0)
}

// FindRoot resolves the sightmap root directory using this precedence:
//
//  1. SIGHTMAP_ROOT env var
//  2. configRoot (from config file, if non-empty)
//  3. Walk up from cwd looking for a directory that contains .sightmap/
//
// The library has no discovery of its own (Load/DirLoader take a directory directly) —
// finding that directory is a CLI concern, so it stays here.
func FindRoot(cwd, configRoot string) (string, error) {
	if v := os.Getenv("SIGHTMAP_ROOT"); v != "" {
		return v, nil
	}
	if configRoot != "" {
		return configRoot, nil
	}
	d := cwd
	for {
		if fi, err := os.Stat(filepath.Join(d, ".sightmap")); err == nil && fi.IsDir() {
			return d, nil
		}
		parent := filepath.Dir(d)
		if parent == d {
			break
		}
		d = parent
	}
	return "", fmt.Errorf("no .sightmap/ directory found (searched from %s)", cwd)
}
