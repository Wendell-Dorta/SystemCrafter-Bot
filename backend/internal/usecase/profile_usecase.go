package usecase

import (
	"context"
	"strings"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/repository"
)

// ProfileUseCase manages learned cross-session architectural memory.
type ProfileUseCase struct {
	profileRepo repository.ProfileRepository
}

// NewProfileUseCase creates an instance of ProfileUseCase.
func NewProfileUseCase(repo repository.ProfileRepository) *ProfileUseCase {
	return &ProfileUseCase{profileRepo: repo}
}

func (uc *ProfileUseCase) GetProfile(ctx context.Context) (*entity.ArchitectProfile, error) {
	return uc.profileRepo.GetProfile(ctx)
}

func (uc *ProfileUseCase) UpdateProfile(ctx context.Context, profile entity.ArchitectProfile) error {
	return uc.profileRepo.UpdateProfile(ctx, profile)
}

func (uc *ProfileUseCase) ResetProfile(ctx context.Context) error {
	return uc.profileRepo.ResetProfile(ctx)
}

func (uc *ProfileUseCase) AddNote(ctx context.Context, note string) error {
	return uc.profileRepo.AddNote(ctx, strings.TrimSpace(note))
}

// LearnFromMessage inspects user messages to automatically capture architectural preferences.
func (uc *ProfileUseCase) LearnFromMessage(ctx context.Context, text string) {
	if strings.TrimSpace(text) == "" {
		return
	}

	profile, err := uc.profileRepo.GetProfile(ctx)
	if err != nil || profile == nil {
		return
	}

	lower := strings.ToLower(text)
	updated := false

	// Detect Cloud Preference
	if strings.Contains(lower, "usamos aws") || strings.Contains(lower, "nossa nuvem é aws") || strings.Contains(lower, "preferimos aws") {
		if profile.PreferredCloud != "AWS" {
			profile.PreferredCloud = "AWS"
			updated = true
		}
	} else if strings.Contains(lower, "usamos gcp") || strings.Contains(lower, "google cloud") || strings.Contains(lower, "preferimos gcp") {
		if profile.PreferredCloud != "GCP" {
			profile.PreferredCloud = "GCP"
			updated = true
		}
	} else if strings.Contains(lower, "usamos azure") || strings.Contains(lower, "preferimos azure") {
		if profile.PreferredCloud != "Azure" {
			profile.PreferredCloud = "Azure"
			updated = true
		}
	}

	// Detect Languages
	langChecks := map[string]string{
		"go":         "Go",
		"golang":     "Go",
		"typescript": "TypeScript",
		"python":     "Python",
		"java":       "Java",
		"c#":         "C#",
		"rust":       "Rust",
	}
	for keyword, standardLang := range langChecks {
		if (strings.Contains(lower, "usamos "+keyword) || strings.Contains(lower, "preferimos "+keyword) || strings.Contains(lower, "nossa stack é "+keyword)) && !containsString(profile.PrimaryLanguages, standardLang) {
			profile.PrimaryLanguages = append(profile.PrimaryLanguages, standardLang)
			updated = true
		}
	}

	// Detect Databases
	dbChecks := map[string]string{
		"postgres":   "PostgreSQL",
		"postgresql": "PostgreSQL",
		"mongodb":    "MongoDB",
		"mongo":      "MongoDB",
		"mysql":      "MySQL",
		"redis":      "Redis",
		"dynamodb":   "DynamoDB",
		"clickhouse": "ClickHouse",
	}
	for keyword, standardDB := range dbChecks {
		if (strings.Contains(lower, "usamos "+keyword) || strings.Contains(lower, "banco "+keyword) || strings.Contains(lower, "preferimos "+keyword)) && !containsString(profile.PreferredDatabases, standardDB) {
			profile.PreferredDatabases = append(profile.PreferredDatabases, standardDB)
			updated = true
		}
	}

	// Detect Compliance
	complianceChecks := map[string]string{
		"lgpd":    "LGPD",
		"gdpr":    "GDPR",
		"pci":     "PCI-DSS",
		"pci-dss": "PCI-DSS",
		"hipaa":   "HIPAA",
		"soc2":    "SOC2",
	}
	for keyword, standardComp := range complianceChecks {
		if (strings.Contains(lower, keyword) && (strings.Contains(lower, "compliance") || strings.Contains(lower, "conformidade") || strings.Contains(lower, "obrigatório"))) && !containsString(profile.ComplianceRules, standardComp) {
			profile.ComplianceRules = append(profile.ComplianceRules, standardComp)
			updated = true
		}
	}

	if updated {
		_ = uc.profileRepo.UpdateProfile(ctx, *profile)
	}
}

func containsString(slice []string, val string) bool {
	for _, item := range slice {
		if strings.EqualFold(item, val) {
			return true
		}
	}
	return false
}
