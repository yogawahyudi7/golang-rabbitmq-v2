# 🎨 RabbitMQ Monitoring Dashboard - Quick Guide

## 🚀 Quick Start

### Method 1: Using Makefile (Recommended)

```bash
# Start backend
make docker-up

# Open dashboard (in new terminal)
make dashboard
```

### Method 2: Manual

```bash
# Start backend
docker-compose up -d

# Open dashboard
cd web
./open-dashboard.bat    # Windows
./open-dashboard.sh     # Linux/Mac
```

### Method 3: Direct Open (If CORS OK)

```bash
# Just double-click
web/dashboard.html
```

---

## 📊 Dashboard Overview

### 🎯 Main Components

#### 1. **Header Section**
```
┌─────────────────────────────────────┐
│ 🐰 RabbitMQ Monitoring Dashboard   │
│ Status: ● Connected | Uptime: 2h   │
└─────────────────────────────────────┘
```

#### 2. **Metrics Cards (Top Row)**
```
┌──────────┬──────────┬──────────┬──────────────┐
│ Published│ Consumed │  Failed  │ Success Rate │
│   1,523  │  1,521   │    2     │   99.87%     │
└──────────┴──────────┴──────────┴──────────────┘
```

#### 3. **Performance Cards (Second Row)**
```
┌──────────────┬──────────────┬──────────┬──────────────┐
│ Publish Rate │ Consume Rate │In Progress│ Avg Process  │
│   5.2 msg/s  │   5.0 msg/s  │     3    │    42ms      │
└──────────────┴──────────────┴──────────┴──────────────┘
```

#### 4. **Real-time Chart**
```
📊 Real-time Rates
   ┌────────────────────────────────────┐
   │                  ╱╲                │
   │              ╱╲  ╱  ╲   ╱╲         │
   │         ╱╲  ╱  ╲╱    ╲ ╱  ╲        │
   │    ╱╲  ╱  ╲╱          ╲╱    ╲      │
   │   ╱  ╲╱                      ╲╱╲   │
   └────────────────────────────────────┘
     — Publish Rate    — Consume Rate
```

#### 5. **Events Table**
```
📜 Recent Events    [All][Publish][Consume][Failed]
┌──────────┬────────┬─────────────┬─────────┬────────┐
│   Time   │  Type  │ Routing Key │Duration │ Status │
├──────────┼────────┼─────────────┼─────────┼────────┤
│ 10:15:29 │Consume │app.message  │  42ms   │Success │
│ 10:15:28 │Publish │app.message  │   5ms   │Success │
│ 10:15:27 │Failed  │app.message  │  30ms   │ Error  │
└──────────┴────────┴─────────────┴─────────┴────────┘
```

---

## 🎮 Features in Action

### ✅ Real-time Updates
- **Auto-refresh**: Every 5 seconds
- **Live charts**: Smooth animations
- **Instant updates**: No page reload needed

### 🔍 Event Filtering
Click filter buttons to see specific events:
- **All** - Show all events
- **Publish** - Only publish operations
- **Consume** - Only consume operations
- **Failed** - Only failed operations (for debugging)

### 📈 Metrics Tracking
Dashboard automatically tracks:
- Total messages (published/consumed)
- Error counts and success rates
- Real-time rates (messages/second)
- Processing performance (average time)
- Active operations (in progress)

---

## 🎨 Visual Indicators

### Status Badges

**Connection Status:**
- 🟢 `Connected` - API reachable, system healthy
- 🔴 `Disconnected` - Cannot reach API

**Event Status:**
- 🟢 `Success` - Operation completed successfully
- 🔴 `Failed` - Operation failed

**Event Types:**
- 🔵 `Publish` - Message published
- 🟣 `Consume` - Message consumed
- 🔴 `Failed` - Operation failed

### Color Scheme
```
Purple Gradient Background: #667eea → #764ba2
Cards: White with shadow
Hover Effect: Cards lift up
Chart Lines: Blue (publish), Purple (consume)
Success: Green badges
Error: Red badges
```

---

## 🔧 Customization

### Change API URL
If backend runs on different port:

**Edit `dashboard.js`:**
```javascript
const API_BASE_URL = 'http://localhost:8080/api/v1/monitoring';
// Change to your port:
const API_BASE_URL = 'http://localhost:9000/api/v1/monitoring';
```

### Change Refresh Rate
Default is 5 seconds. To change:

**Edit `dashboard.js`:**
```javascript
const REFRESH_INTERVAL = 5000; // 5 seconds
// Change to:
const REFRESH_INTERVAL = 3000;  // 3 seconds (faster)
const REFRESH_INTERVAL = 10000; // 10 seconds (slower)
```

### Show More Events
Default shows last 20 events. To show more:

**Edit `dashboard.js`, line ~75:**
```javascript
const endpoint = currentFilter === 'all'
    ? `${API_BASE_URL}/activity?limit=20`   // Change to 50
    : `${API_BASE_URL}/events/${currentFilter}?limit=20`;
```

### Change Colors
**Edit `dashboard.html` CSS:**
```css
/* Background gradient */
background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);

/* Change to green theme */
background: linear-gradient(135deg, #10b981 0%, #059669 100%);

/* Or blue theme */
background: linear-gradient(135deg, #3b82f6 0%, #1e40af 100%);
```

---

## 📱 Responsive Design

Dashboard adapts to screen size:

**Desktop (1920px+):**
- 4 cards per row
- Full chart width
- Large fonts

**Tablet (768px-1920px):**
- 2-3 cards per row
- Adjusted chart
- Medium fonts

**Mobile (320px-768px):**
- 1-2 cards per row
- Compact chart
- Small fonts

---

## 🐛 Troubleshooting

### Problem: Dashboard shows "Disconnected"

**Solution:**
```bash
# Check backend is running
curl http://localhost:8080/api/v1/monitoring/health

# If not running, start it
make docker-up
# or
make run
```

### Problem: Events not loading

**Solution:**
```bash
# Generate some traffic first
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","email":"test@example.com","phone":"123"}'

# Check if events exist
curl http://localhost:8080/api/v1/monitoring/events
```

### Problem: CORS errors in browser console

**Solution 1:** Use HTTP server (recommended)
```bash
cd web
python -m http.server 3000
# Open: http://localhost:3000/dashboard.html
```

**Solution 2:** CORS already configured in backend
```bash
# Make sure you're running the latest code
make docker-rebuild
```

### Problem: Chart not displaying

**Solution:**
```bash
# Check browser console (F12)
# - Look for Chart.js loading errors
# - Check if CDN is accessible

# Try clearing cache:
# Ctrl+Shift+R (Chrome/Firefox)
# Cmd+Shift+R (Mac)
```

---

## 💡 Tips & Tricks

### 1. Keep Dashboard Open While Testing
```bash
# Terminal 1: Run backend
make docker-up

# Terminal 2: Open dashboard
make dashboard

# Terminal 3: Generate traffic
make test-api
```

### 2. Quick Health Check
Watch dashboard while:
```bash
# Send test messages
for i in {1..100}; do
  curl -X POST http://localhost:8080/api/v1/messages \
    -H "Content-Type: application/json" \
    -d '{"user_id":1,"content":"Test '$i'","type":"test"}'
done
```

### 3. Debug Failed Messages
1. Click **"Failed"** filter button
2. Look at error messages in events
3. Check duration - slow operations might timeout
4. Review routing keys - wrong keys cause failures

### 4. Monitor Performance
Watch these metrics:
- **Success Rate** - Should be >99%
- **Avg Processing** - Should be <100ms
- **Rates** - Should match your expected load
- **In Progress** - Should be low (not stuck)

### 5. Keyboard Shortcuts
```
F5          - Manual refresh
F12         - Open dev tools (see console logs)
Ctrl + +/-  - Zoom in/out
Ctrl + 0    - Reset zoom
```

---

## 🎯 Common Use Cases

### Use Case 1: Development Monitoring
```bash
# Keep dashboard open while coding
make dashboard

# Watch for errors in real-time
# Click "Failed" filter when errors occur
```

### Use Case 2: Demo/Presentation
```bash
# Start everything
make docker-up

# Open dashboard in fullscreen
make dashboard
# Press F11 for fullscreen

# Generate demo traffic
make test-api
```

### Use Case 3: Load Testing
```bash
# Open dashboard
make dashboard

# In another terminal, run load test
# Watch metrics to ensure system handles load

# Check:
# - Rates stay stable
# - Success rate stays high
# - Processing time reasonable
# - No failed messages pile up
```

### Use Case 4: Debugging
```bash
# When you see errors:
# 1. Click "Failed" filter
# 2. Check error messages
# 3. Note the routing keys
# 4. Check processing times
# 5. Look for patterns (same error repeated?)
```

---

## 📚 Next Steps

After getting familiar with the dashboard:

1. **Customize it** - Change colors, refresh rate, etc.
2. **Monitor production** - Deploy behind nginx/apache
3. **Add alerts** - Integrate with monitoring tools
4. **Export data** - Add CSV export feature
5. **Add more charts** - Show more metrics visually

---

## 🔗 Related Resources

- **API Documentation**: `MONITORING_API.md`
- **Dashboard README**: `web/README.md`
- **Project README**: `README.md`

---

**Happy Monitoring!** 🎉

For issues or questions, check the logs:
```bash
make docker-logs
```
