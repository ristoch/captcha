package http

import (
	"captcha-service/internal/config"
	"captcha-service/internal/domain/entity"
	"captcha-service/internal/service"
	"captcha-service/internal/usecase"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type DemoHandler struct {
	usecase        *usecase.DemoUsecase
	captchaService *service.CaptchaService
	tmpl           *template.Template
	config         *config.DemoConfig
}

func NewDemoHandler(usecase *usecase.DemoUsecase, captchaService *service.CaptchaService, tmpl *template.Template, config *config.DemoConfig) *DemoHandler {
	return &DemoHandler{
		usecase:        usecase,
		captchaService: captchaService,
		tmpl:           tmpl,
		config:         config,
	}
}

func (h *DemoHandler) HandleDemo(w http.ResponseWriter, r *http.Request) {
	complexityStr := r.URL.Query().Get("complexity")
	complexity := int(h.config.DefaultComplexity)
	if c, err := strconv.Atoi(complexityStr); err == nil {
		if c >= 0 && c <= 100 {
			complexity = c
		} else {
			log.Printf("Invalid complexity %d, using default %d", c, complexity)
		}
	}

	userID := "demo_user"
	if uid := r.URL.Query().Get("user_id"); uid != "" {
		userID = uid
	}

	isBlocked, err := h.usecase.IsUserBlocked(userID)
	if err != nil {
		log.Printf("Error checking user block status: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if isBlocked {
		http.Error(w, "User is blocked", http.StatusTooManyRequests)
		return
	}

	session, err := h.usecase.GetOrCreateSession(userID)
	if err != nil {
		log.Printf("Error getting/creating session: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.setSessionCookie(w, session.SessionID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	challenge, err := h.captchaService.CreateChallenge(ctx, entity.ChallengeTypeSliderPuzzle, int32(complexity), userID)
	if err != nil {
		log.Printf("Error creating challenge: %v", err)
		http.Error(w, "Failed to create challenge", http.StatusInternalServerError)
		return
	}

	captchaHTML := h.createRealChallengeHTML(challenge, userID)

	tmpl, err := template.ParseFiles("./templates/demo.html")
	if err != nil {
		log.Printf("Error loading demo template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		ChallengeID string
		Complexity  int
		HTMLSize    int
		UserID      string
		HTML        string
		MaxAttempts int32
	}{
		ChallengeID: challenge.ID,
		Complexity:  complexity,
		HTMLSize:    len(captchaHTML),
		UserID:      userID,
		HTML:        captchaHTML,
		MaxAttempts: challenge.MaxAttempts,
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *DemoHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"version":   "1.0.0",
		"service":   "demo",
	}

	json.NewEncoder(w).Encode(health)
}

func (h *DemoHandler) setSessionCookie(w http.ResponseWriter, sessionID string) {
	cookie := &http.Cookie{
		Name:     "captcha_session",
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}

func (h *DemoHandler) createRealChallengeHTML(challenge *entity.Challenge, userID string) string {
	rand.Seed(time.Now().UnixNano())
	backgroundImage := entity.BackgroundImages[rand.Intn(len(entity.BackgroundImages))]
	puzzleShape := entity.PuzzleShapes[rand.Intn(len(entity.PuzzleShapes))]

	targetX := int(h.config.DefaultComplexity)
	targetY := 50
	puzzleWidth := 60
	puzzleHeight := 60

	if sliderData, ok := challenge.Data.(entity.SliderPuzzleData); ok {
		targetX = sliderData.ChallengeData.TargetX
		targetY = sliderData.ChallengeData.TargetY

		complexity := int(challenge.Complexity)
		minSize := 40
		maxSize := 80
		puzzleSize := minSize + (complexity-1)*(maxSize-minSize)/99
		puzzleWidth = puzzleSize
		puzzleHeight = puzzleSize
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="ru">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width,initial-scale=1" />
  <title>Капча — Слайдер-пазл</title>
  <style>
    :root { 
      --w: 400px; 
      --h: 300px; 
      --pz: %dpx; 
      --gap-top: 50px; 
    }
    * { box-sizing: border-box; }
    body { margin: 0; font-family: system-ui, -apple-system, Arial, sans-serif; background: #f5f5f5; color:#222; }
    .wrap { max-width: 480px; margin: 24px auto; background: #fff; border-radius: 10px; padding: 16px; box-shadow: 0 6px 20px rgba(0,0,0,.08); }
    h1 { font-size: 18px; margin: 0 0 10px; text-align:center; }
    .canvas-box { position: relative; width: var(--w); height: var(--h); margin: 12px auto; border: 1px solid #e5e5e5; border-radius: 8px; overflow: hidden; background: #fafafa; }
    canvas { display:block; width:100%%; height:100%%; }
    .piece { position:absolute; top: var(--gap-top); left: 0; width: var(--pz); height: var(--pz); background:transparent; box-shadow: 0 4px 10px rgba(0,0,0,.15); display:grid; place-items:center; }
    .ctrl { width:100%%; margin: 12px auto 0; }
    .slider-container { margin: 8px 0; }
    .slider-container label { display: block; margin-bottom: 4px; font-size: 14px; }
    input[type="range"] { width:100%%; height: 34px; -webkit-appearance:none; appearance:none; background:#e9ecef; border-radius: 999px; outline: none; }
    input[type="range"]::-webkit-slider-thumb { -webkit-appearance:none; width:34px; height:34px; border-radius:50%%; background:#1976d2; border:3px solid #fff; box-shadow: 0 2px 6px rgba(0,0,0,.25); cursor:pointer; }
    input[type="range"]::-moz-range-thumb { width:34px; height:34px; border-radius:50%%; background:#1976d2; border:3px solid #fff; box-shadow: 0 2px 6px rgba(0,0,0,.25); cursor:pointer; }
    .msg { text-align:center; margin:10px 0 0; font-weight:600; }
    .ok { color:#1b5e20; }
    .bad { color:#b00020; }
    .hint { text-align:center; font-size:13px; color:#666; margin-top:6px; }
    .noselect { user-select: none; -webkit-user-select:none; }
  </style>
</head>
<body>
  <div class="wrap noselect">
    <h1>Переместите слайдер для решения капчи</h1>
    <div class="canvas-box">
      <canvas id="cv" width="400" height="300" aria-label="captcha"></canvas>
      <div id="piece" class="piece" aria-hidden="true"></div>
    </div>
           <div class="ctrl">
             <div class="slider-container">
               <label>X: <span id="x-value">0</span></label>
               <input id="slider-x" type="range" min="0" max="340" value="0" />
             </div>
           </div>
    <div id="msg" class="msg"></div>
  </div>

  <script>
    const challengeData = {
      challenge_id: "%s",
      user_id: "%s",
      canvas_width: 400,
      canvas_height: 300,
      puzzle_width: %d,
      puzzle_height: %d,
      target_x: %d,
      target_y: %d,
      tolerance: 15,
      background_image: "http://localhost:8081/backgrounds/%s",
      puzzle_shape: "%s"
    };

    let attempts = 0;
    const MAX_ATTEMPTS = %d;
    let isBlocked = false;

    const cv = document.getElementById("cv");
    const ctx = cv.getContext("2d");
    const pieceEl = document.getElementById("piece");
           const sliderX = document.getElementById("slider-x");
           const xValue = document.getElementById("x-value");
    const msg = document.getElementById("msg");

    let bgImage = null;

    function loadImage() {
      const img = new Image();
      img.crossOrigin = "anonymous";
      img.referrerPolicy = "no-referrer";

      img.onload = () => {
        bgImage = img;
        drawAll(0);
      };
      img.onerror = () => {
        bgImage = null;
        drawAll(0);
      };
      
      img.src = challengeData.background_image;
    }

    function drawBackground() {
      ctx.clearRect(0, 0, cv.width, cv.height);
      if (bgImage) {
        ctx.drawImage(bgImage, 0, 0, bgImage.width, bgImage.height, 0, 0, cv.width, cv.height);
      } else {
        ctx.fillStyle = "#ffffff";
        ctx.fillRect(0, 0, cv.width, cv.height);
        
        ctx.fillStyle = "#e0e0e0";
        for (let i = 0; i < cv.width; i += 20) {
          for (let j = 0; j < cv.height; j += 20) {
            if ((i + j) %% 40 === 0) {
              ctx.fillRect(i, j, 10, 10);
            }
          }
        }
      }
    }

    function drawTargetHole() {
      const { puzzle_width: w, puzzle_height: h, target_x: x, target_y: y } = challengeData;
      
      ctx.fillStyle = "rgba(0,0,0,.28)";
      ctx.fillRect(x, y, w, h);
      
      ctx.lineWidth = 2;
      ctx.strokeStyle = "#1976d2";
      ctx.strokeRect(x + 1, y + 1, w - 2, h - 2);
    }

    function drawPieceAt(x) {
      pieceEl.style.left = x + "px";
      pieceEl.style.top = challengeData.target_y + "px";

      pieceEl.replaceChildren();
      const pz = document.createElement("canvas");
      pz.width = challengeData.puzzle_width;
      pz.height = challengeData.puzzle_height;
      const pctx = pz.getContext("2d");

      if (bgImage && bgImage.naturalWidth > 0) {
        const scaleX = bgImage.width / cv.width;
        const scaleY = bgImage.height / cv.height;
        const srcX = challengeData.target_x * scaleX;
        const srcY = challengeData.target_y * scaleY;
        const srcW = challengeData.puzzle_width * scaleX;
        const srcH = challengeData.puzzle_height * scaleY;

        pctx.drawImage(bgImage, srcX, srcY, srcW, srcH, 0, 0, pz.width, pz.height);
      } else {
        pctx.fillStyle = "#ffffff";
        pctx.fillRect(0, 0, pz.width, pz.height);
        
        pctx.fillStyle = "#e0e0e0";
        for (let i = 0; i < pz.width; i += 20) {
          for (let j = 0; j < pz.height; j += 20) {
            if ((i + j) %% 40 === 0) {
              pctx.fillRect(i, j, 10, 10);
            }
          }
        }
      }

      pctx.globalCompositeOperation = 'destination-in';
      pctx.fillRect(0, 0, pz.width, pz.height);
      pctx.globalCompositeOperation = 'source-over';

      pctx.lineWidth = 2;
      pctx.strokeStyle = "#1976d2";
      pctx.strokeRect(0, 0, pz.width, pz.height);

      pieceEl.appendChild(pz);
    }

    function drawAll(x) {
      drawBackground();
      drawTargetHole();
      drawPieceAt(x);
    }

    function setMsg(text, kind) {
      msg.textContent = text || "";
      msg.className = "msg" + (kind ? " " + kind : "");
    }

    function createNewChallenge() {
      if (isBlocked) return;
      
      challengeData.challenge_id = "mock_challenge_" + Date.now();
      challengeData.target_x = Math.floor(Math.random() * (400 - 60)) + 30;
      
      sliderX.value = 0;
      updateDisplay();
      
      setMsg("Новая капча загружена. Попробуйте снова.", "hint");
      
      if (window.top && window.top !== window) {
        window.top.postMessage({
          type: 'captcha:sendData',
          challengeId: challengeData.challenge_id,
          userId: challengeData.user_id,
          eventType: 'newChallenge',
          data: {
            challengeId: challengeData.challenge_id,
            targetX: challengeData.target_x,
            attempts: attempts,
            timestamp: Date.now()
          }
        }, '*');
      }
    }

    function check() {
      if (isBlocked) {
        setMsg("Вы заблокированы за превышение лимита попыток.", "bad");
        return;
      }

      const x = parseInt(sliderX.value);
      const distX = Math.abs(x - challengeData.target_x);
      const isCorrect = distX <= challengeData.tolerance;
      
      attempts++;
      
      if (isCorrect) {
        setMsg("Успешно! Капча пройдена.", "ok");
        if (window.top && window.top !== window) {
          window.top.postMessage({
            type: 'captcha:sendData',
            challengeId: challengeData.challenge_id,
            userId: challengeData.user_id,
            eventType: 'captchaSolved',
              data: {
                positionX: x,
                distanceX: distX,
                attempts: attempts,
                timestamp: Date.now()
              }
          }, '*');
        }
      } else {
        if (attempts >= MAX_ATTEMPTS) {
          isBlocked = true;
          setMsg("Превышен лимит попыток. Вы заблокированы.", "bad");
          if (window.top && window.top !== window) {
            window.top.postMessage({
              type: 'captcha:sendData',
              challengeId: challengeData.challenge_id,
              userId: challengeData.user_id,
              eventType: 'userBlocked',
              data: {
                positionX: x,
                distanceX: distX,
                targetX: challengeData.target_x,
                attempts: attempts,
                maxAttempts: MAX_ATTEMPTS,
                timestamp: Date.now()
              }
            }, '*');
          }
        } else {
          setMsg("Неверно. Попытка " + attempts + "/" + MAX_ATTEMPTS + ". Попробуйте ещё раз.", "bad");
          if (window.top && window.top !== window) {
            window.top.postMessage({
              type: 'captcha:sendData',
              challengeId: challengeData.challenge_id,
              userId: challengeData.user_id,
              eventType: 'captchaFailed',
              data: {
                positionX: x,
                distanceX: distX,
                targetX: challengeData.target_x,
                attempts: attempts,
                maxAttempts: MAX_ATTEMPTS,
                timestamp: Date.now()
              }
            }, '*');
          }
          setTimeout(createNewChallenge, 1000);
        }
      }
    }

    function updateDisplay() {
      const x = parseInt(sliderX.value) || 0;
      xValue.textContent = x;
      drawAll(x);
      setMsg("");
    }

    sliderX.addEventListener("input", updateDisplay);
    sliderX.addEventListener("change", check);

    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', () => {
        loadImage();
        updateDisplay();
      });
    } else {
      loadImage();
      updateDisplay();
    }
    
    if (window.top && window.top !== window) {
      window.top.postMessage({
        type: 'captcha:sendData',
        challengeId: challengeData.challenge_id,
        userId: challengeData.user_id,
        eventType: 'captchaLoaded',
        data: {
          timestamp: Date.now()
        }
      }, '*');
    }

    window.addEventListener('message', function(event) {
      if (event.data.type === 'new_challenge_data') {
        console.log('Received new challenge data:', event.data);
        
        if (event.data.background_image) {
          challengeData.background_image = event.data.background_image;
        }
        if (event.data.puzzle_shape) {
          challengeData.puzzle_shape = event.data.puzzle_shape;
        }
        if (event.data.target_x !== undefined) {
          challengeData.target_x = event.data.target_x;
        }
        if (event.data.challenge_id) {
          challengeData.challenge_id = event.data.challenge_id;
        }
        
        sliderX.value = 0;
        updateDisplay();
        setMsg("Новая капча загружена. Попробуйте снова.", "hint");
        
        loadImage();
      }
    });
  </script>
  </body>
  </html>`, puzzleWidth, challenge.ID, userID, puzzleWidth, puzzleHeight, targetX, targetY, backgroundImage, puzzleShape, challenge.MaxAttempts)
}
