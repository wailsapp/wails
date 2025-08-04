package changelog

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// ProcessorResult contains the results of processing a changelog file
type ProcessorResult struct {
	Entry            *ChangelogEntry
	ValidationResult ValidationResult
	HasContent       bool
}

// Processor handles the complete changelog processing pipeline
type Processor struct {
	parser    *Parser
	validator *Validator
}

// NewProcessor creates a new changelog processor with parser and validator
func NewProcessor() *Processor {
	return &Processor{
		parser:    NewParser(),
		validator: NewValidator(),
	}
}

// ProcessFile processes a changelog file and returns the parsed and validated entry
func (p *Processor) ProcessFile(filePath string) (*ProcessorResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open changelog file %s: %w", filePath, err)
	}
	defer file.Close()

	return p.ProcessReader(file)
}

// ProcessReader processes changelog content from a reader
func (p *Processor) ProcessReader(reader io.Reader) (*ProcessorResult, error) {
	// Parse the content
	entry, err := p.parser.ParseContent(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse changelog content: %w", err)
	}

	// Validate the parsed entry
	validationResult := p.validator.ValidateEntry(entry)

	result := &ProcessorResult{
		Entry:            entry,
		ValidationResult: validationResult,
		HasContent:       entry.HasContent(),
	}

	return result, nil
}

// ProcessString processes changelog content from a string
func (p *Processor) ProcessString(content string) (*ProcessorResult, error) {
	return p.ProcessReader(strings.NewReader(content))
}

// ValidateFile validates a changelog file without parsing it into an entry
func (p *Processor) ValidateFile(filePath string) (ValidationResult, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return ValidationResult{Valid: false}, fmt.Errorf("failed to read changelog file %s: %w", filePath, err)
	}

	return p.validator.ValidateContent(string(content)), nil
}

// ValidateString validates changelog content from a string
func (p *Processor) ValidateString(content string) ValidationResult {
	return p.validator.ValidateContent(content)
}
