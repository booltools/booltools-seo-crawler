package analyzer

type Registry struct {
	pageCheckers []PageRuleChecker
	siteCheckers []SiteRuleChecker
}

func NewRegistry() *Registry {
	return &Registry{
		pageCheckers: make([]PageRuleChecker, 0),
		siteCheckers: make([]SiteRuleChecker, 0),
	}
}

func (r *Registry) RegisterPageChecker(checker PageRuleChecker) {
	r.pageCheckers = append(r.pageCheckers, checker)
}

func (r *Registry) RegisterSiteChecker(checker SiteRuleChecker) {
	r.siteCheckers = append(r.siteCheckers, checker)
}

func (r *Registry) PageCheckers() []PageRuleChecker {
	return r.pageCheckers
}

func (r *Registry) SiteCheckers() []SiteRuleChecker {
	return r.siteCheckers
}
