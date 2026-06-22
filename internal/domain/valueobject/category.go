package valueobject

type Category string

const (
	CategoryTechnical            Category = "technical"
	CategoryOnPage               Category = "on_page"
	CategoryLinks                Category = "links"
	CategoryPerformance          Category = "performance"
	CategoryStructuredData       Category = "structured_data"
	CategorySecurity             Category = "security"
	CategoryAccessibility        Category = "accessibility"
	CategorySocial               Category = "social"
	CategoryMobile               Category = "mobile"
	CategoryURLStructure         Category = "url_structure"
	CategoryInternationalization Category = "internationalization"
	CategoryEEAT                 Category = "eeat"
	CategoryContent              Category = "content"
	CategoryDuplicateContent     Category = "duplicate_content"
	CategoryGEO                  Category = "geo"
)

func (c Category) IsValid() bool {
	switch c {
	case CategoryTechnical, CategoryOnPage, CategoryLinks, CategoryPerformance,
		CategoryStructuredData, CategorySecurity, CategoryAccessibility,
		CategorySocial, CategoryMobile, CategoryURLStructure,
		CategoryInternationalization, CategoryEEAT, CategoryContent,
		CategoryDuplicateContent, CategoryGEO:
		return true
	}
	return false
}

func (c Category) Label() string {
	labels := map[Category]string{
		CategoryTechnical:            "Technical SEO",
		CategoryOnPage:               "On-Page SEO",
		CategoryLinks:                "Links",
		CategoryPerformance:          "Performance",
		CategoryStructuredData:       "Structured Data",
		CategorySecurity:             "Security",
		CategoryAccessibility:        "Accessibility",
		CategorySocial:               "Social / Open Graph",
		CategoryMobile:               "Mobile",
		CategoryURLStructure:         "URL Structure",
		CategoryInternationalization: "Internationalization",
		CategoryEEAT:                 "E-E-A-T",
		CategoryContent:              "Content Quality",
		CategoryDuplicateContent:     "Duplicate Content",
		CategoryGEO:                  "GEO (AI Search)",
	}
	if label, exists := labels[c]; exists {
		return label
	}
	return string(c)
}

func AllCategories() []Category {
	return []Category{
		CategoryTechnical, CategoryOnPage, CategoryLinks, CategoryPerformance,
		CategoryStructuredData, CategorySecurity, CategoryAccessibility,
		CategorySocial, CategoryMobile, CategoryURLStructure,
		CategoryInternationalization, CategoryEEAT, CategoryContent,
		CategoryDuplicateContent, CategoryGEO,
	}
}

var CategoryWeights = map[Category]float64{
	CategoryTechnical:            0.15,
	CategoryOnPage:               0.13,
	CategoryPerformance:          0.13,
	CategoryLinks:                0.09,
	CategoryStructuredData:       0.06,
	CategorySecurity:             0.07,
	CategoryAccessibility:        0.04,
	CategorySocial:               0.04,
	CategoryMobile:               0.04,
	CategoryURLStructure:         0.03,
	CategoryInternationalization: 0.02,
	CategoryEEAT:                 0.05,
	CategoryContent:              0.10,
	CategoryDuplicateContent:     0.05,
}
