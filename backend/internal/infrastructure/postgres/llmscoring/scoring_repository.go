package postgres

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ruveydagundogan/cvmatcher/backend/internal/domain/llmscoring/model"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/domain/llmscoring/repository"
)

type ScoringRepository struct {
	pool *pgxpool.Pool
}

func NewScoringRepository(pool *pgxpool.Pool) *ScoringRepository {
	return &ScoringRepository{pool: pool}
}

func nullableUserID(userID string) interface{} {
	if userID == "" {
		return nil
	}
	return userID
}

func (r *ScoringRepository) Save(ctx context.Context, score *model.ScoreRequest) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO scoring_requests (id, user_id, prompt, response, score, model, inference_ms, word_count, char_count, category, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		score.ID, nullableUserID(score.UserID), score.Prompt, score.Response, score.Score,
		score.Model, score.InferenceMs, score.WordCount, score.CharCount,
		score.Category, score.CreatedAt,
	)
	return err
}

func (r *ScoringRepository) FindByID(ctx context.Context, id string) (*model.ScoreRequest, error) {
	score := &model.ScoreRequest{}
	var userID sql.NullString
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, prompt, response, score, model, inference_ms, word_count, char_count, category, created_at
		 FROM scoring_requests WHERE id = $1`, id,
	).Scan(&score.ID, &userID, &score.Prompt, &score.Response, &score.Score,
		&score.Model, &score.InferenceMs, &score.WordCount, &score.CharCount,
		&score.Category, &score.CreatedAt)
	if err != nil {
		return nil, err
	}
	if userID.Valid {
		score.UserID = userID.String
	}
	return score, nil
}

func (r *ScoringRepository) FindByUserID(ctx context.Context, userID string, offset, limit int) ([]*model.ScoreRequest, int, error) {
	var total int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM scoring_requests WHERE user_id = $1`, nullableUserID(userID),
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, prompt, response, score, model, inference_ms, word_count, char_count, category, created_at
		 FROM scoring_requests WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		nullableUserID(userID), limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	scores := make([]*model.ScoreRequest, 0, min(limit, 64))
	for rows.Next() {
		score := &model.ScoreRequest{}
		var userIDCol sql.NullString
		if err := rows.Scan(&score.ID, &userIDCol, &score.Prompt, &score.Response, &score.Score,
			&score.Model, &score.InferenceMs, &score.WordCount, &score.CharCount,
			&score.Category, &score.CreatedAt); err != nil {
			return nil, 0, err
		}
		if userIDCol.Valid {
			score.UserID = userIDCol.String
		}
		scores = append(scores, score)
	}

	return scores, total, nil
}

func (r *ScoringRepository) DeleteByUserID(ctx context.Context, userID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM scoring_requests WHERE user_id = $1`, nullableUserID(userID))
	return err
}

func (r *ScoringRepository) GetStats(ctx context.Context, userID string) (*repository.ScoringStats, error) {
	stats := &repository.ScoringStats{}
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*), COALESCE(AVG(score), 0), COALESCE(SUM(word_count), 0), COALESCE(SUM(char_count), 0)
		 FROM scoring_requests WHERE user_id = $1`, nullableUserID(userID),
	).Scan(&stats.TotalRequests, &stats.AverageScore, &stats.TotalWords, &stats.TotalChars)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
