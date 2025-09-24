package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"captcha-service/internal/domain/dto"
	"captcha-service/internal/domain/entity"
	"captcha-service/internal/domain/interfaces"
)

type SliderPuzzleGenerator struct {
	config         *entity.Config
	repo           interfaces.ChallengeRepository
	templateEngine interfaces.TemplateEngine
	rand           *rand.Rand
}

func NewSliderPuzzleGenerator(config *entity.Config, repo interfaces.ChallengeRepository, templateEngine interfaces.TemplateEngine) *SliderPuzzleGenerator {
	return &SliderPuzzleGenerator{
		config:         config,
		repo:           repo,
		templateEngine: templateEngine,
		rand:           rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (g *SliderPuzzleGenerator) Generate(ctx context.Context, complexity int32, userID string) (*entity.Challenge, error) {
	if complexity < 0 || complexity > 100 {
		complexity = g.config.ComplexityMedium
	}

	canvasWidth := entity.CanvasWidth

	var puzzleSize int

	if complexity <= g.config.ComplexityLow {
		puzzleSize = int(g.config.PuzzleSizeLow)
	} else if complexity <= g.config.ComplexityMedium {
		puzzleSize = int(g.config.PuzzleSizeMedium)
	} else {
		puzzleSize = int(g.config.PuzzleSizeHigh)
	}

	// Generate random position for the puzzle piece
	_ = g.rand.Intn(canvasWidth - puzzleSize)
	_ = g.rand.Intn(entity.CanvasHeight - puzzleSize - entity.PuzzleGapTop) // Avoid top area

	challengeID := fmt.Sprintf("slider_%d", time.Now().UnixNano())

	challenge := &entity.Challenge{
		ID:         challengeID,
		Type:       entity.ChallengeTypeSliderPuzzle,
		UserID:     userID,
		Complexity: complexity,
		Data: entity.SliderPuzzleData{
			ChallengeData: dto.ChallengeData{},
			CanvasWidth:   entity.CanvasWidth,
			CanvasHeight:  entity.CanvasHeight,
		},
		ExpiresAt:          time.Now().Add(time.Duration(g.config.ExpirationTimeMedium) * time.Second),
		CreatedAt:          time.Now(),
		Attempts:           0,
		MaxAttempts:        g.config.MaxAttempts,
		MinTime:            int64(g.config.MinTimeMs),
		MaxTime:            int64(g.config.MaxTimeMs),
		IsBlocked:          false,
		BlockReason:        "",
		TimeoutAttempts:    0,
		MaxTimeoutAttempts: g.config.MaxTimeoutAttempts,
	}

	return challenge, nil
}

func (g *SliderPuzzleGenerator) Validate(answer interface{}, data interface{}) (bool, int32, error) {
	answerMap, ok := answer.(map[string]interface{})
	if !ok {
		return false, 0, fmt.Errorf("неверный формат ответа")
	}

	// Извлекаем координаты ответа
	answeredX, xOk := answerMap["x"].(float64)
	answeredY, yOk := answerMap["y"].(float64)
	if !xOk || !yOk {
		return false, 0, fmt.Errorf("неверные координаты ответа")
	}

	// Извлекаем данные челленджа
	_, ok = data.(entity.SliderPuzzleData)
	if !ok {
		return false, 0, fmt.Errorf("неверный формат данных челленджа")
	}

	// Для демо, используем фиксированные значения
	targetX := 200
	targetY := 150
	tolerance := 10

	// Проверяем, находится ли ответ в пределах допуска
	diffX := int(answeredX) - targetX
	diffY := int(answeredY) - targetY

	isValid := (diffX >= -tolerance && diffX <= tolerance) && (diffY >= -tolerance && diffY <= tolerance)

	// Для демо, уверенность фиксированная
	confidence := int32(85)

	return isValid, confidence, nil
}

func (g *SliderPuzzleGenerator) GenerateHTML(challenge *entity.Challenge) (string, error) {
	if g.templateEngine == nil {
		return "", fmt.Errorf("template engine not initialized")
	}

	return g.templateEngine.Render("slider_puzzle", challenge)
}

func (g *SliderPuzzleGenerator) convertToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", value)
	}
}
