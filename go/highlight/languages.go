package highlight

import (
	"path/filepath"
	"strings"
	"sync"

	"github.com/mooss/bagend/go/fun/eager/mae"
	"github.com/mooss/bagend/go/obj/omap"
	"github.com/mooss/bagend/go/obj/osl"
)

type Language struct {
	Name       string
	Shortnames []string
	Patterns   []string
}

func (l Language) MDName() string {
	if len(l.Shortnames) == 0 {
		return ""
	}
	return l.Shortnames[0]
}

type LanguageList []Language

func (ll LanguageList) Matcher() FileMatcher {
	extension := omap.NewDefault[string](osl.EmptyP[Language])
	wholename := omap.NewDefault[string](osl.EmptyP[Language])
	var globbed []globbedPattern

	for _, lang := range ll {
		for _, pattern := range lang.Patterns {
			switch {
			case !strings.ContainsAny(pattern, `*?[`):
				wholename.Get(pattern).Push(lang)
			case strings.HasPrefix(pattern, "*.") && !strings.ContainsAny(pattern[2:], `*?[.]`):
				// Simple extension like "*.go", matchable via filepath.Ext.
				extension.Get(pattern[1:]).Push(lang)
			default:
				// Anything else ("*.x[bp]m", "Makefile.*", "Kconfig*", "*.tar.gz")
				// needs glob matching against the basename.
				globbed = append(globbed, globbedPattern{pattern: pattern, language: lang})
			}
		}
	}

	deref := func(ptr *osl.Slice[Language]) []Language { return *ptr }
	return FileMatcher{
		extension: mae.MapV(deref, extension.Map),
		wholename: mae.MapV(deref, wholename.Map),
		globbed:   globbed,
	}
}

// globbedPattern is a pattern with wildcards that must be glob-matched against the basename.
type globbedPattern struct {
	pattern  string
	language Language
}

type FileMatcher struct {
	extension map[string][]Language
	wholename map[string][]Language
	globbed   []globbedPattern
}

func (fm FileMatcher) DetectLanguage(path string) Language {
	basename := filepath.Base(path)
	languages, exists := fm.wholename[basename]
	if exists {
		return languages[0]
	}

	ext := filepath.Ext(basename)
	languages, exists = fm.extension[ext]
	if exists {
		return languages[0]
	}

	for _, gp := range fm.globbed {
		if matched, _ := filepath.Match(gp.pattern, basename); matched {
			return gp.language
		}
	}

	return Language{
		Name:       "",
		Shortnames: []string{""},
		Patterns:   []string{""},
	}
}

var LanguageMatcher = sync.OnceValue(KnownLanguages.Matcher)
