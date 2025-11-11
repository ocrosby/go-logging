package logging

import "regexp"

type Redactor interface {
	Redact(input string) string
}

type RegexRedactor struct {
	pattern     *regexp.Regexp
	replacement string
}

func NewRegexRedactor(pattern *regexp.Regexp, replacement string) *RegexRedactor {
	return &RegexRedactor{
		pattern:     pattern,
		replacement: replacement,
	}
}

func (r *RegexRedactor) Redact(input string) string {
	return r.pattern.ReplaceAllString(input, r.replacement)
}

type RedactorChain struct {
	redactors []Redactor
}

func NewRedactorChain(patterns ...*regexp.Regexp) *RedactorChain {
	redactors := make([]Redactor, len(patterns))
	for i, pattern := range patterns {
		redactors[i] = NewRegexRedactor(pattern, `$1$2...<REDACTED>$7`)
	}
	return &RedactorChain{
		redactors: redactors,
	}
}

func (rc *RedactorChain) AddRedactor(redactor Redactor) {
	rc.redactors = append(rc.redactors, redactor)
}

func (rc *RedactorChain) Redact(input string) string {
	result := input
	for _, redactor := range rc.redactors {
		result = redactor.Redact(result)
	}
	return result
}

var DefaultAPIKeyPattern = regexp.MustCompile(`(?i)([?&])((apiKey)(=)([a-z0-9]{7}))([^&]+)`)

func RedactAPIKeys(input string) string {
	redactor := NewRegexRedactor(DefaultAPIKeyPattern, `$1$2...<REDACTED>$7`)
	return redactor.Redact(input)
}
