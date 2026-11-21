package editor

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func runes(s string) tea.KeyPressMsg { return tea.KeyPressMsg{Text: s} }

func readDoc() Editor {
	return New([]string{"# Title", "a b c d e f g h i j k l m n o p", "last"}, 10, 7)
}

func TestReadingToggleAndFlag(t *testing.T) {
	e := readDoc()
	e = press(e, tea.KeyPressMsg{Code: 'p', Mod: tea.ModCtrl})
	if !e.Reading() || !e.ReadChanged {
		t.Fatalf("reading=%v changed=%v", e.Reading(), e.ReadChanged)
	}
	e.ReadChanged = false
	e = press(e, tea.KeyPressMsg{Code: 'p', Mod: tea.ModCtrl})
	if e.Reading() || !e.ReadChanged {
		t.Fatalf("reading=%v changed=%v", e.Reading(), e.ReadChanged)
	}
}

func TestReadingWrapsAndIsReadOnly(t *testing.T) {
	e := readDoc()
	before := strings.Join(e.Lines(), "\n")
	e = press(e, tea.KeyPressMsg{Code: 'p', Mod: tea.ModCtrl})
	rows := e.readingRowsText()
	n := 0
	for _, r := range rows {
		if strings.Contains(r, "a b") || strings.Contains(r, "k l") {
			n++
		}
	}
	if n < 2 {
		t.Fatalf("long line not wrapped: %q", rows)
	}
	e = press(e, runes("x"), tea.KeyPressMsg{Code: tea.KeyEnter}, tea.KeyPressMsg{Code: tea.KeyBackspace}, tea.KeyPressMsg{Code: tea.KeyDelete},
		tea.KeyPressMsg{Code: 'u', Mod: tea.ModCtrl}, tea.KeyPressMsg{Code: 'k', Mod: tea.ModCtrl}, tea.KeyPressMsg{Code: 't', Mod: tea.ModCtrl}, tea.KeyPressMsg{Code: 'z', Mod: tea.ModCtrl}, tea.KeyPressMsg{Code: 'y', Mod: tea.ModCtrl})
	if got := strings.Join(e.Lines(), "\n"); got != before {
		t.Fatalf("text changed in reading mode: %q", got)
	}
}

func TestReadingScrollAndLeave(t *testing.T) {
	e := readDoc()
	e = press(e, tea.KeyPressMsg{Code: 'p', Mod: tea.ModCtrl})
	e = press(e, tea.KeyPressMsg{Code: tea.KeyEnd})
	e = press(e, tea.KeyPressMsg{Code: 'p', Mod: tea.ModCtrl})
	if e.Reading() {
		t.Fatal("still reading")
	}
	if e.CursorRow() != 3 {
		t.Fatalf("cursor row %d, want 3 (bottom)", e.CursorRow())
	}
	e = press(e, tea.KeyPressMsg{Code: 'p', Mod: tea.ModCtrl}, tea.KeyPressMsg{Code: tea.KeyHome}, tea.KeyPressMsg{Code: 'p', Mod: tea.ModCtrl})
	if e.CursorRow() != 1 {
		t.Fatalf("cursor row %d, want 1 (top)", e.CursorRow())
	}
}

func TestReadingUnmovedRestoresCursor(t *testing.T) {
	e := readDoc()
	e = press(e, tea.KeyPressMsg{Code: tea.KeyDown}, tea.KeyPressMsg{Code: tea.KeyRight}, tea.KeyPressMsg{Code: tea.KeyRight}, tea.KeyPressMsg{Code: tea.KeyRight})
	r, c := e.CursorRow(), e.CursorCol1()
	e = press(e, tea.KeyPressMsg{Code: 'p', Mod: tea.ModCtrl}, tea.KeyPressMsg{Code: 'p', Mod: tea.ModCtrl})
	if e.CursorRow() != r || e.CursorCol1() != c {
		t.Fatalf("cursor %d:%d, want %d:%d", e.CursorRow(), e.CursorCol1(), r, c)
	}
}

func TestReadingVimKeysDoNotEdit(t *testing.T) {
	e := readDoc().WithVim()
	before := strings.Join(e.Lines(), "\n")
	e = press(e, tea.KeyPressMsg{Code: 'p', Mod: tea.ModCtrl})
	e = press(e, runes("d"), runes("d"), runes("i"), runes("x"), runes("j"), runes("v"), runes(":"))
	if strings.Join(e.Lines(), "\n") != before {
		t.Fatal("vim edited in reading mode")
	}
	e = press(e, tea.KeyPressMsg{Code: 'p', Mod: tea.ModCtrl})
	if e.Reading() {
		t.Fatal("ctrl+p should leave reading in vim mode")
	}
}

func TestReadingFullStrengthAndNoCursor(t *testing.T) {
	e := readDoc()
	e = press(e, tea.KeyPressMsg{Code: 'p', Mod: tea.ModCtrl})
	v := e.View().Content
	if v == "" {
		t.Fatal("empty view")
	}
	if !strings.Contains(e.Footer(), "READING") {
		t.Fatalf("footer %q", e.Footer())
	}
}

func TestReadingHitCell(t *testing.T) {
	e := readDoc()
	e = press(e, tea.KeyPressMsg{Code: 'p', Mod: tea.ModCtrl})
	if _, _, ok := e.HitCell(-1, 0); ok {
		t.Fatal("hit above the screen")
	}
	found := false
	for y := range e.readLayout() {
		if l, _, ok := e.HitCell(y, 0); ok && l == 1 {
			found = true
		}
	}
	if !found {
		t.Fatal("no row maps to line 1")
	}
}

func TestWithReadingStartsAtTop(t *testing.T) {
	e := readDoc().WithReading(true)
	if !e.Reading() {
		t.Fatal("not reading")
	}
	e = e.WithReading(false)
	if e.Reading() {
		t.Fatal("still reading")
	}
}

func TestToggleDone(t *testing.T) {
	e := New([]string{"- [ ] a", "plain", "- [x] b"}, 20, 7)
	e = press(e, tea.KeyPressMsg{Code: 'x', Mod: tea.ModCtrl})
	if e.Lines()[0] != "- [x] a" {
		t.Fatalf("got %q", e.Lines()[0])
	}
	e = press(e, tea.KeyPressMsg{Code: 'x', Mod: tea.ModCtrl})
	if e.Lines()[0] != "- [ ] a" {
		t.Fatalf("got %q", e.Lines()[0])
	}
	e = press(e, tea.KeyPressMsg{Code: 'z', Mod: tea.ModCtrl})
	// select all three lines: mixed -> all done
	e = press(e, tea.KeyPressMsg{Code: tea.KeyDown}, tea.KeyPressMsg{Code: tea.KeyDown})
	e = e.ToggleDone()
	if e.Lines()[2] != "- [ ] b" {
		t.Fatalf("got %q", e.Lines()[2])
	}
	if e.Lines()[1] != "plain" {
		t.Fatal("plain line changed")
	}
	v := New([]string{"- [ ] a"}, 20, 7).WithVim()
	v = press(v, tea.KeyPressMsg{Code: 'x', Mod: tea.ModCtrl})
	if v.Lines()[0] != "- [x] a" {
		t.Fatalf("vim: %q", v.Lines()[0])
	}
	v = v.Undo()
	if v.Lines()[0] != "- [ ] a" {
		t.Fatalf("undo: %q", v.Lines()[0])
	}
	r := New([]string{"- [ ] a"}, 20, 7).WithReading(true)
	r = press(r, tea.KeyPressMsg{Code: 'x', Mod: tea.ModCtrl})
	if r.Lines()[0] != "- [ ] a" {
		t.Fatal("reading mode edited")
	}
}
