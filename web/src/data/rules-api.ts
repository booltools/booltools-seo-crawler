export interface RuleDefinition {
  key: string;
  label: string;
  severity: string;
}

export interface RuleCategory {
  key: string;
  label: string;
  rules: RuleDefinition[];
}

export interface RulesApiResponse {
  categories: RuleCategory[];
  presets: unknown[];
}

export interface CategoryNavItem {
  href: string;
  label: string;
  key: string;
  rules: { href: string; label: string; key: string }[];
}

const apiUrl = import.meta.env.PUBLIC_API_URL || 'http://localhost:8080';

export async function fetchRuleCategories(): Promise<{ categories: RuleCategory[]; error: string }> {
  try {
    const response = await fetch(`${apiUrl}/api/rules`);
    if (response.ok) {
      const data: RulesApiResponse = await response.json();
      return { categories: data.categories, error: '' };
    }
    return { categories: [], error: `API returned ${response.status}` };
  } catch {
    return { categories: [], error: 'Could not connect to the API server. Start the backend with `make dev` to see rules.' };
  }
}

export function buildCategoryNav(categories: RuleCategory[]): CategoryNavItem[] {
  return categories.map((cat) => ({
    href: `/docs/rules/${cat.key}`,
    label: cat.label,
    key: cat.key,
    rules: cat.rules.map((rule) => ({
      href: `/docs/rules/${cat.key}/${rule.key}`,
      label: rule.label,
      key: rule.key,
    })),
  }));
}
