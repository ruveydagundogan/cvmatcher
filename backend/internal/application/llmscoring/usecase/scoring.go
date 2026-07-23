package usecase

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/ruveydagundogan/cvmatcher/backend/internal/application/llmscoring/dto"
	auditmodel "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/audit/model"
	auditrepo "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/audit/repository"
	scoringmodel "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/llmscoring/model"
	scoringrepo "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/llmscoring/repository"
	apperrors "github.com/ruveydagundogan/cvmatcher/backend/internal/shared/errors"
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
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

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
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

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
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

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
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

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

	lowerResp := strings.ToLower(response)
	respLen := len(lowerResp)

	switch {
	case respLen >= 800:
		score += 25
	case respLen >= 500:
		score += 22
	case respLen >= 300:
		score += 18
	case respLen >= 150:
		score += 14
	case respLen >= 80:
		score += 10
	default:
		score += 5
	}

	lowerPrompt := strings.ToLower(prompt)
	validWords, matchedWords := countMatchingWords(lowerPrompt, lowerResp)

	if validWords > 0 {
		score += (matchedWords * 30) / validWords
	}

	var hasDot, hasNewline, hasFormatting, hasComma, hasColon, hasSemicolon, hasParen bool
	dotCount := 0
	for i := 0; i < respLen; i++ {
		switch lowerResp[i] {
		case '.':
			hasDot = true
			dotCount++
		case '\n':
			hasNewline = true
		case '*', '-':
			hasFormatting = true
		case ',':
			hasComma = true
		case ':':
			hasColon = true
		case ';':
			hasSemicolon = true
		case '(':
			hasParen = true
		}
	}
	if hasDot {
		score += 5
	}
	if hasNewline {
		score += 5
	}
	if hasFormatting {
		score += 5
	}
	if dotCount >= 3 {
		score += 5
	}
	if hasComma {
		score += 3
	}
	if hasColon {
		score += 3
	}
	if hasSemicolon {
		score += 3
	}
	if hasParen {
		score += 3
	}

	wordCount := countWords(lowerResp)
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

func countWords(s string) int {
	count := 0
	inWord := false
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' || s[i] == '\n' || s[i] == '\t' {
			if inWord {
				count++
				inWord = false
			}
		} else {
			inWord = true
		}
	}
	if inWord {
		count++
	}
	return count
}

func countMatchingWords(lowerPrompt, lowerResp string) (valid, matched int) {
	start := -1
	for i := 0; i <= len(lowerPrompt); i++ {
		isSpace := i == len(lowerPrompt) || lowerPrompt[i] == ' ' || lowerPrompt[i] == '\n' || lowerPrompt[i] == '\t'
		if !isSpace && start == -1 {
			start = i
			continue
		}
		if isSpace && start != -1 {
			word := lowerPrompt[start:i]
			start = -1
			for len(word) > 0 && isPunct(word[0]) {
				word = word[1:]
			}
			for len(word) > 0 && isPunct(word[len(word)-1]) {
				word = word[:len(word)-1]
			}
			if len(word) <= 2 {
				continue
			}
			valid++
			if strings.Contains(lowerResp, word) {
				matched++
			}
		}
	}
	return
}

func isPunct(b byte) bool {
	switch b {
	case '.', ',', '!', '?', ';', ':', '"', '\'', '(', ')', '[', ']', '{', '}':
		return true
	}
	return false
}
