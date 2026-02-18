package graph

import (
	"strings"
	"testing"
)

func TestNewModel(t *testing.T) {
	m := NewModel()
	if m.tasks == nil {
		t.Error("tasks map should be initialized")
	}
	if m.roots == nil {
		t.Error("roots slice should be initialized")
	}
	if len(m.roots) != 0 {
		t.Errorf("expected 0 roots, got %d", len(m.roots))
	}
}

func TestAddTask(t *testing.T) {
	m := NewModel()
	m.AddTask("task-001", "executor", "Test task")

	if len(m.roots) != 1 {
		t.Errorf("expected 1 root, got %d", len(m.roots))
	}
	if m.roots[0] != "task-001" {
		t.Errorf("expected root task-001, got %s", m.roots[0])
	}

	task := m.tasks["task-001"]
	if task == nil {
		t.Fatal("task should exist")
	}
	if task.TaskID != "task-001" {
		t.Errorf("expected TaskID task-001, got %s", task.TaskID)
	}
	if task.AgentID != "executor" {
		t.Errorf("expected AgentID executor, got %s", task.AgentID)
	}
	if task.State != "active" {
		t.Errorf("expected State active, got %s", task.State)
	}
}

func TestAddChildTask(t *testing.T) {
	m := NewModel()
	m.AddTask("task-001", "executor", "Parent task")
	m.AddChildTask("task-001", "task-002", "reviewer", "Child task")

	parent := m.tasks["task-001"]
	if len(parent.Children) != 1 {
		t.Errorf("expected 1 child, got %d", len(parent.Children))
	}
	if parent.Children[0] != "task-002" {
		t.Errorf("expected child task-002, got %s", parent.Children[0])
	}

	// task-002 should not be in roots
	for _, root := range m.roots {
		if root == "task-002" {
			t.Error("child task should not be in roots")
		}
	}

	child := m.tasks["task-002"]
	if child == nil {
		t.Fatal("child task should exist")
	}
	if child.TaskID != "task-002" {
		t.Errorf("expected TaskID task-002, got %s", child.TaskID)
	}
}

func TestUpdateTaskState(t *testing.T) {
	m := NewModel()
	m.AddTask("task-001", "executor", "Test task")
	m.UpdateTaskState("task-001", "done")

	task := m.tasks["task-001"]
	if task.State != "done" {
		t.Errorf("expected State done, got %s", task.State)
	}
}

func TestSetSize(t *testing.T) {
	m := NewModel()
	m.SetSize(100, 50)

	if m.width != 100 {
		t.Errorf("expected width 100, got %d", m.width)
	}
	if m.height != 50 {
		t.Errorf("expected height 50, got %d", m.height)
	}
}

func TestViewEmpty(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 20)
	view := m.View()

	if !strings.Contains(view, "[No tasks]") {
		t.Error("empty view should contain '[No tasks]'")
	}
}

func TestViewSingleTask(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 20)
	m.AddTask("task-001", "executor", "Test task")

	view := m.View()

	if !strings.Contains(view, "task-001") {
		t.Error("view should contain task-001")
	}
	if !strings.Contains(view, "[executor]") {
		t.Error("view should contain [executor]")
	}
	if !strings.Contains(view, "active") {
		t.Error("view should contain active state")
	}
}

func TestViewTreeStructure(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 30)
	m.AddTask("task-001", "executor", "Root task")
	m.AddChildTask("task-001", "task-002", "reviewer", "Child 1")
	m.AddChildTask("task-001", "task-003", "tester", "Child 2")
	m.AddChildTask("task-003", "task-004", "executor", "Grandchild")

	view := m.View()

	// Check all tasks are present
	if !strings.Contains(view, "task-001") {
		t.Error("view should contain task-001")
	}
	if !strings.Contains(view, "task-002") {
		t.Error("view should contain task-002")
	}
	if !strings.Contains(view, "task-003") {
		t.Error("view should contain task-003")
	}
	if !strings.Contains(view, "task-004") {
		t.Error("view should contain task-004")
	}

	// Check tree structure characters
	if !strings.Contains(view, "├──") || !strings.Contains(view, "└──") {
		t.Error("view should contain tree structure characters")
	}
}

func TestViewStateColors(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 20)
	m.AddTask("task-001", "executor", "Task 1")
	m.AddTask("task-002", "executor", "Task 2")
	m.AddTask("task-003", "executor", "Task 3")
	m.AddTask("task-004", "executor", "Task 4")

	m.UpdateTaskState("task-001", "active")
	m.UpdateTaskState("task-002", "done")
	m.UpdateTaskState("task-003", "failed")
	m.UpdateTaskState("task-004", "cancelled")

	view := m.View()

	// Verify all states are rendered
	if !strings.Contains(view, "active") {
		t.Error("view should contain active state")
	}
	if !strings.Contains(view, "done") {
		t.Error("view should contain done state")
	}
	if !strings.Contains(view, "failed") {
		t.Error("view should contain failed state")
	}
	if !strings.Contains(view, "cancelled") {
		t.Error("view should contain cancelled state")
	}

	// Note: Color escape codes are present in the view,
	// but we don't test for exact ANSI codes as they may vary
}

func TestViewZeroSize(t *testing.T) {
	m := NewModel()
	m.AddTask("task-001", "executor", "Test task")

	// Without SetSize, width and height are 0
	view := m.View()
	if view != "" {
		t.Error("view should be empty when size is 0")
	}
}

func TestMultipleRoots(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 30)
	m.AddTask("task-001", "executor", "Root 1")
	m.AddTask("task-002", "executor", "Root 2")

	if len(m.roots) != 2 {
		t.Errorf("expected 2 roots, got %d", len(m.roots))
	}

	view := m.View()
	if !strings.Contains(view, "task-001") || !strings.Contains(view, "task-002") {
		t.Error("view should contain both root tasks")
	}
}

func TestAddChildTaskNonexistentParent(t *testing.T) {
	m := NewModel()
	m.AddChildTask("nonexistent", "task-001", "executor", "Child task")

	// Child should still be created
	if _, exists := m.tasks["task-001"]; !exists {
		t.Error("child task should be created even if parent doesn't exist")
	}
}

func TestAddDuplicateChild(t *testing.T) {
	m := NewModel()
	m.AddTask("task-001", "executor", "Parent task")
	m.AddChildTask("task-001", "task-002", "reviewer", "Child task")
	m.AddChildTask("task-001", "task-002", "reviewer", "Child task")

	parent := m.tasks["task-001"]
	if len(parent.Children) != 1 {
		t.Errorf("expected 1 child (no duplicates), got %d", len(parent.Children))
	}
}
