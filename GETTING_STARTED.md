# 🚀 Getting Started - RabbitMQ Monitoring Dashboard

Panduan lengkap untuk memulai menggunakan RabbitMQ monitoring dashboard dari awal hingga berjalan.

## 📋 Prerequisites

Pastikan sudah terinstall:
- ✅ **Docker & Docker Compose** - Untuk menjalankan services
- ✅ **Git** (optional) - Untuk clone repository
- ✅ **Go 1.21+** (optional) - Jika ingin run tanpa Docker
- ✅ **Python** (optional) - Untuk HTTP server dashboard

### Check Prerequisites

```bash
# Check Docker
docker --version
docker-compose --version

# Check Python (optional)
python --version

# Check Go (optional)
go version
```

---

## 🎯 Quick Start (5 Menit!)

### **Step 1: Clone & Navigate**

```bash
# Clone repository (jika belum)
git clone <repository-url>
cd golang-rabbitmq-v2

# Atau jika sudah ada, navigate ke folder
cd C:\Users\yogaw\playgroud\golang-rabbitmq-v2
```

### **Step 2: Start Backend Services**

```bash
# Start semua services (PostgreSQL, RabbitMQ, App)
make docker-up

# Tunggu ~30 detik untuk services siap
# Output yang diharapkan:
# ✅ Services started
# Container Status:
# rabbitmq_app    Up    0.0.0.0:8080->8080/tcp
# postgres        Up    0.0.0.0:5432->5432/tcp
# rabbitmq        Up    0.0.0.0:5672->5672/tcp, 0.0.0.0:15672->15672/tcp
```

### **Step 3: Verify Services Running**

```bash
# Check container status
make docker-ps

# Test health endpoint
make health

# Output yang diharapkan:
# {
#   "status": "healthy",
#   "uptime_seconds": 45.2,
#   "rabbitmq": {
#     "connected": true,
#     "messages_published": 0,
#     "messages_consumed": 0
#   }
# }
```

### **Step 4: Open Dashboard**

**Windows:**
```bash
make dashboard
```

**Linux/Mac:**
```bash
make dashboard
```

**Atau manual:**
```bash
cd web
# Windows
.\open-dashboard.bat

# Linux/Mac
./open-dashboard.sh
```

### **Step 5: Generate Sample Data**

Di terminal baru:

```bash
# Test API dan generate traffic
make test-api

# Atau manual create messages
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","phone":"1234567890"}'

curl -X POST http://localhost:8080/api/v1/messages \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"content":"Hello RabbitMQ","type":"general"}'
```

### **Step 6: Watch Dashboard Update!** 🎉

Dashboard akan auto-refresh setiap 5 detik dan menampilkan:
- ✅ Metrics cards terupdate
- ✅ Chart menampilkan rates
- ✅ Events table menampilkan aktivitas

**Selesai!** Dashboard sudah berjalan! 🎊

---

## 🔍 Detailed Step-by-Step Guide

### **Option A: Using Docker (Recommended)**

#### Step 1: Setup Project

```bash
# Navigate to project
cd golang-rabbitmq-v2

# (Optional) Setup environment
make dev-setup

# Ini akan create .env file dari .env.example
# Edit .env jika perlu custom configuration
```

#### Step 2: Start All Services

```bash
# Start PostgreSQL, RabbitMQ, and Application
make docker-up

# Wait for services to initialize (~30 seconds)
# You can watch logs:
make docker-logs
```

**What's happening:**
- 🐘 PostgreSQL starts on port 5432
- 🐰 RabbitMQ starts on port 5672 (AMQP) and 15672 (Management UI)
- 🚀 Go Application starts on port 8080

#### Step 3: Verify Everything is Running

```bash
# Check containers
docker ps

# Should see 3 containers:
# - golang-rabbitmq-v2_app
# - golang-rabbitmq-v2_postgres
# - golang-rabbitmq-v2_rabbitmq

# Test application health
curl http://localhost:8080/api/v1/monitoring/health

# Test RabbitMQ management (browser)
start http://localhost:15672
# Login: guest / guest
```

#### Step 4: Open Monitoring Dashboard

**Method 1: Using Makefile**
```bash
make dashboard
```

**Method 2: Using Launcher Script**
```bash
cd web
# Windows
.\open-dashboard.bat

# Linux/Mac
chmod +x open-dashboard.sh
./open-dashboard.sh
```

**Method 3: Direct Open**
```bash
# Just double-click or:
start web/dashboard.html         # Windows
open web/dashboard.html          # Mac
xdg-open web/dashboard.html      # Linux
```

**Method 4: Python HTTP Server (No CORS issues)**
```bash
cd web
python -m http.server 3000
# Open: http://localhost:3000/dashboard.html
```

#### Step 5: Generate Test Traffic

Open new terminal:

```bash
# Create test user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com","phone":"0812345678"}'

# Create messages
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"content":"Test message 1","type":"general"}'

curl -X POST http://localhost:8080/api/v1/messages \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"content":"Test message 2","type":"urgent"}'

# Generate bulk traffic
for i in {1..20}; do
  curl -X POST http://localhost:8080/api/v1/messages \
    -H "Content-Type: application/json" \
    -d '{"user_id":1,"content":"Message '$i'","type":"test"}'
  echo "Message $i sent"
  sleep 1
done
```

#### Step 6: Monitor in Real-time

Watch the dashboard:
- ✅ **Published count** increases
- ✅ **Consumed count** increases
- ✅ **Chart** shows real-time rates
- ✅ **Events table** shows each publish/consume

Try filtering:
- Click **"Publish"** - See only publish events
- Click **"Consume"** - See only consume events
- Click **"All"** - See everything

---

### **Option B: Running Without Docker**

#### Step 1: Start PostgreSQL & RabbitMQ

```bash
# Start only PostgreSQL and RabbitMQ
docker-compose up -d postgres rabbitmq

# Wait for services to be ready
sleep 10
```

#### Step 2: Setup Environment

```bash
# Copy env file
cp .env.example .env

# Edit .env if needed
# Make sure these are set:
# DB_HOST=localhost
# RABBITMQ_URL=amqp://guest:guest@localhost:5672/
```

#### Step 3: Install Dependencies

```bash
# Download Go modules
go mod download
go mod tidy
```

#### Step 4: Run Application

```bash
# Run directly
go run cmd/api/main.go

# Or build first
make build-local
./bin/rabbitmq-app
```

**Application should start on port 8080**

#### Step 5: Open Dashboard & Test

Same as Docker option (Step 4-6)

---

## 🎯 Multiple Terminal Setup (Recommended Workflow)

### **Terminal 1: Backend Services**
```bash
cd golang-rabbitmq-v2
make docker-up
make docker-logs
# Keep this open to watch logs
```

### **Terminal 2: Dashboard**
```bash
cd golang-rabbitmq-v2
make dashboard
# Dashboard opens and HTTP server runs
# Keep this open
```

### **Terminal 3: Testing & Traffic Generation**
```bash
cd golang-rabbitmq-v2

# Quick tests
make test-api

# Or manual commands
curl http://localhost:8080/api/v1/monitoring/health
curl http://localhost:8080/api/v1/monitoring/metrics
curl http://localhost:8080/api/v1/monitoring/activity

# Generate traffic
for i in {1..50}; do
  curl -X POST http://localhost:8080/api/v1/messages \
    -H "Content-Type: application/json" \
    -d '{"user_id":1,"content":"Bulk test '$i'","type":"load-test"}'
done
```

### **Browser: Dashboard**
```
Auto-opened by Terminal 2:
http://localhost:3000/dashboard.html

Watch real-time updates!
```

---

## 🔧 Common Scenarios

### **Scenario 1: Fresh Start (Baru Install)**

```bash
# 1. Clone & navigate
git clone <repo-url>
cd golang-rabbitmq-v2

# 2. Quick start
make quick-start

# 3. Open dashboard
make dashboard

# 4. Test
make test-api

# Done! 🎉
```

---

### **Scenario 2: Restart After Shutdown**

```bash
# 1. Navigate to project
cd golang-rabbitmq-v2

# 2. Start services
make docker-up

# 3. Open dashboard
make dashboard

# Dashboard shows previous data if containers not removed
```

---

### **Scenario 3: Clean Restart (Reset Everything)**

```bash
# 1. Stop and remove everything
make docker-down

# 2. Clean volumes (removes all data)
docker-compose down -v

# 3. Fresh start
make docker-up

# 4. Dashboard
make dashboard

# Fresh empty metrics!
```

---

### **Scenario 4: Development Mode (Code Changes)**

```bash
# Terminal 1: Run with live reload
make run-watch
# Requires: make install-tools (installs air)

# Terminal 2: Dashboard
make dashboard

# Terminal 3: Test changes
make test-api

# Code changes auto-reload!
```

---

### **Scenario 5: Load Testing**

```bash
# Terminal 1: Services
make docker-up
make docker-logs

# Terminal 2: Dashboard
make dashboard

# Terminal 3: Generate high load
for i in {1..1000}; do
  curl -X POST http://localhost:8080/api/v1/messages \
    -H "Content-Type: application/json" \
    -d '{"user_id":1,"content":"Load test '$i'","type":"stress"}' &
done
wait

# Watch dashboard handle the load!
# Monitor:
# - Success rate stays high
# - Processing time reasonable
# - No errors accumulate
```

---

## 🎨 Dashboard Access URLs

Once started, you can access:

| Service | URL | Credentials |
|---------|-----|-------------|
| **Monitoring Dashboard** | `http://localhost:3000/dashboard.html` | - |
| **API Health** | `http://localhost:8080/api/v1/monitoring/health` | - |
| **API Metrics** | `http://localhost:8080/api/v1/monitoring/metrics` | - |
| **API Activity** | `http://localhost:8080/api/v1/monitoring/activity` | - |
| **RabbitMQ Management** | `http://localhost:15672` | guest/guest |

---

## ✅ Verification Checklist

After starting, verify:

- [ ] **Docker containers running**
  ```bash
  make docker-ps
  # Should show 3 containers: app, postgres, rabbitmq
  ```

- [ ] **Health endpoint responds**
  ```bash
  curl http://localhost:8080/api/v1/monitoring/health
  # Should return: {"status":"healthy"...}
  ```

- [ ] **Dashboard opens**
  ```bash
  make dashboard
  # Browser opens with dashboard
  ```

- [ ] **Dashboard shows "Connected"**
  - Green badge in header
  - Uptime displayed

- [ ] **Metrics displayed**
  - Cards show values (may be 0 initially)

- [ ] **Chart visible**
  - Empty at first
  - Updates after traffic

- [ ] **Events table visible**
  - May be empty initially
  - Shows "No events found"

- [ ] **Can create messages**
  ```bash
  curl -X POST http://localhost:8080/api/v1/messages \
    -H "Content-Type: application/json" \
    -d '{"user_id":1,"content":"Test","type":"test"}'
  ```

- [ ] **Dashboard updates**
  - Published count increases
  - Events appear in table

All green? **You're ready to go!** ✅

---

## 🐛 Troubleshooting Startup Issues

### **Issue 1: Port Already in Use**

**Error:**
```
Error: bind: address already in use
```

**Solution:**
```bash
# Check what's using the port
# Windows
netstat -ano | findstr :8080
netstat -ano | findstr :5432
netstat -ano | findstr :5672

# Linux/Mac
lsof -i :8080
lsof -i :5432
lsof -i :5672

# Stop the process or change ports in docker-compose.yml
```

---

### **Issue 2: Docker Not Running**

**Error:**
```
Cannot connect to the Docker daemon
```

**Solution:**
```bash
# Start Docker Desktop
# Windows: Open Docker Desktop
# Linux: sudo systemctl start docker
# Mac: Open Docker app

# Verify
docker ps
```

---

### **Issue 3: Containers Not Starting**

**Error:**
```
Container exited with code 1
```

**Solution:**
```bash
# Check logs
make docker-logs

# Or specific container
docker-compose logs app
docker-compose logs postgres
docker-compose logs rabbitmq

# Common fixes:
# 1. Remove old containers
make docker-clean

# 2. Rebuild
make docker-rebuild

# 3. Check .env file
cat .env
```

---

### **Issue 4: Dashboard Shows "Disconnected"**

**Problem:** Dashboard can't reach API

**Solution:**
```bash
# 1. Verify backend is running
curl http://localhost:8080/api/v1/monitoring/health

# If fails:
# 2. Check container status
make docker-ps

# 3. Check logs
make docker-logs

# 4. Restart services
make docker-down
make docker-up
```

---

### **Issue 5: Dashboard Not Opening**

**Problem:** `make dashboard` doesn't open browser

**Solution:**
```bash
# Option 1: Open manually
start http://localhost:3000/dashboard.html  # Windows
open http://localhost:3000/dashboard.html   # Mac
xdg-open http://localhost:3000/dashboard.html # Linux

# Option 2: Direct file open
start web/dashboard.html  # Windows
open web/dashboard.html   # Mac
```

---

### **Issue 6: Events Not Showing**

**Problem:** Events table is empty

**Solution:**
```bash
# 1. Generate some traffic first
make test-api

# 2. Check if events endpoint works
curl http://localhost:8080/api/v1/monitoring/events

# 3. Check browser console (F12)
# Look for errors

# 4. Try filtering
# Click "All" button in dashboard
```

---

### **Issue 7: CORS Errors in Console**

**Problem:** Browser console shows CORS errors

**Solution:**
```bash
# 1. Use HTTP server (recommended)
cd web
python -m http.server 3000
# Open: http://localhost:3000/dashboard.html

# 2. Or check CORS middleware
# Should already be configured in:
# internal/middleware/cors.go

# 3. Restart backend
make docker-rebuild
```

---

## 🎓 Next Steps

After successfully starting:

1. **Read Documentation**
   - `DASHBOARD_GUIDE.md` - Dashboard usage guide
   - `MONITORING_API.md` - API documentation
   - `README.md` - Project overview

2. **Explore Features**
   - Try event filtering
   - Watch real-time charts
   - Generate load tests

3. **Customize**
   - Edit colors in `web/dashboard.html`
   - Adjust refresh rate in `web/dashboard.js`
   - Add more metrics

4. **Integrate**
   - Use API in your own tools
   - Export data to monitoring systems
   - Add alerts

---

## 📞 Need Help?

**Quick Reference:**
```bash
# Show all available commands
make help

# Show project info
make info

# Check health
make health

# Show metrics
make metrics

# View logs
make docker-logs

# Clean and restart
make docker-clean
make docker-up
```

**Check Logs:**
```bash
# Application logs
make docker-logs

# All services
make docker-logs-all

# Log files (if configured)
ls logs/
```

**Test Endpoints:**
```bash
# Health
curl http://localhost:8080/api/v1/monitoring/health

# Metrics
curl http://localhost:8080/api/v1/monitoring/metrics | jq

# Activity
curl http://localhost:8080/api/v1/monitoring/activity | jq

# Events
curl http://localhost:8080/api/v1/monitoring/events | jq
```

---

## 🎉 Success!

If you see:
- ✅ Dashboard showing "Connected"
- ✅ Metrics cards with values
- ✅ Charts updating
- ✅ Events appearing

**Congratulations!** You're all set up! 🎊

Now you can:
- Monitor your RabbitMQ in real-time
- Debug issues quickly
- Track performance
- Demo to stakeholders

**Happy Monitoring!** 🚀

---

**Quick Commands Cheat Sheet:**

```bash
make docker-up      # Start everything
make dashboard      # Open dashboard
make test-api       # Generate test traffic
make docker-logs    # Watch logs
make health         # Check health
make docker-down    # Stop everything
```
