package generator

import (
	"path/filepath"

	"github.com/wailsapp/wails/v3/internal/generator/collect"
)

func (generator *Generator) generateEvents(events *collect.EventMap) {
	// Generate event data table.
	generator.scheduler.Schedule(func() {
		file, err := generator.creator.Create(filepath.Join(events.Imports.Self, generator.renderer.EventDataFile()))
		if err != nil {
			generator.logger.Errorf("%v", err)
			generator.logger.Errorf("event data table generation failed")
			return
		}
		defer func() {
			if err := file.Close(); err != nil {
				generator.logger.Errorf("%v", err)
				generator.logger.Errorf("event data table generation failed")
			}
		}()

		err = generator.renderer.EventData(file, events)
		if err != nil {
			generator.logger.Errorf("%v", err)
			generator.logger.Errorf("event data table generation failed")
		}
	})

	// Generate event creation code.
	generator.scheduler.Schedule(func() {
		file, err := generator.creator.Create(filepath.Join(events.Imports.Self, generator.renderer.EventCreateFile()))
		if err != nil {
			generator.logger.Errorf("%v", err)
			generator.logger.Errorf("event creation code generation failed")
			return
		}
		defer func() {
			if err := file.Close(); err != nil {
				generator.logger.Errorf("%v", err)
				generator.logger.Errorf("event creation code generation failed")
			}
		}()

		err = generator.renderer.EventCreate(file, events)
		if err != nil {
			generator.logger.Errorf("%v", err)
			generator.logger.Errorf("event creation code generation failed")
		}
	})
}
