package memory

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/ruveydagundogan/cvmatcher/backend/internal/domain/matching/model"
)

type InMemoryMatchingRepo struct {
	mu      sync.RWMutex
	matches map[string]*model.MatchResult
}

func NewInMemoryMatchingRepo() *InMemoryMatchingRepo {
	return &InMemoryMatchingRepo{matches: make(map[string]*model.MatchResult)}
}

func (r *InMemoryMatchingRepo) Save(ctx context.Context, m *model.MatchResult) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.matches[m.ID] = m
	return nil
}

func (r *InMemoryMatchingRepo) FindByID(ctx context.Context, id string) (*model.MatchResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.matches[id]
	if !ok {
		return nil, fmt.Errorf("match result not found")
	}
	return m, nil
}

func (r *InMemoryMatchingRepo) FindByCVAndJD(ctx context.Context, cvID, jdID string) (*model.MatchResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, m := range r.matches {
		if m.CVID == cvID && m.JDID == jdID {
			return m, nil
		}
	}
	return nil, fmt.Errorf("match not found")
}

func (r *InMemoryMatchingRepo) FindByUserID(ctx context.Context, userID string, offset, limit int) ([]*model.MatchResult, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var matches []*model.MatchResult
	for _, m := range r.matches {
		if m.UserID == userID {
			matches = append(matches, m)
		}
	}
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].CreatedAt.After(matches[j].CreatedAt)
	})

	total := len(matches)
	if offset >= total {
		return []*model.MatchResult{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return matches[offset:end], total, nil
}

func (r *InMemoryMatchingRepo) GetDashboardStats(ctx context.Context, userID string) (*model.DashboardStats, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stats := &model.DashboardStats{
		RecentMatches: []*model.MatchResult{},
	}

	var matches []*model.MatchResult
	for _, m := range r.matches {
		if m.UserID == userID {
			matches = append(matches, m)
		}
	}

	stats.TotalMatches = len(matches)

	var totalScore float64
	for _, m := range matches {
		totalScore += m.OverallScore
	}
	if len(matches) > 0 {
		stats.AverageScore = totalScore / float64(len(matches))
	}

	// CV ve JD sayılarını da hesapla (diğer repo'lardan gelebilir ama in-memory'de basit)
	// Not: Bu repo sadece match tutar, CV/JD repo'ları ayrı. 
	// İleride birleştirilebilir. Şimdilik 0 kalsın.
	stats.TotalCVs = 0
	stats.TotalJDs = 0

	// Top skill hesapla
	skillCounts := make(map[string]int)
	for _, m := range matches {
		for _, s := range m.MatchedSkills {
			skillCounts[strings.ToLower(s)]++
		}
	}
	var topSkill string
	var topSkillCount int
	for skill, count := range skillCounts {
		if count > topSkillCount {
			topSkill = skill
			topSkillCount = count
		}
	}
	stats.TopSkill = topSkill
	stats.TopSkillCount = topSkillCount

	// Match rate (basit: total matches / total CVs veya 0)
	stats.MatchRate = 0

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].CreatedAt.After(matches[j].CreatedAt)
	})
	n := 5
	if len(matches) < n {
		n = len(matches)
	}
	if n > 0 {
		stats.RecentMatches = matches[:n]
	}

	return stats, nil
}

func (r *InMemoryMatchingRepo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.matches, id)
	return nil
}
