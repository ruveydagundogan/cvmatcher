package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/ruveydagundogan/llm-decision-score/backend/internal/domain/llmscoring/model"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/domain/llmscoring/repository"
)

type InMemoryScoringRepo struct {
	mu     sync.RWMutex
	scores []*model.ScoreRequest
}

func NewInMemoryScoringRepo() *InMemoryScoringRepo {
	return &InMemoryScoringRepo{}
}

func (r *InMemoryScoringRepo) Save(_ context.Context, score *model.ScoreRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.scores = append(r.scores, score)
	return nil
}

func (r *InMemoryScoringRepo) FindByID(_ context.Context, id string) (*model.ScoreRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, s := range r.scores {
		if s.ID == id {
			return s, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (r *InMemoryScoringRepo) FindByUserID(_ context.Context, userID string, offset, limit int) ([]*model.ScoreRequest, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*model.ScoreRequest
	for _, s := range r.scores {
		if s.UserID == userID {
			filtered = append(filtered, s)
		}
	}

	total := len(filtered)
	if offset >= total {
		return nil, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return filtered[offset:end], total, nil
}

func (r *InMemoryScoringRepo) DeleteByUserID(_ context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	var remaining []*model.ScoreRequest
	for _, s := range r.scores {
		if s.UserID != userID {
			remaining = append(remaining, s)
		}
	}
	r.scores = remaining
	return nil
}

func (r *InMemoryScoringRepo) GetStats(_ context.Context, userID string) (*repository.ScoringStats, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stats := &repository.ScoringStats{}
	for _, s := range r.scores {
		if s.UserID == userID {
			stats.TotalRequests++
			stats.TotalWords += s.WordCount
			stats.TotalChars += s.CharCount
		}
	}
	if stats.TotalRequests > 0 {
		var totalScore float64
		for _, s := range r.scores {
			if s.UserID == userID {
				totalScore += float64(s.Score)
			}
		}
		stats.AverageScore = totalScore / float64(stats.TotalRequests)
	}
	return stats, nil
}
