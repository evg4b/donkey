package donkey

import (
	"fmt"

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
	suffix  *string

	task string
}

// A command that waits for the activity on a channel.
func waitForActivity(sub <-chan store.Event) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

func InitialModel(store *store.Store, pattern, promt, suffix string) tea.Model {
	donkey := donkeyApp{
		store:   store,
		promt:   &promt,
		pattern: &pattern,
		suffix:  &suffix,
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Dot),
		),
	}

	donkey.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter files mask").
				Description("Use glob pattern to specify files").
				Prompt("").
				Value(donkey.pattern).
				Validate(fileglob.ValidPattern),
			huh.NewText().
				Title("Promnt").
				Description("Promt for the model").
				Value(donkey.promt),
			huh.NewInput().
				Title("Suffix").
				Description("Suffix for the output files").
				Prompt("").
				Value(donkey.suffix),
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
	case store.Event:
		switch event.Type {
		case store.FileProcessing:
			m.task = event.InputFileName
		case store.FileProcessed:
			if event.HasDifferentOutput() {
				cmds = append(cmds, tea.Printf("📄 Processed %s -> %s", event.InputFileName, event.OutputFileName))
			} else {
				cmds = append(cmds, tea.Printf("📄 Processed %s", event.InputFileName))
			}
		case store.MemoryCleared:
			cmds = append(cmds, tea.Printf("🧠 Memory cleared"))
		}
		cmds = append(cmds, waitForActivity(m.store.Events()))
		return m, tea.Batch(cmds...)
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
				m.store.Generate(*m.promt, *m.pattern, *m.suffix)
				return finishLoading{}
			},
			waitForActivity(m.store.Events()),
		)
	}

	return m, tea.Batch(cmds...)
}

func (m donkeyApp) View() string {
	if m.loading {
		if m.task == "" {
			return m.spinner.View() + " 🫏 does his job..."
		}

		return fmt.Sprintf(
			"%s 🫏 processing %s file...",
			m.spinner.View(),
			m.task,
		)
	}

	return m.form.View()
}
