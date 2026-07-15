package conventional

import (
	"regexp"
	"strings"
)

var titleRe = regexp.MustCompile(`^([a-z]+)(\([^)]+\))?(!)?: (.+)$`)

// jiraKeyRe matches a Jira-style ticket key anywhere in a string, e.g. TBUD-190.
var jiraKeyRe = regexp.MustCompile(`\b([A-Z][A-Z0-9]+-\d+)\b`)

var validTypes = map[string]bool{
	"feat":     true,
	"fix":      true,
	"docs":     true,
	"style":    true,
	"refactor": true,
	"perf":     true,
	"test":     true,
	"build":    true,
	"ci":       true,
	"chore":    true,
	"revert":   true,
}

const ReleaseTypeRule = `- If the change has any user-facing product impact, the type must use feat or fix so release automation can pick it up. Use feat for a new user-visible capability and fix for a user-visible correction or behavior improvement. Use docs, refactor, chore, test, build, or ci only when the change has no user-facing product behavior impact.`

// ExtractJiraKey returns the first Jira ticket key found in branch, e.g.
// "TBUD-190" from "feature/TBUD-190-add-retry". Returns "" if none found.
func ExtractJiraKey(branch string) string {
	m := jiraKeyRe.FindStringSubmatch(branch)
	if len(m) < 2 {
		return ""
	}
	return m[1]
}

// InjectScope rewrites a conventional commit title to use scope as its scope,
// replacing any scope the agent may have chosen. The title must already be in
// conventional commit format; if it is not, it is returned unchanged.
//
// Examples:
//
//	InjectScope("feat: add retry", "TBUD-190")      → "feat(TBUD-190): add retry"
//	InjectScope("feat(pipeline): add retry", "TBUD-190") → "feat(TBUD-190): add retry"
func InjectScope(title, scope string) string {
	if scope == "" {
		return title
	}
	m := titleRe.FindStringSubmatch(strings.TrimSpace(title))
	if len(m) == 0 || !validTypes[m[1]] {
		return title
	}
	typ := m[1]
	bang := m[3]
	description := m[4]
	return typ + "(" + scope + ")" + bang + ": " + description
}

func IsTitle(title string) bool {
	m := titleRe.FindStringSubmatch(strings.TrimSpace(title))
	return len(m) > 0 && validTypes[m[1]]
}

func TightenTitle(title string) string {
	title = strings.TrimSpace(title)
	if title == "" {
		return ""
	}

	m := titleRe.FindStringSubmatch(title)
	if len(m) == 0 || !validTypes[m[1]] {
		return inferType(title) + ": " + title
	}
	return title
}

func inferType(text string) string {
	lower := strings.ToLower(strings.TrimSpace(text))
	switch {
	case hasDocumentationLanguage(lower):
		return "docs"
	case hasProductImpactLanguage(lower) || isFeatureLanguage(lower) || isFixLanguage(lower):
		return inferReleaseType(lower)
	default:
		return "chore"
	}
}

func inferReleaseType(text string) string {
	if isFeatureLanguage(text) {
		return "feat"
	}
	return "fix"
}

func isFixLanguage(text string) bool {
	lower := strings.ToLower(strings.TrimSpace(text))
	fixPrefixes := []string{
		"fix ", "fixes ", "fixed ", "resolve ", "resolves ", "resolved ",
		"correct ", "corrects ", "corrected ", "repair ", "repairs ", "repaired ",
	}
	for _, prefix := range fixPrefixes {
		if strings.HasPrefix(lower, prefix) {
			return true
		}
	}
	return false
}

func isFeatureLanguage(text string) bool {
	lower := strings.ToLower(strings.TrimSpace(text))
	featurePrefixes := []string{
		"add ", "adds ", "added ", "introduce ", "introduces ", "introduced ",
		"create ", "creates ", "created ", "implement ", "implements ", "implemented ",
		"support ", "supports ", "supported ", "enable ", "enables ", "enabled ",
		"allow ", "allows ", "allowed ",
	}
	for _, prefix := range featurePrefixes {
		if strings.HasPrefix(lower, prefix) {
			return true
		}
	}
	return strings.Contains(lower, " new ") || strings.HasPrefix(lower, "new ")
}

func hasProductImpactLanguage(text string) bool {
	lower := strings.ToLower(text)
	terms := []string{
		"user-facing", "user visible", "user-visible", "user experience", " ux", "ux ",
		" ui", "ui ", "cli", "command", "output", "behavior", "workflow",
		"prompt", "flag", "error message",
	}
	for _, term := range terms {
		if strings.Contains(lower, term) {
			return true
		}
	}
	return false
}

func hasDocumentationLanguage(text string) bool {
	lower := strings.ToLower(text)
	return strings.Contains(lower, "readme") || strings.Contains(lower, "documentation") || strings.Contains(lower, "docs")
}
