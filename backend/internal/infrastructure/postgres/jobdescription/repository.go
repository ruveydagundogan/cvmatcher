package jobdescription

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/domain/jobdescription/model"
)

type JobDescriptionRepository struct {
	pool *pgxpool.Pool
}

func NewJobDescriptionRepository(pool *pgxpool.Pool) *JobDescriptionRepository {
	return &JobDescriptionRepository{pool: pool}
}

func (r *JobDescriptionRepository) Save(ctx context.Context, jd *model.JobDescription) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO job_descriptions (id, user_id, title, content, required_skills, preferred_skills, experience_level, employment_type, location, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		jd.ID, jd.UserID, jd.Title, jd.Content, jd.RequiredSkills, jd.PreferredSkills, jd.ExperienceLevel, jd.EmploymentType, jd.Location, jd.CreatedAt, jd.UpdatedAt,
	)
	return err
}

func (r *JobDescriptionRepository) FindByID(ctx context.Context, id string) (*model.JobDescription, error) {
	jd := &model.JobDescription{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, title, content, required_skills, preferred_skills, experience_level, employment_type, location, created_at, updated_at FROM job_descriptions WHERE id=$1`, id,
	).Scan(&jd.ID, &jd.UserID, &jd.Title, &jd.Content, &jd.RequiredSkills, &jd.PreferredSkills, &jd.ExperienceLevel, &jd.EmploymentType, &jd.Location, &jd.CreatedAt, &jd.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return jd, nil
}

func (r *JobDescriptionRepository) FindByUserID(ctx context.Context, userID string, offset, limit int) ([]*model.JobDescription, int, error) {
	var total int
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM job_descriptions WHERE user_id=$1`, userID).Scan(&total)

	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, title, content, required_skills, preferred_skills, experience_level, employment_type, location, created_at, updated_at
		 FROM job_descriptions WHERE user_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, userID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var jds []*model.JobDescription
	for rows.Next() {
		jd := &model.JobDescription{}
		if err := rows.Scan(&jd.ID, &jd.UserID, &jd.Title, &jd.Content, &jd.RequiredSkills, &jd.PreferredSkills, &jd.ExperienceLevel, &jd.EmploymentType, &jd.Location, &jd.CreatedAt, &jd.UpdatedAt); err != nil {
			return nil, 0, err
		}
		jds = append(jds, jd)
	}
	if jds == nil {
		jds = []*model.JobDescription{}
	}
	return jds, total, nil
}

func (r *JobDescriptionRepository) Update(ctx context.Context, jd *model.JobDescription) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE job_descriptions SET title=$1, content=$2, required_skills=$3, preferred_skills=$4, experience_level=$5, employment_type=$6, location=$7, updated_at=NOW() WHERE id=$8`,
		jd.Title, jd.Content, jd.RequiredSkills, jd.PreferredSkills, jd.ExperienceLevel, jd.EmploymentType, jd.Location, jd.ID,
	)
	return err
}

func (r *JobDescriptionRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM job_descriptions WHERE id=$1`, id)
	return err
}
