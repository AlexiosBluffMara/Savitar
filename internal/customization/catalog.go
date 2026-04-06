package customization

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Agent struct {
	Name          string
	Description   string
	Path          string
	Tools         []string
	UserInvocable bool
}

type Skill struct {
	Name         string
	Description  string
	ArgumentHint string
	Path         string
}

func DiscoverAgents(rootDir string) ([]Agent, error) {
	dirPath := filepath.Join(rootDir, ".github", "agents")
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	agents := make([]Agent, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".agent.md") {
			continue
		}

		path := filepath.Join(dirPath, entry.Name())
		frontmatter, err := parseFrontmatter(path)
		if err != nil {
			return nil, err
		}

		name := frontmatter["name"]
		if name == "" {
			name = strings.TrimSuffix(entry.Name(), ".agent.md")
		}

		agents = append(agents, Agent{
			Name:          name,
			Description:   frontmatter["description"],
			Path:          path,
			Tools:         parseStringList(frontmatter["tools"]),
			UserInvocable: parseBool(frontmatter["user-invocable"], true),
		})
	}

	sort.Slice(agents, func(left, right int) bool {
		return agents[left].Name < agents[right].Name
	})

	return agents, nil
}

func DiscoverSkills(rootDir string) ([]Skill, error) {
	dirPath := filepath.Join(rootDir, ".github", "skills")
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	skills := make([]Skill, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		path := filepath.Join(dirPath, entry.Name(), "SKILL.md")
		frontmatter, err := parseFrontmatter(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}

			return nil, err
		}

		name := frontmatter["name"]
		if name == "" {
			name = entry.Name()
		}

		skills = append(skills, Skill{
			Name:         name,
			Description:  frontmatter["description"],
			ArgumentHint: frontmatter["argument-hint"],
			Path:         path,
		})
	}

	sort.Slice(skills, func(left, right int) bool {
		return skills[left].Name < skills[right].Name
	})

	return skills, nil
}
