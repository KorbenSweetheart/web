package postgres

import (
	"context"
	"fmt"
	"match-me-api/internal/domain"

	"gorm.io/gorm/clause"
)

func (s *Storage) SeedDummyUsers(ctx context.Context) error {
	const op = "storage.postgres.SeedDummyUsers"

	// Seed users
	var count int64
	if err := s.db.WithContext(ctx).Model(&domain.Account{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count existing accounts: op: %s, error: %w", op, err)
	}

	if count > 0 {
		return nil
	}

	users := generateDummyUsers()
	const batchLimit = 50

	if err := s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoNothing: true,
	}).CreateInBatches(&users, batchLimit).Error; err != nil {
		return fmt.Errorf("failed to seed users, op: %s, error: %w", op, err)
	}

	return nil
}

// generateSeedUsers creates fake users for testing purposes
func generateDummyUsers() []domain.Account {
	users := make([]domain.Account, 0, 100)

	user1 := domain.Account{
		Email:        "obiwan@matchme.com",
		PasswordHash: "$2a$12$15dw2.nyH6xOf10DjQezcOIDY.PL.Jkr6ZjJOjpmqcL3xHtVeTWIq", // 12345678
		Profile: domain.Profile{
			Name: "Obi-Wan Kenobi", // TODO: maybe add it during registration
			Age:  35,
			Bio: `A disciplined mind, a patient approach, and a good cup of tea are my essentials.
			I value loyalty, strategy, and staying calm in chaos. Always down for a witty debate or a long walk.`,
			MaxRadius:       30,
			InteractionMode: 2,
			Activities: []domain.ProfileActivity{
				{
					ActivityID:    1, // Running
					Experience:    3, // 1-5 levels: "Beginner", "Active Novice", "Intermediate", "Advanced", "Professional"
					InterestLevel: 5,
				},
				{
					ActivityID:    14,
					Experience:    5, // 1-5 levels: "Beginner", "Active Novice", "Intermediate", "Advanced", "Professional"
					InterestLevel: 4, // 1-5: "Not interested", "Open to it" , "Interested" , "Highly interested", "Actively looking"
				},
			},
			Lat: 60.170856,
			Lon: 24.941492,
		},
	}

	user2 := domain.Account{
		Email:        "anakin@matchme.com",
		PasswordHash: "$2a$12$15dw2.nyH6xOf10DjQezcOIDY.PL.Jkr6ZjJOjpmqcL3xHtVeTWIq", // 12345678
		Profile: domain.Profile{
			Name: "Anakin Skywalker", // TODO: maybe add it during registration
			Age:  19,
			Bio: `I live for speed, high stakes, and pushing limits.
			I trust my gut, speak my mind, and never back down from a challenge.
			If it's fast, intense, or "impossible", count me in.`,
			MaxRadius:       50,
			InteractionMode: 2,
			Activities: []domain.ProfileActivity{
				{
					ActivityID:    13, // Running
					Experience:    4,  // 1-5 levels: "Beginner", "Active Novice", "Intermediate", "Advanced", "Professional"
					InterestLevel: 4,
				},
				{
					ActivityID:    1,
					Experience:    3, // 1-5 levels: "Beginner", "Active Novice", "Intermediate", "Advanced", "Professional"
					InterestLevel: 5, // 1-5: "Not interested", "Open to it" , "Interested" , "Highly interested", "Actively looking"
				},
			},
			Lat: 60.204758,
			Lon: 24.656968,
		},
	}

	// TODO: add 100 randomly generated users

	users = append(users, user1, user2)

	return users
}
