package conventional

import "testing"

func TestTightenTitleKeepsReleaseTypes(t *testing.T) {
	t.Parallel()

	tests := []string{
		"feat(cli): add onboarding wizard",
		"fix: improve command output",
		"fix(api)!: require auth token",
	}

	for _, tc := range tests {
		t.Run(tc, func(t *testing.T) {
			t.Parallel()
			if got := TightenTitle(tc); got != tc {
				t.Fatalf("TightenTitle(%q) = %q", tc, got)
			}
		})
	}
}

func TestTightenTitleKeepsConventionalNonReleaseTypes(t *testing.T) {
	t.Parallel()

	tests := []string{
		"refactor: improve CLI output",
		"docs: add user-facing export command",
		"chore(cli)!: improve UI behavior",
	}

	for _, tc := range tests {
		t.Run(tc, func(t *testing.T) {
			t.Parallel()
			if got := TightenTitle(tc); got != tc {
				t.Fatalf("TightenTitle(%q) = %q", tc, got)
			}
		})
	}
}

func TestTightenTitleKeepsNonProductImpactTypes(t *testing.T) {
	t.Parallel()

	tests := []string{
		"docs: update README",
		"docs: update CLI command documentation",
		"refactor: simplify internal retry loop",
		"test: cover config parsing",
	}

	for _, tc := range tests {
		t.Run(tc, func(t *testing.T) {
			t.Parallel()
			if got := TightenTitle(tc); got != tc {
				t.Fatalf("TightenTitle(%q) = %q", tc, got)
			}
		})
	}
}

func TestTightenTitlePrefixesNonConventionalTitles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		title string
		want  string
	}{
		{name: "new feature", title: "add export command", want: "feat: add export command"},
		{name: "direct fix verb", title: "fix login redirect", want: "fix: fix login redirect"},
		{name: "direct correction verb", title: "correct cache invalidation", want: "fix: correct cache invalidation"},
		{name: "user-facing fix", title: "Improve pipeline header UX", want: "fix: Improve pipeline header UX"},
		{name: "documentation", title: "update README", want: "docs: update README"},
		{name: "generic internal", title: "tidy retry helper", want: "chore: tidy retry helper"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := TightenTitle(tc.title); got != tc.want {
				t.Fatalf("TightenTitle(%q) = %q, want %q", tc.title, got, tc.want)
			}
		})
	}
}

func TestExtractJiraKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		branch string
		want   string
	}{
		{branch: "TBUD-190", want: "TBUD-190"},
		{branch: "feature/TBUD-190", want: "TBUD-190"},
		{branch: "TBUD-190-add-retry", want: "TBUD-190"},
		{branch: "feature/TBUD-190-add-retry-logic", want: "TBUD-190"},
		{branch: "main", want: ""},
		{branch: "fix-some-bug", want: ""},
		{branch: "", want: ""},
	}

	for _, tc := range tests {
		t.Run(tc.branch, func(t *testing.T) {
			t.Parallel()
			if got := ExtractJiraKey(tc.branch); got != tc.want {
				t.Fatalf("ExtractJiraKey(%q) = %q, want %q", tc.branch, got, tc.want)
			}
		})
	}
}

func TestInjectScope(t *testing.T) {
	t.Parallel()

	tests := []struct {
		title string
		scope string
		want  string
	}{
		{title: "feat: add retry logic", scope: "TBUD-190", want: "feat(TBUD-190): add retry logic"},
		{title: "fix(pipeline): correct timeout", scope: "TBUD-190", want: "fix(TBUD-190): correct timeout"},
		{title: "chore!: drop legacy api", scope: "TBUD-190", want: "chore(TBUD-190)!: drop legacy api"},
		{title: "feat(cli): add flag", scope: "", want: "feat(cli): add flag"},
		{title: "not a conventional title", scope: "TBUD-190", want: "not a conventional title"},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()
			if got := InjectScope(tc.title, tc.scope); got != tc.want {
				t.Fatalf("InjectScope(%q, %q) = %q, want %q", tc.title, tc.scope, got, tc.want)
			}
		})
	}
}

func TestIsTitle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		title string
		want  bool
	}{
		{title: "feat: add export", want: true},
		{title: "fix(cli)!: change output", want: true},
		{title: "add export", want: false},
		{title: "Feat: add export", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()
			if got := IsTitle(tc.title); got != tc.want {
				t.Fatalf("IsTitle(%q) = %v, want %v", tc.title, got, tc.want)
			}
		})
	}
}
