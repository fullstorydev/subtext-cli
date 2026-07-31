package sightmap

import (
	"fmt"
	"os"
	"path/filepath"

	sm "github.com/sightmap/sightmap/go/sightmap"
)

// Component is a single named element with fully-resolved CSS selectors, the wire shape
// this package's Collect builds and Upload sends.
type Component struct {
	Name      string   `json:"name"`
	Selectors []string `json:"selectors"`
	Source    string   `json:"source"`
	Memory    []string `json:"memory"`
	Tags      []string `json:"tags,omitempty"`
}

// Payload is the JSON body sent to the sightmap upload endpoint.
type Payload struct {
	Sightmap []Component `json:"sightmap"`
	Memory   []string    `json:"memory"`
}

// Collect loads root/.sightmap/ via the shared sightmap library and returns the flat
// upload payload: every component in the corpus (github.com/sightmap/sightmap/go/sightmap's
// Corpus.AllComponents — globals plus every view's, deduped by name) projected onto the
// wire shape, plus the corpus's file-level memory. Returns an empty payload (not an error)
// when no .sightmap/ directory exists — sightmap.Load has no such tolerance (it walks the
// directory and errors if it's missing), so that contract is preserved here explicitly.
func Collect(root string) (Payload, error) {
	sdir := filepath.Join(root, ".sightmap")
	if _, err := os.Stat(sdir); os.IsNotExist(err) {
		return Payload{}, nil
	}

	corpus, err := sm.Load(sdir)
	if err != nil {
		return Payload{}, fmt.Errorf("sightmap: load %s: %w", sdir, err)
	}

	p := Payload{Memory: corpus.Memory}
	for _, c := range corpus.AllComponents() {
		p.Sightmap = append(p.Sightmap, Component{
			Name:      c.Name,
			Selectors: c.Selectors,
			Source:    c.Source,
			Memory:    c.Memory,
			Tags:      c.Tags,
		})
	}
	return p, nil
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
