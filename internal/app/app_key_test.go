package app

import (
	"image/color"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/Unisay/Nook/editor"
	"github.com/Unisay/Nook/theme"
	"github.com/Unisay/Nook/themefile"
)

// TestHelpKeyToggle checks that f1 opens the help popup without making
// the application quit.
func TestHelpKeyToggle(t *testing.T) {
	m := Model{width: 80, height: 24}

	m2, _, quit := m.updateKey(tea.KeyPressMsg{Code: tea.KeyF1})

	if quit {
		t.Fatal("quit = true, want false")
	}
	if !m2.helpOpen {
		t.Fatal("helpOpen = false, want true")
	}
}

// TestEditorSpaceFlowsRunes checks that a space sent as KeyRunes
// reaches the editor of the active screen and inserts a space in the text.
// The real case (KeySpace) is in TestEditorSpaceFlows, in keys_test.go.
func TestEditorSpaceFlowsRunes(t *testing.T) {
	m := Model{width: 80, height: 24, screen: screenEditor}
	m.editor = editorScreen{
		ed:     editor.New([]string{"a"}, contentW(80), edHeight(contentH(24))),
		width:  80,
		height: 24,
	}

	m2, _, _ := m.updateKey(tea.KeyPressMsg{Text: " "})

	if got := m2.editor.ed.Lines()[0]; got != "a " {
		t.Fatalf("editor.ed.Lines()[0] = %q, want %q", got, "a ")
	}
}

func TestAutoThemeBackground(t *testing.T) {
	dark, _ := themefile.Builtin("dark")
	light, _ := themefile.Builtin("light")
	t.Cleanup(func() { theme.Use(dark) })
	theme.Use(light)
	lightBanner := theme.BannerAcc.Render("banner")
	theme.Use(dark)
	darkBanner := theme.BannerAcc.Render("banner")
	if darkBanner == lightBanner {
		t.Fatal("test requires distinct banner styles")
	}
	auto := Model{}.WithAutoTheme(true)
	if auto.Init() == nil {
		t.Fatal("missing background request")
	}
	auto.Update(tea.BackgroundColorMsg{Color: color.White})
	if theme.BannerAcc.Render("banner") != lightBanner {
		t.Fatal("white background did not select light theme")
	}
	auto.Update(tea.BackgroundColorMsg{Color: color.Black})
	if theme.BannerAcc.Render("banner") != darkBanner {
		t.Fatal("black background did not select dark theme")
	}
	fixed := Model{}
	if fixed.Init() != nil {
		t.Fatal("fixed theme queried background")
	}
	fixed.Update(tea.BackgroundColorMsg{Color: color.White})
	if theme.BannerAcc.Render("banner") != darkBanner {
		t.Fatal("fixed palette changed")
	}
}
