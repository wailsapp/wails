package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"testing"
)

var loadingModelResult Model
var loadingViewResult string
var loadingUpdateResult tea.Model
var loadingCommandResult tea.Cmd

func BenchmarkDoctorNewModel(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		loadingModelResult = NewModel()
	}
}

func BenchmarkDoctorLoadingView(b *testing.B) {
	model := NewModel()
	b.ReportAllocs()
	for b.Loop() {
		loadingViewResult = model.viewLoading()
	}
}

func BenchmarkDoctorLoadingUpdate(b *testing.B) {
	model := NewModel()
	commands := model.Init()().(tea.BatchMsg)
	message := commands[0]() // Execute only the clock, never a real system scan.
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		loadingUpdateResult, loadingCommandResult = model.Update(message)
	}
}
