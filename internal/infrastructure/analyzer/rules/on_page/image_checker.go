package on_page

import (
	"fmt"
	"path"
	"strings"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

type ImageChecker struct{}

func (c *ImageChecker) Check(page crawler.PageData) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	images := page.Images
	if len(images) == 0 {
		return rules
	}

	var missingAltURLs []string
	var nonDescriptiveAltURLs []string
	var longAltURLs []string
	var missingDimensionURLs []string
	var nonModernFormatURLs []string
	var missingLazyLoadURLs []string

	for _, image := range images {
		if image.Alt == "" {
			missingAltURLs = append(missingAltURLs, image.URL)
		} else {
			if isNonDescriptiveAlt(image.Alt) {
				nonDescriptiveAltURLs = append(nonDescriptiveAltURLs, fmt.Sprintf("%s (alt=\"%s\")", image.URL, image.Alt))
			}
			if len(image.Alt) > 125 {
				longAltURLs = append(longAltURLs, image.URL)
			}
		}

		if image.Width == "" || image.Height == "" {
			missingDimensionURLs = append(missingDimensionURLs, image.URL)
		}

		if image.URL != "" && !isModernImageFormat(image.URL) && !image.HasPictureSource {
			nonModernFormatURLs = append(nonModernFormatURLs, image.URL)
		}

		isAboveFold := image.FetchPriority == "high" || image.Loading == "eager"
		if image.Loading != "lazy" && !isAboveFold {
			missingLazyLoadURLs = append(missingLazyLoadURLs, image.URL)
		}
	}

	altRule := valueobject.NewAuditRule("images_alt_text", valueobject.CategoryOnPage, valueobject.SeverityHigh)
	altRule.AffectedURL = page.URL
	if len(missingAltURLs) > 0 {
		altRule.Fail(
			fmt.Sprintf("%d of %d images are missing alt attributes", len(missingAltURLs), len(images)),
			"Add descriptive alt text to all images. Alt text should describe the image content and include relevant keywords when natural.",
		)
		altRule.WithDetails(formatAssetList(missingAltURLs))
	} else {
		altRule.Pass("All images have alt attributes")
	}
	rules = append(rules, altRule)

	descriptiveRule := valueobject.NewAuditRule("images_alt_descriptive", valueobject.CategoryOnPage, valueobject.SeverityLow)
	descriptiveRule.AffectedURL = page.URL
	if len(nonDescriptiveAltURLs) > 0 {
		descriptiveRule.Warn(
			fmt.Sprintf("%d images have non-descriptive alt text (e.g., 'image.jpg', 'photo')", len(nonDescriptiveAltURLs)),
			"Replace generic alt text with descriptive text that explains the image content.",
		)
		descriptiveRule.WithDetails(formatAssetList(nonDescriptiveAltURLs))
	} else {
		descriptiveRule.Pass("Image alt text is descriptive")
	}
	rules = append(rules, descriptiveRule)

	dimensionsRule := valueobject.NewAuditRule("images_dimensions", valueobject.CategoryOnPage, valueobject.SeverityMedium)
	dimensionsRule.AffectedURL = page.URL
	if len(missingDimensionURLs) > 0 {
		dimensionsRule.Warn(
			fmt.Sprintf("%d images are missing width/height attributes", len(missingDimensionURLs)),
			"Add explicit width and height attributes to images to prevent Cumulative Layout Shift (CLS).",
		)
		dimensionsRule.WithDetails(formatAssetList(missingDimensionURLs))
	} else {
		dimensionsRule.Pass("All images have width and height attributes")
	}
	rules = append(rules, dimensionsRule)

	formatRule := valueobject.NewAuditRule("images_modern_format", valueobject.CategoryOnPage, valueobject.SeverityLow)
	formatRule.AffectedURL = page.URL
	if len(nonModernFormatURLs) > 0 {
		formatRule.Warn(
			fmt.Sprintf("%d images use legacy formats based on URL extension", len(nonModernFormatURLs)),
			"Convert images to WebP or AVIF format. Note: frameworks like Next.js Image serve modern formats via content negotiation even if the URL path shows .jpg/.png — this may be a false positive.",
		)
		formatRule.WithDetails(formatAssetList(nonModernFormatURLs))
	} else {
		formatRule.Pass("All images use modern formats")
	}
	rules = append(rules, formatRule)

	lazyRule := valueobject.NewAuditRule("images_lazy_loading", valueobject.CategoryOnPage, valueobject.SeverityLow)
	lazyRule.AffectedURL = page.URL
	if len(missingLazyLoadURLs) > len(images)/2 && len(images) > 2 {
		lazyRule.Warn(
			fmt.Sprintf("%d of %d images are not lazy-loaded", len(missingLazyLoadURLs), len(images)),
			"Add loading=\"lazy\" to below-fold images to improve initial page load performance.",
		)
		lazyRule.WithDetails(formatAssetList(missingLazyLoadURLs))
	} else {
		lazyRule.Pass("Images use lazy loading appropriately")
	}
	rules = append(rules, lazyRule)

	return rules
}

func isNonDescriptiveAlt(alt string) bool {
	lowered := strings.ToLower(strings.TrimSpace(alt))
	genericTerms := []string{"image", "photo", "picture", "img", "icon", "logo", "banner", "thumbnail"}

	for _, term := range genericTerms {
		if lowered == term {
			return true
		}
	}

	imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg", ".avif", ".bmp"}
	for _, ext := range imageExtensions {
		if strings.HasSuffix(lowered, ext) {
			return true
		}
	}

	return false
}

var imageOptimizationPatterns = []string{
	"/_next/image",
	"res.cloudinary.com",
	"imgix.net",
	"imagekit.io",
	"cdn.shopify.com",
	"images.unsplash.com",
	"img.clerk.com",
	"twimg.com",
	"googleusercontent.com",
	"wp.com/",
	"i0.wp.com",
	"i1.wp.com",
	"i2.wp.com",
	"cdn.sanity.io",
	"images.ctfassets.net",
	"storyblok.com/f/",
	"cloudfront.net",
	"fastly.net",
	"akamaized.net",
	"imagedelivery.net",
}

func isImageOptimizationURL(imageURL string) bool {
	lowered := strings.ToLower(imageURL)
	for _, pattern := range imageOptimizationPatterns {
		if strings.Contains(lowered, pattern) {
			return true
		}
	}
	return false
}

func isModernImageFormat(imageURL string) bool {
	extension := strings.ToLower(path.Ext(imageURL))
	if extension == ".webp" || extension == ".avif" {
		return true
	}
	return isImageOptimizationURL(imageURL)
}

func formatAssetList(assets []string) string {
	return strings.Join(assets, "\n")
}
