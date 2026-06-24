const fs = require('fs');
const path = require('path');

function loadRuleDescriptions() {
  const tsPath = path.join(__dirname, '..', 'web', 'src', 'data', 'rule-descriptions.ts');
  let source = fs.readFileSync(tsPath, 'utf8');

  source = source.replace(/export\s+interface\s+RuleDescription\s*\{[^}]+\}\s*/s, '');
  source = source.replace(/export\s+const\s+ruleDescriptions:\s*Record<[^>]+>\s*=/, 'const ruleDescriptions =');
  source = source.replace(/export\s+function\s+getRuleDescription[\s\S]*$/, '');

  const fn = new Function(source + '\nreturn ruleDescriptions;');
  return fn();
}

const ruleDescriptionsFromSource = loadRuleDescriptions();

const CATEGORIES = {
  on_page:               { label: 'On-Page SEO',          slug: 'on-page-seo' },
  content:               { label: 'Content Quality',       slug: 'content' },
  technical:             { label: 'Technical SEO',         slug: 'technical-seo' },
  links:                 { label: 'Links',                 slug: 'links' },
  performance:           { label: 'Performance',           slug: 'performance' },
  structured_data:       { label: 'Structured Data',       slug: 'structured-data' },
  security:              { label: 'Security',              slug: 'security' },
  accessibility:         { label: 'Accessibility',         slug: 'accessibility' },
  social:                { label: 'Social / Open Graph',   slug: 'social-open-graph' },
  mobile:                { label: 'Mobile',                slug: 'mobile' },
  url_structure:         { label: 'URL Structure',         slug: 'url-structure' },
  internationalization:  { label: 'Internationalization',  slug: 'internationalization' },
  eeat:                  { label: 'E-E-A-T',               slug: 'eeat' },
  duplicate_content:     { label: 'Duplicate Content',     slug: 'duplicate-content' },
  geo:                   { label: 'GEO (AI Search)',       slug: 'geo-ai-search' },
};

const RULES = [
  { key: 'title_exists', category: 'on_page', label: 'Page has a <title> tag', severity: 'critical' },
  { key: 'title_length', category: 'on_page', label: 'Title is between 30-60 characters', severity: 'medium' },
  { key: 'title_multiple', category: 'on_page', label: 'Page has only one title tag', severity: 'high' },
  { key: 'meta_description_exists', category: 'on_page', label: 'Page has a meta description', severity: 'high' },
  { key: 'meta_description_length', category: 'on_page', label: 'Description is 120-160 characters', severity: 'medium' },
  { key: 'meta_description_multiple', category: 'on_page', label: 'Page has only one meta description', severity: 'medium' },
  { key: 'h1_count', category: 'on_page', label: 'Page has exactly one H1', severity: 'high' },
  { key: 'h1_not_empty', category: 'on_page', label: 'H1 has non-empty text', severity: 'high' },
  { key: 'heading_hierarchy', category: 'on_page', label: 'Headings follow proper hierarchy', severity: 'medium' },
  { key: 'heading_not_empty', category: 'on_page', label: 'Headings have non-empty text', severity: 'low' },
  { key: 'images_alt_text', category: 'on_page', label: 'All images have alt attributes', severity: 'high' },
  { key: 'images_alt_descriptive', category: 'on_page', label: 'Image alt text is descriptive', severity: 'low' },
  { key: 'images_dimensions', category: 'on_page', label: 'Images have width/height attributes', severity: 'medium' },
  { key: 'images_modern_format', category: 'on_page', label: 'Images use WebP/AVIF format', severity: 'low' },
  { key: 'images_lazy_loading', category: 'on_page', label: 'Below-fold images use lazy loading', severity: 'low' },
  { key: 'content_word_count', category: 'content', label: 'Page has sufficient word count', severity: 'medium' },
  { key: 'content_text_html_ratio', category: 'content', label: 'Text to HTML ratio is healthy', severity: 'medium' },
  { key: 'canonical_exists', category: 'technical', label: 'Page has a canonical link', severity: 'high' },
  { key: 'canonical_absolute', category: 'technical', label: 'Canonical URL is absolute', severity: 'medium' },
  { key: 'canonical_self_ref', category: 'technical', label: 'Canonical is self-referencing', severity: 'low' },
  { key: 'canonical_conflict', category: 'technical', label: 'No conflicting canonical tags', severity: 'high' },
  { key: 'meta_robots_noindex', category: 'technical', label: 'Page has a noindex directive', severity: 'info' },
  { key: 'meta_robots_conflict', category: 'technical', label: 'No conflicting robots directives', severity: 'high' },
  { key: 'http_status_ok', category: 'technical', label: 'Page returns HTTP 200 status', severity: 'critical' },
  { key: 'crawl_depth', category: 'technical', label: 'Page crawl depth is reasonable', severity: 'medium' },
  { key: 'pagination_rel_tags', category: 'technical', label: 'Pagination uses rel=next/prev', severity: 'low' },
  { key: 'robots_txt_exists', category: 'technical', label: 'robots.txt file exists', severity: 'high' },
  { key: 'robots_txt_size', category: 'technical', label: 'robots.txt is not too large', severity: 'low' },
  { key: 'robots_txt_sitemap', category: 'technical', label: 'robots.txt references sitemap', severity: 'medium' },
  { key: 'robots_txt_syntax', category: 'technical', label: 'robots.txt has valid syntax', severity: 'medium' },
  { key: 'sitemap_exists', category: 'technical', label: 'Sitemap.xml exists', severity: 'high' },
  { key: 'sitemap_valid_xml', category: 'technical', label: 'Sitemap has valid XML', severity: 'high' },
  { key: 'sitemap_index', category: 'technical', label: 'Sitemap index detected', severity: 'info' },
  { key: 'sitemap_size', category: 'technical', label: 'Sitemap is within size limits', severity: 'medium' },
  { key: 'sitemap_freshness', category: 'technical', label: 'Sitemap URLs have recent lastmod', severity: 'low' },
  { key: 'sitemap_coverage', category: 'technical', label: 'Crawled URLs appear in sitemap', severity: 'medium' },
  { key: 'sitemap_orphan_urls', category: 'technical', label: 'No sitemap-only orphan URLs', severity: 'medium' },
  { key: 'sitemap_broken_urls', category: 'technical', label: 'No broken URLs in sitemap', severity: 'high' },
  { key: 'sitemap_redirect_urls', category: 'technical', label: 'No redirect URLs in sitemap', severity: 'medium' },
  { key: 'sitemap_robots_conflict', category: 'technical', label: 'No sitemap/robots.txt conflicts', severity: 'high' },
  { key: 'sitemap_image', category: 'technical', label: 'Sitemap includes image entries', severity: 'low' },
  { key: 'sitemap_video', category: 'technical', label: 'Sitemap includes video entries', severity: 'info' },
  { key: 'redirect_chains', category: 'technical', label: 'No redirect chains detected', severity: 'medium' },
  { key: 'temporary_redirects', category: 'technical', label: 'Temporary redirects are intentional', severity: 'low' },
  { key: 'internal_links_present', category: 'links', label: 'Page has internal links', severity: 'medium' },
  { key: 'internal_links_count', category: 'links', label: 'Internal link count is reasonable', severity: 'low' },
  { key: 'internal_links_anchor_text', category: 'links', label: 'Internal links have descriptive anchor text', severity: 'low' },
  { key: 'external_links_rel', category: 'links', label: 'External links have rel attributes', severity: 'medium' },
  { key: 'broken_internal_links', category: 'links', label: 'No broken internal links', severity: 'high' },
  { key: 'broken_external_links', category: 'links', label: 'No broken external links', severity: 'medium' },
  { key: 'page_size', category: 'performance', label: 'Page size is within limits', severity: 'medium' },
  { key: 'html_size', category: 'performance', label: 'HTML document size is reasonable', severity: 'medium' },
  { key: 'ttfb', category: 'performance', label: 'Time to First Byte is fast', severity: 'high' },
  { key: 'compression', category: 'performance', label: 'Response uses compression', severity: 'medium' },
  { key: 'js_file_count', category: 'performance', label: 'JavaScript file count is reasonable', severity: 'medium' },
  { key: 'css_file_count', category: 'performance', label: 'CSS file count is reasonable', severity: 'medium' },
  { key: 'render_blocking', category: 'performance', label: 'No render-blocking resources', severity: 'high' },
  { key: 'total_requests', category: 'performance', label: 'Total requests are optimized', severity: 'medium' },
  { key: 'cache_headers', category: 'performance', label: 'Cache headers are configured', severity: 'low' },
  { key: 'broken_scripts', category: 'performance', label: 'All script files are reachable', severity: 'high' },
  { key: 'broken_stylesheets', category: 'performance', label: 'All CSS files are reachable', severity: 'high' },
  { key: 'broken_images', category: 'performance', label: 'All images are reachable', severity: 'medium' },
  { key: 'jsonld_exists', category: 'structured_data', label: 'Page has JSON-LD structured data', severity: 'medium' },
  { key: 'jsonld_valid', category: 'structured_data', label: 'JSON-LD is valid JSON', severity: 'high' },
  { key: 'jsonld_breadcrumb', category: 'structured_data', label: 'BreadcrumbList schema is present', severity: 'low' },
  { key: 'uses_https', category: 'security', label: 'Page uses HTTPS', severity: 'critical' },
  { key: 'mixed_content', category: 'security', label: 'No mixed content detected', severity: 'high' },
  { key: 'hsts_header', category: 'security', label: 'HSTS header is present', severity: 'medium' },
  { key: 'security_xcto', category: 'security', label: 'X-Content-Type-Options header', severity: 'low' },
  { key: 'security_xfo', category: 'security', label: 'X-Frame-Options header', severity: 'low' },
  { key: 'security_csp', category: 'security', label: 'Content-Security-Policy header', severity: 'low' },
  { key: 'security_referrer', category: 'security', label: 'Referrer-Policy header', severity: 'low' },
  { key: 'security_permissions', category: 'security', label: 'Permissions-Policy header', severity: 'low' },
  { key: 'security_server_disclosure', category: 'security', label: 'Server version not disclosed', severity: 'low' },
  { key: 'html_lang', category: 'accessibility', label: 'HTML has lang attribute', severity: 'high' },
  { key: 'viewport_meta', category: 'accessibility', label: 'Viewport meta tag is present', severity: 'high' },
  { key: 'charset_meta', category: 'accessibility', label: 'Character encoding declared', severity: 'medium' },
  { key: 'empty_links', category: 'accessibility', label: 'Links have accessible text', severity: 'medium' },
  { key: 'aria_landmarks', category: 'accessibility', label: 'Page uses ARIA landmarks', severity: 'low' },
  { key: 'og_title', category: 'social', label: 'og:title meta tag is present', severity: 'medium' },
  { key: 'og_description', category: 'social', label: 'og:description meta tag is present', severity: 'medium' },
  { key: 'og_image', category: 'social', label: 'og:image meta tag is present', severity: 'medium' },
  { key: 'og_url', category: 'social', label: 'og:url meta tag is present', severity: 'low' },
  { key: 'og_type', category: 'social', label: 'og:type meta tag is present', severity: 'low' },
  { key: 'og_site_name', category: 'social', label: 'og:site_name meta tag is present', severity: 'low' },
  { key: 'og_locale', category: 'social', label: 'og:locale meta tag is present', severity: 'low' },
  { key: 'twitter_card', category: 'social', label: 'twitter:card meta tag is present', severity: 'medium' },
  { key: 'twitter_title', category: 'social', label: 'twitter:title meta tag is present', severity: 'low' },
  { key: 'twitter_description', category: 'social', label: 'twitter:description meta tag is present', severity: 'low' },
  { key: 'twitter_image', category: 'social', label: 'twitter:image meta tag is present', severity: 'low' },
  { key: 'twitter_site', category: 'social', label: 'twitter:site meta tag is present', severity: 'low' },
  { key: 'mobile_viewport', category: 'mobile', label: 'Viewport meta tag exists', severity: 'critical' },
  { key: 'mobile_viewport_config', category: 'mobile', label: 'Viewport uses device-width config', severity: 'medium' },
  { key: 'url_lowercase', category: 'url_structure', label: 'URL uses lowercase', severity: 'low' },
  { key: 'url_hyphens', category: 'url_structure', label: 'URL uses hyphens for separation', severity: 'low' },
  { key: 'url_length', category: 'url_structure', label: 'URL length is reasonable', severity: 'low' },
  { key: 'url_special_chars', category: 'url_structure', label: 'URL has no special characters', severity: 'low' },
  { key: 'url_double_slash', category: 'url_structure', label: 'URL has no double slashes', severity: 'medium' },
  { key: 'url_parameters', category: 'url_structure', label: 'URL parameters are minimal', severity: 'low' },
  { key: 'hreflang_valid', category: 'internationalization', label: 'Hreflang tags are valid', severity: 'medium' },
  { key: 'hreflang_x_default', category: 'internationalization', label: 'Hreflang x-default is present', severity: 'low' },
  { key: 'eeat_author', category: 'eeat', label: 'Content has author attribution', severity: 'medium' },
  { key: 'eeat_copyright', category: 'eeat', label: 'Page has copyright notice', severity: 'low' },
  { key: 'eeat_about_page', category: 'eeat', label: 'Site has an About page', severity: 'medium' },
  { key: 'eeat_contact_page', category: 'eeat', label: 'Site has a Contact page', severity: 'medium' },
  { key: 'eeat_privacy_policy', category: 'eeat', label: 'Site has a Privacy Policy', severity: 'medium' },
  { key: 'eeat_terms', category: 'eeat', label: 'Site has Terms of Service', severity: 'low' },
  { key: 'duplicate_content_body', category: 'duplicate_content', label: 'No duplicate body content', severity: 'high' },
  { key: 'duplicate_titles', category: 'duplicate_content', label: 'No duplicate page titles', severity: 'medium' },
  { key: 'duplicate_descriptions', category: 'duplicate_content', label: 'No duplicate meta descriptions', severity: 'medium' },
  { key: 'geo_crawler_oai-searchbot', category: 'geo', label: 'OAI-SearchBot is allowed', severity: 'high' },
  { key: 'geo_crawler_perplexitybot', category: 'geo', label: 'PerplexityBot is allowed', severity: 'high' },
  { key: 'geo_crawler_claude-searchbot', category: 'geo', label: 'Claude-SearchBot is allowed', severity: 'high' },
  { key: 'geo_crawler_google-extended', category: 'geo', label: 'Google-Extended is allowed', severity: 'high' },
  { key: 'geo_crawler_applebot', category: 'geo', label: 'Applebot is allowed', severity: 'high' },
  { key: 'geo_crawler_block_gptbot', category: 'geo', label: 'GPTBot (training) is blocked', severity: 'low' },
  { key: 'geo_crawler_block_ccbot', category: 'geo', label: 'CCBot (training) is blocked', severity: 'low' },
  { key: 'geo_llms_txt_exists', category: 'geo', label: 'llms.txt file exists', severity: 'medium' },
  { key: 'geo_llms_full_txt', category: 'geo', label: 'llms-full.txt file exists', severity: 'low' },
  { key: 'geo_llms_txt_h1', category: 'geo', label: 'llms.txt has H1 heading', severity: 'medium' },
  { key: 'geo_llms_txt_blockquote', category: 'geo', label: 'llms.txt has blockquote summary', severity: 'medium' },
  { key: 'geo_llms_txt_sections', category: 'geo', label: 'llms.txt has H2 sections', severity: 'low' },
  { key: 'geo_llms_txt_links', category: 'geo', label: 'llms.txt has markdown links', severity: 'low' },
  { key: 'geo_entity_org_schema', category: 'geo', label: 'Organization schema is present', severity: 'high' },
  { key: 'geo_entity_social', category: 'geo', label: 'Social profile links are present', severity: 'medium' },
  { key: 'geo_citability_statistics', category: 'geo', label: 'Content has citable statistics', severity: 'medium' },
  { key: 'geo_citability_faq', category: 'geo', label: 'FAQ schema is present', severity: 'low' },
  { key: 'geo_citability_tables', category: 'geo', label: 'Content has data tables', severity: 'low' },
  { key: 'geo_citability_lists', category: 'geo', label: 'Content has structured lists', severity: 'low' },
  { key: 'geo_citability_question_headings', category: 'geo', label: 'Content has question-style headings', severity: 'medium' },
  { key: 'geo_ai_descriptive_headings', category: 'geo', label: 'Headings are descriptive for AI', severity: 'medium' },
  { key: 'geo_ai_freshness', category: 'geo', label: 'Content has freshness indicators', severity: 'medium' },
  { key: 'geo_ai_semantic_html', category: 'geo', label: 'Page uses semantic HTML elements', severity: 'medium' },
];

function escapeHtml(text) {
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}

const OG_IMAGE_URL = 'https://raw.githubusercontent.com/booltools/booltools-seo-crawler/master/web/public/screenshots/home.png';
const TWITTER_SITE = '@booltools';
const COPYRIGHT_YEAR = new Date().getFullYear();
const LAST_UPDATED_ISO = new Date().toISOString().split('T')[0];

function getDescription(key, ruleLabel) {
  const source = ruleDescriptionsFromSource[key];
  if (source) {
    return {
      short: source.shortDescription,
      full: source.fullDescription,
      why: source.whyImportant,
      fix: source.howToFix,
      bad: source.badExample,
      good: source.goodExample,
      snippet: source.agentSnippet,
    };
  }
  const label = ruleLabel || key;
  return {
    short: `${label}. This audit rule is part of the Booltools SEO Crawler rule set for optimizing search engine visibility.`,
    full: `This rule verifies that your page follows the "${label}" best practice for search engine optimization and generative engine optimization (GEO). It is checked automatically by the Booltools SEO Crawler.`,
    why: `Following the "${label}" best practice improves your page's visibility in search engines and AI-powered search systems. Failing this check can negatively impact your SEO score.`,
    fix: 'Check each affected URL in the audit report and apply the recommendation provided. See the examples below for guidance.',
    bad: '<!-- See the audit report for specific examples -->',
    good: '<!-- Follow the rule label guidance for correct implementation -->',
    snippet: 'Review the affected URLs in the audit report and apply the fix recommendation.',
  };
}

function groupRulesByCategory() {
  const grouped = {};
  for (const rule of RULES) {
    if (!grouped[rule.category]) grouped[rule.category] = [];
    grouped[rule.category].push(rule);
  }
  return grouped;
}

function buildSidebar(activeCategorySlug, activeRuleKey, docsBase) {
  const grouped = groupRulesByCategory();
  const categoryOrder = Object.keys(CATEGORIES);

  let sidebarHtml = `    <aside class="docs-sidebar">
      <nav aria-label="Documentation navigation">
        <ul>
          <li><a href="${docsBase}/index.html">Overview</a></li>
          <li><a href="${docsBase}/getting-started.html">Getting Started</a></li>
          <li><a href="${docsBase}/web-ui.html">Web UI Guide</a></li>
          <li><a href="${docsBase}/sdk.html">CI/CD SDK</a></li>
          <li>
            <a href="${docsBase}/rules.html">Rules Reference</a>
            <ul class="docs-subnav">`;

  for (const catKey of categoryOrder) {
    const catInfo = CATEGORIES[catKey];
    const isActive = catInfo.slug === activeCategorySlug;
    const activeClass = isActive ? ' class="active"' : '';
    sidebarHtml += `
              <li><a href="${docsBase}/rules/${catInfo.slug}.html"${activeClass}>${catInfo.label}</a>`;

    if (isActive && grouped[catKey]) {
      sidebarHtml += `
                <ul class="docs-subnav-rules">`;
      for (const rule of grouped[catKey]) {
        const ruleActive = rule.key === activeRuleKey ? ' class="active"' : '';
        sidebarHtml += `
                  <li><a href="${docsBase}/rules/${catInfo.slug}/${rule.key}.html"${ruleActive} title="${escapeHtml(rule.label)}">${rule.key}</a></li>`;
      }
      sidebarHtml += `
                </ul>`;
    }

    sidebarHtml += `
              </li>`;
  }

  sidebarHtml += `
            </ul>
          </li>
          <li><a href="${docsBase}/presets.html">Rule Presets</a></li>
          <li><a href="${docsBase}/localhost.html">Localhost Guide</a></li>
          <li><a href="${docsBase}/api.html">API Reference</a></li>
        </ul>
      </nav>
    </aside>`;

  return sidebarHtml;
}

function buildMobileMenu() {
  return `    <button id="docs-menu-toggle" class="docs-menu-toggle" aria-label="Toggle menu">
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <line x1="3" y1="12" x2="21" y2="12" />
        <line x1="3" y1="6" x2="21" y2="6" />
        <line x1="3" y1="18" x2="21" y2="18" />
      </svg>
      <span>Menu</span>
    </button>
    <div id="docs-drawer-overlay" class="docs-drawer-overlay"></div>`;
}

function buildMobileScript() {
  return `  <script>
    (function() {
      var toggle = document.getElementById('docs-menu-toggle');
      var sidebar = document.querySelector('.docs-sidebar');
      var overlay = document.getElementById('docs-drawer-overlay');
      function toggleDrawer() {
        var isOpen = sidebar.classList.contains('open');
        sidebar.classList.toggle('open', !isOpen);
        overlay.classList.toggle('visible', !isOpen);
        document.body.style.overflow = !isOpen ? 'hidden' : '';
      }
      toggle.addEventListener('click', toggleDrawer);
      overlay.addEventListener('click', toggleDrawer);
      sidebar.querySelectorAll('a').forEach(function(link) {
        link.addEventListener('click', function() {
          if (sidebar.classList.contains('open')) toggleDrawer();
        });
      });
    })();
  </script>`;
}

function buildFooter() {
  return `  <footer class="docs-footer" role="contentinfo">
    <p>&copy; ${COPYRIGHT_YEAR} Booltools. All rights reserved. Open-source under MIT License.</p>
  </footer>`;
}

function buildCopyScript() {
  return `  <script>
    document.querySelectorAll('.copy-btn').forEach(function(button) {
      button.addEventListener('click', function() {
        var targetId = button.getAttribute('data-snippet-target');
        if (!targetId) return;
        var el = document.getElementById(targetId);
        if (!el) return;
        navigator.clipboard.writeText(el.textContent || '').then(function() {
          var span = button.querySelector('.copy-text');
          if (span) {
            span.textContent = 'Copied!';
            setTimeout(function() { span.textContent = 'Copy'; }, 2000);
          }
        });
      });
    });
  </script>`;
}

function buildPageHead(title, description, canonicalPath, cssPath, jsonLd) {
  const baseUrl = 'https://booltools.github.io/booltools-seo-crawler';
  const fullUrl = `${baseUrl}/${canonicalPath}`;
  const fullTitle = `${escapeHtml(title)} — Booltools`;
  const escapedDescription = escapeHtml(description);
  const jsonLdTag = jsonLd ? `\n  <script type="application/ld+json">${JSON.stringify(jsonLd)}</script>` : '';
  return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>${fullTitle}</title>
  <meta name="description" content="${escapedDescription}" />
  <link rel="canonical" href="${fullUrl}" />
  <meta name="author" content="Booltools" />
  <meta property="og:title" content="${fullTitle}" />
  <meta property="og:description" content="${escapedDescription}" />
  <meta property="og:type" content="article" />
  <meta property="og:url" content="${fullUrl}" />
  <meta property="og:locale" content="en_US" />
  <meta property="og:site_name" content="Booltools SEO Crawler" />
  <meta property="og:image" content="${OG_IMAGE_URL}" />
  <meta name="twitter:card" content="summary_large_image" />
  <meta name="twitter:site" content="${TWITTER_SITE}" />
  <meta name="twitter:title" content="${fullTitle}" />
  <meta name="twitter:description" content="${escapedDescription}" />
  <meta name="twitter:image" content="${OG_IMAGE_URL}" />
  <link rel="preconnect" href="https://fonts.googleapis.com" />
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap" rel="stylesheet" />
  <link rel="stylesheet" href="${cssPath}" />${jsonLdTag}
</head>`;
}

function buildHeader(backPath) {
  return `<body>
  <header class="docs-header">
    <a href="${backPath}" class="logo"><svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg> <span>Booltools SEO Crawler</span></a>
    <a href="https://github.com/booltools/booltools-seo-crawler" class="back-link" target="_blank" rel="noopener noreferrer">GitHub</a>
  </header>`;
}

function generateCategoryPage(categoryKey) {
  const catInfo = CATEGORIES[categoryKey];
  const grouped = groupRulesByCategory();
  const rules = grouped[categoryKey] || [];
  const title = `${catInfo.label} Rules`;
  const description = `Browse all ${rules.length} ${catInfo.label} audit rules in the Booltools SEO Crawler. Each rule includes severity, description, examples, and an AI-agent fix snippet.`;
  const sidebar = buildSidebar(catInfo.slug, null, '..');

  const baseUrl = 'https://booltools.github.io/booltools-seo-crawler';
  const canonicalPath = `docs/rules/${catInfo.slug}.html`;
  const jsonLd = {
    '@context': 'https://schema.org',
    '@type': 'CollectionPage',
    name: title,
    description: description,
    url: `${baseUrl}/${canonicalPath}`,
    isPartOf: { '@type': 'WebSite', name: 'Booltools SEO Crawler', url: baseUrl },
    breadcrumb: {
      '@type': 'BreadcrumbList',
      itemListElement: [
        { '@type': 'ListItem', position: 1, name: 'Rules', item: `${baseUrl}/docs/rules.html` },
        { '@type': 'ListItem', position: 2, name: catInfo.label, item: `${baseUrl}/${canonicalPath}` },
      ],
    },
    dateModified: LAST_UPDATED_ISO,
  };

  let rulesListHtml = '';
  for (const rule of rules) {
    const desc = getDescription(rule.key, rule.label);
    rulesListHtml += `
          <a href="${catInfo.slug}/${rule.key}.html" class="rule-card-link">
            <div class="rule-card">
              <div class="rule-card-header">
                <code>${rule.key}</code>
                <span class="severity-badge sev-${rule.severity}">${rule.severity}</span>
              </div>
              <h2 class="rule-card-title">${escapeHtml(rule.label)}</h2>
              <p>${escapeHtml(desc.short)}</p>
            </div>
          </a>`;
  }

  return `${buildPageHead(title, description, canonicalPath, '../docs.css', jsonLd)}
${buildHeader('../../')}

  <div class="docs-page">
${buildMobileMenu()}

${sidebar}
    <main class="docs-content" role="main">
      <nav class="breadcrumb" aria-label="Breadcrumb">
        <a href="../rules.html">Rules</a>
        <span class="sep">/</span>
        <span>${catInfo.label}</span>
      </nav>

      <h1>${catInfo.label}</h1>
      <p>${rules.length} audit rules in this category. Click any rule to see the full description, examples, and an AI-agent fix snippet.</p>
      <time datetime="${LAST_UPDATED_ISO}" class="freshness-date">Last updated: ${LAST_UPDATED_ISO}</time>

      <div class="rules-grid">${rulesListHtml}
      </div>
    </main>
  </div>

${buildFooter()}
${buildMobileScript()}
</body>
</html>`;
}

function generateRulePage(rule, categoryKey) {
  const catInfo = CATEGORIES[categoryKey];
  const grouped = groupRulesByCategory();
  const rulesInCategory = grouped[categoryKey] || [];
  const currentIndex = rulesInCategory.findIndex(r => r.key === rule.key);
  const previousRule = currentIndex > 0 ? rulesInCategory[currentIndex - 1] : null;
  const nextRule = currentIndex < rulesInCategory.length - 1 ? rulesInCategory[currentIndex + 1] : null;

  const desc = getDescription(rule.key, rule.label);
  const title = escapeHtml(rule.label);
  const description = `${escapeHtml(desc.short)} Learn what this rule checks, why it matters, and how to fix it with examples and an AI-agent snippet.`;
  const sidebar = buildSidebar(catInfo.slug, rule.key, '../..');

  const baseUrl = 'https://booltools.github.io/booltools-seo-crawler';
  const canonicalPath = `docs/rules/${catInfo.slug}/${rule.key}.html`;
  const jsonLd = {
    '@context': 'https://schema.org',
    '@type': 'TechArticle',
    headline: rule.label,
    description: desc.short,
    url: `${baseUrl}/${canonicalPath}`,
    author: { '@type': 'Organization', name: 'Booltools' },
    publisher: { '@type': 'Organization', name: 'Booltools' },
    dateModified: LAST_UPDATED_ISO,
    isPartOf: { '@type': 'WebSite', name: 'Booltools SEO Crawler', url: baseUrl },
    breadcrumb: {
      '@type': 'BreadcrumbList',
      itemListElement: [
        { '@type': 'ListItem', position: 1, name: 'Rules', item: `${baseUrl}/docs/rules.html` },
        { '@type': 'ListItem', position: 2, name: catInfo.label, item: `${baseUrl}/docs/rules/${catInfo.slug}.html` },
        { '@type': 'ListItem', position: 3, name: rule.label, item: `${baseUrl}/${canonicalPath}` },
      ],
    },
  };

  let prevNextHtml = '';
  if (previousRule || nextRule) {
    prevNextHtml = `
      <nav class="rule-nav" aria-label="Rule navigation">
        <div class="nav-links-row">`;
    if (previousRule) {
      prevNextHtml += `
          <a href="${previousRule.key}.html" class="nav-prev">&larr; ${escapeHtml(previousRule.label)}</a>`;
    }
    if (nextRule) {
      prevNextHtml += `
          <a href="${nextRule.key}.html" class="nav-next">${escapeHtml(nextRule.label)} &rarr;</a>`;
    }
    prevNextHtml += `
        </div>
      </nav>`;
  }

  return `${buildPageHead(title, description, canonicalPath, '../../docs.css', jsonLd)}
${buildHeader('../../../')}

  <div class="docs-page">
${buildMobileMenu()}

${sidebar}
    <main class="docs-content" role="main">
      <article>
      <nav class="breadcrumb" aria-label="Breadcrumb">
        <a href="../../rules.html">Rules</a>
        <span class="sep">/</span>
        <a href="../${catInfo.slug}.html">${catInfo.label}</a>
        <span class="sep">/</span>
        <span>${escapeHtml(rule.label)}</span>
      </nav>

      <div class="rule-title-row">
        <h1>${escapeHtml(rule.label)}</h1>
        <span class="severity-badge sev-${rule.severity}">${rule.severity}</span>
      </div>

      <p class="rule-key-display"><code>${rule.key}</code></p>
      <time datetime="${LAST_UPDATED_ISO}" class="freshness-date">Last updated: ${LAST_UPDATED_ISO}</time>

      <p class="short-desc">${escapeHtml(desc.short)}</p>

      <h2>What does this rule check?</h2>
      <p>${escapeHtml(desc.full)}</p>

      <h2>Why is this important?</h2>
      <p>${escapeHtml(desc.why)}</p>

      <h2>How to fix</h2>
      <p>${escapeHtml(desc.fix)}</p>

      <h2>Examples</h2>
      <div class="example-grid">
        <div class="example-card example-bad">
          <div class="example-label">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>
            Incorrect
          </div>
          <pre><code>${escapeHtml(desc.bad)}</code></pre>
        </div>
        <div class="example-card example-good">
          <div class="example-label">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="20 6 9 17 4 12"/></svg>
            Correct
          </div>
          <pre><code>${escapeHtml(desc.good)}</code></pre>
        </div>
      </div>

      <h2>Agent fix snippet</h2>
      <p>Copy this snippet and paste it into your AI coding agent to fix this issue automatically.</p>
      <div class="snippet-container">
        <button class="copy-btn" data-snippet-target="agent-snippet" aria-label="Copy snippet">
          <svg class="copy-icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>
          <span class="copy-text">Copy</span>
        </button>
        <pre id="agent-snippet"><code>${escapeHtml(desc.snippet)}</code></pre>
      </div>
${prevNextHtml}
      </article>
    </main>
  </div>

${buildFooter()}
${buildMobileScript()}
${buildCopyScript()}
</body>
</html>`;
}

function main() {
  const docsDir = path.join(__dirname, 'docs');
  const rulesDir = path.join(docsDir, 'rules');

  fs.mkdirSync(rulesDir, { recursive: true });

  const grouped = groupRulesByCategory();
  let categoryCount = 0;
  let ruleCount = 0;

  for (const [categoryKey, catInfo] of Object.entries(CATEGORIES)) {
    const categorySlugDir = path.join(rulesDir, catInfo.slug);
    fs.mkdirSync(categorySlugDir, { recursive: true });

    const categoryHtml = generateCategoryPage(categoryKey);
    fs.writeFileSync(path.join(rulesDir, `${catInfo.slug}.html`), categoryHtml);
    categoryCount++;

    const rules = grouped[categoryKey] || [];
    for (const rule of rules) {
      const ruleHtml = generateRulePage(rule, categoryKey);
      fs.writeFileSync(path.join(categorySlugDir, `${rule.key}.html`), ruleHtml);
      ruleCount++;
    }
  }

  console.log(`Generated ${categoryCount} category pages and ${ruleCount} rule pages.`);
}

main();
