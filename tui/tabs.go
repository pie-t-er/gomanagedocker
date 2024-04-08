package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tab int

const (
	images tab = iota
	containers
	volumes
)

type Model struct {
	Tabs       []string
	TabContent []listModel
	activeTab  int
	width      int
	height     int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func NewModel(tabs []string) Model {
	contents := make([]listModel, 3)

	for i := range contents {
		contents[i] = InitList(fmt.Sprintf("title %d", i), 100, 30)
	}
	return Model{
		Tabs:       tabs,
		TabContent: contents,
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		windowStyle = windowStyle.
			Width(m.width - listDocStyle.GetHorizontalFrameSize() - 2).
			Height(m.height - listDocStyle.GetVerticalFrameSize() - 3)

		//change list dimentions when window size changes
		//TODO: change width
		m.getActiveList().SetWidth(msg.Width)
		m.getActiveList().SetHeight(msg.Height - 7)

		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "l", "n", "tab":
			m.activeTab = min(m.activeTab+1, len(m.Tabs)-1)
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.TabContent[m.activeTab].list, cmd = m.TabContent[m.activeTab].list.Update(msg)

	return m, cmd
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func (m Model) View() string {
	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range m.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == m.activeTab
		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "└"
		} else if isLast && !isActive {
			border.BottomRight = "┴"
		}

		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	var row string
	row = lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	fillerStringLen := windowStyle.GetWidth() - lipgloss.Width(row)
	if fillerStringLen > 0 {
		fillerString := strings.Repeat("─", fillerStringLen+1)
		fillerString += "┐"
		row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, fillerStyle.Render(fillerString))
	}

	body := windowStyle.Render(m.TabContent[m.activeTab].View())
	doc.WriteString(row)
	doc.WriteString("\n")

	doc.WriteString(body)
	return docStyle.Render(doc.String())
}

//Util

func (m Model) getActiveTab() listModel {
	return m.TabContent[m.activeTab]
}

func (m Model) getActiveList() *list.Model {
	return &m.TabContent[m.activeTab].list
}
