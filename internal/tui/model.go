package tui

import (
	"encoding/json"
	"time"

	"github.com/chamdom/omc-agent-tui/internal/store"
	"github.com/chamdom/omc-agent-tui/internal/tui/arena"
	"github.com/chamdom/omc-agent-tui/internal/tui/footer"
	"github.com/chamdom/omc-agent-tui/internal/tui/graph"
	"github.com/chamdom/omc-agent-tui/internal/tui/inspector"
	"github.com/chamdom/omc-agent-tui/internal/tui/timeline"
	"github.com/chamdom/omc-agent-tui/pkg/schema"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const panelCount = 4

// EventMsg is sent when a new event arrives from the pipeline.
type EventMsg schema.CanonicalEvent

// tickMsg is sent periodically to update the UI.
type tickMsg time.Time

// Model is the root Bubbletea model for the TUI.
type Model struct {
	store      *store.Store
	arena      arena.Model
	timeline   timeline.Model
	graph      graph.Model
	inspector  inspector.Model
	footer     footer.Model
	agentTasks map[string]string // agentID -> active taskID

	width  int
	height int

	focused int // 0=arena, 1=timeline, 2=graph, 3=inspector
}

// NewModel creates a new TUI model with the given store.
func NewModel(s *store.Store) Model {
	return Model{
		store:      s,
		arena:      arena.NewModel(),
		timeline:   timeline.NewModel(),
		graph:      graph.NewModel(),
		inspector:  inspector.NewModel(),
		footer:     footer.NewModel(),
		agentTasks: make(map[string]string),
		focused:    0,
	}
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
	)
}

// Update handles messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.focused = (m.focused + 1) % panelCount
		default:
			// Route scroll keys to focused viewport panel
			switch m.focused {
			case 1: // timeline
				m.timeline, cmd = m.timeline.Update(msg)
				cmds = append(cmds, cmd)
			case 3: // inspector
				m.inspector, cmd = m.inspector.Update(msg)
				cmds = append(cmds, cmd)
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateLayout()

	case EventMsg:
		m.addEvent(schema.CanonicalEvent(msg))

	case tickMsg:
		if m.store != nil {
			metrics := m.store.GetMetrics()
			m.footer.SetMetrics(
				metrics.TotalLatency,
				metrics.TotalTokensIn,
				metrics.TotalTokensOut,
				metrics.TotalCostUSD,
			)
		}
		cmds = append(cmds, tickCmd())
	}

	return m, tea.Batch(cmds...)
}

// View renders the entire TUI.
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	// Top: Arena (full width)
	arenaView := m.arena.View()

	// Middle left: Timeline
	timelineView := m.timeline.View()

	// Middle right: Graph + Inspector stacked vertically
	graphView := m.graph.View()
	inspectorView := m.inspector.View()
	rightPanel := lipgloss.JoinVertical(lipgloss.Left, graphView, inspectorView)

	// Middle section: Timeline | Right panel
	middleSection := lipgloss.JoinHorizontal(lipgloss.Top, timelineView, rightPanel)

	// Bottom: Footer
	footerView := m.footer.View()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		arenaView,
		middleSection,
		footerView,
	)
}

// updateLayout recalculates panel sizes based on window dimensions.
func (m *Model) updateLayout() {
	width := m.width
	height := m.height
	if width < 80 {
		width = 80
	}
	if height < 24 {
		height = 24
	}

	footerHeight := 1
	arenaHeight := int(float64(height-footerHeight) * 0.25)
	middleHeight := height - arenaHeight - footerHeight

	leftWidth := int(float64(width) * 0.55)
	rightWidth := width - leftWidth

	graphHeight := middleHeight / 2
	inspectorHeight := middleHeight - graphHeight

	m.arena.SetSize(width, arenaHeight)
	m.timeline.SetSize(leftWidth, middleHeight)
	m.graph.SetSize(rightWidth, graphHeight)
	m.inspector.SetSize(rightWidth, inspectorHeight)
	m.footer.SetSize(width)
}

// addEvent processes a new event, updating store and all sub-models.
func (m *Model) addEvent(event schema.CanonicalEvent) {
	// Store event
	if m.store != nil {
		m.store.AddEvent(event)
	}

	// Update timeline
	m.timeline.AddEvent(event)

	// Update arena
	m.arena.UpdateAgent(event.AgentID, event.Role, event.State)

	// Update footer counters
	m.footer.IncrementEvents()
	if event.Type == schema.TypeError {
		m.footer.IncrementErrors()
	}
	if event.Mode != "" {
		m.footer.SetMode(event.Mode)
	}

	// Update graph for task lifecycle events
	if event.TaskID != "" {
		switch event.Type {
		case schema.TypeTaskSpawn:
			title := extractTaskTitle(event)
			parentTaskID := ""
			if event.ParentAgentID != "" {
				if ptid, ok := m.agentTasks[event.ParentAgentID]; ok {
					parentTaskID = ptid
				}
			}
			if parentTaskID != "" {
				m.graph.AddChildTask(parentTaskID, event.TaskID, event.AgentID, title)
			} else {
				m.graph.AddTask(event.TaskID, event.AgentID, title)
			}
			m.agentTasks[event.AgentID] = event.TaskID

		case schema.TypeTaskDone:
			state := extractTaskDoneState(event)
			m.graph.UpdateTaskState(event.TaskID, state)
			delete(m.agentTasks, event.AgentID)
		}
	}

	// Update inspector with latest event
	m.inspector.SetEvent(&event)

	// Refresh metrics from store
	if m.store != nil {
		metrics := m.store.GetMetrics()
		m.footer.SetMetrics(
			metrics.TotalLatency,
			metrics.TotalTokensIn,
			metrics.TotalTokensOut,
			metrics.TotalCostUSD,
		)
	}
}

// AddEvent is the public API for adding events externally.
func (m *Model) AddEvent(event schema.CanonicalEvent) {
	m.addEvent(event)
}

// extractTaskTitle gets the title from a TaskSpawn payload.
func extractTaskTitle(event schema.CanonicalEvent) string {
	if len(event.Payload) > 0 {
		var payload schema.TaskSpawnPayload
		if err := json.Unmarshal(event.Payload, &payload); err == nil && payload.Title != "" {
			return payload.Title
		}
	}
	return event.TaskID
}

// extractTaskDoneState maps TaskDone result to graph state string.
func extractTaskDoneState(event schema.CanonicalEvent) string {
	if len(event.Payload) > 0 {
		var payload schema.TaskDonePayload
		if err := json.Unmarshal(event.Payload, &payload); err == nil {
			switch payload.Result {
			case "success":
				return "done"
			case "failure":
				return "failed"
			case "cancelled":
				return "cancelled"
			}
		}
	}
	return "done"
}

// tickCmd returns a command that sends a tick message every second.
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
