let ws = null;
let currentChallengeId = '';
let userId = '{{.UserID}}'; // Set by server

// Binary data utilities for optimization
const BinaryUtils = {
    // Pack coordinates into single int (max 8192x8192 as per requirements)
    packCoordinates: (x, y) => {
        return (y << 13) | x; // 13 bits for x (max 8191), 13 bits for y (max 8191)
    },
    
    // Unpack coordinates from int
    unpackCoordinates: (packed) => {
        return {
            x: packed & 0x1FFF, // 13 bits
            y: (packed >> 13) & 0x1FFF // 13 bits
        };
    },
    
    // Pack event data for efficient transmission
    packEventData: (eventType, data) => {
        const buffer = new ArrayBuffer(8);
        const view = new DataView(buffer);
        
        // First 4 bytes: event type (0-255)
        view.setUint8(0, eventType);
        
        // Next 4 bytes: packed coordinates or other data
        if (data.x !== undefined && data.y !== undefined) {
            const packed = BinaryUtils.packCoordinates(data.x, data.y);
            view.setUint32(1, packed, true); // little-endian
        }
        
        return buffer;
    }
};

// window.postMessage API for iframe interaction (as per requirements)
function setupPostMessageAPI() {
    // Listen for messages from iframe captcha
    window.addEventListener("message", (e) => {
        if (e.data?.type === "captcha:sendData") {
            handleCaptchaData(e.data.data);
        }
    });
    
    // Send data to iframe captcha
    function sendToCaptcha(data) {
        const iframe = document.getElementById('captcha-iframe');
        if (iframe && iframe.contentWindow) {
            iframe.contentWindow.postMessage({
                type: "captcha:serverData",
                data: data
            }, '*');
        }
    }
    
    // Handle data from captcha iframe
    function handleCaptchaData(data) {
        console.log('Received from captcha:', data);
        
        // Convert to binary for efficient transmission
        if (data.eventType === 'slider_move' && data.data) {
            const binaryData = BinaryUtils.packEventData(1, data.data); // 1 = slider_move
            sendWebSocketEvent('slider_move', binaryData);
        } else if (data.eventType === 'validation' && data.data) {
            const binaryData = BinaryUtils.packEventData(2, data.data); // 2 = validation
            sendWebSocketEvent('validation', binaryData);
        }
    }
    
    // Make functions globally available
    window.sendToCaptcha = sendToCaptcha;
    window.handleCaptchaData = handleCaptchaData;
}

// Send WebSocket event with binary data support
function sendWebSocketEvent(eventType, data) {
    if (!ws || ws.readyState !== WebSocket.OPEN) {
        console.error('WebSocket not connected');
        return;
    }
    
    const message = {
        type: eventType,
        challenge_id: currentChallengeId,
        user_id: userId,
        data: data
    };
    
    // Send as binary if data is ArrayBuffer, otherwise as JSON
    if (data instanceof ArrayBuffer) {
        // Convert message to binary format
        const jsonStr = JSON.stringify({
            type: eventType,
            challenge_id: currentChallengeId,
            user_id: userId
        });
        const jsonBytes = new TextEncoder().encode(jsonStr);
        
        // Combine JSON length + JSON data + binary data
        const totalLength = 4 + jsonBytes.length + data.byteLength;
        const buffer = new ArrayBuffer(totalLength);
        const view = new DataView(buffer);
        
        // Write JSON length (4 bytes)
        view.setUint32(0, jsonBytes.length, true);
        
        // Write JSON data
        new Uint8Array(buffer, 4, jsonBytes.length).set(jsonBytes);
        
        // Write binary data
        new Uint8Array(buffer, 4 + jsonBytes.length, data.byteLength).set(new Uint8Array(data));
        
        ws.send(buffer);
    } else {
        ws.send(JSON.stringify(message));
    }
}

function showStatus(message) {
    console.log('Status:', message);
}

function showError(message) {
    console.error('Error:', message);
}

function hideError() {
    console.log('Error hidden');
}

function showBlockedMessage(message) {
    console.error('User blocked:', message);
    
    // Create blocking overlay
    const overlay = document.createElement('div');
    overlay.style.cssText = `
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: rgba(0, 0, 0, 0.8);
        z-index: 10000;
        display: flex;
        justify-content: center;
        align-items: center;
        font-family: system-ui, -apple-system, Arial, sans-serif;
    `;
    
    const blockDialog = document.createElement('div');
    blockDialog.style.cssText = `
        background: white;
        padding: 30px;
        border-radius: 10px;
        box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
        text-align: center;
        max-width: 400px;
        margin: 20px;
    `;
    
    blockDialog.innerHTML = `
        <div style="font-size: 24px; color: #d32f2f; margin-bottom: 15px;">🚫</div>
        <h2 style="color: #d32f2f; margin: 0 0 15px 0; font-size: 20px;">Доступ заблокирован</h2>
        <p style="color: #666; margin: 0 0 20px 0; line-height: 1.4;">${message}</p>
        <p style="color: #999; font-size: 14px; margin: 0;">Попробуйте снова через несколько минут</p>
    `;
    
    overlay.appendChild(blockDialog);
    document.body.appendChild(overlay);
    
    // Disable all captcha interactions
    const captchaElements = document.querySelectorAll('canvas, input[type="range"], .piece');
    captchaElements.forEach(el => {
        el.style.pointerEvents = 'none';
        el.style.opacity = '0.5';
    });
}

function connectWebSocket() {
    const wsUrl = 'ws://localhost:8081/ws';
    showStatus('Connecting to WebSocket...');
    
    try {
        ws = new WebSocket(wsUrl);
        
        ws.onopen = function() {
            showStatus('WebSocket connected');
            hideError();
            
            // Send initial challenge request
            const message = {
                type: 'challenge_request',
                challenge_id: currentChallengeId,
                user_id: userId
            };
            ws.send(JSON.stringify(message));
        };
        
        ws.onmessage = function(event) {
            try {
                const data = JSON.parse(event.data);
                console.log('Received WebSocket message:', data);
                
                if (data.type === 'new_challenge') {
                    showStatus('New challenge received!');
                    console.log('New challenge received, updating captcha...');
                    
                    // Update the captcha container with new HTML
                    if (data.data && data.data.html) {
                        const captchaContainer = document.querySelector('.captcha-container');
                        if (captchaContainer) {
                            captchaContainer.innerHTML = data.data.html;
                            // Re-initialize the captcha after updating HTML
                            if (typeof initializeCaptcha === 'function') {
                                initializeCaptcha();
                            }
                        }
                    } else {
                        // Fallback: reload the page
                        setTimeout(() => {
                            window.location.reload();
                        }, 1000);
                    }
                } else if (data.type === 'challenge_created') {
                    showStatus('Challenge created successfully');
                } else if (data.type === 'grpc_response') {
                    // Handle gRPC responses, including blocking
                    if (data.data && data.data.blocked) {
                        showBlockedMessage(data.data.error || 'You are blocked due to too many attempts');
                    } else if (data.data && data.data.error) {
                        showError(data.data.error);
                    }
                } else if (data.type === 'error') {
                    showError(data.message || 'Unknown error');
                }
            } catch (e) {
                console.error('Error parsing WebSocket message:', e);
            }
        };
        
        ws.onclose = function() {
            showStatus('WebSocket disconnected');
            // Try to reconnect after 3 seconds
            setTimeout(connectWebSocket, 3000);
        };
        
        ws.onerror = function(error) {
            showError('WebSocket error: ' + error);
        };
        
    } catch (e) {
        showError('Failed to connect to WebSocket: ' + e.message);
    }
}

// Handle messages from captcha iframe
window.addEventListener('message', function(event) {
    if (event.data && event.data.type === 'captcha:sendData') {
        console.log('Captcha event:', event.data);
        
        if (ws && ws.readyState === WebSocket.OPEN) {
            // Forward captcha events to WebSocket
            const message = {
                type: 'captcha_event',
                challenge_id: event.data.challengeId || currentChallengeId,
                user_id: event.data.userId || userId,
                event_type: event.data.eventType,
                data: event.data.data
            };
            ws.send(JSON.stringify(message));
        } else {
            showError('WebSocket not connected');
        }
    }
});

// Start WebSocket connection when page loads
window.addEventListener('load', function() {
    setupPostMessageAPI(); // Initialize postMessage API first
    connectWebSocket();
});
