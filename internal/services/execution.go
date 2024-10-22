package services

import (
	"fmt"

	"github.com/Iluhander/currency-project-backend/internal/model"
	"github.com/Iluhander/currency-project-backend/internal/repository/pipelines"
	"github.com/google/uuid"
)

type ExecutionService struct {
	pipeRepo *pipelines.PipelineRepository
}

func Init(pipeRepo *pipelines.PipelineRepository) *ExecutionService {
	return &ExecutionService {
		pipeRepo,
	}
}

func (s *ExecutionService) GetPipeline(pluginType int) []*model.Plugin {
	mergedArr := make([]*model.Plugin, 0)

	pipeline := s.pipeRepo.GetPipeline()

	switch pluginType {
	case model.TAuthPlugin:
		mergedArr = append(mergedArr, pipeline.Auth...)
	case model.TPaymentPlugin:
		mergedArr = append(mergedArr, pipeline.Payment...)
	case model.TStatisticsPlugin:
		mergedArr = append(mergedArr, pipeline.Statistics...)
	default:
		mergedArr = append(mergedArr, pipeline.Auth...)
		mergedArr = append(mergedArr, pipeline.Payment...)
		mergedArr = append(mergedArr, pipeline.Statistics...)
	}

	return mergedArr
}

func (s *ExecutionService) AddPlugin(newData *model.Plugin) (*model.Plugin, error) {
	pipeline := s.pipeRepo.GetPipeline()
	newData.Id = uuid.New().String()

	if newData.Type == model.TAuthPlugin {
		pipeline.Auth = append(pipeline.Auth, newData)
	} else if newData.Type == model.TPaymentPlugin {
		pipeline.Auth = append(pipeline.Payment, newData)
	} else if newData.Type == model.TStatisticsPlugin {
		pipeline.Auth = append(pipeline.Statistics, newData)
	} else {
		return nil, fmt.Errorf("unknown plugin type %d", newData.Type)
	}

	s.pipeRepo.UpdatePipeline(pipeline)
	return newData, nil
}

func (s *ExecutionService) UpdatePlugin(newData *model.Plugin) (*model.Plugin, error) {
	pipeline := s.pipeRepo.GetPipeline()

	for i, v := range pipeline.Auth {
		if v.Id == newData.Id {
			pipeline.Auth[i] = newData
			s.pipeRepo.UpdatePipeline(pipeline)

			return newData, nil
		}
	}

	for i, v := range pipeline.Payment {
		if v.Id == newData.Id {
			pipeline.Payment[i] = newData
			s.pipeRepo.UpdatePipeline(pipeline)

			return newData, nil
		}
	}

	for i, v := range pipeline.Statistics {
		if v.Id == newData.Id {
			pipeline.Statistics[i] = newData
			s.pipeRepo.UpdatePipeline(pipeline)

			return newData, nil
		}
	}

	return nil, fmt.Errorf("plugin not found")
}

func (s *ExecutionService) DeletePlugin(newData *model.Plugin) (error) {
	pipeline := s.pipeRepo.GetPipeline()

	for i, v := range pipeline.Auth {
		if v.Id == newData.Id {
			pipeline.Auth = append(pipeline.Auth[:i], pipeline.Auth[i+1:]...)
			s.pipeRepo.UpdatePipeline(pipeline)

			return nil
		}
	}

	for i, v := range pipeline.Payment {
		if v.Id == newData.Id {
			pipeline.Payment = append(pipeline.Payment[:i], pipeline.Payment[i+1:]...)
			s.pipeRepo.UpdatePipeline(pipeline)

			return nil
		}
	}

	for i, v := range pipeline.Statistics {
		if v.Id == newData.Id {
			pipeline.Statistics = append(pipeline.Statistics[:i], pipeline.Statistics[i+1:]...)
			s.pipeRepo.UpdatePipeline(pipeline)

			return nil
		}
	}

	return fmt.Errorf("plugin not found")
}
