package store

import (
	"context"
	"fmt"
	"log"
	"os"
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

func (s *Store) process(file string, prompt string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	promptObject := gollm.NewPrompt(
		fmt.Sprintf(`%s. Input file: %s`, prompt, string(data)),
		gollm.WithContext(
			"You are applcation wich process input file using the instructions."+
				"Your output will be used as content file.",
		),
		gollm.WithDirectives(
			"Do not add or remove content unless specifically stated in the instructions.",
			"Follow the instructions strictly",
			"Don't remove empty lines",
			"Avoid wrapping the result in ```",
		),
		gollm.WithOutput("Return only the text. Do not include a preamble. Do not wrap yout in markdown tags"),
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
