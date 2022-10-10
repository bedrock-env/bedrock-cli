package helpers

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

const COLORINFO = lipgloss.Color("#00cdff")
const COLORWARN = lipgloss.Color("#ffc200")
const COLORERROR = lipgloss.Color("ff3800")

var InfoStyleBold = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#23aae5"))
var InfoStyle = lipgloss.NewStyle().Foreground(COLORINFO)
var WarnStyleBold = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#e5c223"))
var WarnStyle = lipgloss.NewStyle().Foreground(COLORWARN)
var ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#e64949"))
var BasicStyle = lipgloss.NewStyle()

var style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	PaddingTop(2).
	PaddingLeft(4).
	Width(22)

func Render() {
	fmt.Println(style.Render("Hello, kitty."))
}
