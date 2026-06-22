# SEO/GEO Checks for CI/CD Pipeline

These checks can be run in development mode against local/staging HTML output before deployment.
Designed for integration into a CI/CD SDK that prevents SEO regressions.

## How it works

The SDK will render each page to HTML (or read built HTML files) and run these checks statically.
Checks marked **requires network** need HTTP access to the deployed/preview site.

---

## HTML-Based Checks (No Network Required)

These checks analyze the raw HTML output and can run against built files locally.

### On-Page SEO

| Rule Key | Check | Severity |
|---|---|---|
| `title_exists` | Page has a `<title>` tag | High |
| `title_length` | Title is between 30-60 characters | Medium |
| `meta_description_exists` | Page has a meta description | High |
| `meta_description_length` | Description is 120-160 characters | Medium |
| `h1_count` | Page has exactly one `<h1>` | High |
| `h1_not_empty` | `<h1>` has non-empty text | High |
| `heading_hierarchy` | Headings follow proper hierarchy (no skipped levels) | Medium |
| `images_alt_text` | All images have alt attributes | High |
| `images_alt_descriptive` | Image alt text is not generic ("image", "photo") | Low |
| `images_dimensions` | Images have width/height attributes | Medium |
| `images_modern_format` | Images use WebP/AVIF format | Low |
| `images_lazy_loading` | Below-fold images have `loading="lazy"` | Low |

### Accessibility

| Rule Key | Check | Severity |
|---|---|---|
| `html_lang` | `<html>` has a `lang` attribute | High |
| `viewport_meta` | Viewport meta tag is present | High |
| `charset_meta` | Character encoding is declared | Medium |
| `empty_links` | Links have accessible text/aria-label | Medium |
| `aria_landmarks` | Page uses `<nav>`, `<main>`, or ARIA roles | Low |

### Mobile

| Rule Key | Check | Severity |
|---|---|---|
| `mobile_viewport` | Viewport meta tag exists | Critical |
| `mobile_viewport_config` | Viewport uses `width=device-width, initial-scale=1` | Medium |

### Structured Data

| Rule Key | Check | Severity |
|---|---|---|
| `jsonld_exists` | Page has JSON-LD structured data | Medium |
| `jsonld_valid` | JSON-LD blocks contain valid JSON | High |
| `jsonld_breadcrumb` | BreadcrumbList schema is present | Low |

### Social / Open Graph

| Rule Key | Check | Severity |
|---|---|---|
| `og_title` | `og:title` meta tag is present | Medium |
| `og_description` | `og:description` meta tag is present | Medium |
| `og_image` | `og:image` meta tag is present | Medium |
| `og_url` | `og:url` meta tag is present | Low |
| `og_type` | `og:type` meta tag is present | Low |
| `og_site_name` | `og:site_name` meta tag is present | Low |
| `og_locale` | `og:locale` meta tag is present | Low |
| `twitter_card` | `twitter:card` meta tag is present | Medium |
| `twitter_title` | `twitter:title` meta tag is present | Low |
| `twitter_description` | `twitter:description` meta tag is present | Low |
| `twitter_image` | `twitter:image` meta tag is present | Low |
| `twitter_site` | `twitter:site` meta tag is present | Low |

### Technical SEO

| Rule Key | Check | Severity |
|---|---|---|
| `canonical_exists` | Page has a canonical link | High |
| `canonical_absolute` | Canonical URL is absolute | Medium |
| `canonical_self_ref` | Canonical points to self | Low |
| `canonical_conflict` | No conflicting canonical declarations | High |
| `meta_robots_noindex` | Page is not `noindex` (unless intended) | Critical |
| `meta_robots_conflict` | No conflicting robots directives | High |
| `pagination_rel_tags` | Paginated pages have `rel="next"` / `rel="prev"` | Low |

### Content Quality

| Rule Key | Check | Severity |
|---|---|---|
| `content_word_count` | Page has adequate content (300+ words) | Medium |
| `content_text_html_ratio` | Text-to-HTML ratio is above 10% | Medium |

### Performance (Static)

| Rule Key | Check | Severity |
|---|---|---|
| `js_file_count` | Page loads a reasonable number of JS files (≤10) | Medium |
| `css_file_count` | Page loads a reasonable number of CSS files (≤5) | Medium |
| `render_blocking` | No render-blocking scripts in `<head>` | High |
| `total_requests` | Total resource count is reasonable (≤50) | Medium |
| `page_size` | Total page size is reasonable | Medium |
| `html_size` | HTML document size is reasonable | Medium |

### Links (Static)

| Rule Key | Check | Severity |
|---|---|---|
| `internal_links_present` | Page has internal links | Medium |
| `internal_links_count` | Not excessive internal links (≤100) | Low |
| `internal_links_anchor_text` | Links use descriptive anchor text | Low |
| `external_links_rel` | External `_blank` links have `rel="noopener noreferrer"` | Medium |

### URL Structure

| Rule Key | Check | Severity |
|---|---|---|
| `url_lowercase` | URL uses lowercase characters | Low |
| `url_hyphens` | URL uses hyphens instead of underscores | Low |
| `url_length` | URL is under 100 characters | Low |
| `url_special_chars` | URL has no special characters | Low |
| `url_double_slash` | URL path has no double slashes | Medium |
| `url_parameters` | URL has minimal query parameters | Low |

### Internationalization

| Rule Key | Check | Severity |
|---|---|---|
| `hreflang_valid` | Hreflang tags have valid language codes | Medium |
| `hreflang_x_default` | `x-default` hreflang is present | Low |

### E-E-A-T

| Rule Key | Check | Severity |
|---|---|---|
| `eeat_author` | Content pages show author information | Medium |
| `eeat_copyright` | Page has copyright notice | Low |

### GEO (AI Readiness - Static)

| Rule Key | Check | Severity |
|---|---|---|
| `geo_ai_descriptive_headings` | Pages have descriptive H2 headings | Medium |
| `geo_ai_freshness` | Pages have content freshness signals | Medium |
| `geo_ai_semantic_html` | Pages use semantic HTML elements | Medium |
| `geo_entity_org_schema` | Organization schema is complete | High |
| `geo_entity_social` | Social profile links in Organization schema | Medium |
| `geo_citability_statistics` | Pages contain specific data points | Medium |
| `geo_citability_faq` | FAQ-style content exists | Low |
| `geo_citability_tables` | Comparison tables exist | Low |
| `geo_citability_lists` | Numbered/ordered lists exist | Low |
| `geo_citability_question_headings` | Question-format headings exist | Medium |

---

## Network-Required Checks (Requires Deployed/Preview Site)

These checks need HTTP access to the running site and cannot run on static HTML alone.

### Technical SEO

| Rule Key | Check | Severity |
|---|---|---|
| `http_status_ok` | Page returns 200 status | Critical |
| `robots_txt_exists` | robots.txt file exists and is accessible | High |
| `robots_txt_size` | robots.txt is not excessively large | Low |
| `robots_txt_sitemap` | robots.txt references sitemap | Medium |
| `robots_txt_syntax` | robots.txt has valid syntax | Medium |
| `sitemap_exists` | sitemap.xml exists | High |
| `sitemap_valid_xml` | Sitemap contains valid XML | High |
| `sitemap_size` | Sitemap has reasonable number of URLs | Medium |
| `sitemap_freshness` | Sitemap has lastmod dates | Low |
| `sitemap_coverage` | Sitemap covers crawled pages | Medium |
| `sitemap_orphan_urls` | No orphan URLs in sitemap | Medium |
| `sitemap_broken_urls` | No broken URLs in sitemap | High |
| `sitemap_redirect_urls` | No redirect URLs in sitemap | Medium |
| `sitemap_robots_conflict` | No robots.txt/sitemap conflicts | High |
| `redirect_chains` | No multi-hop redirect chains | Medium |
| `temporary_redirects` | No unnecessary 302/307 redirects | Low |
| `crawl_depth` | Pages are reachable within reasonable depth | Medium |

### Performance (Runtime)

| Rule Key | Check | Severity |
|---|---|---|
| `ttfb` | Server response time (TTFB) is fast | High |
| `compression` | Response uses gzip/brotli compression | Medium |
| `cache_headers` | Cache-Control headers are set | Low |

### Security (Runtime)

| Rule Key | Check | Severity |
|---|---|---|
| `uses_https` | Page is served over HTTPS | Critical |
| `mixed_content` | No HTTP resources on HTTPS page | High |
| `hsts_header` | HSTS header is present | Medium |
| `security_xcto` | X-Content-Type-Options header present | Low |
| `security_xfo` | X-Frame-Options header present | Low |
| `security_csp` | Content-Security-Policy header present | Low |
| `security_referrer` | Referrer-Policy header present | Low |
| `security_permissions` | Permissions-Policy header present | Low |
| `security_server_disclosure` | Server version is not disclosed | Low |

### Links (Runtime)

| Rule Key | Check | Severity |
|---|---|---|
| `broken_internal_links` | No broken internal links | High |
| `broken_external_links` | No broken external links | Medium |

### Content (Cross-Page)

| Rule Key | Check | Severity |
|---|---|---|
| `duplicate_content_body` | No duplicate page content | High |
| `duplicate_titles` | No duplicate titles across pages | Medium |
| `duplicate_descriptions` | No duplicate meta descriptions | Medium |

### E-E-A-T (Site-Level)

| Rule Key | Check | Severity |
|---|---|---|
| `eeat_about_page` | About page exists | Medium |
| `eeat_contact_page` | Contact page exists | Medium |
| `eeat_privacy_policy` | Privacy policy page exists | Medium |
| `eeat_terms` | Terms of service page exists | Low |

### GEO (Runtime)

| Rule Key | Check | Severity |
|---|---|---|
| `geo_llms_txt_exists` | /llms.txt file exists | Medium |
| `geo_llms_txt_h1` | llms.txt has H1 heading | Medium |
| `geo_llms_txt_blockquote` | llms.txt has summary blockquote | Medium |
| `geo_llms_txt_sections` | llms.txt has H2 sections | Low |
| `geo_llms_txt_links` | llms.txt has resource links | Low |
| `geo_llms_full_txt` | /llms-full.txt exists | Low |
| `geo_crawler_*` | AI search crawlers are not blocked | High |

---

## SDK Integration Notes

### Recommended CI Pipeline Usage

```yaml
# Example GitHub Actions step
- name: SEO/GEO Check
  run: seo-crawler check --mode=static --dir=./dist --fail-on=high
```

### Exit Codes

- `0` — All checks pass
- `1` — Failures found at or above the configured severity threshold
- `2` — Configuration error

### Configuration

The SDK should support:
- `--mode=static` — Run HTML-based checks only (no network)
- `--mode=full` — Run all checks against a live URL
- `--fail-on=critical|high|medium|low` — Set the severity threshold for CI failure
- `--ignore=rule_key,rule_key` — Skip specific rules
- `--config=.seo-crawler.yml` — Load rules configuration from file

### Total Checks

- **HTML-based (static)**: 68 checks
- **Network-required**: 42 checks
- **Total**: 110 unique rule keys
