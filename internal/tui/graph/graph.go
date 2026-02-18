package graph

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// TaskNode represents a task in the dependency graph.
type TaskNode struct {
	TaskID   string
	AgentID  string
	Title    string
	State    string // "active"|"done"|"failed"|"cancelled"
	Children []string
}

// Model represents the Graph panel state.
type Model struct {
	tasks  map[string]*TaskNode
	roots  []string // root task IDs (tasks without parents)
	width  int
	height int
}

// NewModel creates a new Graph model.
func NewModel() Model {
	return Model{
		tasks: make(map[string]*TaskNode),
		roots: []string{},
	}
}

// AddTask adds a root task to the graph.
func (m *Model) AddTask(taskID, agentID, title string) {
	if _, exists := m.tasks[taskID]; !exists {
		m.tasks[taskID] = &TaskNode{
			TaskID:   taskID,
			AgentID:  agentID,
			Title:    title,
			State:    "active",
			Children: []string{},
		}
		m.roots = append(m.roots, taskID)
	}
}

// AddChildTask adds a child task under a parent.
func (m *Model) AddChildTask(parentID, childID, agentID, title string) {
	// Create child node if it doesn't exist
	if _, exists := m.tasks[childID]; !exists {
		m.tasks[childID] = &TaskNode{
			TaskID:   childID,
			AgentID:  agentID,
			Title:    title,
			State:    "active",
			Children: []string{},
		}
	}

	// Add to parent's children
	if parent, exists := m.tasks[parentID]; exists {
		// Check if child already in list
		for _, c := range parent.Children {
			if c == childID {
				return
			}
		}
		parent.Children = append(parent.Children, childID)
	}

	// Remove from roots if it was there
	for i, root := range m.roots {
		if root == childID {
			m.roots = append(m.roots[:i], m.roots[i+1:]...)
			break
		}
	}
}

// UpdateTaskState updates the state of a task.
func (m *Model) UpdateTaskState(taskID string, state string) {
	if task, exists := m.tasks[taskID]; exists {
		task.State = state
	}
}

// SetSize updates the panel dimensions.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// View renders the Graph panel.
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	content := m.renderTree()

	style := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	if content == "" {
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Align(lipgloss.Center, lipgloss.Center)
		return style.Render(emptyStyle.Render("[No tasks]"))
	}

	return style.Render(content)
}

// renderTree renders the task tree structure.
func (m Model) renderTree() string {
	if len(m.roots) == 0 {
		return ""
	}

	var lines []string
	for i, rootID := range m.roots {
		isLast := i == len(m.roots)-1
		lines = append(lines, m.renderNode(rootID, "", isLast, true)...)
	}

	return strings.Join(lines, "\n")
}

// renderNode renders a single task node and its children recursively.
func (m Model) renderNode(taskID string, prefix string, isLast bool, isRoot bool) []string {
	task, exists := m.tasks[taskID]
	if !exists {
		return []string{}
	}

	var lines []string

	// Render current node
	nodePrefix := prefix
	if !isRoot {
		if isLast {
			nodePrefix += "└── "
		} else {
			nodePrefix += "├── "
		}
	}

	// State color mapping
	stateStyle := m.getStateStyle(task.State)
	nodeLine := fmt.Sprintf("%s%s [%s] %s",
		nodePrefix,
		task.TaskID,
		task.AgentID,
		task.State,
	)
	lines = append(lines, stateStyle.Render(nodeLine))

	// Render children
	childPrefix := prefix
	if !isRoot {
		if isLast {
			childPrefix += "    "
		} else {
			childPrefix += "│   "
		}
	}

	for i, childID := range task.Children {
		isChildLast := i == len(task.Children)-1
		lines = append(lines, m.renderNode(childID, childPrefix, isChildLast, false)...)
	}

	return lines
}

// getStateStyle returns the lipgloss style for a given state.
func (m Model) getStateStyle(state string) lipgloss.Style {
	switch state {
	case "active":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#4FC3F7"))
	case "done":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#66BB6A"))
	case "failed":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#EF5350"))
	case "cancelled":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#9E9E9E"))
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	}
}
