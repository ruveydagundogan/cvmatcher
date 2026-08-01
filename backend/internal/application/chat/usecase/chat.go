package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	chatmodel "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/chat/model"
	chatrepo "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/chat/repository"
	cvmodel "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/cv/model"
	cvrepo "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/cv/repository"
	jdrepo "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/jobdescription/repository"
	matchmodel "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/matching/model"
	matchrepo "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/matching/repository"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/infrastructure/llm"
	apperrors "github.com/ruveydagundogan/cvmatcher/backend/internal/shared/errors"
)

const (
	maxHistoryMessages = 10
	maxChatTokens      = 1024
)

type ConversationDetail struct {
	Conversation *chatmodel.Conversation `json:"conversation"`
	Messages     []*chatmodel.Message    `json:"messages"`
}

type SendMessageResult struct {
	Conversation *chatmodel.Conversation `json:"conversation"`
	UserMessage  *chatmodel.Message      `json:"user_message"`
	Assistant    *chatmodel.Message      `json:"assistant_message"`
}

type ChatUseCase struct {
	chatRepo   chatrepo.ChatRepository
	cvRepo     cvrepo.CVRepository
	jdRepo     jdrepo.JobDescriptionRepository
	matchRepo  matchrepo.MatchingRepository
	llmClient  *llm.Client
	log        *slog.Logger
}

func NewChatUseCase(chatRepo chatrepo.ChatRepository, cvRepo cvrepo.CVRepository, jdRepo jdrepo.JobDescriptionRepository, matchRepo matchrepo.MatchingRepository, llmClient *llm.Client, log *slog.Logger) *ChatUseCase {
	return &ChatUseCase{chatRepo: chatRepo, cvRepo: cvRepo, jdRepo: jdRepo, matchRepo: matchRepo, llmClient: llmClient, log: log}
}

func (uc *ChatUseCase) CreateConversation(ctx context.Context, userID, title, cvID, jdID, matchID string) (*chatmodel.Conversation, error) {
	// A match reference carries its own CV and JD; if the user only provides a
	// match, derive the other references from it.
	if matchID != "" {
		match, err := uc.matchRepo.FindByID(ctx, matchID)
		if err != nil {
			return nil, apperrors.Internal("failed to load match", err)
		}
		if match == nil {
			return nil, apperrors.NotFound("match not found")
		}
		if match.UserID != userID {
			return nil, apperrors.Forbidden("match does not belong to user")
		}
		if cvID == "" {
			cvID = match.CVID
		}
		if jdID == "" {
			jdID = match.JDID
		}
	}

	if cvID != "" {
		cv, err := uc.cvRepo.FindByID(ctx, cvID)
		if err != nil {
			return nil, apperrors.Internal("failed to load cv", err)
		}
		if cv == nil {
			return nil, apperrors.NotFound("cv not found")
		}
		if cv.UserID != userID {
			return nil, apperrors.Forbidden("cv does not belong to user")
		}
	}

	if jdID != "" {
		jd, err := uc.jdRepo.FindByID(ctx, jdID)
		if err != nil {
			return nil, apperrors.Internal("failed to load job description", err)
		}
		if jd == nil {
			return nil, apperrors.NotFound("job description not found")
		}
		if jd.UserID != userID {
			return nil, apperrors.Forbidden("job description does not belong to user")
		}
	}

	conv := chatmodel.NewConversation(userID, title, cvID, jdID, matchID)
	if err := uc.chatRepo.SaveConversation(ctx, conv); err != nil {
		return nil, apperrors.Internal("failed to save conversation", err)
	}
	return conv, nil
}

func (uc *ChatUseCase) ListConversations(ctx context.Context, userID string) ([]*chatmodel.Conversation, error) {
	return uc.chatRepo.ListByUser(ctx, userID)
}

func (uc *ChatUseCase) GetConversation(ctx context.Context, userID, convID string) (*ConversationDetail, error) {
	conv, err := uc.chatRepo.FindByID(ctx, convID)
	if err != nil {
		return nil, apperrors.Internal("failed to load conversation", err)
	}
	if conv == nil {
		return nil, apperrors.NotFound("conversation not found")
	}
	if conv.UserID != userID {
		return nil, apperrors.Forbidden("conversation does not belong to user")
	}
	msgs, err := uc.chatRepo.ListMessages(ctx, convID)
	if err != nil {
		return nil, apperrors.Internal("failed to load messages", err)
	}
	return &ConversationDetail{Conversation: conv, Messages: msgs}, nil
}

func (uc *ChatUseCase) DeleteConversation(ctx context.Context, userID, convID string) error {
	conv, err := uc.chatRepo.FindByID(ctx, convID)
	if err != nil {
		return apperrors.Internal("failed to load conversation", err)
	}
	if conv == nil {
		return apperrors.NotFound("conversation not found")
	}
	if conv.UserID != userID {
		return apperrors.Forbidden("conversation does not belong to user")
	}
	return uc.chatRepo.DeleteConversation(ctx, convID)
}

func (uc *ChatUseCase) SendMessage(ctx context.Context, userID, convID, content string) (*SendMessageResult, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, apperrors.Validation("message is empty")
	}

	conv, err := uc.chatRepo.FindByID(ctx, convID)
	if err != nil {
		return nil, apperrors.Internal("failed to load conversation", err)
	}
	if conv == nil {
		return nil, apperrors.NotFound("conversation not found")
	}
	if conv.UserID != userID {
		return nil, apperrors.Forbidden("conversation does not belong to user")
	}

	// Save user message
	userMsg := chatmodel.NewMessage(convID, "user", content, 0)
	if err := uc.chatRepo.SaveMessage(ctx, userMsg); err != nil {
		return nil, apperrors.Internal("failed to save message", err)
	}

	// Load history for context
	history, err := uc.chatRepo.ListMessages(ctx, convID)
	if err != nil {
		return nil, apperrors.Internal("failed to load history", err)
	}

	// Build messages: system (with CV/JD/match context) + last N history + new user message
	chatMessages := []llm.ChatMessage{
		{Role: "system", Content: uc.buildSystemPrompt(ctx, userID, conv.CVID, conv.JDID, conv.MatchID)},
	}
	start := 0
	if len(history) > maxHistoryMessages {
		start = len(history) - maxHistoryMessages
	}
	for _, m := range history[start:] {
		chatMessages = append(chatMessages, llm.ChatMessage{Role: m.Role, Content: m.Content})
	}
	chatMessages = append(chatMessages, llm.ChatMessage{Role: "user", Content: content})

	startTime := time.Now()
	response, err := uc.llmClient.ChatCompletionWithModel(ctx, llm.ModelCVCoach, chatMessages, maxChatTokens)
	durationMs := time.Since(startTime).Milliseconds()
	if err != nil {
		uc.log.Warn("cv-coach chat failed, falling back to base model", "error", err)
		response, err = uc.llmClient.ChatCompletionWithModel(ctx, llm.ModelBase, chatMessages, maxChatTokens)
		if err != nil {
			return nil, apperrors.Internal("LLM request failed", err)
		}
	}

	assistantMsg := chatmodel.NewMessage(convID, "assistant", response, 0)
	if err := uc.chatRepo.SaveMessage(ctx, assistantMsg); err != nil {
		return nil, apperrors.Internal("failed to save assistant message", err)
	}

	// Auto-title from first user message
	if conv.Title == "New Chat" && len(history) <= 1 {
		conv.Title = deriveTitle(content)
	}

	if err := uc.chatRepo.TouchConversation(ctx, convID); err != nil {
		uc.log.Warn("failed to touch conversation", "error", err)
	}

	uc.log.Info("chat message sent", "conversation_id", convID, "duration_ms", durationMs, "model_used", "cv-coach/fallback")

	return &SendMessageResult{
		Conversation: conv,
		UserMessage:  userMsg,
		Assistant:    assistantMsg,
	}, nil
}

func (uc *ChatUseCase) buildSystemPrompt(ctx context.Context, userID, cvID, jdID, matchID string) string {
	var b strings.Builder
	b.WriteString("You are the CV Coach, an expert career assistant. Help the user improve their CV, prepare for interviews, and tailor applications. Be concrete, practical and encouraging. Use short paragraphs and give examples with numbers when possible.")

	var match *matchmodel.MatchResult
	if matchID != "" {
		m, err := uc.matchRepo.FindByID(ctx, matchID)
		if err == nil && m != nil && m.UserID == userID {
			match = m
			if cvID == "" {
				cvID = m.CVID
			}
			if jdID == "" {
				jdID = m.JDID
			}
		}
	}

	// Load JD context if linked
	if jdID != "" {
		jd, err := uc.jdRepo.FindByID(ctx, jdID)
		if err == nil && jd != nil && jd.UserID == userID {
			b.WriteString("\n\nThe target job description (the user is applying to this role):\n")
			b.WriteString(fmt.Sprintf("Title: %s\n", jd.Title))
			if jd.Content != "" {
				b.WriteString(fmt.Sprintf("Description: %s\n", jd.Content))
			}
			if len(jd.RequiredSkills) > 0 {
				b.WriteString(fmt.Sprintf("Required skills: %s\n", strings.Join(jd.RequiredSkills, ", ")))
			}
			if len(jd.PreferredSkills) > 0 {
				b.WriteString(fmt.Sprintf("Preferred skills: %s\n", strings.Join(jd.PreferredSkills, ", ")))
			}
			if jd.ExperienceLevel != "" {
				b.WriteString(fmt.Sprintf("Experience level: %s\n", jd.ExperienceLevel))
			}
		}
	}

	// Load match context if linked
	if match != nil {
		b.WriteString("\n\nThe AI match result between the user's CV and this job description (the user wants to understand and improve this score):\n")
		b.WriteString(fmt.Sprintf("Overall score: %.0f/100\n", match.OverallScore))
		b.WriteString(fmt.Sprintf("Skill match: %.0f/100, Experience: %.0f/100, Education: %.0f/100\n", match.SkillMatchScore, match.ExperienceScore, match.EducationScore))
		if len(match.MatchedSkills) > 0 {
			b.WriteString(fmt.Sprintf("Matched skills: %s\n", strings.Join(match.MatchedSkills, ", ")))
		}
		if len(match.MissingSkills) > 0 {
			b.WriteString(fmt.Sprintf("Missing skills: %s\n", strings.Join(match.MissingSkills, ", ")))
		}
		if match.LLMAnalysis != "" {
			b.WriteString(fmt.Sprintf("Analysis: %s\n", match.LLMAnalysis))
		}
	}

	// Load CV context if linked
	var cv *cvmodel.CV
	var err error
	if cvID != "" {
		cv, err = uc.cvRepo.FindByID(ctx, cvID)
	} else {
		cv, err = uc.findLatestCV(ctx, userID)
	}
	if err == nil && cv != nil && cv.UserID == userID {
		b.WriteString("\n\nThe user's CV context (use it to give personalized advice):\n")
		b.WriteString(fmt.Sprintf("Title: %s\n", cv.Title))
		if len(cv.ParsedSkills) > 0 {
			b.WriteString(fmt.Sprintf("Skills: %s\n", strings.Join(cv.ParsedSkills, ", ")))
		}
		if cv.ParsedSummary != "" {
			b.WriteString(fmt.Sprintf("Summary: %s\n", cv.ParsedSummary))
		}
		if len(cv.ParsedExperience) > 0 {
			b.WriteString("Experience:\n")
			for _, e := range cv.ParsedExperience {
				b.WriteString(fmt.Sprintf("- %s @ %s (%s - %s)\n", e.Title, e.Company, e.StartDate, e.EndDate))
			}
		}
		if len(cv.ParsedEducation) > 0 {
			b.WriteString("Education:\n")
			for _, e := range cv.ParsedEducation {
				b.WriteString(fmt.Sprintf("- %s, %s\n", e.Degree, e.Institution))
			}
		}
	}
	return b.String()
}

func (uc *ChatUseCase) findLatestCV(ctx context.Context, userID string) (*cvmodel.CV, error) {
	cvs, _, err := uc.cvRepo.FindByUserID(ctx, userID, 0, 1)
	if err != nil {
		return nil, err
	}
	if len(cvs) == 0 {
		return nil, nil
	}
	return cvs[0], nil
}

func deriveTitle(content string) string {
	content = strings.TrimSpace(content)
	words := strings.Fields(content)
	if len(words) == 0 {
		return "New Chat"
	}
	title := strings.Join(words[:min(len(words), 8)], " ")
	if len(words) > 8 {
		title += "..."
	}
	if len(title) > 60 {
		title = title[:60] + "..."
	}
	return title
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
