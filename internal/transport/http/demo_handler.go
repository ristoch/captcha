package http

import (
	"captcha-service/internal/usecase"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type DemoHandler struct {
	usecase *usecase.DemoUsecase
	tmpl    *template.Template
}

func NewDemoHandler(usecase *usecase.DemoUsecase, tmpl *template.Template) *DemoHandler {
	return &DemoHandler{
		usecase: usecase,
		tmpl:    tmpl,
	}
}

func (h *DemoHandler) HandleDemo(w http.ResponseWriter, r *http.Request) {
	complexityStr := r.URL.Query().Get("complexity")
	complexity := 50
	if c, err := strconv.Atoi(complexityStr); err == nil {
		complexity = c
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

	challenge, err := h.createMockChallenge(userID, complexity)
	if err != nil {
		log.Printf("Error creating challenge: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(challenge))
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

func (h *DemoHandler) createMockChallenge(userID string, complexity int) (string, error) {
	challengeID := fmt.Sprintf("mock_challenge_%d", time.Now().UnixNano())
	targetX := 200

	html := fmt.Sprintf(`
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
      --pz: 60px; 
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
      <input id="slider" type="range" min="0" max="340" value="0" />
    </div>
    <div id="msg" class="msg"></div>
  </div>

  <script>
    const challengeData = {
      challenge_id: "%s",
      user_id: "%s",
      canvas_width: 400,
      canvas_height: 300,
      puzzle_width: 60,
      puzzle_height: 60,
      target_x: %d,
      tolerance: 15,
      background_image: "http://localhost:8081/backgrounds/background1.png",
      puzzle_shape: "square"
    };

    let attempts = 0;
    const MAX_ATTEMPTS = 3;
    let isBlocked = false;

    const cv = document.getElementById("cv");
    const ctx = cv.getContext("2d");
    const pieceEl = document.getElementById("piece");
    const slider = document.getElementById("slider");
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
      const y = 50;
      const { puzzle_width: w, puzzle_height: h, target_x: x } = challengeData;
      
      ctx.fillStyle = "rgba(0,0,0,.28)";
      ctx.fillRect(x, y, w, h);
      
      ctx.lineWidth = 2;
      ctx.strokeStyle = "#1976d2";
      ctx.strokeRect(x + 1, y + 1, w - 2, h - 2);
    }

    function drawPieceAt(x) {
      pieceEl.style.left = x + "px";
      pieceEl.style.top = "50px";

      pieceEl.replaceChildren();
      const pz = document.createElement("canvas");
      pz.width = challengeData.puzzle_width;
      pz.height = challengeData.puzzle_height;
      const pctx = pz.getContext("2d");

      if (bgImage && bgImage.naturalWidth > 0) {
        const scaleX = bgImage.width / cv.width;
        const scaleY = bgImage.height / cv.height;
        const srcX = challengeData.target_x * scaleX;
        const srcY = 50 * scaleY;
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
      
      slider.value = 0;
      drawAll(0);
      
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

      const x = parseInt(slider.value);
      const dist = Math.abs(x - challengeData.target_x);
      const isCorrect = dist <= challengeData.tolerance;
      
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
              position: x,
              distance: dist,
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
                position: x,
                distance: dist,
                target: challengeData.target_x,
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
                position: x,
                distance: dist,
                target: challengeData.target_x,
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

    slider.addEventListener("input", () => {
      const x = parseInt(slider.value) || 0;
      drawAll(x);
      setMsg("");
    });

    slider.addEventListener("change", () => {
      check();
    });

    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', loadImage);
    } else {
      loadImage();
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
  </script>
</body>
</html>`, challengeID, userID, targetX)

	return html, nil
}
