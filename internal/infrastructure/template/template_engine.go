package template

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path/filepath"

	"captcha-service/internal/domain/entity"
	"captcha-service/pkg/logger"

	"go.uber.org/zap"
)

type NewTemplateData interface {
	GetData() map[string]interface{}
}

type TemplateEngine interface {
	Render(templateName string, data interface{}) (string, error)
}

type TemplateEngineService struct {
	templates map[string]*template.Template
	basePath  string
}

func NewTemplateEngineService(basePath string) *TemplateEngineService {
	tes := &TemplateEngineService{
		templates: make(map[string]*template.Template),
		basePath:  basePath,
	}

	tes.loadTemplates()

	return tes
}

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

func (tes *TemplateEngineService) RenderSliderPuzzle(ctx context.Context, challenge *entity.Challenge) (string, error) {
	tmpl, exists := tes.templates["slider_puzzle"]
	if !exists {
		return "", fmt.Errorf("slider_puzzle template not found")
	}

	data := map[string]interface{}{
		"challenge": challenge,
	}

	data["PuzzleWidth"] = 60
	data["PuzzleHeight"] = 60
	data["PuzzleShape"] = "square"

	return tes.renderTemplate(tmpl, data)
}

func (tes *TemplateEngineService) RenderDragDrop(ctx context.Context, challenge *entity.Challenge) (string, error) {
	tmpl, exists := tes.templates["drag_drop"]
	if !exists {
		return "", fmt.Errorf("drag_drop template not found")
	}

	data := map[string]interface{}{
		"challenge": challenge,
	}

	data["PuzzleWidth"] = 40
	data["PuzzleHeight"] = 40
	data["PuzzleShape"] = "circle"

	return tes.renderTemplate(tmpl, data)
}

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

func (tes *TemplateEngineService) RenderDemoPage(ctx context.Context, data *entity.DemoData) (string, error) {
	tmpl, exists := tes.templates["demo"]
	if !exists {
		return "", fmt.Errorf("demo template not found")
	}

	templateData := map[string]interface{}{
		"user_id":               data.UserID,
		"session_id":            data.SessionID,
		entity.FieldChallengeID: data.ChallengeID,
		"html":                  data.HTML,
	}

	return tes.renderTemplate(tmpl, templateData)
}

func (tes *TemplateEngineService) renderTemplate(tmpl *template.Template, data interface{}) (string, error) {
	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, data); err != nil {
		logger.Error("Failed to execute template", zap.Error(err))
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func (tes *TemplateEngineService) ReloadTemplates() error {
	logger.Info("Reloading templates")
	tes.templates = make(map[string]*template.Template)
	tes.loadTemplates()
	return nil
}

func (tes *TemplateEngineService) GetTemplateNames() []string {
	names := make([]string, 0, len(tes.templates))
	for name := range tes.templates {
		names = append(names, name)
	}
	return names
}

func (tes *TemplateEngineService) Render(templateName string, data interface{}) (string, error) {
	tmpl, exists := tes.templates[templateName]
	if !exists {
		return "", fmt.Errorf("template %s not found", templateName)
	}
	return tes.renderTemplate(tmpl, data)
}

var _ TemplateEngine = (*TemplateEngineService)(nil)
