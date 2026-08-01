package chat

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	chatmodel "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/chat/model"
)

type ChatRepository struct {
	pool *pgxpool.Pool
}

func NewChatRepository(pool *pgxpool.Pool) *ChatRepository {
	return &ChatRepository{pool: pool}
}

func (r *ChatRepository) ListByUser(ctx context.Context, userID string) ([]*chatmodel.Conversation, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, title, cv_id, jd_id, match_id, created_at, updated_at
		 FROM conversations WHERE user_id=$1 ORDER BY updated_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	convs := make([]*chatmodel.Conversation, 0)
	for rows.Next() {
		c := &chatmodel.Conversation{}
		var cvID, jdID, matchID sql.NullString
		if err := rows.Scan(&c.ID, &c.UserID, &c.Title, &cvID, &jdID, &matchID, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		if cvID.Valid {
			c.CVID = cvID.String
		}
		if jdID.Valid {
			c.JDID = jdID.String
		}
		if matchID.Valid {
			c.MatchID = matchID.String
		}
		convs = append(convs, c)
	}
	return convs, rows.Err()
}

func (r *ChatRepository) FindByID(ctx context.Context, id string) (*chatmodel.Conversation, error) {
	c := &chatmodel.Conversation{}
	var cvID, jdID, matchID sql.NullString
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, title, cv_id, jd_id, match_id, created_at, updated_at FROM conversations WHERE id=$1`, id,
	).Scan(&c.ID, &c.UserID, &c.Title, &cvID, &jdID, &matchID, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if cvID.Valid {
		c.CVID = cvID.String
	}
	if jdID.Valid {
		c.JDID = jdID.String
	}
	if matchID.Valid {
		c.MatchID = matchID.String
	}
	return c, nil
}

func (r *ChatRepository) SaveConversation(ctx context.Context, conv *chatmodel.Conversation) error {
	var cvID, jdID, matchID any
	if conv.CVID != "" {
		cvID = conv.CVID
	}
	if conv.JDID != "" {
		jdID = conv.JDID
	}
	if conv.MatchID != "" {
		matchID = conv.MatchID
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO conversations (id, user_id, title, cv_id, jd_id, match_id, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		conv.ID, conv.UserID, conv.Title, cvID, jdID, matchID, conv.CreatedAt, conv.UpdatedAt,
	)
	return err
}

func (r *ChatRepository) DeleteConversation(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM conversations WHERE id=$1`, id)
	return err
}

func (r *ChatRepository) TouchConversation(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `UPDATE conversations SET updated_at=NOW() WHERE id=$1`, id)
	return err
}

func (r *ChatRepository) ListMessages(ctx context.Context, conversationID string) ([]*chatmodel.Message, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, conversation_id, role, content, token_count, created_at
		 FROM messages WHERE conversation_id=$1 ORDER BY created_at ASC`, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	msgs := make([]*chatmodel.Message, 0)
	for rows.Next() {
		m := &chatmodel.Message{}
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.Role, &m.Content, &m.TokenCount, &m.CreatedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, rows.Err()
}

func (r *ChatRepository) SaveMessage(ctx context.Context, msg *chatmodel.Message) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO messages (id, conversation_id, role, content, token_count, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		msg.ID, msg.ConversationID, msg.Role, msg.Content, msg.TokenCount, msg.CreatedAt,
	)
	return err
}
