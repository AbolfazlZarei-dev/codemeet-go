package codemeet

const loginHTML = `<!DOCTYPE html>
<html lang="en" dir="ltr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login - CodeMeet Dashboard</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;600;800&display=swap" rel="stylesheet">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.1/css/all.min.css">
    <style>
        body {
            margin: 0;
            padding: 0;
            background-color: #0b1120;
            background-image: radial-gradient(circle at top right, rgba(59, 130, 246, 0.1), transparent 40%), radial-gradient(circle at bottom left, rgba(139, 92, 246, 0.1), transparent 40%);
            color: #f1f5f9;
            font-family: 'Inter', 'Segoe UI', sans-serif;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
        }
        .login-card {
            background: rgba(30, 41, 59, 0.6);
            padding: 40px;
            border-radius: 20px;
            width: 100%;
            max-width: 400px;
            text-align: center;
            box-shadow: 0 10px 40px rgba(0,0,0,0.4);
            backdrop-filter: blur(12px);
            border: 1px solid rgba(51, 65, 85, 0.5);
        }
        .login-card h1 { 
            margin: 0 0 10px 0; 
            font-size: 28px; 
            color: #3b82f6; 
        }
        .login-card p { 
            margin: 0 0 30px 0; 
            color: #94a3b8; 
            font-size: 14px;
        }
        .input-group {
            position: relative;
            margin-bottom: 20px;
        }
        .input-group i {
            position: absolute;
            left: 15px;
            top: 50%;
            transform: translateY(-50%);
            color: #94a3b8;
        }
        .input-group input {
            width: 100%;
            padding: 15px 20px 15px 45px;
            border-radius: 10px;
            border: 1px solid #334155;
            background: #0f172a;
            color: #f1f5f9;
            box-sizing: border-box;
            outline: none;
            transition: all 0.2s;
            font-size: 14px;
        }
        .input-group input:focus {
            border-color: #3b82f6;
            box-shadow: 0 0 10px rgba(59, 130, 246, 0.4);
        }
        .btn-login {
            width: 100%;
            padding: 15px;
            background: linear-gradient(135deg, #3b82f6, #8b5cf6);
            border: none;
            color: white;
            border-radius: 10px;
            font-weight: 700;
            font-size: 16px;
            cursor: pointer;
            transition: all 0.2s;
            margin-top: 10px;
        }
        .btn-login:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 15px rgba(59, 130, 246, 0.4);
        }
    </style>
</head>
<body>
    <div class="login-card">
        <h1><i class="fas fa-shield-halved"></i> Login</h1>
        <p>Access to CodeMeet Dashboard is restricted. Please login.</p>
        {{ERR_MSG}}
        <form method="POST" action="/login">
            <div class="input-group">
                <i class="fas fa-user"></i>
                <input type="text" name="username" placeholder="Username" required autofocus>
            </div>
            <div class="input-group">
                <i class="fas fa-key"></i>
                <input type="password" name="password" placeholder="Password" required>
            </div>
            <button type="submit" class="btn-login">Login</button>
        </form>
    </div>
</body>
</html>`

const dashboardHTML = `<!DOCTYPE html>
<html lang="en" dir="ltr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>CodeMeet Bot Dashboard</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;600;800&family=Vazirmatn:wght@400;700&family=JetBrains+Mono:wght@400;700&display=swap" rel="stylesheet">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.1/css/all.min.css">
    <style>
        :root {
            --bg-color: #0b1120;
            --card-bg: rgba(30, 41, 59, 0.6);
            --card-bg-solid: #1e293b;
            --text-color: #f1f5f9;
            --text-muted: #94a3b8;
            --accent-color: #3b82f6;
            --accent-glow: rgba(59, 130, 246, 0.4);
            --success-color: #10b981;
            --warn-color: #facc15;
            --error-color: #ef4444;
            --border-color: rgba(51, 65, 85, 0.5);
            --log-bg: rgba(2, 6, 23, 0.9);
            --shadow: 0 10px 40px rgba(0,0,0,0.4);
        }

        body.light-mode {
            --bg-color: #f1f5f9;
            --card-bg: rgba(255, 255, 255, 0.9);
            --card-bg-solid: #ffffff;
            --text-color: #0f172a;
            --text-muted: #475569;
            --border-color: rgba(203, 213, 225, 0.8);
            --log-bg: rgba(248, 250, 252, 0.9);
            --shadow: 0 10px 30px rgba(0,0,0,0.05);
        }

        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            background-color: var(--bg-color); 
            background-image: radial-gradient(circle at top right, rgba(59, 130, 246, 0.1), transparent 40%), radial-gradient(circle at bottom left, rgba(139, 92, 246, 0.1), transparent 40%);
            color: var(--text-color); 
            padding: 20px; 
            font-family: 'Inter', 'Vazirmatn', 'Segoe UI', sans-serif; 
            transition: background-color 0.3s ease, color 0.3s ease;
            min-height: 100vh;
        }
        body[dir="rtl"] { font-family: 'Vazirmatn', 'Inter', sans-serif; }

        .container { max-width: 1200px; margin: 0 auto; }
        
        .header { 
            display: flex; justify-content: space-between; align-items: center; margin-bottom: 30px; 
            background: var(--card-bg); padding: 20px 30px; border-radius: 20px; 
            border: 1px solid var(--border-color); 
            backdrop-filter: blur(12px);
            box-shadow: var(--shadow);
        }
        .header h1 { font-size: 24px; font-weight: 800; display: flex; align-items: center; gap: 15px; }
        .header h1 i { color: var(--accent-color); font-size: 28px; }
        .logo { 
            width: 45px; height: 45px; background: linear-gradient(135deg, var(--accent-color), #8b5cf6); 
            border-radius: 12px; display: flex; align-items: center; justify-content: center; 
            font-weight: bold; color: white; font-size: 18px;
            box-shadow: 0 4px 15px var(--accent-glow);
        }
        .header-controls { display: flex; align-items: center; gap: 15px; }
        
        select, .theme-toggle, .btn-action {
            background: var(--card-bg-solid); color: var(--text-color); border: 1px solid var(--border-color);
            padding: 8px 16px; border-radius: 8px; cursor: pointer; font-family: inherit; font-size: 14px;
            transition: all 0.2s; display: flex; align-items: center; justify-content: center; gap: 8px;
        }
        select:hover, .theme-toggle:hover, .btn-action:hover { 
            border-color: var(--accent-color); 
            box-shadow: 0 0 10px var(--accent-glow);
            transform: translateY(-1px);
        }
        .theme-toggle { width: 38px; height: 38px; padding: 0; }
        .btn-action { height: 32px; padding: 0 12px; font-size: 12px; font-weight: 600; }

        .status-indicator { 
            display: flex; align-items: center; gap: 8px; font-size: 14px; color: var(--success-color); 
            font-weight: bold; background: rgba(16, 185, 129, 0.1); padding: 8px 16px; border-radius: 20px; 
            border: 1px solid rgba(16, 185, 129, 0.3);
        }
        .dot { width: 8px; height: 8px; background-color: var(--success-color); border-radius: 50%; animation: pulse 2s infinite; }
        @keyframes pulse { 0% { box-shadow: 0 0 0 0 rgba(16, 185, 129, 0.7); } 70% { box-shadow: 0 0 0 10px rgba(16, 185, 129, 0); } 100% { box-shadow: 0 0 0 0 rgba(16, 185, 129, 0); } }

        .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); gap: 20px; margin-bottom: 30px; }
        .card { 
            background: var(--card-bg); padding: 25px; border-radius: 16px; border: 1px solid var(--border-color); 
            transition: all 0.3s ease; backdrop-filter: blur(12px); box-shadow: var(--shadow);
            display: flex; flex-direction: column; justify-content: center; position: relative; overflow: hidden;
        }
        .card::before { content: ''; position: absolute; top: 0; left: 0; width: 4px; height: 100%; background: var(--accent-color); opacity: 0; transition: opacity 0.3s; }
        .card:hover { transform: translateY(-5px); border-color: var(--accent-color); }
        .card:hover::before { opacity: 1; }
        .card h3 { font-size: 13px; text-transform: uppercase; color: var(--text-muted); margin-bottom: 12px; letter-spacing: 1px; display: flex; align-items: center; gap: 10px; }
        .card h3 i { font-size: 16px; width: 30px; height: 30px; background: rgba(59, 130, 246, 0.1); border-radius: 8px; display: flex; align-items: center; justify-content: center; color: var(--accent-color); }
        .card .value { font-size: 32px; font-weight: 800; font-family: 'Inter', sans-serif; }

        .panel { 
            background: var(--card-bg); padding: 25px; border-radius: 16px; border: 1px solid var(--border-color); 
            margin-bottom: 30px; backdrop-filter: blur(12px); box-shadow: var(--shadow);
        }
        .panel h2 { font-size: 20px; margin-bottom: 20px; color: var(--text-color); border-bottom: 1px solid var(--border-color); padding-bottom: 15px; font-weight: 700; display: flex; align-items: center; gap: 10px; }
        .panel h2 i { color: var(--accent-color); }
        
        .status-banner { background: linear-gradient(135deg, rgba(16, 185, 129, 0.1), rgba(59, 130, 246, 0.1)); border: 1px solid rgba(16, 185, 129, 0.3); padding: 25px; border-radius: 12px; display: flex; justify-content: space-between; align-items: center; margin-bottom: 25px; }
        .status-banner .left { display: flex; align-items: center; gap: 15px; }
        .status-banner .left h3 { color: var(--success-color); font-size: 24px; font-weight: 700; }
        .status-banner .left p { color: var(--text-muted); font-size: 14px; }
        .status-banner .left i { font-size: 40px; color: var(--success-color); }
        .status-banner .right { font-size: 48px; font-weight: 900; color: var(--success-color); }
        
        .summary-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 15px; margin-bottom: 30px; }
        .summary-box { background: var(--card-bg-solid); padding: 20px; border-radius: 8px; border: 1px solid var(--border-color); text-align: center; }
        .summary-box span { display: block; color: var(--text-muted); font-size: 13px; margin-bottom: 8px; text-transform: uppercase; }
        .summary-box strong { font-size: 28px; color: var(--text-color); font-weight: 800; }

        .status-group { background: var(--card-bg-solid); border: 1px solid var(--border-color); border-radius: 12px; margin-bottom: 20px; overflow: hidden; }
        .status-group-header { padding: 15px 20px; border-bottom: 1px solid var(--border-color); display: flex; justify-content: space-between; align-items: center; }
        .status-group-header h4 { font-size: 16px; color: var(--text-color); font-weight: 700; display: flex; align-items: center; gap: 10px; }
        .status-group-header h4 i { color: var(--accent-color); }
        .status-tag { display: flex; align-items: center; gap: 8px; color: var(--success-color); font-size: 14px; font-weight: 600; }
        .status-item { padding: 15px 20px; display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid var(--border-color); }
        .status-item:last-child { border-bottom: none; }
        .status-item .name { color: var(--text-muted); font-size: 14px; display: flex; align-items: center; gap: 10px; }
        .status-item .name i { width: 20px; text-align: center; color: var(--text-muted); }
        .status-item .status-ok { color: var(--success-color); font-size: 13px; font-weight: bold; display: flex; align-items: center; gap: 6px; }
        .small-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--success-color); }

        .info-row { display: flex; justify-content: space-between; align-items: center; padding: 15px 0; border-bottom: 1px solid var(--border-color); }
        .info-row:last-child { border-bottom: none; }
        .info-label { color: var(--text-muted); font-weight: 600; font-size: 14px; display: flex; align-items: center; gap: 10px; }
        .info-label i { width: 20px; text-align: center; }
        .info-value { color: var(--text-color); font-weight: 700; text-align: right; font-size: 14px; }
        a { color: var(--accent-color); text-decoration: none; }
        a:hover { text-decoration: underline; }
        .tag { display: inline-block; background: rgba(59, 130, 246, 0.1); color: var(--accent-color); padding: 6px 14px; border-radius: 20px; font-size: 12px; margin: 3px; border: 1px solid rgba(59, 130, 246, 0.2); font-weight: 600; }

        .terminal { background: var(--log-bg); border-radius: 12px; border: 1px solid var(--border-color); overflow: hidden; direction: ltr; text-align: left; box-shadow: var(--shadow); }
        .terminal-header { background: var(--card-bg-solid); padding: 12px 15px; display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid var(--border-color); }
        .terminal-header-left { display: flex; gap: 8px; align-items: center; }
        .dot-red { width: 12px; height: 12px; border-radius: 50%; background: #ef4444; }
        .dot-yellow { width: 12px; height: 12px; border-radius: 50%; background: #facc15; }
        .dot-green { width: 12px; height: 12px; border-radius: 50%; background: #4ade80; }
        .terminal-title { margin-left: 10px; font-size: 12px; color: var(--text-muted); font-family: 'JetBrains Mono', monospace; }
        .terminal-header-right { display: flex; gap: 10px; align-items: center; }
        .search-wrapper { position: relative; display: flex; align-items: center; }
        .search-wrapper i { position: absolute; left: 10px; color: var(--text-muted); font-size: 12px; }
        .log-search { background: transparent; border: 1px solid var(--border-color); color: var(--text-color); padding: 5px 10px 5px 30px; border-radius: 4px; font-size: 12px; outline: none; width: 200px; transition: all 0.2s; }
        .log-search:focus { border-color: var(--accent-color); box-shadow: 0 0 10px var(--accent-glow); }

        .terminal-body { padding: 15px; height: 500px; overflow-y: auto; font-family: 'JetBrains Mono', 'Courier New', monospace; font-size: 13px; }
        .terminal-body::-webkit-scrollbar { width: 8px; }
        .terminal-body::-webkit-scrollbar-track { background: transparent; }
        .terminal-body::-webkit-scrollbar-thumb { background: #475569; border-radius: 4px; }
        .terminal-body::-webkit-scrollbar-thumb:hover { background: #64748b; }
        
        .log-line { margin-bottom: 8px; padding: 8px 12px; border-radius: 6px; white-space: pre-wrap; word-wrap: break-word; border-left: 4px solid transparent; background: rgba(255,255,255,0.02); transition: all 0.2s; }
        .log-line:hover { background: rgba(255,255,255,0.05); transform: translateX(2px); }
        .log-ts { color: #64748b; margin-right: 10px; }
        .log-level { font-weight: bold; margin-right: 12px; padding: 2px 8px; border-radius: 4px; text-transform: uppercase; font-size: 11px; }
        .level-info { color: #4ade80; border-left-color: #4ade80; background: rgba(74, 222, 128, 0.1); }
        .level-warn { color: #facc15; border-left-color: #facc15; background: rgba(250, 204, 21, 0.1); }
        .level-error, .level-fatal { color: #ef4444; border-left-color: #ef4444; background: rgba(239, 68, 68, 0.1); }
        .level-debug { color: #38bdf8; border-left-color: #38bdf8; background: rgba(56, 189, 248, 0.1); }
        .level-unknown { color: #cbd5e1; border-left-color: #334155; }
        
        .highlight { color: #facc15; font-weight: bold; }
        .file-highlight { color: #94a3b8; font-style: italic; }

        /* Toast Notification */
        .toast {
            visibility: hidden;
            min-width: 250px;
            margin-left: -125px;
            background-color: var(--success-color);
            color: #fff;
            text-align: center;
            border-radius: 8px;
            padding: 12px;
            position: fixed;
            z-index: 1;
            left: 50%;
            bottom: 30px;
            font-weight: 600;
            box-shadow: 0 4px 15px rgba(0,0,0,0.3);
            opacity: 0;
            transition: opacity 0.3s, bottom 0.3s;
        }
        .toast.show {
            visibility: visible;
            opacity: 1;
            bottom: 50px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>
                <div class="logo">CM</div> 
                <i class="fas fa-chart-line"></i>
                <span id="t-dashboard">CodeMeet Dashboard</span>
            </h1>
            <div class="header-controls">
                <select id="lang-select" onchange="changeLanguage(this.value)">
                    <option value="en">English</option>
                    <option value="fa">فارسی</option>
                    <option value="ar">العربية</option>
                </select>
                <button class="theme-toggle" onclick="toggleTheme()" title="Toggle Theme">
                    <i class="fas fa-moon" id="theme-icon"></i>
                </button>
                <button class="theme-toggle" onclick="window.location.href='/logout'" title="Logout" style="color: var(--error-color);">
                    <i class="fas fa-right-from-bracket"></i>
                </button>
                <div class="status-indicator">
                    <div class="dot"></div>
                    <i class="fas fa-bolt"></i>
                    <span id="t-status">RUNNING</span>
                </div>
            </div>
        </div>
        
        <div class="grid">
            <div class="card">
                <h3><i class="fas fa-arrow-trend-up"></i> <span id="t-total-req">Total Requests</span></h3>
                <div class="value" id="stat-requests">0</div>
            </div>
            <div class="card">
                <h3><i class="fas fa-circle-check"></i> <span id="t-success">Success</span></h3>
                <div class="value" id="stat-success" style="color: var(--success-color);">0</div>
            </div>
            <div class="card">
                <h3><i class="fas fa-circle-exclamation"></i> <span id="t-errors">Errors</span></h3>
                <div class="value" id="stat-errors" style="color: var(--error-color);">0</div>
            </div>
            <div class="card">
                <h3><i class="fas fa-stopwatch"></i> <span id="t-latency">Avg Latency</span></h3>
                <div class="value" id="stat-latency">0.00ms</div>
            </div>
        </div>

        <div class="panel" style="margin-top: 0;">
            <h2><i class="fas fa-heart-pulse"></i> <span id="t-sys-status">System Status</span></h2>
            <div class="status-banner">
                <div class="left">
                    <i class="fas fa-shield-halved"></i>
                    <div>
                        <h3 id="t-all-op">All Systems Operational</h3>
                        <p id="t-sys-smooth">Bot is running smoothly. Updated just now.</p>
                    </div>
                </div>
                <div class="right">100%</div>
            </div>
            <div class="summary-grid">
                <div class="summary-box"><span id="t-comp">Components</span><strong>9</strong></div>
                <div class="summary-box"><span id="t-op">Operational</span><strong>9</strong></div>
                <div class="summary-box"><span id="t-issues">Issues</span><strong>0</strong></div>
            </div>

            <div class="status-group">
                <div class="status-group-header">
                    <h4><i class="fas fa-microchip"></i> <span id="t-core">Core Services</span></h4>
                    <div class="status-tag"><div class="small-dot"></div> <span id="t-op2">Operational</span></div>
                </div>
                <div class="status-item">
                    <span class="name"><i class="fas fa-satellite-dish"></i> <span id="t-update-rec">Update Receiver (Polling/Webhook)</span></span>
                    <span class="status-ok"><div class="small-dot"></div> <span id="t-op3">Operational</span></span>
                </div>
                <div class="status-item">
                    <span class="name"><i class="fas fa-shuffle"></i> <span id="t-disp">Update Dispatcher</span></span>
                    <span class="status-ok"><div class="small-dot"></div> <span id="t-op4">Operational</span></span>
                </div>
                <div class="status-item">
                    <span class="name"><i class="fas fa-gauge-high"></i> <span id="t-rl">Rate Limiter</span></span>
                    <span class="status-ok"><div class="small-dot"></div> <span id="t-op5">Operational</span></span>
                </div>
            </div>

            <div class="status-group">
                <div class="status-group-header">
                    <h4><i class="fas fa-robot"></i> <span id="t-bot-api">Bot API Methods</span></h4>
                    <div class="status-tag"><div class="small-dot"></div> <span id="t-op6">Operational</span></div>
                </div>
                <div class="status-item">
                    <span class="name"><i class="fas fa-paper-plane"></i> <span id="t-send-msg">Send Message</span></span>
                    <span class="status-ok"><div class="small-dot"></div> <span id="t-op7">Operational</span></span>
                </div>
                <div class="status-item">
                    <span class="name"><i class="fas fa-cloud-arrow-up"></i> <span id="t-media-up">Media Upload</span></span>
                    <span class="status-ok"><div class="small-dot"></div> <span id="t-op8">Operational</span></span>
                </div>
                <div class="status-item">
                    <span class="name"><i class="fas fa-comments"></i> <span id="t-chat-mgmt">Chat Management</span></span>
                    <span class="status-ok"><div class="small-dot"></div> <span id="t-op9">Operational</span></span>
                </div>
            </div>

            <div class="status-group">
                <div class="status-group-header">
                    <h4><i class="fas fa-server"></i> <span id="t-infra">Infrastructure</span></h4>
                    <div class="status-tag"><div class="small-dot"></div> <span id="t-op10">Operational</span></div>
                </div>
                <div class="status-item">
                    <span class="name"><i class="fas fa-database"></i> <span id="t-cache">Memory Cache</span></span>
                    <span class="status-ok"><div class="small-dot"></div> <span id="t-op11">Operational</span></span>
                </div>
                <div class="status-item">
                    <span class="name"><i class="fas fa-rotate-right"></i> <span id="t-retry">Retry Policy</span></span>
                    <span class="status-ok"><div class="small-dot"></div> <span id="t-op12">Operational</span></span>
                </div>
                <div class="status-item">
                    <span class="name"><i class="fas fa-desktop"></i> <span id="t-dash-srv">Dashboard Server</span></span>
                    <span class="status-ok"><div class="small-dot"></div> <span id="t-op13">Operational</span></span>
                </div>
            </div>
        </div>

        <div class="panel">
            <h2><i class="fas fa-circle-info"></i> <span id="t-sys-info">System Information</span></h2>
            <div class="info-row">
                <span class="info-label"><i class="fas fa-user-pen"></i> <span id="t-author">Author</span></span>
                <span class="info-value" id="info-author">Loading...</span>
            </div>
            <div class="info-row">
                <span class="info-label"><i class="fab fa-github"></i> <span id="t-github">GitHub Profile</span></span>
                <span class="info-value"><a id="info-github" href="#" target="_blank">Loading...</a></span>
            </div>
            <div class="info-row">
                <span class="info-label"><i class="fas fa-code-branch"></i> <span id="t-repo">Repository</span></span>
                <span class="info-value"><a id="info-repo" href="#" target="_blank">Loading...</a></span>
            </div>
            <div class="info-row">
                <span class="info-label"><i class="fas fa-tag"></i> <span id="t-version">Library Version</span></span>
                <span class="info-value" id="info-version">Loading...</span>
            </div>
            <div class="info-row">
                <span class="info-label"><i class="fas fa-play"></i> <span id="t-runmode">Run Mode</span></span>
                <span class="info-value" id="info-runmode">Loading...</span>
            </div>
            <div class="info-row">
                <span class="info-label"><i class="fas fa-puzzle-piece"></i> <span id="t-features">Active Features</span></span>
                <span class="info-value" id="info-features">Loading...</span>
            </div>
        </div>

        <div class="panel" style="padding: 0; background: transparent; border: none; box-shadow: none;">
            <div class="terminal">
                <div class="terminal-header">
                    <div class="terminal-header-left">
                        <div class="dot-red"></div><div class="dot-yellow"></div><div class="dot-green"></div>
                        <div class="terminal-title" id="t-logs-title">bot@codemeet: ~ (live logs)</div>
                    </div>
                    <div class="terminal-header-right">
                        <div class="search-wrapper">
                            <i class="fas fa-magnifying-glass"></i>
                            <input type="text" class="log-search" id="log-search" placeholder="Search logs..." oninput="fetchLogs()">
                        </div>
                        <button class="btn-action" id="pause-btn" onclick="togglePause()" title="Pause Scroll">
                            <i class="fas fa-pause"></i>
                        </button>
                        <button class="btn-action" id="toggle-log-btn" onclick="toggleLogs()" title="Start/Stop Logging" style="color: var(--error-color);">
                            <i class="fas fa-stop"></i> Stop Logs
                        </button>
                        <button class="btn-action" onclick="copyLogs()" title="Copy All Logs">
                            <i class="fas fa-copy"></i> <span id="t-copy">Copy</span>
                        </button>
                    </div>
                </div>
                <div class="terminal-body" id="log-container"></div>
            </div>
        </div>
    </div>

    <div id="toast" class="toast">Logs copied to clipboard!</div>

    <script>
        const translations = {
            en: { dashboard: "CodeMeet Dashboard", status: "RUNNING", total_req: "Total Requests", success: "Success", errors: "Errors", latency: "Avg Latency", sys_status: "System Status", all_op: "All Systems Operational", sys_smooth: "Bot is running smoothly. Updated just now.", comp: "Components", op: "Operational", issues: "Issues", core: "Core Services", bot_api: "Bot API Methods", infra: "Infrastructure", update_rec: "Update Receiver (Polling/Webhook)", disp: "Update Dispatcher", rl: "Rate Limiter", send_msg: "Send Message", media_up: "Media Upload", chat_mgmt: "Chat Management", cache: "Memory Cache", retry: "Retry Policy", dash_srv: "Dashboard Server", sys_info: "System Information", author: "Author", github: "GitHub Profile", repo: "Repository", version: "Library Version", runmode: "Run Mode", features: "Active Features", logs_title: "bot@codemeet: ~ (live logs)", search_log: "Search logs...", copy: "Copy", toast_copied: "Logs copied to clipboard!" },
            fa: { dashboard: "داشبورد کدمیت", status: "در حال اجرا", total_req: "کل درخواست‌ها", success: "موفقیت", errors: "خطاها", latency: "میانگین تأخیر", sys_status: "وضعیت سیستم", all_op: "تمام سیستم‌ها عملیاتی هستند", sys_smooth: "ربات بدون مشکل در حال اجراست.", comp: "اجزا", op: "عملیاتی", issues: "مشکلات", core: "سرویس‌های اصلی", bot_api: "متدهای بات API", infra: "زیرساخت", update_rec: "دریافت‌کننده آپدیت", disp: "مدیریت آپدیت‌ها", rl: "محدودکننده نرخ", send_msg: "ارسال پیام", media_up: "آپلود رسانه", chat_mgmt: "مدیریت چت", cache: "کش حافظه", retry: "سیاست تلاش مجدد", dash_srv: "سرور داشبورد", sys_info: "اطلاعات سیستم", author: "سازنده", github: "گیت‌هاب", repo: "مخزن", version: "نسخه", runmode: "حالت اجرا", features: "امکانات فعال", logs_title: "ربات@کدمیت: ~ (لاگ‌های زنده)", search_log: "جستجو...", copy: "کپی", toast_copied: "لاگ‌ها کپی شدند!" },
            ar: { dashboard: "لوحة تحكم كودميت", status: "يعمل", total_req: "إجمالي الطلبات", success: "نجاح", errors: "أخطاء", latency: "متوسط التأخير", sys_status: "حالة النظام", all_op: "جميع الأنظمة تعمل", sys_smooth: "يعمل البوت بسلاسة.", comp: "المكونات", op: "تعمل", issues: "مشاكل", core: "الخدمات الأساسية", bot_api: "طرق بوت API", infra: "البنية التحتية", update_rec: "مستقبل التحديثات", disp: "موزع التحديثات", rl: "محدد المعدل", send_msg: "إرسال رسالة", media_up: "رفع الوسائط", chat_mgmt: "إدارة الدردشة", cache: "ذاكرة التخزين المؤقت", retry: "إعادة المحاولة", dash_srv: "خادم لوحة التحكم", sys_info: "معلومات النظام", author: "المؤلف", github: "GitHub", repo: "المستودع", version: "الإصدار", runmode: "وضع التشغيل", features: "الميزات النشطة", logs_title: "بوت@كودميت: ~ (سجلات مباشرة)", search_log: "بحث...", copy: "نسخ", toast_copied: "تم نسخ السجلات!" }
        };

        let isPaused = false;
        let logsEnabled = true;

        function applyLanguage(lang) {
            const t = translations[lang] || translations.en;
            document.documentElement.lang = lang;
            document.documentElement.dir = (lang === 'fa' || lang === 'ar') ? 'rtl' : 'ltr';
            
            const ids = ['dashboard', 'status', 'total_req', 'success', 'errors', 'latency', 'sys_status', 'all_op', 'sys_smooth', 'comp', 'op', 'issues', 'core', 'bot_api', 'infra', 'update_rec', 'disp', 'rl', 'send_msg', 'media_up', 'chat_mgmt', 'cache', 'retry', 'dash_srv', 'sys_info', 'author', 'github', 'repo', 'version', 'runmode', 'features', 'logs_title', 'search_log', 'copy', 'toast_copied'];
            
            ids.forEach(id => {
                const el = document.getElementById('t-' + id);
                if(el) el.innerText = t[id];
            });

            document.getElementById('log-search').placeholder = t.search_log;
            document.getElementById('lang-select').value = lang;
            localStorage.setItem('cm-lang', lang);
        }

        function changeLanguage(lang) { applyLanguage(lang); }

        function toggleTheme() {
            document.body.classList.toggle('light-mode');
            const isLight = document.body.classList.contains('light-mode');
            document.getElementById('theme-icon').className = isLight ? 'fas fa-sun' : 'fas fa-moon';
            localStorage.setItem('cm-theme', isLight ? 'light' : 'dark');
        }

        function togglePause() {
            isPaused = !isPaused;
            const btn = document.getElementById('pause-btn');
            btn.innerHTML = isPaused ? '<i class="fas fa-play"></i>' : '<i class="fas fa-pause"></i>';
            btn.style.color = isPaused ? 'var(--warn-color)' : 'var(--text-color)';
        }

        async function toggleLogs() {
            try {
                const res = await fetch('/api/logs/toggle');
                const data = await res.json();
                logsEnabled = data.enabled;
                const btn = document.getElementById('toggle-log-btn');
                if (logsEnabled) {
                    btn.innerHTML = '<i class="fas fa-stop"></i> Stop Logs';
                    btn.style.color = 'var(--error-color)';
                } else {
                    btn.innerHTML = '<i class="fas fa-play"></i> Resume Logs';
                    btn.style.color = 'var(--success-color)';
                }
                fetchLogs(); // Fetch immediately to see the toggle message
            } catch (e) {
                console.error('Error toggling logs:', e);
            }
        }

        function showToast(message) {
            const toast = document.getElementById("toast");
            toast.innerText = message;
            toast.className = "toast show";
            setTimeout(() => { toast.className = toast.className.replace("show", ""); }, 3000);
        }

        async function copyLogs() {
            try {
                const res = await fetch('/api/logs');
                const logs = await res.json();
                const text = logs.join('\n');
                
                if (navigator.clipboard && window.isSecureContext) {
                    await navigator.clipboard.writeText(text);
                } else {
                    // Fallback
                    const textArea = document.createElement('textarea');
                    textArea.value = text;
                    textArea.style.position = 'fixed';
                    textArea.style.opacity = '0';
                    document.body.appendChild(textArea);
                    textArea.focus();
                    textArea.select();
                    document.execCommand('copy');
                    document.body.removeChild(textArea);
                }
                
                const lang = localStorage.getItem('cm-lang') || 'en';
                showToast(translations[lang].toast_copied);
            } catch (err) {
                console.error('Copy failed:', err);
            }
        }

        function escapeHtml(text) {
            var div = document.createElement('div');
            div.innerText = text;
            return div.innerHTML;
        }

        function highlightLog(text) {
            let html = escapeHtml(text);
            // Highlight key=value pairs
            html = html.replace(/(\b\w+=\S+)/g, '<span class="highlight">$1</span>');
            // Highlight file paths like (file.go:123)
            html = html.replace(/(\([\w\d\_\-\/\.]+:\d+\))/g, '<span class="file-highlight">$1</span>');
            return html;
        }

        async function fetchData() {
            try {
                const infoRes = await fetch('/api/info');
                const info = await infoRes.json();
                
                document.getElementById('info-author').innerText = info.author;
                var ghLink = document.getElementById('info-github');
                ghLink.href = 'https://' + info.github;
                ghLink.innerText = info.github;
                var repoLink = document.getElementById('info-repo');
                repoLink.href = 'https://' + info.repo;
                repoLink.innerText = info.repo;
                document.getElementById('info-version').innerText = info.version;
                document.getElementById('info-runmode').innerText = info.runMode;
                
                var featuresDiv = document.getElementById('info-features');
                featuresDiv.innerHTML = '';
                info.features.forEach(function(f) {
                    var tag = document.createElement('span');
                    tag.className = 'tag';
                    tag.innerText = f;
                    featuresDiv.appendChild(tag);
                });

                // Sync log button state
                logsEnabled = info.logs_enabled;
                const btn = document.getElementById('toggle-log-btn');
                if (logsEnabled) {
                    btn.innerHTML = '<i class="fas fa-stop"></i> Stop Logs';
                    btn.style.color = 'var(--error-color)';
                } else {
                    btn.innerHTML = '<i class="fas fa-play"></i> Resume Logs';
                    btn.style.color = 'var(--success-color)';
                }

                const statsRes = await fetch('/api/stats');
                const stats = await statsRes.json();
                document.getElementById('stat-requests').innerText = stats.api.Requests || 0;
                document.getElementById('stat-success').innerText = stats.api.SuccessCount || 0;
                document.getElementById('stat-errors').innerText = stats.api.ErrorCount || 0;
                document.getElementById('stat-latency').innerText = ((stats.api.AvgLatency / 1000000) || 0).toFixed(2) + 'ms';

            } catch (e) {
                console.error('Error fetching data:', e);
                document.getElementById('t-status').innerText = 'ERROR';
                document.getElementById('t-status').style.color = 'var(--error-color)';
            }
        }

        async function fetchLogs() {
            try {
                const res = await fetch('/api/logs');
                const logs = await res.json();
                var container = document.getElementById('log-container');
                var searchTerm = document.getElementById('log-search').value.toLowerCase();
                
                var isScrolledToBottom = container.scrollHeight - container.clientHeight <= container.scrollTop + 1;
                
                container.innerHTML = '';
                
                var filteredLogs = logs.filter(log => log.toLowerCase().includes(searchTerm));
                
                filteredLogs.forEach(function(log) {
                    var div = document.createElement('div');
                    var match = log.match(/^\[(.*?)\]\s+(DEBUG|INFO |WARN |ERROR|FATAL)\s+([\s\S]*)$/);
                    if (match) {
                        var ts = match[1];
                        var level = match[2].trim();
                        var rest = match[3];
                        div.className = 'log-line level-' + level.toLowerCase();
                        var tsHtml = '<span class="log-ts">[' + ts + ']</span>';
                        var levelHtml = '<span class="log-level">' + level + '</span>';
                        div.innerHTML = tsHtml + levelHtml + highlightLog(rest);
                    } else {
                        div.className = 'log-line level-unknown';
                        div.innerHTML = highlightLog(log);
                    }
                    container.appendChild(div);
                });

                if (isScrolledToBottom && !isPaused) {
                    container.scrollTop = container.scrollHeight;
                }
            } catch (e) {
                console.error('Error fetching logs:', e);
            }
        }

        document.addEventListener('DOMContentLoaded', function() {
            const savedLang = localStorage.getItem('cm-lang') || 'en';
            applyLanguage(savedLang);
            
            const savedTheme = localStorage.getItem('cm-theme');
            if (savedTheme === 'light') {
                toggleTheme();
            }

            setInterval(fetchData, 2000);
            setInterval(fetchLogs, 1000);
            fetchData();
            fetchLogs();
        });
    </script>
</body>
</html>`
