package cv

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/domain/cv/model"
)

type CVRepository struct {
	pool *pgxpool.Pool
}

func NewCVRepository(pool *pgxpool.Pool) *CVRepository {
	return &CVRepository{pool: pool}
}

func (r *CVRepository) Save(ctx context.Context, cv *model.CV) error {
	exp, _ := json.Marshal(cv.ParsedExperience)
	edu, _ := json.Marshal(cv.ParsedEducation)
	_, err := r.pool.Exec(ctx,
		`INSERT INTO cvs (id, user_id, title, content, parsed_skills, parsed_experience, parsed_education, parsed_summary, status, file_name, file_size, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		cv.ID, cv.UserID, cv.Title, cv.Content, cv.ParsedSkills, string(exp), string(edu), cv.ParsedSummary,
		cv.Status, cv.FileName, cv.FileSize, cv.CreatedAt, cv.UpdatedAt,
	)
	return err
}

func (r *CVRepository) FindByID(ctx context.Context, id string) (*model.CV, error) {
	var expJSON, eduJSON string
	cv := &model.CV{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, title, content, parsed_skills, parsed_experience, parsed_education, parsed_summary, status, file_name, file_size, created_at, updated_at FROM cvs WHERE id=$1`, id,
	).Scan(&cv.ID, &cv.UserID, &cv.Title, &cv.Content, &cv.ParsedSkills, &expJSON, &eduJSON, &cv.ParsedSummary, &cv.Status, &cv.FileName, &cv.FileSize, &cv.CreatedAt, &cv.UpdatedAt)
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(expJSON), &cv.ParsedExperience)
	json.Unmarshal([]byte(eduJSON), &cv.ParsedEducation)
	return cv, nil
}

func (r *CVRepository) FindByUserID(ctx context.Context, userID string, offset, limit int) ([]*model.CV, int, error) {
	var total int
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM cvs WHERE user_id=$1`, userID).Scan(&total)

	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, title, content, parsed_skills, parsed_experience, parsed_education, parsed_summary, status, file_name, file_size, created_at, updated_at
		 FROM cvs WHERE user_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, userID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var cvs []*model.CV
	for rows.Next() {
		var expJSON, eduJSON string
		cv := &model.CV{}
		if err := rows.Scan(&cv.ID, &cv.UserID, &cv.Title, &cv.Content, &cv.ParsedSkills, &expJSON, &eduJSON, &cv.ParsedSummary, &cv.Status, &cv.FileName, &cv.FileSize, &cv.CreatedAt, &cv.UpdatedAt); err != nil {
			return nil, 0, err
		}
		json.Unmarshal([]byte(expJSON), &cv.ParsedExperience)
		json.Unmarshal([]byte(eduJSON), &cv.ParsedEducation)
		cvs = append(cvs, cv)
	}
	if cvs == nil {
		cvs = []*model.CV{}
	}
	return cvs, total, nil
}

func (r *CVRepository) Update(ctx context.Context, cv *model.CV) error {
	exp, _ := json.Marshal(cv.ParsedExperience)
	edu, _ := json.Marshal(cv.ParsedEducation)
	_, err := r.pool.Exec(ctx,
		`UPDATE cvs SET title=$1, content=$2, parsed_skills=$3, parsed_experience=$4, parsed_education=$5, parsed_summary=$6, status=$7, file_name=$8, file_size=$9, updated_at=NOW() WHERE id=$10`,
		cv.Title, cv.Content, cv.ParsedSkills, string(exp), string(edu), cv.ParsedSummary, cv.Status, cv.FileName, cv.FileSize, cv.ID,
	)
	return err
}

func (r *CVRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM cvs WHERE id=$1`, id)
	return err
}
