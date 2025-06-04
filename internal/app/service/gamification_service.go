package service

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type GamificationService interface {
	ApplyLessonCompletion(userID, lessonID uuid.UUID) error
}

type gamificationService struct {
	db *sql.DB
}

func NewGamificationService(db *sql.DB) GamificationService {
	return &gamificationService{db: db}
}

func (s *gamificationService) ApplyLessonCompletion(userID, lessonID uuid.UUID) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				fmt.Printf("rollback failed: %v\n", rollbackErr)
			}
		} else {
			_ = tx.Commit() // move Commit here, guarded
		}
	}()

	// 1. Get XP and crown reward from lesson
	var xpReward, crownsReward int
	err = tx.QueryRow(`
		SELECT xp_reward, crowns_reward
		FROM lessons
		WHERE id = $1
	`, lessonID).Scan(&xpReward, &crownsReward)
	if err != nil {
		return fmt.Errorf("lesson not found: %w", err)
	}

	// 2. Update user_stats (XP, streaks, hearts)
	_, err = tx.Exec(`
		INSERT INTO user_stats (user_id, total_xp, current_streak, max_streak, hearts, last_active)
		VALUES ($1, $2, 1, 1, 5, CURRENT_DATE)
		ON CONFLICT (user_id) DO UPDATE
		SET
			total_xp = user_stats.total_xp + EXCLUDED.total_xp,
			current_streak = CASE
				WHEN user_stats.last_active = CURRENT_DATE - INTERVAL '1 day' THEN user_stats.current_streak + 1
				WHEN user_stats.last_active = CURRENT_DATE THEN user_stats.current_streak
				ELSE 1
			END,
			max_streak = GREATEST(
				user_stats.max_streak,
				CASE
					WHEN user_stats.last_active = CURRENT_DATE - INTERVAL '1 day' THEN user_stats.current_streak + 1
					WHEN user_stats.last_active = CURRENT_DATE THEN user_stats.current_streak
					ELSE 1
				END
			),
			last_active = CURRENT_DATE
	`, userID, xpReward)
	if err != nil {
		return fmt.Errorf("failed to update user_stats: %w", err)
	}

	// 3. Update user_crowns (lesson mastery)
	_, err = tx.Exec(`
		INSERT INTO user_crowns (user_id, lesson_id, crown_level)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, lesson_id) DO UPDATE
		SET crown_level = user_crowns.crown_level + EXCLUDED.crown_level
	`, userID, lessonID, crownsReward)
	if err != nil {
		return fmt.Errorf("failed to update user_crowns: %w", err)
	}

	// 4. Update user_wallet (gems)
	_, err = tx.Exec(`
		INSERT INTO user_wallet (user_id, gems, updated_at)
		VALUES ($1, 5, now())
		ON CONFLICT (user_id) DO UPDATE
		SET gems = user_wallet.gems + 5, updated_at = now()
	`, userID)
	if err != nil {
		return fmt.Errorf("failed to update user_wallet: %w", err)
	}

	return tx.Commit()
}

