package store

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/goreleaser/fileglob"
	"github.com/teilomillet/gollm"
	"github.com/teilomillet/gollm/utils"
)

type Store struct {
	llm gollm.LLM
}

func NewStore(provider string, model string, timeout time.Duration) (*Store, error) {
	llm, err := gollm.NewLLM(
		gollm.SetProvider(provider),
		gollm.SetModel(model),
		gollm.SetTimeout(timeout),
		gollm.SetLogLevel(utils.LogLevelOff),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create LLM: %v", err)
	}

	return &Store{llm: llm}, nil
}

func (s *Store) Generate(promt string, pattern string) {
	files, err := fileglob.Glob(pattern)
	if err != nil {
		panic(err)
	}

	if len(files) == 0 {
		panic("No files")
	}

	for _, file := range files {
		s.process(file, promt)
	}
}

var template = `Process input file using the following instructions: "%s",
Return only proccessed text. Do not include a preamble.
Input file: %s`

func (s *Store) process(file string, prompt string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	promptObject := gollm.NewPrompt(
		fmt.Sprintf(template, strings.ReplaceAll(prompt, "\"", "\\\""), string(data)),
		gollm.WithOutput("Return only the text. Do not include a preamble."),
	)

	ctx := context.Background()
	response, err := s.llm.Generate(ctx, promptObject)
	if err != nil {
		log.Fatalf("Failed to generate response: %v", err)
	}

	err = os.WriteFile(file, []byte(response), 0664)
	if err != nil {
		return err
	}

	return nil
}
