package store

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/goreleaser/fileglob"
	"github.com/teilomillet/gollm"
	"github.com/teilomillet/gollm/utils"
)

type Store struct {
	llm       gollm.LLM
	eventChan chan Event
}

func NewStore(provider, model string, timeout time.Duration) (*Store, error) {
	llm, err := gollm.NewLLM(
		gollm.SetProvider(provider),
		gollm.SetModel(model),
		gollm.SetTimeout(timeout),
		gollm.SetLogLevel(utils.LogLevelOff),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create LLM: %v", err)
	}
	return &Store{
		llm:       llm,
		eventChan: make(chan Event, 100),
	}, nil
}

func (s *Store) Generate(promt, pattern, suffix string) {
	files, err := fileglob.Glob(pattern)
	if err != nil {
		panic(err)
	}

	if len(files) == 0 {
		panic("No files")
	}

	for _, file := range files {
		s.eventChan <- Event{Type: FileProcessing, InputFileName: file}
		outFile := createOutputFileName(file, suffix)
		s.process(file, promt, outFile)
		s.eventChan <- Event{Type: FileProcessed, InputFileName: file, OutputFileName: outFile}
	}
	close(s.eventChan)
}

func (s *Store) Events() <-chan Event {
	return s.eventChan
}

func createOutputFileName(file, suffix string) string {
	if suffix == "" {
		return file
	}
	ext := filepath.Ext(file)
	return fmt.Sprintf(
		"%s.%s.%s",
		strings.TrimSuffix(file, ext),
		strings.Trim(suffix, "."),
		strings.Trim(ext, "."),
	)
}

func (s *Store) process(inputFile, prompt, outFile string) error {
	defer func() {
		if memoryLLM, ok := s.llm.(interface{ ClearMemory() }); ok {
			memoryLLM.ClearMemory()
			log.Println("Memory cleared. Starting a new conversation.")
			s.eventChan <- Event{Type: MemoryCleared}
		}
	}()

	data, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	promptObject := gollm.NewPrompt(
		fmt.Sprintf(`Process input file folowing instructions: %s.
Input file: %s`, prompt, string(data)),
		gollm.WithContext(
			"You are applcation wich process input file using the instructions."+
				"Your output will be used as content file.",
		),
		gollm.WithDirectives(
			"Do not add or remove content unless specifically stated in the instructions.",
			"Do not add or remove content unless specifically stated in the instructions.",
			"Follow the instructions strictly",
			"Don't remove empty lines",
			"Avoid wrapping the result in ```",
			"Do not add preambles or postambles",
		),
		gollm.WithOutput("Return full content of the file with the changes."),
	)

	ctx := context.Background()
	response, err := s.llm.Generate(ctx, promptObject)
	if err != nil {
		log.Fatalf("Failed to generate response: %v", err)
	}

	err = os.WriteFile(outFile, []byte(response), 0664)
	if err != nil {
		return err
	}

	return nil
}
