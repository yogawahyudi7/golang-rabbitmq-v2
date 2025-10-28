# 🚀 Quick Start Test Log

## Prerequisites Check

### 1. Docker Desktop Status
**Issue:** Docker Desktop tidak running
```
Error: cannot connect to Docker daemon
The system cannot find the file specified.
```

**Solution:** Start Docker Desktop
- Windows: Buka Docker Desktop dari Start Menu
- Wait hingga Docker icon di system tray menunjukkan "Docker Desktop is running"

---

## Quick Start Steps (After Docker Running)

### Step 1: Start Services
```bash
cd C:\Users\yogaw\playgroud\golang-rabbitmq-v2
make docker-up
```

**Expected Output:**
```
Starting Docker services...
Creating network "golang-rabbitmq-v2_default" with the default driver
Creating golang-rabbitmq-v2_postgres_1 ... done
Creating golang-rabbitmq-v2_rabbitmq_1 ... done
Creating golang-rabbitmq-v2_app_1      ... done
✅ Services started

Container Status:
NAME                           STATUS              PORTS
golang-rabbitmq-v2_app_1      Up 2 seconds        0.0.0.0:8080->8080/tcp
golang-rabbitmq-v2_postgres_1 Up 3 seconds        0.0.0.0:5432->5432/tcp
golang-rabbitmq-v2_rabbitmq_1 Up 3 seconds        0.0.0.0:5672->5672/tcp, 0.0.0.0:15672->15672/tcp
```

**Wait ~30 seconds** untuk services fully ready.

---

### Step 2: Verify Services

```bash
# Check containers
make docker-ps

# Check health
make health
```

**Expected Health Response:**
```json
{
  "status": "healthy",
  "timestamp": "45.2s",
  "uptime_seconds": 45.2,
  "rabbitmq": {
    "connected": true,
    "messages_published": 0,
    "messages_consumed": 0
  }
}
```

---

### Step 3: Open Dashboard

```bash
make dashboard
```

**What happens:**
1. Starts Python HTTP server on port 3000
2. Opens browser to `http://localhost:3000/dashboard.html`
3. Dashboard shows "Connected" status

**Expected Dashboard View:**
```
┌─────────────────────────────────────┐
│ 🐰 RabbitMQ Monitoring Dashboard   │
│ ● Connected | Uptime: 0h 0m 45s    │
└─────────────────────────────────────┘

┌──────────┬──────────┬──────────┬──────────────┐
│ Published│ Consumed │  Failed  │ Success Rate │
│    0     │    0     │    0     │     0%       │
└──────────┴──────────┴──────────┴──────────────┘
```

---

### Step 4: Generate Test Traffic

Open **new terminal**:

```bash
cd C:\Users\yogaw\playgroud\golang-rabbitmq-v2
make test-api
```

**Expected Output:**
```
Testing API endpoints...
1. Health Check:
"healthy"

2. Create Test User:
{
  "id": 1,
  "name": "Test User",
  "email": "test@example.com",
  "phone": "1234567890",
  "created_at": "2025-10-28T13:57:00Z"
}

3. Message Stats:
{
  "total_messages": 0,
  "pending_messages": 0,
  "failed_messages": 0
}
```

---

### Step 5: Watch Dashboard Update!

Go back to dashboard browser tab.

**Watch these update (auto-refresh every 5 seconds):**

1. **Metrics Cards:**
   - Published: 0 → 1 → 2 → ...
   - Consumed: 0 → 1 → 2 → ...
   - Success Rate: 0% → 100%

2. **Chart:**
   - Lines start appearing
   - Blue line (publish rate)
   - Purple line (consume rate)

3. **Events Table:**
   - New rows appear
   - Shows publish events
   - Shows consume events
   - All marked as "Success" (green)

---

## ✅ Success Indicators

You know it's working when:

### Backend
- ✅ `make docker-ps` shows 3 containers running
- ✅ `make health` returns `"status": "healthy"`
- ✅ No errors in `make docker-logs`

### Dashboard
- ✅ Shows "● Connected" (green)
- ✅ Uptime counter increases
- ✅ Metrics cards show numbers
- ✅ Chart has data lines
- ✅ Events table populated

### API
- ✅ Can create users
- ✅ Can create messages
- ✅ Stats show correct counts

---

## 🎯 Complete Test Sequence

Run this complete sequence:

### Terminal 1: Backend & Logs
```bash
cd C:\Users\yogaw\playgroud\golang-rabbitmq-v2

# Start services
make docker-up

# Wait 30 seconds
timeout 30

# Check logs
make docker-logs
```

### Terminal 2: Dashboard
```bash
cd C:\Users\yogaw\playgroud\golang-rabbitmq-v2

# Open dashboard
make dashboard

# Keep this terminal open (HTTP server running)
```

### Terminal 3: Testing
```bash
cd C:\Users\yogaw\playgroud\golang-rabbitmq-v2

# Test health
curl http://localhost:8080/api/v1/monitoring/health

# Test API
make test-api

# Generate more traffic
for i in {1..20}; do
  curl -X POST http://localhost:8080/api/v1/messages \
    -H "Content-Type: application/json" \
    -d "{\"user_id\":1,\"content\":\"Test message $i\",\"type\":\"test\"}"
  echo "Message $i sent"
  sleep 1
done
```

### Browser: Dashboard
```
URL: http://localhost:3000/dashboard.html

Watch:
- Connection status: ● Connected
- Published count increasing
- Consumed count increasing
- Chart lines moving
- Events appearing in table
```

---

## 📊 Expected Results After 5 Minutes

### Metrics
```
Published:     ~20-25
Consumed:      ~20-25
Failed:        0
Success Rate:  100%

Publish Rate:  ~0.3 msg/s (depending on test speed)
Consume Rate:  ~0.3 msg/s
In Progress:   0
Avg Process:   ~40-60ms
```

### Chart
- Should show smooth lines
- Blue line (publish) and purple line (consume) close together
- No gaps or drops

### Events Table
```
Recent Events showing:
- Mix of "Publish" and "Consume" events
- All with "Success" status (green)
- Durations 5-50ms
- Routing key: "app.message"
```

---

## 🐛 Troubleshooting

### Issue 1: Docker Not Running

**Error:**
```
Error: cannot connect to Docker daemon
```

**Fix:**
1. Open Docker Desktop
2. Wait for "Docker Desktop is running"
3. Run `docker ps` to verify
4. Then retry `make docker-up`

---

### Issue 2: Port Already in Use

**Error:**
```
Bind for 0.0.0.0:8080 failed: port is already allocated
```

**Fix:**
```bash
# Find process using port
netstat -ano | findstr :8080

# Kill process (replace PID)
taskkill /PID <PID> /F

# Or use different port in docker-compose.yml
```

---

### Issue 3: Dashboard Shows "Disconnected"

**Symptoms:**
- Red badge: "● Disconnected"
- All metrics show 0
- Chart empty

**Fix:**
```bash
# 1. Check backend running
make docker-ps

# 2. Check health
curl http://localhost:8080/api/v1/monitoring/health

# If fails:
make docker-down
make docker-up

# Wait 30 seconds
# Refresh dashboard (F5)
```

---

### Issue 4: No Events in Table

**Symptoms:**
- Table shows "No events found"
- Events count: 0

**Fix:**
```bash
# Generate traffic first
make test-api

# Wait 5 seconds for auto-refresh
# Or manually refresh (F5)
```

---

### Issue 5: Python HTTP Server Issues

**Error:**
```
'python' is not recognized as an internal or external command
```

**Fix Options:**

**Option 1: Direct open (may have CORS)**
```bash
start web/dashboard.html
```

**Option 2: Install Python**
- Download from python.org
- Install with "Add to PATH" checked
- Retry `make dashboard`

**Option 3: Use Node.js**
```bash
cd web
npx http-server -p 3000
# Open: http://localhost:3000/dashboard.html
```

---

## 📝 Test Checklist

Run through this checklist:

- [ ] Docker Desktop running
- [ ] `make docker-up` successful
- [ ] 3 containers running (`make docker-ps`)
- [ ] Health check passes (`make health`)
- [ ] Dashboard opens (`make dashboard`)
- [ ] Dashboard shows "Connected"
- [ ] Test API works (`make test-api`)
- [ ] Metrics cards update
- [ ] Chart shows data
- [ ] Events appear in table
- [ ] Can filter events (All/Publish/Consume/Failed)
- [ ] Auto-refresh working (wait 5 seconds)

**All checked?** ✅ Perfect! Everything working!

---

## 🎉 Success!

If you can:
1. ✅ See 3 containers running
2. ✅ Dashboard shows "Connected"
3. ✅ Metrics update after `make test-api`
4. ✅ Chart displays lines
5. ✅ Events show in table

**Congratulations!** 🎊 Quick start berhasil 100%!

---

## 🔄 Stop Everything

When done testing:

### Stop HTTP Server
In Terminal 2 (dashboard):
```
Press Ctrl+C
```

### Stop Docker Services
```bash
make docker-down
```

**Expected Output:**
```
Stopping Docker services...
Stopping golang-rabbitmq-v2_app_1      ... done
Stopping golang-rabbitmq-v2_rabbitmq_1 ... done
Stopping golang-rabbitmq-v2_postgres_1 ... done
Removing golang-rabbitmq-v2_app_1      ... done
Removing golang-rabbitmq-v2_rabbitmq_1 ... done
Removing golang-rabbitmq-v2_postgres_1 ... done
✅ Services stopped
```

---

## 📚 Next Steps

After successful quick start:

1. **Learn Dashboard Features**
   - Read: `DASHBOARD_GUIDE.md`
   - Try all filter buttons
   - Customize colors/rates

2. **Try API Endpoints**
   - Read: `MONITORING_API.md`
   - Test with curl/Postman
   - Build custom tools

3. **Load Testing**
   - Generate bulk traffic
   - Watch performance metrics
   - Check success rate

4. **Development**
   - Read: `GETTING_STARTED.md`
   - Setup local development
   - Make code changes

---

**Happy Monitoring!** 🚀
