package template

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path/filepath"

	"captcha-service/internal/domain/entity"
	"captcha-service/internal/domain/interfaces"
	"captcha-service/pkg/logger"

	"go.uber.org/zap"
)

// TemplateEngineService implements the TemplateEngine interface
type TemplateEngineService struct {
	templates map[string]*template.Template
	basePath  string
}

// NewTemplateEngineService creates a new template engine service
func NewTemplateEngineService(basePath string) *TemplateEngineService {
	tes := &TemplateEngineService{
		templates: make(map[string]*template.Template),
		basePath:  basePath,
	}

	// Load templates
	tes.loadTemplates()

	return tes
}

// loadTemplates loads all templates from the base path
func (tes *TemplateEngineService) loadTemplates() {
	templateFiles := map[string]string{
		"slider_puzzle": "slider_puzzle.html",
		"drag_drop":     "drag_drop.html",
		"blocked":       "blocked.html",
		"demo":          "demo.html",
	}

	for name, filename := range templateFiles {
		tmpl, err := template.ParseFiles(filepath.Join(tes.basePath, filename))
		if err != nil {
			logger.Error("Failed to load template",
				zap.String("template", name),
				zap.String("file", filename),
				zap.Error(err))
			continue
		}

		tes.templates[name] = tmpl
		logger.Debug("Loaded template", zap.String("template", name))
	}
}

// RenderSliderPuzzle renders the slider puzzle template
func (tes *TemplateEngineService) RenderSliderPuzzle(ctx context.Context, challenge *entity.Challenge) (string, error) {
	tmpl, exists := tes.templates["slider_puzzle"]
	if !exists {
		return "", fmt.Errorf("slider_puzzle template not found")
	}

	data := map[string]interface{}{
		"challenge": challenge,
	}

	// Customize data for slider puzzle
	data["PuzzleWidth"] = 60
	data["PuzzleHeight"] = 60
	data["PuzzleShape"] = "square"

	return tes.renderTemplate(tmpl, data)
}

// RenderDragDrop renders the drag and drop template
func (tes *TemplateEngineService) RenderDragDrop(ctx context.Context, challenge *entity.Challenge) (string, error) {
	tmpl, exists := tes.templates["drag_drop"]
	if !exists {
		return "", fmt.Errorf("drag_drop template not found")
	}

	data := map[string]interface{}{
		"challenge": challenge,
	}

	// Customize data for drag and drop
	data["PuzzleWidth"] = 40
	data["PuzzleHeight"] = 40
	data["PuzzleShape"] = "circle"

	return tes.renderTemplate(tmpl, data)
}

// RenderBlockedPage renders the blocked user page
func (tes *TemplateEngineService) RenderBlockedPage(ctx context.Context, userID, reason string) (string, error) {
	tmpl, exists := tes.templates["blocked"]
	if !exists {
		return "", fmt.Errorf("blocked template not found")
	}

	data := map[string]interface{}{
		"user_id": userID,
		"reason":  reason,
	}

	return tes.renderTemplate(tmpl, data)
}

// RenderDemoPage renders the demo page
func (tes *TemplateEngineService) RenderDemoPage(ctx context.Context, data *entity.DemoData) (string, error) {
	tmpl, exists := tes.templates["demo"]
	if !exists {
		return "", fmt.Errorf("demo template not found")
	}

	templateData := map[string]interface{}{
		"user_id":      data.UserID,
		"session_id":   data.SessionID,
		"challenge_id": data.ChallengeID,
		"html":         data.HTML,
	}

	return tes.renderTemplate(tmpl, templateData)
}

// renderTemplate renders a template with the given data
func (tes *TemplateEngineService) renderTemplate(tmpl *template.Template, data interface{}) (string, error) {
	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, data); err != nil {
		logger.Error("Failed to execute template", zap.Error(err))
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// ReloadTemplates reloads all templates from disk
func (tes *TemplateEngineService) ReloadTemplates() error {
	logger.Info("Reloading templates")
	tes.templates = make(map[string]*template.Template)
	tes.loadTemplates()
	return nil
}

// GetTemplateNames returns the names of all loaded templates
func (tes *TemplateEngineService) GetTemplateNames() []string {
	names := make([]string, 0, len(tes.templates))
	for name := range tes.templates {
		names = append(names, name)
	}
	return names
}

// Render renders a template with the given data
func (tes *TemplateEngineService) Render(templateName string, data interface{}) (string, error) {
	tmpl, exists := tes.templates[templateName]
	if !exists {
		return "", fmt.Errorf("template %s not found", templateName)
	}
	return tes.renderTemplate(tmpl, data)
}

// Ensure TemplateEngineService implements the interface
var _ interfaces.TemplateEngine = (*TemplateEngineService)(nil)
