package analysis

import (
	"fmt"
	"strings"

	"github.com/azr4e1/educationalsp/lsp"
)

type State struct {
	// Map of file names to contents
	Documents map[string]string
}

func NewState() State {
	return State{Documents: map[string]string{}}
}

func getDiagnostics(text string) []lsp.Diagnostic {
	diagnostics := []lsp.Diagnostic{}
	for row, line := range strings.Split(text, "\n") {
		idx := strings.Index(line, "VS Code")
		if idx >= 0 {
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range:    LineRange(row, idx, idx+len("VS Code")),
				Source:   "Common Sense",
				Message:  "Fix to a better editor :(",
				Severity: 1,
			})
		}
		idx = strings.Index(line, "Helix")
		if idx >= 0 {
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range:    LineRange(row, idx, idx+len("Helix")),
				Source:   "Common Sense",
				Message:  "Great Choice :)",
				Severity: 4,
			})
		}
	}
	return diagnostics
}

func (s *State) OpenDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text
	return getDiagnostics(text)
}

func (s *State) Update(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text
	return getDiagnostics(text)
}

func (s *State) Hover(id int, uri string, position lsp.Position) lsp.HoverResponse {
	document := s.Documents[uri]

	return lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.HoverResult{
			Contents: fmt.Sprintf("File: %s, Characters: %d", uri, len(document)),
		},
	}
}

func (s *State) Definition(id int, uri string, position lsp.Position) lsp.DefinitionResponse {
	return lsp.DefinitionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.Location{
			URI: uri,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      position.Line - 1,
					Character: 0,
				},
				End: lsp.Position{
					Line:      position.Line - 1,
					Character: 1,
				},
			},
		},
	}
}

func (s *State) CodeAction(id int, uri string) lsp.CodeActionResponse {
	text := s.Documents[uri]
	actions := []lsp.CodeAction{}
	for row, line := range strings.Split(text, "\n") {
		idx := strings.Index(line, "VS Code")
		if idx >= 0 {
			replaceChange := map[string][]lsp.TextEdit{}
			replaceChange[uri] = []lsp.TextEdit{
				{
					Range:   LineRange(row, idx, idx+len("VS Code")),
					NewText: "Helix",
				},
			}
			actions = append(actions, lsp.CodeAction{
				Title: "Replace VS C*de with a superior editor",
				Edit:  &lsp.WorkspaceEdit{Changes: replaceChange},
			})

			censorChange := map[string][]lsp.TextEdit{}
			censorChange[uri] = []lsp.TextEdit{
				{
					Range:   LineRange(row, idx, idx+len("VS Code")),
					NewText: "VS C*de",
				},
			}
			actions = append(actions, lsp.CodeAction{
				Title: "Censor VS C*de",
				Edit:  &lsp.WorkspaceEdit{Changes: censorChange},
			})
		}
	}
	response := lsp.CodeActionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: actions,
	}
	return response
}

func (s *State) Completion(id int, uri string) lsp.CompletionResponse {
	document := s.Documents[uri]
	completionWords := strings.Fields(document)
	candidates := []lsp.CompletionItem{}
	wordSet := make(map[string]struct{})
	for _, word := range completionWords {
		if _, ok := wordSet[word]; ok {
			continue
		}
		wordSet[word] = struct{}{}
		candidates = append(candidates, lsp.CompletionItem{
			Label:  word,
			Detail: fmt.Sprintf("\"%s\" is a word in this document", word),
		})
	}

	return lsp.CompletionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: candidates,
	}
}

// func (s *State) Diagnostics() lsp.PublishDiagnosticsNotification {
// 	params :=
// 	for uri, document := range s {
// 		text := s.Documents[uri]
// 		diagnostics := []lsp.Diagnostic{}
// 		for row, line := range strings.Split(text, "\n") {
// 			idx := strings.Index(line, "VS Code")
// 			if idx >= 0
// 	}
// }

func LineRange(row, start, end int) lsp.Range {
	return lsp.Range{
		Start: lsp.Position{
			Line:      uint(row),
			Character: uint(start),
		},
		End: lsp.Position{
			Line:      uint(row),
			Character: uint(end),
		},
	}
}

func GetLastEl[E any](s []E) (E, bool) {
	if len(s) == 0 {
		var zero E
		return zero, false
	}
	return s[len(s)-1], true
}
