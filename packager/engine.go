package packager

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type CapsuleManifest struct {
	CapsuleVersion string            `json:"capsule_version"`
	AppID          string            `json:"app_id"`
	AppName        string            `json:"app_name"`
	CreatedAt      string            `json:"created_at"`
	Runtime        string            `json:"runtime"`
	RuntimeVersion string            `json:"runtime_version"`
	Entrypoint     string            `json:"entrypoint"`
	DocumentRoot   string            `json:"document_root"`
	EnvVariables   map[string]string `json:"env_variables"`
	SecretRefs     []string          `json:"secret_refs"`
}

type Packager struct {
	SourceDir string
}

func NewPackager(sourceDir string) *Packager {
	return &Packager{SourceDir: sourceDir}
}

func (p *Packager) Inspect() (*CapsuleManifest, error) {
	runtime := "static"
	version := "1.0"
	entry := "index.html"
	docRoot := "."

	if _, err := os.Stat(filepath.Join(p.SourceDir, "composer.json")); err == nil {
		runtime = "php"
		version = "8.2"
		entry = "index.php"
		if _, err := os.Stat(filepath.Join(p.SourceDir, "public")); err == nil {
			docRoot = "public"
		}
	} else if _, err := os.Stat(filepath.Join(p.SourceDir, "package.json")); err == nil {
		runtime = "nodejs"
		version = "20"
		entry = "npm start"
	} else if _, err := os.Stat(filepath.Join(p.SourceDir, "requirements.txt")); err == nil {
		runtime = "python"
		version = "3.11"
		entry = "main.py"
	}

	appName := filepath.Base(filepath.Clean(p.SourceDir))
	if appName == "." || appName == "/" {
		appName = "my-app"
	}

	manifest := &CapsuleManifest{
		CapsuleVersion: "1.0",
		AppID:          fmt.Sprintf("cap-%s-%d", strings.ToLower(appName), time.Now().Unix()),
		AppName:        appName,
		CreatedAt:      time.Now().UTC().Format(time.RFC3339),
		Runtime:        runtime,
		RuntimeVersion: version,
		Entrypoint:     entry,
		DocumentRoot:   docRoot,
		EnvVariables:   make(map[string]string),
		SecretRefs:     make([]string, 0),
	}

	// Secret redaction from .env
	envPath := filepath.Join(p.SourceDir, ".env")
	if data, err := os.ReadFile(envPath); err == nil {
		secretPattern := regexp.MustCompile(`(?i)(password|secret|key|token|auth|credential|api_key)`)
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			k := strings.TrimSpace(parts[0])
			v := strings.TrimSpace(parts[1])
			if secretPattern.MatchString(k) {
				manifest.EnvVariables[k] = fmt.Sprintf("<KENPANEL_SECRET_REF:%s>", k)
				manifest.SecretRefs = append(manifest.SecretRefs, k)
			} else {
				manifest.EnvVariables[k] = v
			}
		}
	}

	return manifest, nil
}
