package postgres

import (
	"context"
	"fmt"
	"match-me-api/internal/domain"
	"math/rand"

	"gorm.io/gorm/clause"
)

type cityLocation struct {
	Name string
	Lat  float64
	Lon  float64
}

const (
	defaultPasswordHash = "$2a$12$15dw2.nyH6xOf10DjQezcOIDY.PL.Jkr6ZjJOjpmqcL3xHtVeTWIq" // 12345678
	dummyUsersAmount    = 1000
)

func (s *Storage) SeedDummyUsers(ctx context.Context) error {
	const op = "storage.postgres.SeedDummyUsers"

	// Seed users
	var count int64
	if err := s.db.WithContext(ctx).Model(&domain.Account{}).Count(&count).Error; err != nil {
		return fmt.Errorf("%s: failed to count existing accounts: %w", op, err)
	}

	if count > 0 {
		return nil
	}

	users := generateDummyUsers(dummyUsersAmount)
	const batchLimit = 50

	if err := s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoNothing: true,
	}).CreateInBatches(&users, batchLimit).Error; err != nil {
		return fmt.Errorf("%s: failed to seed users: %w", op, err)
	}

	connections := generateDummyConnections()
	if err := s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "from_user_id"}, {Name: "to_user_id"}},
		DoNothing: true,
	}).Create(&connections).Error; err != nil {
		return fmt.Errorf("%s: failed to seed connections: %w", op, err)
	}

	return nil
}

// generateSeedUsers creates fake users for testing purposes
func generateDummyUsers(n int) []domain.Account {
	users := make([]domain.Account, 0, n)

	obiwan := domain.Account{
		Email:        "obiwan@matchme.com",
		PasswordHash: defaultPasswordHash, // 12345678
		Profile: domain.Profile{
			Email: "obiwan@matchme.com",
			Name:  "Obi-Wan Kenobi",
			Age:   35,
			Bio: `A disciplined mind, a patient approach, and a good cup of tea are my essentials.
			I value loyalty, strategy, and staying calm in chaos. Always down for a witty debate or a long walk.`,
			MaxRadius:       30,
			InteractionMode: domain.Social,
			Activities: []domain.ProfileActivity{
				{
					ActivityID:    1,                   // Running
					Experience:    domain.Intermediate, // 1-5 levels: "Beginner", "Active Novice", "Intermediate", "Advanced", "Professional"
					InterestLevel: domain.ActivelyLook, // 1-5: "Not interested", "Open to it" , "Interested" , "Highly interested", "Actively looking"
				},
				{
					ActivityID:    14,                  // Aikido
					Experience:    domain.Professional, // 1-5 levels: "Beginner", "Active Novice", "Intermediate", "Advanced", "Professional"
					InterestLevel: domain.Highly,       // 1-5: "Not interested", "Open to it" , "Interested" , "Highly interested", "Actively looking"
				},
				{
					ActivityID:    9, // Yoga
					Experience:    domain.Advanced,
					InterestLevel: domain.Interested,
				},
			},
			Lat: 60.170856,
			Lon: 24.941492,
		},
	}

	anakin := domain.Account{
		Email:        "anakin@matchme.com",
		PasswordHash: defaultPasswordHash, // 12345678
		Profile: domain.Profile{
			Email: "anakin@matchme.com",
			Name:  "Anakin Skywalker",
			Age:   19,
			Bio: `I live for speed, high stakes, and pushing limits.
			I trust my gut, speak my mind, and never back down from a challenge.
			If it's fast, intense, or "impossible", count me in.`,
			MaxRadius:       50,
			InteractionMode: 2,
			Activities: []domain.ProfileActivity{
				{
					ActivityID:    8, // CrossFit
					Experience:    domain.Advanced,
					InterestLevel: domain.ActivelyLook,
				},
				{
					ActivityID:    13, // MMA
					Experience:    domain.Advanced,
					InterestLevel: domain.Highly,
				},
				{
					ActivityID:    1, // Running
					Experience:    domain.Intermediate,
					InterestLevel: domain.ActivelyLook,
				},
			},
			Lat: 60.204758,
			Lon: 24.656968,
		},
	}

	yoda := domain.Account{
		Email:        "yoda@matchme.com",
		PasswordHash: defaultPasswordHash,
		Profile: domain.Profile{
			Email:           "yoda@matchme.com",
			Name:            "Master Yoda",
			Age:             896,
			Bio:             `Patience you must have. Long path wisdom is. Mindful of the present moment we stay. Meditate and train body and mind I like. Learn continuously, teach others we must.`,
			MaxRadius:       15,
			InteractionMode: domain.Silent,
			Activities: []domain.ProfileActivity{
				{
					ActivityID:    9, // Yoga
					Experience:    domain.Professional,
					InterestLevel: domain.ActivelyLook,
				},
				{
					ActivityID:    14, // Aikido
					Experience:    domain.Professional,
					InterestLevel: domain.Highly,
				},
			},
			Lat: 60.165000,
			Lon: 24.935000,
		},
	}

	windu := domain.Account{
		Email:        "windu@matchme.com",
		PasswordHash: defaultPasswordHash,
		Profile: domain.Profile{
			Email:           "windu@matchme.com",
			Name:            "Mace Windu",
			Age:             53,
			Bio:             `Strict discipline, unwavering principles, and zero tolerance for nonsense. I look for determination and high endurance. Action speaks louder than words.`,
			MaxRadius:       25,
			InteractionMode: domain.Silent,
			Activities: []domain.ProfileActivity{
				{
					ActivityID:    13, // MMA
					Experience:    domain.Professional,
					InterestLevel: domain.ActivelyLook,
				},
				{
					ActivityID:    3, // Gym
					Experience:    domain.Advanced,
					InterestLevel: domain.Highly,
				},
			},
			Lat: 60.180000,
			Lon: 24.950000,
		},
	}

	dooku := domain.Account{
		Email:        "dooku@matchme.com",
		PasswordHash: defaultPasswordHash,
		Profile: domain.Profile{
			Email:           "dooku@matchme.com",
			Name:            "Count Dooku",
			Age:             81,
			Bio:             `Elegance, precision, and refined taste. I appreciate tactical mastery, high-level sportsmanship, and discipline. Mediocrity does not interest me.`,
			MaxRadius:       40,
			InteractionMode: domain.Social,
			Activities: []domain.ProfileActivity{
				{
					ActivityID:    14, // Aikido
					Experience:    domain.Professional,
					InterestLevel: domain.ActivelyLook,
				},
				{
					ActivityID:    2, // Padel
					Experience:    domain.Advanced,
					InterestLevel: domain.Interested,
				},
			},
			Lat: 60.150000,
			Lon: 24.890000,
		},
	}

	maul := domain.Account{
		Email:        "maul@matchme.com",
		PasswordHash: defaultPasswordHash,
		Profile: domain.Profile{
			Email:           "maul@matchme.com",
			Name:            "Darth Maul",
			Age:             34,
			Bio:             `Focused entirely on victory and physical mastery. Intense training sessions only. If you can't keep up with the pace, don't waste my time. Rebuilding stronger every day.`,
			MaxRadius:       60,
			InteractionMode: domain.Silent,
			Activities: []domain.ProfileActivity{
				{
					ActivityID:    15, // Jiu-Jitsu
					Experience:    domain.Professional,
					InterestLevel: domain.ActivelyLook,
				},
				{
					ActivityID:    8, // CrossFit
					Experience:    domain.Professional,
					InterestLevel: domain.Highly,
				},
				{
					ActivityID:    11, // Climbing
					Experience:    domain.Advanced,
					InterestLevel: domain.Interested,
				},
			},
			Lat: 60.210000,
			Lon: 25.030000,
		},
	}

	ventress := domain.Account{
		Email:        "ventress@matchme.com",
		PasswordHash: defaultPasswordHash,
		Profile: domain.Profile{
			Email:           "ventress@matchme.com",
			Name:            "Asajj Ventress",
			Age:             28,
			Bio:             `Independent, sharp-witted, and unpredictable. I value freedom and strength. Looking for sparring partners or companions who aren't afraid of taking risks.`,
			MaxRadius:       35,
			InteractionMode: domain.OpenToAnything,
			Activities: []domain.ProfileActivity{
				{
					ActivityID:    13, // MMA
					Experience:    domain.Advanced,
					InterestLevel: domain.Highly,
				},
				{
					ActivityID:    11, // Climbing
					Experience:    domain.Intermediate,
					InterestLevel: domain.ActivelyLook,
				},
			},
			Lat: 60.190000,
			Lon: 24.910000,
		},
	}

	ahsoka := domain.Account{
		Email:        "ahsoka@matchme.com",
		PasswordHash: defaultPasswordHash,
		Profile: domain.Profile{
			Email:           "ahsoka@matchme.com",
			Name:            "Ahsoka Tano",
			Age:             18,
			Bio:             `Always learning, quick on my feet, and always ready to help. I love outdoor activities, staying active, and meeting genuine people.`,
			MaxRadius:       30,
			InteractionMode: domain.Social,
			Activities: []domain.ProfileActivity{
				{
					ActivityID:    1, // Running
					Experience:    domain.Advanced,
					InterestLevel: domain.ActivelyLook,
				},
				{
					ActivityID:    4, // Cycling
					Experience:    domain.Intermediate,
					InterestLevel: domain.Highly,
				},
				{
					ActivityID:    15, // Jiu-Jitsu
					Experience:    domain.Professional,
					InterestLevel: domain.ActivelyLook,
				},
				{
					ActivityID:    8, // CrossFit
					Experience:    domain.Professional,
					InterestLevel: domain.Highly,
				},
			},
			Lat: 60.175000,
			Lon: 24.925000,
		},
	}

	jarjar := domain.Account{
		Email:        "jarjar@matchme.com",
		PasswordHash: defaultPasswordHash,
		Profile: domain.Profile{
			Email:           "jarjar@matchme.com",
			Name:            "Jar Jar Binks",
			Age:             25,
			Bio:             `Mesa loves making new friends! Mesa big energy, super friendly, and loves swimming and playing ball. Sometimes mesa clumsy, but mesa heart is in right place!`,
			MaxRadius:       100,
			InteractionMode: domain.Social,
			Activities: []domain.ProfileActivity{
				{
					ActivityID:    7, // Swimming
					Experience:    domain.ActiveNovice,
					InterestLevel: domain.ActivelyLook,
				},
				{
					ActivityID:    1, // Running
					Experience:    domain.ActiveNovice,
					InterestLevel: domain.ActivelyLook,
				},
				{
					ActivityID:    5, // Football
					Experience:    domain.Beginner,
					InterestLevel: domain.OpenToIt,
				},
			},
			Lat: 60.140000,
			Lon: 24.960000,
		},
	}

	users = append(users, obiwan, anakin, yoda, windu, dooku, maul, ventress, ahsoka, jarjar)

	// TODO: add 100 randomly generated users
	finlandCities := []cityLocation{
		{Name: "Helsinki", Lat: 60.1699, Lon: 24.9384},
		{Name: "Tampere", Lat: 61.4978, Lon: 23.7610},
		{Name: "Turku", Lat: 60.4518, Lon: 22.2666},
		{Name: "Oulu", Lat: 65.0121, Lon: 25.4651},
		{Name: "Jyväskylä", Lat: 62.2426, Lon: 25.7473},
		{Name: "Lahti", Lat: 60.9827, Lon: 25.6612},
		{Name: "Kuopio", Lat: 62.8924, Lon: 27.6782},
		{Name: "Lappeenranta", Lat: 61.0587, Lon: 28.1887},
		{Name: "Joensuu", Lat: 62.6010, Lon: 29.7636},
		{Name: "Vaasa", Lat: 63.0951, Lon: 21.6165},
		{Name: "Mikkeli", Lat: 61.6887, Lon: 27.2721},
	}

	randomUsersCount := n - len(users)
	for i := 1; i <= randomUsersCount; i++ {
		// Random city
		city := finlandCities[rand.Intn(len(finlandCities))]

		// coordinates shift to add variability
		latOffset := (rand.Float64() - 0.5) * 0.15
		lonOffset := (rand.Float64() - 0.5) * 0.25

		// 1-2 random activities
		activitiesCount := 1 + rand.Intn(2)
		userActivities := make([]domain.ProfileActivity, 0, activitiesCount)
		usedActivityIDs := make(map[int64]bool)

		for len(userActivities) < activitiesCount {
			activityID := int64(1 + rand.Intn(15))
			if !usedActivityIDs[activityID] {
				usedActivityIDs[activityID] = true
				userActivities = append(userActivities, domain.ProfileActivity{
					ActivityID:    activityID,
					Experience:    domain.ExperienceLevel(1 + rand.Intn(5)),
					InterestLevel: domain.InterestLevel(1 + rand.Intn(5)),
				})
			}
		}

		email := fmt.Sprintf("user%d@matchme.com", i)

		dummyAcc := domain.Account{
			Email:        email,
			PasswordHash: defaultPasswordHash,
			Profile: domain.Profile{
				Email:           email,
				Name:            fmt.Sprintf("Athlete %d (%s)", i, city.Name),
				Age:             18 + rand.Intn(43), // 18 - 60
				Bio:             fmt.Sprintf("Hi! I live in %s and love staying active. Looking for sports partners!", city.Name),
				MaxRadius:       float64(10 + rand.Intn(91)), // 10 - 100 km (10 + [0..90])
				InteractionMode: domain.InteractionMode(1 + rand.Intn(3)),
				Activities:      userActivities,
				Lat:             city.Lat + latOffset,
				Lon:             city.Lon + lonOffset,
			},
		}

		users = append(users, dummyAcc)
	}

	return users
}

func generateDummyConnections() []domain.Connection {
	return []domain.Connection{
		// The Jedi
		{FromUserID: 1, ToUserID: 2, Status: domain.Accepted}, // Obi-Wan <-> Anakin
		{FromUserID: 1, ToUserID: 3, Status: domain.Accepted}, // Obi-Wan <-> Yoda
		{FromUserID: 3, ToUserID: 4, Status: domain.Accepted}, // Yoda <-> Windu
		{FromUserID: 2, ToUserID: 8, Status: domain.Accepted}, // Anakin <-> Ahsoka

		// Pending
		{FromUserID: 9, ToUserID: 1, Status: domain.Pending}, // Jar Jar -> Obi-Wan
		{FromUserID: 8, ToUserID: 3, Status: domain.Pending}, // Ahsoka -> Yoda

		// Declined
		{FromUserID: 1, ToUserID: 6, Status: domain.Declined}, // Obi-Wan <-> Maul
		{FromUserID: 2, ToUserID: 5, Status: domain.Declined}, // Anakin <-> Dooku

		// The Sith
		{FromUserID: 5, ToUserID: 7, Status: domain.Accepted}, // Dooku <-> Ventress
		{FromUserID: 6, ToUserID: 5, Status: domain.Pending},  // Maul -> Dooku
	}
}
