package helpers

import (
	"github.com/charmbracelet/lipgloss"
)

// const COLORINFO = lipgloss.Color("#00cdff")
// const COLORWARN = lipgloss.Color("#ffc200")
const COLORERROR = lipgloss.Color("1")
const COLORINFO = lipgloss.Color("4")
const COLORSUCCESS = lipgloss.Color("2")
const COLORWARN = lipgloss.Color("11")

// var InfoStyleBold = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#23aae5"))
var InfoStyleBold = lipgloss.NewStyle().Bold(true).Foreground(COLORINFO)
var InfoStyle = lipgloss.NewStyle().Foreground(COLORINFO)
var WarnStyleBold = lipgloss.NewStyle().Bold(true).Foreground(COLORWARN)
var WarnStyle = lipgloss.NewStyle().Foreground(COLORWARN)
var ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#e64949"))
var BasicStyle = lipgloss.NewStyle()

var ExtensionInstallStep = lipgloss.NewStyle().MarginLeft(2)
