package usecase

import (
	"context"
	"log/slog"
	"strings"

	"github.com/ruveydagundogan/llm-decision-score/backend/internal/application/llmscoring/dto"
	auditmodel "github.com/ruveydagundogan/llm-decision-score/backend/internal/domain/audit/model"
	auditrepo "github.com/ruveydagundogan/llm-decision-score/backend/internal/domain/audit/repository"
	scoringmodel "github.com/ruveydagundogan/llm-decision-score/backend/internal/domain/llmscoring/model"
	scoringrepo "github.com/ruveydagundogan/llm-decision-score/backend/internal/domain/llmscoring/repository"
	apperrors "github.com/ruveydagundogan/llm-decision-score/backend/internal/shared/errors"
)

type ScoreUseCase struct {
	scoringRepo scoringrepo.ScoringRepository
	auditRepo   auditrepo.AuditRepository
	logger      *slog.Logger
}

func NewScoreUseCase(
	scoringRepo scoringrepo.ScoringRepository,
	auditRepo auditrepo.AuditRepository,
	logger *slog.Logger,
) *ScoreUseCase {
	return &ScoreUseCase{
		scoringRepo: scoringRepo,
		auditRepo:   auditRepo,
		logger:      logger,
	}
}

func (uc *ScoreUseCase) Execute(ctx context.Context, userID string, req dto.ScoreRequestDTO) (*dto.ScoreResponseDTO, error) {
	scoreItem := scoringmodel.NewScoreRequest(userID, req.Prompt, req.Response, req.Model, req.InferenceMs)
	scoreItem.Score = calculateDecisionScore(req.Prompt, req.Response)

	if err := uc.scoringRepo.Save(ctx, scoreItem); err != nil {
		return nil, apperrors.Internal("failed to save score", err)
	}

	_ = uc.auditRepo.Save(ctx, auditmodel.NewAuditLog(
		userID, "score", "scoring", scoreItem.ID,
		"", "", nil,
	))

	uc.logger.Info("score calculated",
		"user_id", userID,
		"score_id", scoreItem.ID,
		"score", scoreItem.Score,
	)

	return &dto.ScoreResponseDTO{
		ID:        scoreItem.ID,
		Score:     scoreItem.Score,
		Prompt:    scoreItem.Prompt,
		Response:  scoreItem.Response,
		Model:     scoreItem.Model,
		WordCount: scoreItem.WordCount,
		CharCount: scoreItem.CharCount,
	}, nil
}

type GetHistoryUseCase struct {
	scoringRepo scoringrepo.ScoringRepository
	logger      *slog.Logger
}

func NewGetHistoryUseCase(scoringRepo scoringrepo.ScoringRepository, logger *slog.Logger) *GetHistoryUseCase {
	return &GetHistoryUseCase{scoringRepo: scoringRepo, logger: logger}
}

func (uc *GetHistoryUseCase) Execute(ctx context.Context, userID string, page, limit int) (*dto.HistoryListResponseDTO, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	items, total, err := uc.scoringRepo.FindByUserID(ctx, userID, offset, limit)
	if err != nil {
		return nil, apperrors.Internal("failed to fetch history", err)
	}

	historyItems := make([]dto.HistoryResponseDTO, len(items))
	for i, item := range items {
		historyItems[i] = dto.HistoryResponseDTO{
			ID:        item.ID,
			Prompt:    item.Prompt,
			Response:  item.Response,
			Score:     item.Score,
			Model:     item.Model,
			WordCount: item.WordCount,
			CharCount: item.CharCount,
			CreatedAt: item.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return &dto.HistoryListResponseDTO{
		Items: historyItems,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

type DeleteHistoryUseCase struct {
	scoringRepo scoringrepo.ScoringRepository
	auditRepo   auditrepo.AuditRepository
	logger      *slog.Logger
}

func NewDeleteHistoryUseCase(
	scoringRepo scoringrepo.ScoringRepository,
	auditRepo auditrepo.AuditRepository,
	logger *slog.Logger,
) *DeleteHistoryUseCase {
	return &DeleteHistoryUseCase{
		scoringRepo: scoringRepo,
		auditRepo:   auditRepo,
		logger:      logger,
	}
}

func (uc *DeleteHistoryUseCase) Execute(ctx context.Context, userID string) error {
	if err := uc.scoringRepo.DeleteByUserID(ctx, userID); err != nil {
		return apperrors.Internal("failed to delete history", err)
	}

	_ = uc.auditRepo.Save(ctx, auditmodel.NewAuditLog(
		userID, "delete_history", "scoring", "",
		"", "", nil,
	))

	uc.logger.Info("history deleted", "user_id", userID)
	return nil
}

type GetStatsUseCase struct {
	scoringRepo scoringrepo.ScoringRepository
	logger      *slog.Logger
}

func NewGetStatsUseCase(scoringRepo scoringrepo.ScoringRepository, logger *slog.Logger) *GetStatsUseCase {
	return &GetStatsUseCase{scoringRepo: scoringRepo, logger: logger}
}

func (uc *GetStatsUseCase) Execute(ctx context.Context, userID string) (*dto.StatsResponseDTO, error) {
	stats, err := uc.scoringRepo.GetStats(ctx, userID)
	if err != nil {
		return nil, apperrors.Internal("failed to fetch stats", err)
	}

	return &dto.StatsResponseDTO{
		TotalRequests: stats.TotalRequests,
		AverageScore:  stats.AverageScore,
		TotalWords:    stats.TotalWords,
		TotalChars:    stats.TotalChars,
	}, nil
}

func calculateDecisionScore(prompt, response string) int {
	score := 0

	prompt = strings.ToLower(prompt)
	response = strings.ToLower(response)

	length := len(response)
	switch {
	case length >= 800:
		score += 25
	case length >= 500:
		score += 22
	case length >= 300:
		score += 18
	case length >= 150:
		score += 14
	case length >= 80:
		score += 10
	default:
		score += 5
	}

	promptWords := strings.Fields(prompt)
	validWords := 0
	matchedWords := 0

	for _, word := range promptWords {
		word = strings.Trim(word, ".,!?;:\"'()[]{}")
		if len(word) <= 2 {
			continue
		}
		validWords++
		if strings.Contains(response, word) {
			matchedWords++
		}
	}

	if validWords > 0 {
		score += (matchedWords * 30) / validWords
	}

	if strings.Contains(response, ".") {
		score += 5
	}
	if strings.Contains(response, "\n") {
		score += 5
	}
	if strings.Contains(response, "*") || strings.Contains(response, "-") {
		score += 5
	}
	if strings.Count(response, ".") >= 3 {
		score += 5
	}

	if strings.Contains(response, ",") {
		score += 3
	}
	if strings.Contains(response, ":") {
		score += 3
	}
	if strings.Contains(response, ";") {
		score += 3
	}
	if strings.Contains(response, "(") {
		score += 3
	}
	if len(strings.Fields(response)) >= 80 {
		score += 3
	}

	wordCount := len(strings.Fields(response))
	switch {
	case wordCount >= 120:
		score += 10
	case wordCount >= 80:
		score += 8
	case wordCount >= 50:
		score += 6
	case wordCount >= 30:
		score += 4
	default:
		score += 2
	}

	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}

	return score
}
