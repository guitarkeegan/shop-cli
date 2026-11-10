package main

import (
	feat "github.com/guitarkeegan/shop-cli/internal/feature-service"
)

type SectionCreateOptions struct {
	Name         string
	Settings     string
	Blocks       []string
	useTemplates bool
	Templates    *feat.FeatureService
}

type SnippetCreateOptions struct {
	Name        string
	DocParams   string
	Template    bool
	Description string
	Example     string
}

type SectionCreateOption func(*SectionCreateOptions)

func WithFeatureService(featS *feat.FeatureService) SectionCreateOption {
	return func(sco *SectionCreateOptions) {
		sco.Templates = featS
	}
}

func NewSectionCreateOptions(opts ...SectionCreateOption) *SectionCreateOptions {
	sco := &SectionCreateOptions{}
	for _, opt := range opts {
		opt(sco)
	}
	return sco
}

// Create an empty SnippetCreateOptions, this may be used
// to load a config in the future.
func NewSnippetCreateOptions() *SnippetCreateOptions {
	return &SnippetCreateOptions{}
}
