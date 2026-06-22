package analyzer

import "github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"

type RuleInfo struct {
	Key      string              `json:"key"`
	Category valueobject.Category `json:"category"`
	Label    string              `json:"label"`
	Severity string              `json:"severity"`
}

func AllRuleDefinitions() []RuleInfo {
	return []RuleInfo{
		// On-Page SEO
		{Key: "title_exists", Category: valueobject.CategoryOnPage, Label: "Page has a <title> tag", Severity: "critical"},
		{Key: "title_length", Category: valueobject.CategoryOnPage, Label: "Title is between 30-60 characters", Severity: "medium"},
		{Key: "title_multiple", Category: valueobject.CategoryOnPage, Label: "Page has only one title tag", Severity: "high"},
		{Key: "meta_description_exists", Category: valueobject.CategoryOnPage, Label: "Page has a meta description", Severity: "high"},
		{Key: "meta_description_length", Category: valueobject.CategoryOnPage, Label: "Description is 120-160 characters", Severity: "medium"},
		{Key: "meta_description_multiple", Category: valueobject.CategoryOnPage, Label: "Page has only one meta description", Severity: "medium"},
		{Key: "h1_count", Category: valueobject.CategoryOnPage, Label: "Page has exactly one H1", Severity: "high"},
		{Key: "h1_not_empty", Category: valueobject.CategoryOnPage, Label: "H1 has non-empty text", Severity: "high"},
		{Key: "heading_hierarchy", Category: valueobject.CategoryOnPage, Label: "Headings follow proper hierarchy", Severity: "medium"},
		{Key: "heading_not_empty", Category: valueobject.CategoryOnPage, Label: "Headings have non-empty text", Severity: "low"},
		{Key: "images_alt_text", Category: valueobject.CategoryOnPage, Label: "All images have alt attributes", Severity: "high"},
		{Key: "images_alt_descriptive", Category: valueobject.CategoryOnPage, Label: "Image alt text is descriptive", Severity: "low"},
		{Key: "images_dimensions", Category: valueobject.CategoryOnPage, Label: "Images have width/height attributes", Severity: "medium"},
		{Key: "images_modern_format", Category: valueobject.CategoryOnPage, Label: "Images use WebP/AVIF format", Severity: "low"},
		{Key: "images_lazy_loading", Category: valueobject.CategoryOnPage, Label: "Below-fold images use lazy loading", Severity: "low"},

		// Content
		{Key: "content_word_count", Category: valueobject.CategoryContent, Label: "Page has sufficient word count", Severity: "medium"},
		{Key: "content_text_html_ratio", Category: valueobject.CategoryContent, Label: "Text to HTML ratio is healthy", Severity: "medium"},

		// Technical SEO
		{Key: "canonical_exists", Category: valueobject.CategoryTechnical, Label: "Page has a canonical link", Severity: "high"},
		{Key: "canonical_absolute", Category: valueobject.CategoryTechnical, Label: "Canonical URL is absolute", Severity: "medium"},
		{Key: "canonical_self_ref", Category: valueobject.CategoryTechnical, Label: "Canonical is self-referencing", Severity: "low"},
		{Key: "canonical_conflict", Category: valueobject.CategoryTechnical, Label: "No conflicting canonical tags", Severity: "high"},
		{Key: "meta_robots_noindex", Category: valueobject.CategoryTechnical, Label: "No accidental noindex directive", Severity: "critical"},
		{Key: "meta_robots_conflict", Category: valueobject.CategoryTechnical, Label: "No conflicting robots directives", Severity: "high"},
		{Key: "http_status_ok", Category: valueobject.CategoryTechnical, Label: "Page returns HTTP 200 status", Severity: "critical"},
		{Key: "crawl_depth", Category: valueobject.CategoryTechnical, Label: "Page crawl depth is reasonable", Severity: "medium"},
		{Key: "pagination_rel_tags", Category: valueobject.CategoryTechnical, Label: "Pagination uses rel=next/prev", Severity: "low"},
		{Key: "robots_txt_exists", Category: valueobject.CategoryTechnical, Label: "robots.txt file exists", Severity: "high"},
		{Key: "robots_txt_size", Category: valueobject.CategoryTechnical, Label: "robots.txt is not too large", Severity: "low"},
		{Key: "robots_txt_sitemap", Category: valueobject.CategoryTechnical, Label: "robots.txt references sitemap", Severity: "medium"},
		{Key: "robots_txt_syntax", Category: valueobject.CategoryTechnical, Label: "robots.txt has valid syntax", Severity: "medium"},
		{Key: "sitemap_exists", Category: valueobject.CategoryTechnical, Label: "Sitemap.xml exists", Severity: "high"},
		{Key: "sitemap_valid_xml", Category: valueobject.CategoryTechnical, Label: "Sitemap has valid XML", Severity: "high"},
		{Key: "sitemap_index", Category: valueobject.CategoryTechnical, Label: "Sitemap index detected", Severity: "info"},
		{Key: "sitemap_size", Category: valueobject.CategoryTechnical, Label: "Sitemap is within size limits", Severity: "medium"},
		{Key: "sitemap_freshness", Category: valueobject.CategoryTechnical, Label: "Sitemap URLs have recent lastmod", Severity: "low"},
		{Key: "sitemap_coverage", Category: valueobject.CategoryTechnical, Label: "Crawled URLs appear in sitemap", Severity: "medium"},
		{Key: "sitemap_orphan_urls", Category: valueobject.CategoryTechnical, Label: "No sitemap-only orphan URLs", Severity: "medium"},
		{Key: "sitemap_broken_urls", Category: valueobject.CategoryTechnical, Label: "No broken URLs in sitemap", Severity: "high"},
		{Key: "sitemap_redirect_urls", Category: valueobject.CategoryTechnical, Label: "No redirect URLs in sitemap", Severity: "medium"},
		{Key: "sitemap_robots_conflict", Category: valueobject.CategoryTechnical, Label: "No sitemap/robots.txt conflicts", Severity: "high"},
		{Key: "sitemap_image", Category: valueobject.CategoryTechnical, Label: "Sitemap includes image entries", Severity: "low"},
		{Key: "sitemap_video", Category: valueobject.CategoryTechnical, Label: "Sitemap includes video entries", Severity: "info"},
		{Key: "redirect_chains", Category: valueobject.CategoryTechnical, Label: "No redirect chains detected", Severity: "medium"},
		{Key: "temporary_redirects", Category: valueobject.CategoryTechnical, Label: "Temporary redirects are intentional", Severity: "low"},

		// Links
		{Key: "internal_links_present", Category: valueobject.CategoryLinks, Label: "Page has internal links", Severity: "medium"},
		{Key: "internal_links_count", Category: valueobject.CategoryLinks, Label: "Internal link count is reasonable", Severity: "low"},
		{Key: "internal_links_anchor_text", Category: valueobject.CategoryLinks, Label: "Internal links have descriptive anchor text", Severity: "low"},
		{Key: "external_links_rel", Category: valueobject.CategoryLinks, Label: "External links have rel attributes", Severity: "medium"},
		{Key: "broken_internal_links", Category: valueobject.CategoryLinks, Label: "No broken internal links", Severity: "high"},
		{Key: "broken_external_links", Category: valueobject.CategoryLinks, Label: "No broken external links", Severity: "medium"},

		// Performance
		{Key: "page_size", Category: valueobject.CategoryPerformance, Label: "Page size is within limits", Severity: "medium"},
		{Key: "html_size", Category: valueobject.CategoryPerformance, Label: "HTML document size is reasonable", Severity: "medium"},
		{Key: "ttfb", Category: valueobject.CategoryPerformance, Label: "Time to First Byte is fast", Severity: "high"},
		{Key: "compression", Category: valueobject.CategoryPerformance, Label: "Response uses compression", Severity: "medium"},
		{Key: "js_file_count", Category: valueobject.CategoryPerformance, Label: "JavaScript file count is reasonable", Severity: "medium"},
		{Key: "css_file_count", Category: valueobject.CategoryPerformance, Label: "CSS file count is reasonable", Severity: "medium"},
		{Key: "render_blocking", Category: valueobject.CategoryPerformance, Label: "No render-blocking resources", Severity: "high"},
		{Key: "total_requests", Category: valueobject.CategoryPerformance, Label: "Total requests are optimized", Severity: "medium"},
		{Key: "cache_headers", Category: valueobject.CategoryPerformance, Label: "Cache headers are configured", Severity: "low"},
		{Key: "broken_scripts", Category: valueobject.CategoryPerformance, Label: "All script files are reachable", Severity: "high"},
		{Key: "broken_stylesheets", Category: valueobject.CategoryPerformance, Label: "All CSS files are reachable", Severity: "high"},
		{Key: "broken_images", Category: valueobject.CategoryPerformance, Label: "All images are reachable", Severity: "medium"},

		// Structured Data
		{Key: "jsonld_exists", Category: valueobject.CategoryStructuredData, Label: "Page has JSON-LD structured data", Severity: "medium"},
		{Key: "jsonld_valid", Category: valueobject.CategoryStructuredData, Label: "JSON-LD is valid JSON", Severity: "high"},
		{Key: "jsonld_breadcrumb", Category: valueobject.CategoryStructuredData, Label: "BreadcrumbList schema is present", Severity: "low"},

		// Security
		{Key: "uses_https", Category: valueobject.CategorySecurity, Label: "Page uses HTTPS", Severity: "critical"},
		{Key: "mixed_content", Category: valueobject.CategorySecurity, Label: "No mixed content detected", Severity: "high"},
		{Key: "hsts_header", Category: valueobject.CategorySecurity, Label: "HSTS header is present", Severity: "medium"},
		{Key: "security_xcto", Category: valueobject.CategorySecurity, Label: "X-Content-Type-Options header", Severity: "low"},
		{Key: "security_xfo", Category: valueobject.CategorySecurity, Label: "X-Frame-Options header", Severity: "low"},
		{Key: "security_csp", Category: valueobject.CategorySecurity, Label: "Content-Security-Policy header", Severity: "low"},
		{Key: "security_referrer", Category: valueobject.CategorySecurity, Label: "Referrer-Policy header", Severity: "low"},
		{Key: "security_permissions", Category: valueobject.CategorySecurity, Label: "Permissions-Policy header", Severity: "low"},
		{Key: "security_server_disclosure", Category: valueobject.CategorySecurity, Label: "Server version not disclosed", Severity: "low"},

		// Accessibility
		{Key: "html_lang", Category: valueobject.CategoryAccessibility, Label: "HTML has lang attribute", Severity: "high"},
		{Key: "viewport_meta", Category: valueobject.CategoryAccessibility, Label: "Viewport meta tag is present", Severity: "high"},
		{Key: "charset_meta", Category: valueobject.CategoryAccessibility, Label: "Character encoding declared", Severity: "medium"},
		{Key: "empty_links", Category: valueobject.CategoryAccessibility, Label: "Links have accessible text", Severity: "medium"},
		{Key: "aria_landmarks", Category: valueobject.CategoryAccessibility, Label: "Page uses ARIA landmarks", Severity: "low"},

		// Social / Open Graph
		{Key: "og_title", Category: valueobject.CategorySocial, Label: "og:title meta tag is present", Severity: "medium"},
		{Key: "og_description", Category: valueobject.CategorySocial, Label: "og:description meta tag is present", Severity: "medium"},
		{Key: "og_image", Category: valueobject.CategorySocial, Label: "og:image meta tag is present", Severity: "medium"},
		{Key: "og_url", Category: valueobject.CategorySocial, Label: "og:url meta tag is present", Severity: "low"},
		{Key: "og_type", Category: valueobject.CategorySocial, Label: "og:type meta tag is present", Severity: "low"},
		{Key: "og_site_name", Category: valueobject.CategorySocial, Label: "og:site_name meta tag is present", Severity: "low"},
		{Key: "og_locale", Category: valueobject.CategorySocial, Label: "og:locale meta tag is present", Severity: "low"},
		{Key: "twitter_card", Category: valueobject.CategorySocial, Label: "twitter:card meta tag is present", Severity: "medium"},
		{Key: "twitter_title", Category: valueobject.CategorySocial, Label: "twitter:title meta tag is present", Severity: "low"},
		{Key: "twitter_description", Category: valueobject.CategorySocial, Label: "twitter:description meta tag is present", Severity: "low"},
		{Key: "twitter_image", Category: valueobject.CategorySocial, Label: "twitter:image meta tag is present", Severity: "low"},
		{Key: "twitter_site", Category: valueobject.CategorySocial, Label: "twitter:site meta tag is present", Severity: "low"},

		// Mobile
		{Key: "mobile_viewport", Category: valueobject.CategoryMobile, Label: "Viewport meta tag exists", Severity: "critical"},
		{Key: "mobile_viewport_config", Category: valueobject.CategoryMobile, Label: "Viewport uses device-width config", Severity: "medium"},

		// URL Structure
		{Key: "url_lowercase", Category: valueobject.CategoryURLStructure, Label: "URL uses lowercase", Severity: "low"},
		{Key: "url_hyphens", Category: valueobject.CategoryURLStructure, Label: "URL uses hyphens for separation", Severity: "low"},
		{Key: "url_length", Category: valueobject.CategoryURLStructure, Label: "URL length is reasonable", Severity: "low"},
		{Key: "url_special_chars", Category: valueobject.CategoryURLStructure, Label: "URL has no special characters", Severity: "low"},
		{Key: "url_double_slash", Category: valueobject.CategoryURLStructure, Label: "URL has no double slashes", Severity: "medium"},
		{Key: "url_parameters", Category: valueobject.CategoryURLStructure, Label: "URL parameters are minimal", Severity: "low"},

		// Internationalization
		{Key: "hreflang_valid", Category: valueobject.CategoryInternationalization, Label: "Hreflang tags are valid", Severity: "medium"},
		{Key: "hreflang_x_default", Category: valueobject.CategoryInternationalization, Label: "Hreflang x-default is present", Severity: "low"},

		// E-E-A-T
		{Key: "eeat_author", Category: valueobject.CategoryEEAT, Label: "Content has author attribution", Severity: "medium"},
		{Key: "eeat_copyright", Category: valueobject.CategoryEEAT, Label: "Page has copyright notice", Severity: "low"},
		{Key: "eeat_about_page", Category: valueobject.CategoryEEAT, Label: "Site has an About page", Severity: "medium"},
		{Key: "eeat_contact_page", Category: valueobject.CategoryEEAT, Label: "Site has a Contact page", Severity: "medium"},
		{Key: "eeat_privacy_policy", Category: valueobject.CategoryEEAT, Label: "Site has a Privacy Policy", Severity: "medium"},
		{Key: "eeat_terms", Category: valueobject.CategoryEEAT, Label: "Site has Terms of Service", Severity: "low"},

		// Duplicate Content
		{Key: "duplicate_content_body", Category: valueobject.CategoryDuplicateContent, Label: "No duplicate body content", Severity: "high"},
		{Key: "duplicate_titles", Category: valueobject.CategoryDuplicateContent, Label: "No duplicate page titles", Severity: "medium"},
		{Key: "duplicate_descriptions", Category: valueobject.CategoryDuplicateContent, Label: "No duplicate meta descriptions", Severity: "medium"},

		// GEO
		{Key: "geo_crawler_oai-searchbot", Category: valueobject.CategoryGEO, Label: "OAI-SearchBot is allowed", Severity: "high"},
		{Key: "geo_crawler_perplexitybot", Category: valueobject.CategoryGEO, Label: "PerplexityBot is allowed", Severity: "high"},
		{Key: "geo_crawler_claude-searchbot", Category: valueobject.CategoryGEO, Label: "Claude-SearchBot is allowed", Severity: "high"},
		{Key: "geo_crawler_google-extended", Category: valueobject.CategoryGEO, Label: "Google-Extended is allowed", Severity: "high"},
		{Key: "geo_crawler_applebot", Category: valueobject.CategoryGEO, Label: "Applebot is allowed", Severity: "high"},
		{Key: "geo_crawler_block_gptbot", Category: valueobject.CategoryGEO, Label: "GPTBot (training) is blocked", Severity: "low"},
		{Key: "geo_crawler_block_ccbot", Category: valueobject.CategoryGEO, Label: "CCBot (training) is blocked", Severity: "low"},
		{Key: "geo_llms_txt_exists", Category: valueobject.CategoryGEO, Label: "llms.txt file exists", Severity: "medium"},
		{Key: "geo_llms_full_txt", Category: valueobject.CategoryGEO, Label: "llms-full.txt file exists", Severity: "low"},
		{Key: "geo_llms_txt_h1", Category: valueobject.CategoryGEO, Label: "llms.txt has H1 heading", Severity: "medium"},
		{Key: "geo_llms_txt_blockquote", Category: valueobject.CategoryGEO, Label: "llms.txt has blockquote summary", Severity: "medium"},
		{Key: "geo_llms_txt_sections", Category: valueobject.CategoryGEO, Label: "llms.txt has H2 sections", Severity: "low"},
		{Key: "geo_llms_txt_links", Category: valueobject.CategoryGEO, Label: "llms.txt has markdown links", Severity: "low"},
		{Key: "geo_entity_org_schema", Category: valueobject.CategoryGEO, Label: "Organization schema is present", Severity: "high"},
		{Key: "geo_entity_social", Category: valueobject.CategoryGEO, Label: "Social profile links are present", Severity: "medium"},
		{Key: "geo_citability_statistics", Category: valueobject.CategoryGEO, Label: "Content has citable statistics", Severity: "medium"},
		{Key: "geo_citability_faq", Category: valueobject.CategoryGEO, Label: "FAQ schema is present", Severity: "low"},
		{Key: "geo_citability_tables", Category: valueobject.CategoryGEO, Label: "Content has data tables", Severity: "low"},
		{Key: "geo_citability_lists", Category: valueobject.CategoryGEO, Label: "Content has structured lists", Severity: "low"},
		{Key: "geo_citability_question_headings", Category: valueobject.CategoryGEO, Label: "Content has question-style headings", Severity: "medium"},
		{Key: "geo_ai_descriptive_headings", Category: valueobject.CategoryGEO, Label: "Headings are descriptive for AI", Severity: "medium"},
		{Key: "geo_ai_freshness", Category: valueobject.CategoryGEO, Label: "Content has freshness indicators", Severity: "medium"},
		{Key: "geo_ai_semantic_html", Category: valueobject.CategoryGEO, Label: "Page uses semantic HTML elements", Severity: "medium"},
	}
}
