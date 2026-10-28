package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/fullstorydev/subtext-cli/internal/sightmap"
)

var sightmapCmd = &cobra.Command{
	Use:   "sightmap",
	Short: "Upload a .sightmap/ corpus",
	Long: `Upload a .sightmap/ corpus to a live connection or a session review so
snapshots and signals come back with semantic component names, view names,
request names, and extracted properties.

The whole corpus is sent (components with their properties, views with routes,
requests, messages, memory), the same wire as 'sightmap export --url'.

Upload URLs are single-use and expire quickly; upload right after you get one.

Live connection ('subtext live tunnel', or the sightmap_upload_url line from
'subtext live connect'):
  URL=$(subtext live tunnel --format json | jq -r .data.sightmapUploadUrl)
  subtext sightmap upload --url "$URL"

Session review (the sightmap_upload_url line in 'subtext review open' output):
  subtext sightmap upload --url "<sightmap_upload_url>"`,
}

var sightmapUploadFlags struct {
	url  string
	root string
}

var sightmapUploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a .sightmap/ corpus to a live connection or session review",
	RunE:  runSightmapUpload,
}

func init() {
	f := sightmapUploadCmd.Flags()
	f.StringVar(&sightmapUploadFlags.url, "url", "", "Single-use sightmap upload URL (from live tunnel, live connect, or review open)")
	f.StringVar(&sightmapUploadFlags.root, "root", "", "Directory containing .sightmap/ (auto-detected if omitted)")
	_ = sightmapUploadCmd.MarkFlagRequired("url")

	sightmapCmd.AddCommand(sightmapUploadCmd)
}

func runSightmapUpload(cmd *cobra.Command, _ []string) error {
	root := sightmapUploadFlags.root
	if root == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getwd: %w", err)
		}
		root, err = sightmap.FindRoot(cwd, globalConfig.SightmapRoot)
		if err != nil {
			return err
		}
	}

	corpus, err := sightmap.Collect(root)
	if err != nil {
		return err
	}
	if sightmap.IsEmpty(corpus) {
		if globalFlags.format == "json" {
			return json.NewEncoder(os.Stdout).Encode(map[string]any{"ok": true, "components": 0, "message": "No sightmap definitions found"})
		}
		fmt.Println("No sightmap definitions found")
		return nil
	}

	n, err := sightmap.Upload(cmd.Context(), sightmapUploadFlags.url, corpus)
	if err != nil {
		return err
	}

	if globalFlags.format == "json" {
		return json.NewEncoder(os.Stdout).Encode(map[string]any{"ok": true, "components": n})
	}
	fmt.Printf("Uploaded %d sightmap component(s).\n", n)
	return nil
}
