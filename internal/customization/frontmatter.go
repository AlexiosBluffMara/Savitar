package customization

import (
	"os"
	"strings"
)

func parseFrontmatter(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	frontmatter := map[string]string{}
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return frontmatter, nil
	}

	for _, line := range lines[1:] {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			break
		}

		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		key, value, ok := strings.Cut(trimmed, ":")
		if !ok {
			continue
		}

		frontmatter[strings.TrimSpace(key)] = trimValue(value)
	}

	return frontmatter, nil
}

func trimValue(value string) string {
	trimmed := strings.TrimSpace(value)
	trimmed = strings.Trim(trimmed, `"'`)
	return trimmed
}

func parseStringList(value string) []string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || trimmed == "[]" {
		return nil
	}

	trimmed = strings.TrimPrefix(trimmed, "[")
	trimmed = strings.TrimSuffix(trimmed, "]")
	trimmed = strings.TrimSpace(trimmed)
	if trimmed == "" {
		return nil
	}

	parts := strings.Split(trimmed, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		item := trimValue(part)
		if item == "" {
			continue
		}

		values = append(values, item)
	}

	return values
}

func parseBool(value string, fallback bool) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true":
		return true
	case "false":
		return false
	default:
		return fallback
	}
}
