package types

import "strings"

type LocalizedText map[string]string

func (t LocalizedText) String() string {
	originalLocale, hasOriginal := t["original"]
	if text, originalPresent := t[originalLocale]; hasOriginal && originalPresent {
		return strings.TrimSuffix(text, "\n")
	}

	for _, text := range t {
		return strings.TrimSuffix(text, "\n")
	}

	return ""
}
