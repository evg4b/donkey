package donkey

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/evg4b/donkey/internal/store"
	"github.com/goreleaser/fileglob"
)

type donkeyApp struct {
	store *store.Store

	loading bool
	spinner spinner.Model

	form    *huh.Form
	pattern *string
	promt   *string
}

func InitialModel(store *store.Store) tea.Model {
	var promt, pattern string = "", ""

	donkey := donkeyApp{
		store:   store,
		promt:   &promt,
		pattern: &pattern,
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Spinner{
				Frames: []string{
					"▱▱▱",
					"▰▱▱",
					"▱▰▱",
					"▱▱▰",
				},
				FPS: time.Second / 7,
			}),
		),
	}

	donkey.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter files mask").
				Prompt("").
				Value(donkey.pattern).
				Validate(fileglob.ValidPattern),
			huh.NewText().
				Title("Promnt").
				Value(donkey.promt),
		),
	).
		WithWidth(80).
		WithShowHelp(true).
		WithShowErrors(true).
		WithTheme(huh.ThemeCatppuccin())

	return donkey
}

func (m donkeyApp) Init() tea.Cmd {
	return tea.Batch(
		m.form.Init(),
		tea.SetWindowTitle("🫏 donkey"),
	)
}

func (m donkeyApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch event := msg.(type) {
	case finishLoading:
		return m, tea.Batch(
			tea.SetWindowTitle("🫏 donkey"),
			tea.Quit,
		)
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		if event.String() == "q" || event.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if m.form.State == huh.StateCompleted && !m.loading {
		m.loading = true
		cmds = append(
			cmds,
			m.spinner.Tick,
			func() tea.Msg {
				m.store.Generate(*m.promt, *m.pattern)
				return finishLoading{}
			},
		)
	}

	return m, tea.Batch(cmds...)
}

func (m donkeyApp) View() string {
	if m.loading {
		return m.spinner.View() + " 🫏 the donkey does his job..."
	}

	return m.form.View()
}
