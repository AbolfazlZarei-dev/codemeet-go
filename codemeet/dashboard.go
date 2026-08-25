package codemeet

const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>CodeMeet Bot Dashboard</title>
    <style>
        :root {
            --bg-color: #0f172a;
            --card-bg: #1e293b;
            --text-color: #f1f5f9;
            --accent-color: #3b82f6;
            --success-color: #10b981;
            --warn-color: #facc15;
            --error-color: #ef4444;
            --border-color: #334155;
            --log-bg: #020617;
        }
        * { margin: 0; padding: 0; box-sizing: border-box; font-family: 'Inter', 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; }
        body { background-color: var(--bg-color); color: var(--text-color); padding: 20px; }
        .container { max-width: 1200px; margin: 0 auto; }
        .header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 30px; background: var(--card-bg); padding: 20px; border-radius: 12px; border: 1px solid var(--border-color); }
        .header h1 { font-size: 24px; font-weight: 600; display: flex; align-items: center; gap: 10px; }
        .logo { width: 32px; height: 32px; background: var(--accent-color); border-radius: 8px; display: flex; align-items: center; justify-content: center; font-weight: bold; color: white; }
        .status-indicator { display: flex; align-items: center; gap: 8px; font-size: 14px; color: var(--success-color); font-weight: bold; background: rgba(16, 185, 129, 0.1); padding: 8px 16px; border-radius: 20px; }
        .dot { width: 8px; height: 8px; background-color: var(--success-color); border-radius: 50%; animation: pulse 2s infinite; }
        @keyframes pulse { 0% { box-shadow: 0 0 0 0 rgba(16, 185, 129, 0.7); } 70% { box-shadow: 0 0 0 10px rgba(16, 185, 129, 0); } 100% { box-shadow: 0 0 0 0 rgba(16, 185, 129, 0); } }
        .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 20px; margin-bottom: 30px; }
        .card { background: var(--card-bg); padding: 20px; border-radius: 12px; border: 1px solid var(--border-color); transition: all 0.3s ease; }
        .card:hover { transform: translateY(-2px); box-shadow: 0 10px 20px rgba(0,0,0,0.2); border-color: var(--accent-color); }
        .card h3 { font-size: 12px; text-transform: uppercase; color: #94a3b8; margin-bottom: 10px; letter-spacing: 1px; }
        .card .value { font-size: 28px; font-weight: 700; }
        
        /* استایل‌های بخش وضعیت سیستم (System Status) */
        .status-banner { background: rgba(16, 185, 129, 0.1); border: 1px solid rgba(16, 185, 129, 0.3); padding: 25px; border-radius: 12px; display: flex; justify-content: space-between; align-items: center; margin-bottom: 25px; }
        .status-banner .left h3 { color: var(--success-color); font-size: 22px; margin-bottom: 5px; }
        .status-banner .left p { color: #94a3b8; font-size: 14px; }
        .status-banner .right { font-size: 48px; font-weight: 800; color: var(--success-color); }
        .summary-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 15px; margin-bottom: 30px; }
        .summary-box { background: var(--card-bg); padding: 15px; border-radius: 8px; border: 1px solid var(--border-color); text-align: center; }
        .summary-box span { display: block; color: #94a3b8; font-size: 12px; margin-bottom: 5px; text-transform: uppercase; }
        .summary-box strong { font-size: 24px; color: var(--text-color); }
        .status-group { background: var(--card-bg); border: 1px solid var(--border-color); border-radius: 12px; margin-bottom: 20px; overflow: hidden; }
        .status-group-header { padding: 15px 20px; border-bottom: 1px solid var(--border-color); display: flex; justify-content: space-between; align-items: center; background: #1a2535; }
        .status-group-header h4 { font-size: 16px; color: var(--text-color); }
        .status-tag { display: flex; align-items: center; gap: 8px; color: var(--success-color); font-size: 14px; font-weight: 500; }
        .status-item { padding: 12px 20px; display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid var(--border-color); }
        .status-item:last-child { border-bottom: none; }
        .status-item .name { color: #cbd5e1; font-size: 14px; }
        .status-item .status-ok { color: var(--success-color); font-size: 12px; font-weight: bold; display: flex; align-items: center; gap: 6px; }
        .small-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--success-color); }

        .panel { background: var(--card-bg); padding: 25px; border-radius: 12px; border: 1px solid var(--border-color); margin-bottom: 30px; }
        .panel h2 { font-size: 18px; margin-bottom: 20px; color: var(--text-color); border-bottom: 1px solid var(--border-color); padding-bottom: 10px; }
        .info-row { display: flex; justify-content: space-between; align-items: center; padding: 12px 0; border-bottom: 1px solid var(--border-color); }
        .info-row:last-child { border-bottom: none; }
        .info-label { color: #94a3b8; font-weight: 500; font-size: 14px; }
        .info-value { color: var(--text-color); font-weight: 500; text-align: right; font-size: 14px; }
        a { color: var(--accent-color); text-decoration: none; }
        a:hover { text-decoration: underline; }
        .tag { display: inline-block; background: rgba(59, 130, 246, 0.1); color: var(--accent-color); padding: 5px 12px; border-radius: 6px; font-size: 12px; margin: 3px; border: 1px solid rgba(59, 130, 246, 0.2); }
        .terminal { background: var(--log-bg); border-radius: 12px; border: 1px solid var(--border-color); overflow: hidden; }
        .terminal-header { background: #1e293b; padding: 10px 15px; display: flex; gap: 8px; align-items: center; border-bottom: 1px solid var(--border-color); }
        .dot-red { width: 12px; height: 12px; border-radius: 50%; background: #ef4444; }
        .dot-yellow { width: 12px; height: 12px; border-radius: 50%; background: #facc15; }
        .dot-green { width: 12px; height: 12px; border-radius: 50%; background: #4ade80; }
        .terminal-title { margin-left: 10px; font-size: 12px; color: #94a3b8; font-family: monospace; }
        .terminal-body { padding: 15px; height: 400px; overflow-y: auto; font-family: 'Courier New', monospace; font-size: 13px; }
        .terminal-body::-webkit-scrollbar { width: 8px; }
        .terminal-body::-webkit-scrollbar-track { background: #1e293b; }
        .terminal-body::-webkit-scrollbar-thumb { background: #475569; border-radius: 4px; }
        .log-line { margin-bottom: 8px; white-space: pre-wrap; word-wrap: break-word; border-left: 3px solid transparent; padding-left: 8px; }
        .log-ts { color: #64748b; margin-right: 5px; }
        .log-level { font-weight: bold; margin-right: 5px; }
        .level-info { color: #4ade80; border-left-color: #4ade80; }
        .level-warn { color: #facc15; border-left-color: #facc15; }
        .level-error, .level-fatal { color: #ef4444; border-left-color: #ef4444; }
        .level-debug { color: #38bdf8; border-left-color: #38bdf8; }
        .level-unknown { color: #cbd5e1; border-left-color: #334155; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1><div class="logo">CM</div> CodeMeet Dashboard</h1>
            <div class="status-indicator">
                <div class="dot"></div>
                <span id="status-badge">RUNNING</span>
            </div>
        </div>
        
        <div class="grid">
            <div class="card"><h3>Total Requests</h3><div class="value" id="stat-requests">0</div></div>
            <div class="card"><h3>Success</h3><div class="value" id="stat-success" style="color: var(--success-color);">0</div></div>
            <div class="card"><h3>Errors</h3><div class="value" id="stat-errors" style="color: var(--error-color);">0</div></div>
            <div class="card"><h3>Avg Latency</h3><div class="value" id="stat-latency">0.00ms</div></div>
        </div>

        <!-- بخش وضعیت سیستم -->
        <div class="panel" style="margin-top: 0;">
            <h2>System Status</h2>
            <div class="status-banner">
                <div class="left">
                    <h3>All Systems Operational</h3>
                    <p>Bot is running smoothly. Updated just now.</p>
                </div>
                <div class="right">100%</div>
            </div>
            <div class="summary-grid">
                <div class="summary-box"><span>Total Components</span><strong>9</strong></div>
                <div class="summary-box"><span>Operational</span><strong>9</strong></div>
                <div class="summary-box"><span>Issues Detected</span><strong>0</strong></div>
            </div>

            <div class="status-group">
                <div class="status-group-header">
                    <h4>Core Services</h4>
                    <div class="status-tag"><div class="small-dot"></div> Operational 100%</div>
                </div>
                <div class="status-item"><span class="name">Update Receiver (Polling/Webhook)</span><span class="status-ok"><div class="small-dot"></div> Operational</span></div>
                <div class="status-item"><span class="name">Update Dispatcher</span><span class="status-ok"><div class="small-dot"></div> Operational</span></div>
                <div class="status-item"><span class="name">Rate Limiter</span><span class="status-ok"><div class="small-dot"></div> Operational</span></div>
            </div>

            <div class="status-group">
                <div class="status-group-header">
                    <h4>Bot API Methods</h4>
                    <div class="status-tag"><div class="small-dot"></div> Operational 100%</div>
                </div>
                <div class="status-item"><span class="name">Send Message</span><span class="status-ok"><div class="small-dot"></div> Operational</span></div>
                <div class="status-item"><span class="name">Media Upload</span><span class="status-ok"><div class="small-dot"></div> Operational</span></div>
                <div class="status-item"><span class="name">Chat Management</span><span class="status-ok"><div class="small-dot"></div> Operational</span></div>
            </div>

            <div class="status-group">
                <div class="status-group-header">
                    <h4>Infrastructure</h4>
                    <div class="status-tag"><div class="small-dot"></div> Operational 100%</div>
                </div>
                <div class="status-item"><span class="name">Memory Cache</span><span class="status-ok"><div class="small-dot"></div> Operational</span></div>
                <div class="status-item"><span class="name">Retry Policy</span><span class="status-ok"><div class="small-dot"></div> Operational</span></div>
                <div class="status-item"><span class="name">Dashboard Server</span><span class="status-ok"><div class="small-dot"></div> Operational</span></div>
            </div>
        </div>

        <div class="panel">
            <h2>System Information</h2>
            <div class="info-row"><span class="info-label">Author</span><span class="info-value" id="info-author">Loading...</span></div>
            <div class="info-row"><span class="info-label">GitHub Profile</span><span class="info-value"><a id="info-github" href="#" target="_blank">Loading...</a></span></div>
            <div class="info-row"><span class="info-label">Repository</span><span class="info-value"><a id="info-repo" href="#" target="_blank">Loading...</a></span></div>
            <div class="info-row"><span class="info-label">Library Version</span><span class="info-value" id="info-version">Loading...</span></div>
            <div class="info-row"><span class="info-label">Run Mode</span><span class="info-value" id="info-runmode">Loading...</span></div>
            <div class="info-row"><span class="info-label">Active Features</span><span class="info-value" id="info-features">Loading...</span></div>
        </div>

        <div class="panel" style="padding: 0; background: transparent; border: none;">
            <div class="terminal">
                <div class="terminal-header">
                    <div class="dot-red"></div><div class="dot-yellow"></div><div class="dot-green"></div>
                    <div class="terminal-title">bot@codemeet: ~ (live logs)</div>
                </div>
                <div class="terminal-body" id="log-container"></div>
            </div>
        </div>
    </div>

    <script>
        function escapeHtml(text) {
            var div = document.createElement('div');
            div.innerText = text;
            return div.innerHTML;
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

                const statsRes = await fetch('/api/stats');
                const stats = await statsRes.json();
                document.getElementById('stat-requests').innerText = stats.Requests || 0;
                document.getElementById('stat-success').innerText = stats.SuccessCount || 0;
                document.getElementById('stat-errors').innerText = stats.ErrorCount || 0;
                document.getElementById('stat-latency').innerText = ((stats.AvgLatency / 1000000) || 0).toFixed(2) + 'ms';

            } catch (e) {
                console.error('Error fetching data:', e);
                document.getElementById('status-badge').innerText = 'ERROR';
                document.getElementById('status-badge').style.color = 'var(--error-color)';
            }
        }

        async function fetchLogs() {
            try {
                const res = await fetch('/api/logs');
                const logs = await res.json();
                var container = document.getElementById('log-container');
                
                if(container.childElementCount !== logs.length) {
                    container.innerHTML = '';
                    logs.forEach(function(log) {
                        var div = document.createElement('div');
                        var match = log.match(/^\[(.*?)\]\s+(DEBUG|INFO|WARN|ERROR|FATAL)\s+([\s\S]*)$/);
                        if (match) {
                            var ts = match[1];
                            var level = match[2];
                            var rest = match[3];
                            div.className = 'log-line level-' + level.toLowerCase();
                            var tsHtml = '<span class="log-ts">[' + ts + ']</span>';
                            var levelHtml = '<span class="log-level">' + level + '</span>';
                            div.innerHTML = tsHtml + levelHtml + escapeHtml(rest);
                        } else {
                            div.className = 'log-line level-unknown';
                            div.innerText = log;
                        }
                        container.appendChild(div);
                    });
                    container.scrollTop = container.scrollHeight;
                }
            } catch (e) {
                console.error('Error fetching logs:', e);
            }
        }

        setInterval(fetchData, 2000);
        setInterval(fetchLogs, 1000);
        fetchData();
        fetchLogs();
    </script>
</body>
</html>`
