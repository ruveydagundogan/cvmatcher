package knowledge

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/domain/knowledge/model"
)

type KnowledgeRepository struct {
	pool *pgxpool.Pool
}

func NewKnowledgeRepository(pool *pgxpool.Pool) *KnowledgeRepository {
	return &KnowledgeRepository{pool: pool}
}

func (r *KnowledgeRepository) Save(ctx context.Context, entry *model.KnowledgeEntry) error {
	tags, _ := json.Marshal(entry.Tags)
	_, err := r.pool.Exec(ctx,
		`INSERT INTO knowledge_entries (id, user_id, title, content, tags, category, source, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		entry.ID, entry.UserID, entry.Title, entry.Content, string(tags), entry.Category, entry.Source, entry.CreatedAt, entry.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("save knowledge entry: %w", err)
	}
	return nil
}

func (r *KnowledgeRepository) FindByID(ctx context.Context, id string) (*model.KnowledgeEntry, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, user_id, title, content, tags, category, source, created_at, updated_at
		 FROM knowledge_entries WHERE id = $1`, id,
	)
	var entry model.KnowledgeEntry
	var tagsStr string
	err := row.Scan(&entry.ID, &entry.UserID, &entry.Title, &entry.Content, &tagsStr, &entry.Category, &entry.Source, &entry.CreatedAt, &entry.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("find knowledge entry: %w", err)
	}
	json.Unmarshal([]byte(tagsStr), &entry.Tags)
	return &entry, nil
}

func (r *KnowledgeRepository) FindByUserID(ctx context.Context, userID string, offset, limit int) ([]*model.KnowledgeEntry, int, error) {
	var total int
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM knowledge_entries WHERE user_id = $1`, userID).Scan(&total)

	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, title, content, tags, category, source, created_at, updated_at
		 FROM knowledge_entries WHERE user_id = $1
		 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, userID, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list knowledge entries: %w", err)
	}
	defer rows.Close()

	var entries []*model.KnowledgeEntry
	for rows.Next() {
		var entry model.KnowledgeEntry
		var tagsStr string
		if err := rows.Scan(&entry.ID, &entry.UserID, &entry.Title, &entry.Content, &tagsStr, &entry.Category, &entry.Source, &entry.CreatedAt, &entry.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan knowledge entry: %w", err)
		}
		json.Unmarshal([]byte(tagsStr), &entry.Tags)
		entries = append(entries, &entry)
	}
	return entries, total, nil
}

func (r *KnowledgeRepository) Search(ctx context.Context, query string, tags []string, limit int) ([]*model.KnowledgeSearchResult, error) {
	where := []string{"1=1"}
	args := []any{}
	argIdx := 1

	if query != "" {
		where = append(where, fmt.Sprintf("(to_tsvector('english', title || ' ' || content) @@ plainto_tsquery('english', $%d) OR title ILIKE $%d OR content ILIKE $%d)", argIdx, argIdx+1, argIdx+1))
		likeQuery := "%" + query + "%"
		args = append(args, query, likeQuery)
		argIdx += 2
	}

	if len(tags) > 0 {
		placeholders := []string{}
		for _, tag := range tags {
			placeholders = append(placeholders, fmt.Sprintf("$%d", argIdx))
			args = append(args, tag)
			argIdx++
		}
		where = append(where, fmt.Sprintf("tags && ARRAY[%s]", strings.Join(placeholders, ",")))
	}

	sql := fmt.Sprintf(
		`SELECT id, user_id, title, content, tags, category, source, created_at, updated_at
		 FROM knowledge_entries WHERE %s
		 ORDER BY created_at DESC LIMIT $%d`,
		strings.Join(where, " AND "), argIdx,
	)
	args = append(args, limit)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("search knowledge: %w", err)
	}
	defer rows.Close()

	var results []*model.KnowledgeSearchResult
	for rows.Next() {
		var entry model.KnowledgeEntry
		var tagsStr string
		if err := rows.Scan(&entry.ID, &entry.UserID, &entry.Title, &entry.Content, &tagsStr, &entry.Category, &entry.Source, &entry.CreatedAt, &entry.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan search result: %w", err)
		}
		json.Unmarshal([]byte(tagsStr), &entry.Tags)
		results = append(results, &model.KnowledgeSearchResult{Entry: entry, Score: 1.0})
	}
	return results, nil
}

func (r *KnowledgeRepository) Update(ctx context.Context, entry *model.KnowledgeEntry) error {
	tags, _ := json.Marshal(entry.Tags)
	_, err := r.pool.Exec(ctx,
		`UPDATE knowledge_entries SET title=$1, content=$2, tags=$3, category=$4, source=$5, updated_at=NOW() WHERE id=$6`,
		entry.Title, entry.Content, string(tags), entry.Category, entry.Source, entry.ID,
	)
	return err
}

func (r *KnowledgeRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM knowledge_entries WHERE id=$1`, id)
	return err
}

func (r *KnowledgeRepository) ListCategories(ctx context.Context) ([]string, error) {
	rows, err := r.pool.Query(ctx, `SELECT DISTINCT category FROM knowledge_entries WHERE category != '' ORDER BY category`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var cat string
		if err := rows.Scan(&cat); err != nil {
			continue
		}
		categories = append(categories, cat)
	}
	return categories, nil
}
