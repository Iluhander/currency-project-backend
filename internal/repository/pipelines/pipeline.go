package pipelines

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/Iluhander/currency-project-backend/internal/model/plugins"
)

type PipelineRepository struct {
	pipe *plugins.Pipeline
	file string
	mu *sync.Mutex
}

func Init(persistFileName string) (*PipelineRepository, error) {
	res := PipelineRepository{
		&plugins.Pipeline{
			make([]*plugins.Plugin, 0),
			make([]*plugins.Plugin, 0),
			make([]*plugins.Plugin, 0),
		},
		persistFileName,
		&sync.Mutex{},
	}

	if _, err := os.Stat(persistFileName); errors.Is(err, os.ErrNotExist) {
		return &res, nil
	}

	jsonContents, readErr := os.ReadFile(persistFileName)
	if readErr != nil {
		return nil, fmt.Errorf("%w; %w", fmt.Errorf("pipeline unmarshalling error"), readErr)
	}

	parseErr := json.Unmarshal(jsonContents, &res.pipe)
	if parseErr != nil {
		return nil, fmt.Errorf("%w; %w", fmt.Errorf("pipeline unmarshalling error"), parseErr)
	}

	return &res, nil
}

func (r *PipelineRepository) UpdatePipeline(newPipeline *plugins.Pipeline) error {
	marshaled, marshalErr := json.Marshal(newPipeline)
	if marshalErr != nil {
		return fmt.Errorf("%w; %w", fmt.Errorf("pipeline saving error"), marshalErr)
	}

	writeErr := os.WriteFile(r.file, marshaled, 0644)
	if writeErr != nil {
		return fmt.Errorf("%w; %w", fmt.Errorf("pipeline saving error"), writeErr)
	}

	r.pipe = newPipeline

	return nil
}

func (r *PipelineRepository) GetPipeline() *plugins.Pipeline {
	r.mu.Lock()
	r.mu.Unlock()

	return r.pipe
}
