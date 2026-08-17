package main

import (
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	authModel "github.com/bees/hindu-ritual-platform/internal/auth/model"
	panditModel "github.com/bees/hindu-ritual-platform/internal/pandits/model"
	ritualModel "github.com/bees/hindu-ritual-platform/internal/rituals/model"
	"github.com/bees/hindu-ritual-platform/pkg/security"
)

type panditSeed struct {
	Name           string
	Email          string
	Phone          string
	Bio            string
	Experience     int
	Specialization string
	BasePrice      float64
	ServiceArea    string
}

type ritualSeed struct {
	Name        string
	Slug        string
	Description string
	BasePrice   float64
	Duration    string
	Category    string
	Items       string
	Procedure   string
}

func main() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "ritual_user"),
		getEnv("DB_PASSWORD", "ritual_pass_secure"),
		getEnv("DB_NAME", "hindu_ritual_db"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_SSLMODE", "disable"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto")

	adminHash, _ := security.HashPassword("AdminPass123!")
	customerHash, _ := security.HashPassword("CustomerPass123!")
	panditHash, _ := security.HashPassword("PanditPass123!")

	admin := authModel.User{}
	db.Where("email = ?", "admin@bishalpujasewa.com").FirstOrCreate(&admin, authModel.User{
		Email: "admin@bishalpujasewa.com", PasswordHash: adminHash, FullName: "Admin",
		Role: "admin", IsEmailVerified: true, IsActive: true,
	})
	fmt.Println("Admin:", admin.Email)

	customer := authModel.User{}
	db.Where("email = ?", "customer@example.com").FirstOrCreate(&customer, authModel.User{
		Email: "customer@example.com", PasswordHash: customerHash, FullName: "Hari Prasad",
		Phone: "9876543210", Role: "customer", IsEmailVerified: true, IsActive: true,
	})
	fmt.Println("Customer:", customer.Email)

	categories := []ritualModel.RitualCategory{
		{Name: "Daily Fire Sacrifice", Slug: "daily-fire-sacrifice", Description: "Daily Agnihotra and fire-based rituals"},
		{Name: "Soma Sacrifices", Slug: "soma-sacrifices", Description: "Great Soma Yajna ceremonies"},
		{Name: "Royal Consecration", Slug: "royal-consecration", Description: "Royal sacrifices for kingship and prestige"},
		{Name: "Lunar Observances", Slug: "lunar-observances", Description: "Monthly and seasonal lunar sacrifices"},
		{Name: "Householder Duties", Slug: "householder-duties", Description: "Daily ethical and ritual obligations"},
		{Name: "Life Ceremonies", Slug: "life-ceremonies", Description: "Birth, marriage, and other life rituals"},
	}
	for _, c := range categories {
		existing := ritualModel.RitualCategory{}
		result := db.Where("slug = ?", c.Slug).First(&existing)
		if result.Error != nil {
			db.Create(&c)
		}
	}
	fmt.Println("Categories created")

	pandits := []panditSeed{
		{"Atri", "pandit.atri@example.com", "9841000001", "Descendant of the Atri gotra, specializing in Agnihotra and daily fire rituals with 20+ years of experience.", 20, "Agnihotra, Daily Fires, Griha Pravesh", 1800, "Kathmandu Valley"},
		{"Bharadvaja", "pandit.bharadvaja@example.com", "9841000002", "Expert in Soma Yajna and complex Vedic ceremonies. 18 years of experience in high-level ritual performances.", 18, "Soma Yajna, Agnicayana, Vivah", 2200, "Kathmandu Valley"},
		{"Gautama", "pandit.gautama@example.com", "9841000003", "Scholar of Gautama Dharmashastra, specializing in Griha Pravesh and householder rituals. 15 years experience.", 15, "Griha Pravesh, Sanskar, Puja", 1500, "Lalitpur"},
		{"Jamadagni", "pandit.jamadagni@example.com", "9841000004", "Expert in Aśvamedha and Rājasūya royal ceremonies. Senior priest with 25+ years of Vedic practice.", 25, "Aśvamedha, Rājasūya, Vājapeya", 3000, "Bhaktapur"},
		{"Kashyapa", "pandit.kashyapa@example.com", "9841000005", "Descendant of Kashyapa gotra. Specializes in Pañca Mahāyajñas and daily household rituals. 12 years experience.", 12, "Pañca Mahāyajñas, Daily Rituals, Shraddha", 1300, "Kathmandu Valley"},
		{"Vasistha", "pandit.vasistha@example.com", "9841000006", "Vasistha lineage priest specializing in Darśa-Pūrṇamāsa and lunar observances. Known for precise ritual timing.", 22, "Darśa-Pūrṇamāsa, Vrat, Lunar Rituals", 2000, "Patan"},
		{"Vishvamitra", "pandit.vishvamitra@example.com", "9841000007", "Legendary lineage expert in all major Vedic sacrifices. 30 years of experience across Nepal and India.", 30, "All Vedic Yajnas, Agnicayana, Soma", 3500, "Kathmandu Valley"},
	}

	for _, p := range pandits {
		user := authModel.User{}
		result := db.Where("email = ?", p.Email).First(&user)
		if result.Error != nil {
			user = authModel.User{
				Email: p.Email, PasswordHash: panditHash, FullName: p.Name,
				Phone: p.Phone, Role: "pandit", IsEmailVerified: true, IsActive: true,
			}
			db.Create(&user)
		}

		existingPandit := panditModel.Pandit{}
		if db.Where("user_id = ?", user.ID).First(&existingPandit).Error != nil {
			db.Create(&panditModel.Pandit{
				UserID: user.ID, Bio: p.Bio, ExperienceYears: p.Experience,
				Specialization: p.Specialization,
				Languages:      datatypes.NewJSONSlice([]string{"Nepali", "Sanskrit", "Hindi"}),
				BasePrice:      p.BasePrice, VerificationStatus: panditModel.VerificationApproved,
				IsAvailable: true, ServiceArea: p.ServiceArea, Rating: 4.8,
			})
		}
	}
	fmt.Println("Pandits created:", len(pandits))

	var catMap = make(map[string]string)
	var cats []ritualModel.RitualCategory
	db.Find(&cats)
	for _, c := range cats {
		catMap[c.Slug] = c.ID.String()
	}

	rituals := []ritualSeed{
		{"Agnihotra", "agnihotra", "The foundational daily fire sacrifice performed at sunrise and sunset. Sustains Ṛta (cosmic order) and purifies the household.", 1200, "1 hour", "daily-fire-sacrifice", "Ghee, Samidh (sacred firewood), Rice, Milk", "Purify the space → Light the fire → Recite Agni mantras → Offer ghee and rice → Conclude with peace chant"},
		{"Soma Yajña", "soma-yajna", "The defining great sacrifice of the Vedic religion. Pressing and offering of the sacred Soma plant to the Devas.", 5000, "5 days", "soma-sacrifices", "Soma plant, Ghee, Milk, Grains, Darbha grass", "Preliminary rites → Soma pressing → Offering to Devas → Feast → Concluding sacrifice"},
		{"Agnicayana", "agnicayana", "The most elaborate fire ritual — building the fire altar brick by brick as a symbolic re-creation of the universe.", 8000, "12 days", "soma-sacrifices", "Clay bricks, Gold piece, Ghee, Soma, Animal-shaped pots", "Select sacred space → Shape bricks → Lay altar in five layers → Install sacred fires → Offer oblations"},
		{"Darśa–Pūrṇamāsa", "darsa-purnamasa", "New Moon and Full Moon sacrifices that mark the lunar rhythm. Regular oblations sustaining cosmic and domestic harmony.", 2000, "1 day", "lunar-observances", "Ghee, Rice cakes (Purodāśa), Barley, Milk, Potsherds", "Fasting → Setup fires → Offer Purodāśa → Make oblations → Feed Brahmins"},
		{"Aśvamedha", "asvamedha", "The supreme royal horse sacrifice — a grand ritual demonstrating imperial sovereignty and territorial authority.", 50000, "1 year", "royal-consecration", "Sacrificial horse, Royal chariot, Gold, 100+ priests, Queens", "Horse consecration → Year-long wandering → Horse returns → Grand sacrifice → Queen's ritual"},
		{"Rājasūya", "rajasuya", "Royal consecration ceremony that establishes and legitimizes kingship. Includes the symbolic chariot drive and homage from vassals.", 25000, "2 years", "royal-consecration", "Royal throne, Chariot, Crown, Gold, Ghee, Rare herbs", "Preliminary homa → Anointing → Chariot drive → Vassal homage → Crown installation → Feast"},
		{"Vājapeya", "vajapeya", "The 'drink of strength' sacrifice — a royal ritual for prosperity, prestige, and vitality. Includes a symbolic chariot race.", 15000, "17 days", "royal-consecration", "Chariot, 17 sacrificial animals, Soma, Grains, Gold", "Setup 17 poles → Soma pressing → Chariot race → Royal anointing → Feast"},
		{"Pañca Mahāyajñas", "panca-mahayajnas", "The five great daily sacrifices: Deva Yajña, Pitṛ Yajña, Bhūta Yajña, Manuṣya Yajña, and Brahma Yajña — the householder's daily dharma.", 500, "daily", "householder-duties", "Ghee, Water, Grains, Food, Scriptures", "Brahma Yajña (study) → Deva Yajña (fire offering) → Pitṛ Yajña (ancestors) → Bhūta Yajña (creatures) → Manuṣya Yajña (guests)"},
	}

	for _, r := range rituals {
		existing := ritualModel.Ritual{}
		if db.Where("slug = ?", r.Slug).First(&existing).Error != nil {
			catID, _ := uuid.Parse(catMap[r.Category])
			db.Create(&ritualModel.Ritual{
				CategoryID: catID, Name: r.Name, Slug: r.Slug,
				Description: r.Description, Duration: r.Duration,
				BasePrice: r.BasePrice, IsActive: true,
				RequiredItems: r.Items, Procedure: r.Procedure,
			})
		}
	}
	fmt.Println("Rituals created:", len(rituals))

	fmt.Println("\n========================================")
	fmt.Println("  Seed data created successfully!")
	fmt.Println("========================================")
	fmt.Println("  Admin:    admin@bishalpujasewa.com / AdminPass123!")
	fmt.Println("  Pandit:   pandit.atri@example.com / PanditPass123!")
	fmt.Println("  Customer: customer@example.com / CustomerPass123!")
	fmt.Println("========================================")
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
