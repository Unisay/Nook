package editor

import (
	"slices"
	"strings"
	"testing"
)

func TestStoredControlsRenderAsVisibleText(t *testing.T) {
	const raw = "a\x00\x07\x1b\r\t\n\x7f\u0085\u009f\xffz"
	const visible = "a???? ??¡¡?z"
	doc := func(text string) []string {
		return []string{"# " + text, text, "| name |", "|---|", "| " + text + " |"}
	}
	assertDisplayEquivalent(t, doc(raw), doc(visible))
	for _, cmd := range []bool{false, true} {
		footer := func(text string) string {
			e := New(nil, 80, 10).WithVim()
			e.vimCmdOn, e.vimCmd, e.vimMsg = cmd, text, text
			return e.Footer()
		}
		if footer(raw) != footer(visible) {
			t.Fatalf("unsafe footer, command mode %v: %q", cmd, footer(raw))
		}
	}
}

func TestMarkdownTabsUseOneSpace(t *testing.T) {
	for _, line := range []string{"#\tTitle", ">\tquote", "\t- nested", "\tcode"} {
		doc := []string{line, "end"}
		if line == "\tcode" {
			doc = []string{"```", line, "```", "end"}
		}
		plain := slices.Clone(doc)
		for i, s := range plain {
			plain[i] = strings.ReplaceAll(s, "\t", " ")
		}
		e := New(doc, 40, 15).WithReading(true)
		if got, want := e.View().Content, New(plain, 40, 15).WithReading(true).View().Content; got != want {
			t.Fatalf("%q: got %q, want %q", line, got, want)
		}
	}
}

func assertDisplayEquivalent(t *testing.T, rawDoc, visibleDoc []string) {
	t.Helper()
	for _, width := range []int{7, 40} {
		for line := range rawDoc {
			for _, reading := range []bool{false, true} {
				makeEditor := func(lines []string) Editor {
					e := New(lines, width, 15).SetSelection(0, 1, line, 5)
					return e.WithReading(reading)
				}
				got, want := makeEditor(rawDoc), makeEditor(visibleDoc)
				if got.View().Content != want.View().Content {
					t.Fatalf("line %d, width %d, reading %v: got %q, want %q", line, width, reading, got.View().Content, want.View().Content)
				}
				for y := 0; y < 15; y++ {
					for x := 0; x < width; x++ {
						l, c, ok := got.HitCell(y, x)
						wl, wc, wok := want.HitCell(y, x)
						if l != wl || c != wc || ok != wok {
							t.Fatalf("hit (%d,%d): got (%d,%d,%v), want (%d,%d,%v)", x, y, l, c, ok, wl, wc, wok)
						}
					}
				}
				if !slices.Equal(got.Lines(), rawDoc) {
					t.Fatal("rendering changed stored text")
				}
			}
		}
	}
}

func TestTableControlsRenderAsVisibleText(t *testing.T) {
	for _, tc := range []struct{ raw, visible string }{
		{"\r``` | x |", "?``` | x |"},
		{"# h\n| x |", "# h?| x |"},
	} {
		doc := func(row string) []string { return []string{"| name |", "|---|", row, "end"} }
		assertDisplayEquivalent(t, doc(tc.raw), doc(tc.visible))
	}
}
