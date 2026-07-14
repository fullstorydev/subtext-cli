package cli

import (
	"github.com/spf13/cobra"
)

// namespaceRunE is the shared RunE for namespace commands. It prepends the
// namespace prefix to the first arg so that "subtext live act-click" resolves
// to the "live-act-click" tool. DisableFlagParsing must be set on the parent
// command, so --help is intercepted here before delegating.
func namespaceRunE(cmd *cobra.Command, args []string) error {
	// DisableFlagParsing passes any leading global flag through as args[0];
	// pull them out before building the tool name or it gets glued in.
	args = extractGlobalFlags(args)
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		return cmd.Help()
	}
	toolName := cmd.Use + "-" + args[0]
	return runCall(cmd, append([]string{toolName}, args[1:]...))
}

var liveCmd = &cobra.Command{
	Use:                "live",
	Short:              "Drive a hosted browser via live MCP tools",
	Long:               `Drive a hosted Subtext browser session: open URLs, observe a viewer, take screenshots, run JavaScript in the page, and hand off control.`,
	DisableFlagParsing: true,
	RunE:               namespaceRunE,
}

var commentCmd = &cobra.Command{
	Use:                "comment",
	Short:              "Create and retrieve session comments",
	Long:               `Create and retrieve comments attached to Fullstory sessions. Comments anchor observations to specific moments in a recording.`,
	DisableFlagParsing: true,
	RunE:               namespaceRunE,
}

var docCmd = &cobra.Command{
	Use:                "doc",
	Short:              "Create and manage proof documents",
	Long:               `Create proof documents that embed screenshots, session viewers, and code. Documents are the primary deliverable for sharing agent findings with reviewers.`,
	DisableFlagParsing: true,
	RunE:               namespaceRunE,
}

// tunnelCmd is declared in tunnel.go (hand-written subcommands).

var artifactCmd = &cobra.Command{
	Use:                "artifact",
	Short:              "Upload and retrieve artifacts",
	Long:               `Upload files (screenshots, HAR traces, logs) and retrieve them by ID. Artifacts persist across sessions and can be linked from proof documents.`,
	DisableFlagParsing: true,
	RunE:               namespaceRunE,
}

var privacyCmd = &cobra.Command{
	Use:   "privacy",
	Short: "Detect PII and manage element-block privacy rules",
	Long: `Propose CSS selectors for PII elements in a session, persist them as element-block rules, and manage their lifecycle (list, delete, promote).

Rules are created in preview scope (PREVIEW_SESSIONS_ONLY) and must be explicitly promoted to apply to all sessions. Only mask and exclude rules are supported — unmask rules cannot be created or deleted via this tool.`,
	DisableFlagParsing: true,
	RunE:               namespaceRunE,
}

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Deep-review a recorded Fullstory session",
	Long: `Open a completed Fullstory session for deep agent review: read the map, zoom into the signal stream, and snapshot the screen to understand what happened.

open returns a map — signal counts by kind/tag and page flow — so read that before deciding what to zoom into. zoom takes a resolution map ({scope|kind|tag: grain}, grains digest/standard/machine/detail, finest-wins) to progressively disclose just the signals a hypothesis needs.

Primary use cases:
  - Verify another agent's proof work (BEFORE/AFTER chapter markers as the spine)
  - Diagnose a bug or regression from a captured session
  - Produce a structured summary of what an agent or user did

Sessions are identified by trace_id, session URL, device+session IDs, email address, or user UID. Always close when done to release resources.`,
	DisableFlagParsing: true,
	RunE:               namespaceRunE,
}

func init() {
	liveCmd.SetHelpFunc(namespaceHelpFunc("live-"))
	commentCmd.SetHelpFunc(namespaceHelpFunc("comment-"))
	docCmd.SetHelpFunc(namespaceHelpFunc("doc-"))
	artifactCmd.SetHelpFunc(namespaceHelpFunc("artifact-"))
	privacyCmd.SetHelpFunc(namespaceHelpFunc("privacy-"))
	reviewCmd.SetHelpFunc(namespaceHelpFunc("review-"))
}
