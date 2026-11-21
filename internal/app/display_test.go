package app

import (
	"testing"

	"github.com/Unisay/Nook/internal/store"
)

func TestControlCharactersInAppViews(t *testing.T) {
	const raw = "a\x07\x1b\x7f\u0085b"
	const visible = "a???¡b"
	if toastLine("opening https://example.org/"+raw, 80) != toastLine("opening https://example.org/"+visible, 80) {
		t.Fatal("link toast differs from visible equivalent")
	}
	m, _, _ := openVimNote(t, "")
	m.editor.note.Title = raw + "\nend"
	if got := m.View().WindowTitle; got != visible+"?end" {
		t.Fatalf("window title = %q", got)
	}
	list := func(name string) string {
		s := newListScreen(nil, 80, 24)
		s.path = []store.Folder{{Name: name}}
		s.folders = []store.Folder{{Name: name}}
		s.notes = []store.Note{{Title: name}}
		return s.View()
	}
	if list(raw+"\rend") != list(visible+"?end") {
		t.Fatal("list differs from visible equivalent")
	}
}
