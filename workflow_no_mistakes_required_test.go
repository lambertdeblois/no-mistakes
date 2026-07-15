package main

import (
	"os"
	"strings"
	"testing"
)

// TestNoMistakesRequiredWorkflowExemptsReleaseAutomation pins the exemption
// logic so the release pipeline (release-please via GITHUB_TOKEN) and
// dependabot are never silently blocked by the gate.
func TestNoMistakesRequiredWorkflowExemptsReleaseAutomation(t *testing.T) {
	data, err := os.ReadFile(".github/workflows/no-mistakes-required.yml")
	if err != nil {
		t.Fatalf("read workflow: %v", err)
	}
	content := string(data)

	exempt := []string{
		"github-actions[bot]",
		"dependabot[bot]",
		"release-please[bot]",
	}
	for _, login := range exempt {
		needle := "github.event.pull_request.user.login != '" + login + "'"
		if !strings.Contains(content, needle) {
			t.Errorf("workflow must exempt %q via %q", login, needle)
		}
	}
}

// TestNoMistakesRequiredWorkflowChecksSignatureMarker pins the exact marker
// string the check greps for. It must match the literal section heading produced by
// internal/pipeline/steps/prsummary.go when building the Pipeline section.
func TestNoMistakesRequiredWorkflowChecksSignatureMarker(t *testing.T) {
	data, err := os.ReadFile(".github/workflows/no-mistakes-required.yml")
	if err != nil {
		t.Fatalf("read workflow: %v", err)
	}
	content := string(data)

	marker := "## Pipeline"
	if !strings.Contains(content, marker) {
		t.Fatalf("workflow must grep for the prsummary.go pipeline section marker:\n  %s", marker)
	}

	summary, err := os.ReadFile("internal/pipeline/steps/prsummary.go")
	if err != nil {
		t.Fatalf("read prsummary.go: %v", err)
	}
	if !strings.Contains(string(summary), marker) {
		t.Fatalf("prsummary.go no longer writes the expected marker; update both files in sync")
	}
}

// TestNoMistakesRequiredWorkflowReadsPRBodyViaEnv pins the shell-injection-safe
// pattern: the PR body must be piped through an env var, not interpolated
// directly into the shell script body.
func TestNoMistakesRequiredWorkflowReadsPRBodyViaEnv(t *testing.T) {
	data, err := os.ReadFile(".github/workflows/no-mistakes-required.yml")
	if err != nil {
		t.Fatalf("read workflow: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, "PR_BODY: ${{ github.event.pull_request.body }}") {
		t.Errorf("workflow must expose PR body via the PR_BODY env var")
	}
	if strings.Contains(content, "${{ github.event.pull_request.body }}\n          run:") {
		t.Errorf("workflow must not interpolate PR body directly into run: script (injection risk)")
	}
}

// TestNoMistakesRequiredWorkflowTriggersOnRelevantPREvents ensures the check
// re-runs when the PR body is edited so a contributor cannot bypass by opening
// clean then editing the body.
func TestNoMistakesRequiredWorkflowTriggersOnRelevantPREvents(t *testing.T) {
	data, err := os.ReadFile(".github/workflows/no-mistakes-required.yml")
	if err != nil {
		t.Fatalf("read workflow: %v", err)
	}
	content := string(data)

	for _, typ := range []string{"opened", "edited", "synchronize", "reopened"} {
		if !strings.Contains(content, typ) {
			t.Errorf("workflow must trigger on pull_request type %q", typ)
		}
	}
}
