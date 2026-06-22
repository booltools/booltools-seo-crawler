package handler

import (
	"net/http"

	"github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/analyzer"
)

type RulesHandler struct{}

func NewRulesHandler() *RulesHandler {
	return &RulesHandler{}
}

type RuleCategoryResponse struct {
	Key   string         `json:"key"`
	Label string         `json:"label"`
	Rules []RuleResponse `json:"rules"`
}

type RuleResponse struct {
	Key      string `json:"key"`
	Label    string `json:"label"`
	Severity string `json:"severity"`
}

type RulesListResponse struct {
	Categories []RuleCategoryResponse `json:"categories"`
	Presets    []PresetResponse       `json:"presets"`
}

type PresetResponse struct {
	Key         string   `json:"key"`
	Label       string   `json:"label"`
	Description string   `json:"description"`
	RuleKeys    []string `json:"ruleKeys"`
}

func (h *RulesHandler) ListRules(writer http.ResponseWriter, request *http.Request) {
	definitions := analyzer.AllRuleDefinitions()

	categoryOrder := []valueobject.Category{
		valueobject.CategoryOnPage,
		valueobject.CategoryContent,
		valueobject.CategoryTechnical,
		valueobject.CategoryLinks,
		valueobject.CategoryPerformance,
		valueobject.CategoryStructuredData,
		valueobject.CategorySecurity,
		valueobject.CategoryAccessibility,
		valueobject.CategorySocial,
		valueobject.CategoryMobile,
		valueobject.CategoryURLStructure,
		valueobject.CategoryInternationalization,
		valueobject.CategoryEEAT,
		valueobject.CategoryDuplicateContent,
		valueobject.CategoryGEO,
	}

	grouped := make(map[valueobject.Category][]RuleResponse)
	for _, definition := range definitions {
		grouped[definition.Category] = append(grouped[definition.Category], RuleResponse{
			Key:      definition.Key,
			Label:    definition.Label,
			Severity: definition.Severity,
		})
	}

	var categories []RuleCategoryResponse
	for _, category := range categoryOrder {
		rules, exists := grouped[category]
		if !exists {
			continue
		}
		categories = append(categories, RuleCategoryResponse{
			Key:   string(category),
			Label: category.Label(),
			Rules: rules,
		})
	}

	presetDefs := analyzer.AllPresets()
	presets := make([]PresetResponse, 0, len(presetDefs))
	for _, preset := range presetDefs {
		presets = append(presets, PresetResponse{
			Key:         preset.Key,
			Label:       preset.Label,
			Description: preset.Description,
			RuleKeys:    preset.RuleKeys,
		})
	}

	writeJSON(writer, http.StatusOK, RulesListResponse{
		Categories: categories,
		Presets:    presets,
	})
}
