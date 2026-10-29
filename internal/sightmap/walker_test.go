package sightmap

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	sm "github.com/sightmap/sightmap/go/sightmap"

	"github.com/fullstorydev/subtext-cli/internal/fstesting"
)

// TestCollect_Components exercises the full walk against testdata/components.yaml,
// which has two top-level components (one with string selector, one with list)
// and component-level + top-level memory.
func TestCollect_Components(t *testing.T) {
	p, err := Collect(testdataRoot(t))
	fstesting.Ok(t, err, "Collect")

	byName := indexByName(p.AllComponents())

	hdr, ok := byName["Header"]
	fstesting.Assert(t, ok, "missing Header component")
	fstesting.Equals(t, []string{"header"}, hdr.Selectors, "Header selectors")
	fstesting.Equals(t, "main-nav", hdr.Source, "Header source")
	fstesting.Equals(t, []string{"primary navigation area", "contains logo and search"}, hdr.Memory, "Header memory")

	search, ok := byName["SearchInput"]
	fstesting.Assert(t, ok, "missing SearchInput component")
	fstesting.Equals(t, []string{"input[type=search]", ".search-field"}, search.Selectors, "SearchInput selectors")
}

// TestCollect_TopLevelMemory verifies that top-level memory entries from all
// YAML files are collected into Payload.Memory.
func TestCollect_TopLevelMemory(t *testing.T) {
	p, err := Collect(testdataRoot(t))
	fstesting.Ok(t, err, "Collect")

	memSet := make(map[string]bool, len(p.Memory))
	for _, m := range p.Memory {
		memSet[m] = true
	}
	fstesting.Assert(t, memSet["top-level context about the app layout"], "missing top-level memory from components.yaml")
	fstesting.Assert(t, memSet["dashboard is the landing page after login"], "missing top-level memory from views.yaml")
}

// TestCollect_Nested verifies descendant-combinator flattening via testdata/nested.yaml,
// and that neither source nor tags is inherited by a child that doesn't declare its own
// (the shared sightmap library's convention — only the selector prefix cascades).
func TestCollect_Nested(t *testing.T) {
	dir := writeSightmap(t, "nested.yaml")
	p, err := Collect(dir)
	fstesting.Ok(t, err, "Collect")

	byName := indexByName(p.AllComponents())

	modal, ok := byName["Modal"]
	fstesting.Assert(t, ok, "missing Modal")
	fstesting.Equals(t, []string{".modal"}, modal.Selectors, "Modal selectors")
	fstesting.Equals(t, "dialog", modal.Source, "Modal source")
	fstesting.Equals(t, []string{"overlay dialog container"}, modal.Memory, "Modal memory")
	fstesting.Equals(t, []string{"defect"}, modal.Tags, "Modal tags")

	close_, ok := byName["ModalClose"]
	fstesting.Assert(t, ok, "missing ModalClose")
	fstesting.Equals(t, []string{".modal .modal-close"}, close_.Selectors, "ModalClose selectors")
	fstesting.Equals(t, "", close_.Source, "ModalClose does not inherit parent source")
	fstesting.Equals(t, []string(nil), close_.Tags, "ModalClose does not inherit parent tags")

	body, ok := byName["ModalBody"]
	fstesting.Assert(t, ok, "missing ModalBody")
	fstesting.Equals(t, []string{".modal .modal-body"}, body.Selectors, "ModalBody selectors")
	fstesting.Equals(t, []string{"main content area of the dialog"}, body.Memory, "ModalBody memory")
}

// TestCollect_Views verifies that view-scoped components are collected via
// testdata/views.yaml.
func TestCollect_Views(t *testing.T) {
	dir := writeSightmap(t, "views.yaml")
	p, err := Collect(dir)
	fstesting.Ok(t, err, "Collect")

	byName := indexByName(p.AllComponents())
	_, hasCard := byName["DashboardCard"]
	fstesting.Assert(t, hasCard, "missing DashboardCard")
	_, hasHeader := byName["DashboardHeader"]
	fstesting.Assert(t, hasHeader, "missing DashboardHeader")
}

// TestCollect_WireKeepsViewsAndProperties is the point of sending the corpus: the
// upload body carries the view (name, route, its own components) and component
// properties, which the old flat list dropped.
func TestCollect_WireKeepsViewsAndProperties(t *testing.T) {
	p, err := Collect(writeSightmap(t, "views.yaml"))
	fstesting.Ok(t, err, "Collect")

	blob, err := json.Marshal(p)
	fstesting.Ok(t, err, "marshal")
	var wire sm.Corpus
	fstesting.Ok(t, json.Unmarshal(blob, &wire), "unmarshal wire")

	fstesting.Equals(t, 1, len(wire.Views), "views on the wire")
	v := wire.Views[0]
	fstesting.Equals(t, "Dashboard", v.Name, "view name")
	fstesting.Equals(t, "/dashboard", v.Route, "view route")
	card := indexByName(v.Components)["DashboardCard"]
	fstesting.Equals(t, 1, len(card.Properties), "DashboardCard properties")
	fstesting.Equals(t, "title", card.Properties[0].Name, "property name")
	fstesting.Equals(t, "data-title", card.Properties[0].Extract.Path, "property extract path")
}

// TestCollect_NoSightmapDir returns nothing to upload (not an error) when no
// .sightmap/ directory exists.
func TestCollect_NoSightmapDir(t *testing.T) {
	p, err := Collect(t.TempDir())
	fstesting.Ok(t, err, "Collect")
	fstesting.Assert(t, IsEmpty(p), "nothing to upload expected")
}

func TestFindRoot(t *testing.T) {
	dir := t.TempDir()
	child := filepath.Join(dir, "a", "b", "c")
	fstesting.Ok(t, os.MkdirAll(child, 0755), "MkdirAll child")
	fstesting.Ok(t, os.MkdirAll(filepath.Join(dir, ".sightmap"), 0755), "MkdirAll .sightmap")

	t.Setenv("SIGHTMAP_ROOT", "")
	got, err := FindRoot(child, "")
	fstesting.Ok(t, err, "FindRoot")
	fstesting.Equals(t, dir, got, "root")
}

func TestFindRoot_EnvOverride(t *testing.T) {
	t.Setenv("SIGHTMAP_ROOT", "/custom/root")
	got, err := FindRoot("/some/cwd", "")
	fstesting.Ok(t, err, "FindRoot")
	fstesting.Equals(t, "/custom/root", got, "root from env")
}

func TestFindRoot_ConfigRoot(t *testing.T) {
	t.Setenv("SIGHTMAP_ROOT", "")
	got, err := FindRoot("/some/cwd", "/config/root")
	fstesting.Ok(t, err, "FindRoot")
	fstesting.Equals(t, "/config/root", got, "root from config")
}

func TestFindRoot_EnvBeatsConfig(t *testing.T) {
	t.Setenv("SIGHTMAP_ROOT", "/env/root")
	got, err := FindRoot("/some/cwd", "/config/root")
	fstesting.Ok(t, err, "FindRoot")
	fstesting.Equals(t, "/env/root", got, "env var should beat config")
}

func TestFindRoot_NotFound(t *testing.T) {
	t.Setenv("SIGHTMAP_ROOT", "")
	_, err := FindRoot(t.TempDir(), "")
	fstesting.Assert(t, err != nil, "expected error when .sightmap/ not found")
}

// --- helpers ---

// testdataRoot returns a temp dir whose .sightmap/ contains all testdata files.
func testdataRoot(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	sdir := filepath.Join(dir, ".sightmap")
	fstesting.Ok(t, os.MkdirAll(sdir, 0755), "MkdirAll")
	for _, name := range []string{"components.yaml", "nested.yaml", "views.yaml"} {
		data, err := os.ReadFile(filepath.Join("testdata", name))
		fstesting.Ok(t, err, "read testdata/%s", name)
		fstesting.Ok(t, os.WriteFile(filepath.Join(sdir, name), data, 0644), "write %s", name)
	}
	return dir
}

// writeSightmap returns a temp dir whose .sightmap/ contains only the named testdata file.
func writeSightmap(t *testing.T, name string) string {
	t.Helper()
	dir := t.TempDir()
	sdir := filepath.Join(dir, ".sightmap")
	fstesting.Ok(t, os.MkdirAll(sdir, 0755), "MkdirAll")
	data, err := os.ReadFile(filepath.Join("testdata", name))
	fstesting.Ok(t, err, "read testdata/%s", name)
	fstesting.Ok(t, os.WriteFile(filepath.Join(sdir, name), data, 0644), "write %s", name)
	return dir
}

func indexByName(comps []sm.ComponentDef) map[string]sm.ComponentDef {
	m := make(map[string]sm.ComponentDef, len(comps))
	for _, c := range comps {
		m[c.Name] = c
	}
	return m
}
