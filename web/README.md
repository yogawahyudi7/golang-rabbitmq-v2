# RabbitMQ Monitoring Dashboard

Dashboard web real-time untuk monitoring RabbitMQ pub/sub activity.

## 🎨 Features

### Real-time Monitoring
- ✅ **Live Metrics Cards** - Published, Consumed, Failed, Success Rate
- ✅ **Rate Monitoring** - Publish/Consume rate per second
- ✅ **In-Progress Tracking** - Active operations counter
- ✅ **Performance Stats** - Average processing time

### Visual Charts
- ✅ **Real-time Line Chart** - Publish & Consume rates
- ✅ **Auto-updating** - Refreshes every 5 seconds
- ✅ **Smooth Animations** - Chart.js powered

### Event History
- ✅ **Recent Events Table** - Last 20 events
- ✅ **Event Filtering** - Filter by All, Publish, Consume, Failed
- ✅ **Color-coded Status** - Easy to identify success/failures
- ✅ **Time & Duration** - Track when and how long

### UI/UX
- ✅ **Modern Design** - Gradient background with clean cards
- ✅ **Responsive Layout** - Works on desktop and mobile
- ✅ **Connection Status** - Shows connected/disconnected
- ✅ **Auto-refresh** - Updates every 5 seconds

## 🚀 Quick Start

### 1. Start the Backend API

Pastikan aplikasi Go sudah berjalan:

```bash
# Option 1: Using Docker Compose
make docker-up

# Option 2: Run locally
make run
```

Backend akan berjalan di: `http://localhost:8080`

### 2. Open Dashboard

Buka file HTML di browser:

**Windows:**
```bash
# Langsung buka di browser
start web/dashboard.html

# Atau double-click file dashboard.html
```

**Linux/Mac:**
```bash
# Firefox
firefox web/dashboard.html

# Chrome
google-chrome web/dashboard.html

# Safari
open web/dashboard.html
```

### 3. Atau Serve dengan HTTP Server (Optional)

Untuk menghindari CORS issues, serve via HTTP server:

**Python:**
```bash
cd web
python -m http.server 3000
# Open: http://localhost:3000/dashboard.html
```

**Node.js:**
```bash
cd web
npx http-server -p 3000
# Open: http://localhost:3000/dashboard.html
```

**Go:**
```bash
cd web
go run -m http.server 3000
# Open: http://localhost:3000/dashboard.html
```

## 📊 Dashboard Components

### Header
- Connection status indicator (Connected/Disconnected)
- System uptime display
- Last updated timestamp

### Metrics Cards (Top Row)
1. **Published** - Total messages published
2. **Consumed** - Total messages consumed
3. **Failed** - Total failed operations
4. **Success Rate** - Percentage of successful operations

### Performance Cards (Second Row)
1. **Publish Rate** - Messages/second being published
2. **Consume Rate** - Messages/second being consumed
3. **In Progress** - Current operations in progress
4. **Avg Processing** - Average processing time

### Real-time Chart
- Line chart showing publish and consume rates over time
- Updates every 5 seconds
- Keeps last 20 data points
- Smooth animations

### Events Table
- Shows recent 20 events
- Filterable by type (All/Publish/Consume/Failed)
- Displays: Time, Type, Routing Key, Duration, Status
- Color-coded badges for easy identification

## ⚙️ Configuration

Edit `dashboard.js` untuk customize:

```javascript
// API endpoint (jika backend di port lain)
const API_BASE_URL = 'http://localhost:8080/api/v1/monitoring';

// Auto-refresh interval (ms)
const REFRESH_INTERVAL = 5000; // 5 seconds

// Max data points in chart
const MAX_CHART_POINTS = 20;
```

## 🎨 Customization

### Change Colors

Edit `dashboard.html` di section `<style>`:

```css
/* Main gradient background */
background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);

/* Card hover effect */
.card:hover {
    transform: translateY(-5px);
}

/* Chart colors (in dashboard.js) */
borderColor: 'rgb(59, 130, 246)',  // Publish line
borderColor: 'rgb(147, 51, 234)',  // Consume line
```

### Change Refresh Rate

Edit `dashboard.js`:

```javascript
const REFRESH_INTERVAL = 3000; // 3 seconds (faster)
const REFRESH_INTERVAL = 10000; // 10 seconds (slower)
```

### Show More Events

Edit `dashboard.js`:

```javascript
// In fetchActivity function
const endpoint = currentFilter === 'all'
    ? `${API_BASE_URL}/activity?limit=50`  // Show 50 events
    : `${API_BASE_URL}/events/${currentFilter}?limit=50`;
```

## 🔍 Troubleshooting

### Dashboard shows "Disconnected"

**Problem:** Cannot connect to API

**Solutions:**
1. Check if backend is running: `curl http://localhost:8080/api/v1/monitoring/health`
2. Check API_BASE_URL in `dashboard.js`
3. Check browser console for errors (F12)
4. Verify CORS is enabled in backend

### Events not loading

**Problem:** Events table empty or shows error

**Solutions:**
1. Generate some traffic first (create messages via API)
2. Check backend logs for errors
3. Verify monitoring endpoints are working: `curl http://localhost:8080/api/v1/monitoring/activity`

### Chart not updating

**Problem:** Chart stays static

**Solutions:**
1. Check browser console for JavaScript errors
2. Verify Chart.js loaded: Check Network tab (F12)
3. Make sure backend is receiving traffic
4. Check if rates are actually changing

### CORS errors in browser console

**Problem:** CORS policy blocking requests

**Solutions:**
1. CORS already configured in backend (`middleware/cors.go`)
2. If still issues, serve dashboard via HTTP server (see Quick Start #3)
3. Or use browser extension to disable CORS for development

## 📱 Mobile Responsive

Dashboard is responsive and works on:
- ✅ Desktop (1920x1080+)
- ✅ Laptop (1366x768+)
- ✅ Tablet (768px+)
- ✅ Mobile (320px+)

## 🎯 Use Cases

### 1. Development Monitoring
Monitor your local development in real-time while testing.

### 2. Demo/Presentation
Show live metrics during presentations or demos.

### 3. Debugging
Quickly identify failed messages and performance issues.

### 4. Performance Testing
Monitor rates and processing times during load tests.

## 🔗 API Endpoints Used

Dashboard menggunakan monitoring API endpoints:

- `GET /api/v1/monitoring/metrics/detailed` - Main metrics
- `GET /api/v1/monitoring/activity?limit=20` - Real-time activity
- `GET /api/v1/monitoring/events/publish?limit=20` - Publish events
- `GET /api/v1/monitoring/events/consume?limit=20` - Consume events
- `GET /api/v1/monitoring/events/failed?limit=20` - Failed events

## 🚀 Production Deployment

Untuk production, consider:

1. **Serve via web server** (Nginx, Apache, Caddy)
2. **Add authentication** - Protect dashboard access
3. **HTTPS** - Use SSL certificates
4. **Reverse proxy** - Proxy API calls through same origin
5. **Caching** - Cache static assets
6. **Monitoring** - Add error tracking (Sentry, etc.)

### Example Nginx Config

```nginx
server {
    listen 80;
    server_name dashboard.example.com;

    # Serve dashboard
    location / {
        root /path/to/web;
        index dashboard.html;
    }

    # Proxy API requests
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 📚 Dependencies

- **Chart.js 4.4.0** - Via CDN (no installation needed)
- Modern browser with JavaScript enabled
- Backend API running on port 8080

## 💡 Tips

1. **Keep dashboard open** - Auto-refreshes every 5 seconds
2. **Use filters** - Filter events by type for focused debugging
3. **Watch the chart** - Spot patterns and anomalies in real-time
4. **Check failed events** - Use "Failed" filter to debug errors
5. **Monitor success rate** - Should be close to 100% in healthy system

## 🐛 Known Issues

- Dashboard keeps data in memory only (resets on refresh)
- Max 20 chart points (configurable)
- No historical data export (yet)
- Single page only (no multi-page support)

## 🔄 Future Enhancements

Possible improvements:
- [ ] Export data to CSV
- [ ] Alert notifications
- [ ] Dark mode toggle
- [ ] Multiple chart types
- [ ] Historical data persistence
- [ ] Real-time WebSocket updates
- [ ] User authentication
- [ ] Custom time ranges

---

**Enjoy monitoring!** 🎉
