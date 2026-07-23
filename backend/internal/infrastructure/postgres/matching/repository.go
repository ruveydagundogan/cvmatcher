package matching

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/domain/matching/model"
)

type MatchingRepository struct {
	pool *pgxpool.Pool
}

func NewMatchingRepository(pool *pgxpool.Pool) *MatchingRepository {
	return &MatchingRepository{pool: pool}
}

func (r *MatchingRepository) Save(ctx context.Context, m *model.MatchResult) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO match_results (id, user_id, cv_id, jd_id, overall_score, skill_match_score, experience_score, education_score, llm_analysis, details, matched_skills, missing_skills, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		 ON CONFLICT (cv_id, jd_id) DO UPDATE SET
		   overall_score=EXCLUDED.overall_score,
		   skill_match_score=EXCLUDED.skill_match_score,
		   experience_score=EXCLUDED.experience_score,
		   education_score=EXCLUDED.education_score,
		   llm_analysis=EXCLUDED.llm_analysis,
		   matched_skills=EXCLUDED.matched_skills,
		   missing_skills=EXCLUDED.missing_skills`,
		m.ID, m.UserID, m.CVID, m.JDID, m.OverallScore, m.SkillMatchScore, m.ExperienceScore, m.EducationScore,
		m.LLMAnalysis, "{}", m.MatchedSkills, m.MissingSkills, m.CreatedAt,
	)
	return err
}

func (r *MatchingRepository) FindByID(ctx context.Context, id string) (*model.MatchResult, error) {
	m := &model.MatchResult{}
	err := r.pool.QueryRow(ctx,
		`SELECT mr.id, mr.user_id, mr.cv_id, mr.jd_id, mr.overall_score, mr.skill_match_score, mr.experience_score, mr.education_score, mr.llm_analysis,
		        mr.matched_skills, mr.missing_skills, mr.created_at,
		        c.title, c.file_name, jd.title
		 FROM match_results mr
		 JOIN cvs c ON c.id = mr.cv_id
		 JOIN job_descriptions jd ON jd.id = mr.jd_id
		 WHERE mr.id=$1`, id,
	).Scan(&m.ID, &m.UserID, &m.CVID, &m.JDID, &m.OverallScore, &m.SkillMatchScore, &m.ExperienceScore, &m.EducationScore, &m.LLMAnalysis,
		&m.MatchedSkills, &m.MissingSkills, &m.CreatedAt, &m.CVTitle, &m.CVFileName, &m.JDTitle)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *MatchingRepository) FindByCVAndJD(ctx context.Context, cvID, jdID string) (*model.MatchResult, error) {
	m := &model.MatchResult{}
	err := r.pool.QueryRow(ctx,
		`SELECT mr.id, mr.user_id, mr.cv_id, mr.jd_id, mr.overall_score, mr.skill_match_score, mr.experience_score, mr.education_score, mr.llm_analysis,
		        mr.matched_skills, mr.missing_skills, mr.created_at,
		        c.title, c.file_name, jd.title
		 FROM match_results mr
		 JOIN cvs c ON c.id = mr.cv_id
		 JOIN job_descriptions jd ON jd.id = mr.jd_id
		 WHERE mr.cv_id=$1 AND mr.jd_id=$2`, cvID, jdID,
	).Scan(&m.ID, &m.UserID, &m.CVID, &m.JDID, &m.OverallScore, &m.SkillMatchScore, &m.ExperienceScore, &m.EducationScore, &m.LLMAnalysis,
		&m.MatchedSkills, &m.MissingSkills, &m.CreatedAt, &m.CVTitle, &m.CVFileName, &m.JDTitle)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *MatchingRepository) FindByUserID(ctx context.Context, userID string, offset, limit int) ([]*model.MatchResult, int, error) {
	var total int
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM match_results WHERE user_id=$1`, userID).Scan(&total)

	rows, err := r.pool.Query(ctx,
		`SELECT mr.id, mr.user_id, mr.cv_id, mr.jd_id, mr.overall_score, mr.skill_match_score, mr.experience_score, mr.education_score, mr.llm_analysis,
		        mr.matched_skills, mr.missing_skills, mr.created_at,
		        c.title, c.file_name, jd.title
		 FROM match_results mr
		 JOIN cvs c ON c.id = mr.cv_id
		 JOIN job_descriptions jd ON jd.id = mr.jd_id
		 WHERE mr.user_id=$1 ORDER BY mr.created_at DESC LIMIT $2 OFFSET $3`, userID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var matches []*model.MatchResult
	for rows.Next() {
		m := &model.MatchResult{}
		if err := rows.Scan(&m.ID, &m.UserID, &m.CVID, &m.JDID, &m.OverallScore, &m.SkillMatchScore, &m.ExperienceScore, &m.EducationScore, &m.LLMAnalysis,
			&m.MatchedSkills, &m.MissingSkills, &m.CreatedAt, &m.CVTitle, &m.CVFileName, &m.JDTitle); err != nil {
			return nil, 0, err
		}
		matches = append(matches, m)
	}
	if matches == nil {
		matches = []*model.MatchResult{}
	}
	return matches, total, nil
}

func (r *MatchingRepository) GetDashboardStats(ctx context.Context, userID string) (*model.DashboardStats, error) {
	stats := &model.DashboardStats{}

	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM cvs WHERE user_id=$1`, userID).Scan(&stats.TotalCVs)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM job_descriptions WHERE user_id=$1`, userID).Scan(&stats.TotalJDs)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM match_results WHERE user_id=$1`, userID).Scan(&stats.TotalMatches)
	r.pool.QueryRow(ctx, `SELECT COALESCE(AVG(overall_score),0) FROM match_results WHERE user_id=$1`, userID).Scan(&stats.AverageScore)

	if stats.TotalJDs > 0 {
		stats.MatchRate = float64(stats.TotalMatches) / float64(stats.TotalJDs) * 100
	}

	rows, err := r.pool.Query(ctx,
		`SELECT mr.id, mr.user_id, mr.cv_id, mr.jd_id, mr.overall_score, mr.skill_match_score, mr.experience_score, mr.education_score, mr.llm_analysis,
		        mr.matched_skills, mr.missing_skills, mr.created_at,
		        c.title, c.file_name, jd.title
		 FROM match_results mr
		 JOIN cvs c ON c.id = mr.cv_id
		 JOIN job_descriptions jd ON jd.id = mr.jd_id
		 WHERE mr.user_id=$1 ORDER BY mr.created_at DESC LIMIT 5`, userID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			m := &model.MatchResult{}
			if err := rows.Scan(&m.ID, &m.UserID, &m.CVID, &m.JDID, &m.OverallScore, &m.SkillMatchScore, &m.ExperienceScore, &m.EducationScore, &m.LLMAnalysis,
				&m.MatchedSkills, &m.MissingSkills, &m.CreatedAt, &m.CVTitle, &m.CVFileName, &m.JDTitle); err == nil {
				stats.RecentMatches = append(stats.RecentMatches, m)
			}
		}
	}
	if stats.RecentMatches == nil {
		stats.RecentMatches = []*model.MatchResult{}
	}

	return stats, nil
}
