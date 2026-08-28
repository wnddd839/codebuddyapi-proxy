package models

import "testing"

func TestPublicModelIDStripsCodeBuddyPrefix(t *testing.T) {
	cases := map[string]string{
		"":                          "auto",
		"default":                   "auto",
		"auto":                      "auto",
		"codebuddy/auto":            "auto",
		"codebuddy:auto":            "auto",
		"codebuddy/deepseek-v4-pro": "deepseek-v4-pro",
		"CODEBUDDY/glm-5.3":         "glm-5.3",
		"minimax-m2.5":              "minimax-m2.5",
	}
	for in, want := range cases {
		if got := PublicModelID(in); got != want {
			t.Fatalf("PublicModelID(%q)=%q want %q", in, got, want)
		}
	}
}
