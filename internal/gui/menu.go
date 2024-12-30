package gui

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	confluentKafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/symonk/kui/internal/kafka"
	"golang.org/x/term"
)

var menuStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))

var highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
var tabStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).BorderForeground(highlight).Padding(0, 1)
var activeTab = tabStyle.BorderStyle(lipgloss.ThickBorder()).Padding(0, 1)
var tabGap = tabStyle.BorderTop(false).BorderLeft(false).BorderRight(false)

// Menu is the core view of kui.
type Menu struct {
	client       *kafka.Client
	status       *StatusBar
	connected    bool
	width        int
	height       int
	meta         *confluentKafka.Metadata
	brokers      []*confluentKafka.BrokerMetadata
	tabs         []string
	topicsView   tea.Model
	groupsView   tea.Model
	settingsView tea.Model
	configView   tea.Model
	debugView    tea.Model
}

func NewMenu(client *kafka.Client) *Menu {
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))
	dView := NewDebugView()
	return &Menu{
		client:    client,
		status:    NewStatusBar(client),
		width:     w,
		height:    h,
		tabs:      []string{"Overview", "Topics", "Consumer Groups", "Settings"},
		debugView: dView}
}

func (m *Menu) View() string {
	var doc strings.Builder
	row := lipgloss.JoinHorizontal(lipgloss.Top, activeTab.Render("Overview"), tabStyle.Render("Topics"), tabStyle.Render("Consumer Groups"), tabStyle.Render("Settings"), tabStyle.Render("Config"), tabStyle.Render("Debug"))
	gap := tabGap.Render(strings.Repeat(" ", max(0, m.width-lipgloss.Width(row))))
	newRow := lipgloss.JoinHorizontal(lipgloss.Top, row, gap)
	doc.WriteString(newRow + "\n\n")
	doc.WriteString(m.debugView.View())
	return doc.String() + m.status.View()
}

func (m *Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	case connMessage:
		m.status.Update(msg)
	case *MetaData:
		m.meta = msg.meta
		return m, kafkaMetaDataCommand(m.client)
	}
	m.debugView.Update(msg)
	return m, nil
}

func (m *Menu) Init() tea.Cmd {
	return kafkaMetaDataCommand(m.client)
}

// OverviewTabView encapsulates the view when the application has the
// `overview` tab selected.
type OverviewTabView struct {
	client *kafka.Client
	meta   *confluentKafka.Metadata
	table  table.Model
}

func (o *OverviewTabView) View() string {
	return "Overview"
}

func (o *OverviewTabView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return o, nil
}

func (o *OverviewTabView) Init() tea.Cmd {
	return kafkaMetaDataCommand(o.client)
}
