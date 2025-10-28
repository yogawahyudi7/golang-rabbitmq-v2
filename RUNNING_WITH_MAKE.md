# 🚀 Running dengan Make Commands

Panduan lengkap running project ini menggunakan Make commands.

---

## ⚡ Super Quick Start (1 Command!)

### **Cara Tercepat:**

```bash
make start
```

**Itu saja!** 🎉

Command ini akan:
1. ✅ Start semua Docker services (PostgreSQL, RabbitMQ, App)
2. ✅ Wait 30 detik untuk services ready
3. ✅ Test API endpoints
4. ✅ Show instructions untuk buka dashboard

**Total waktu:** ~1-2 menit (including Docker startup)

---

## 📋 What `make start` Does

```bash
make start
```

**Output yang akan Anda lihat:**

```
🚀 Starting Full Quick Start Demo...

Step 1/4: Starting Docker services...
Creating network "golang-rabbitmq-v2_default"...
Creating postgres... done
Creating rabbitmq... done
Creating app... done
✅ Services started

Container Status:
NAME                           STATUS              PORTS
golang-rabbitmq-v2_app         Up 2 seconds        0.0.0.0:8080->8080/tcp
golang-rabbitmq-v2_postgres    Up 3 seconds        0.0.0.0:5432->5432/tcp
golang-rabbitmq-v2_rabbitmq    Up 3 seconds        0.0.0.0:5672->5672/tcp

Step 2/4: Waiting for services to be ready...
⏳ Please wait 30 seconds...

Step 3/4: Testing API...
Testing API endpoints...
1. Health Check:
"healthy"

2. Create Test User:
{
  "id": 1,
  "name": "Test User",
  "email": "test@example.com"
}

3. Message Stats:
{
  "total_messages": 0
}

Step 4/4: Opening dashboard...
🎨 Dashboard will open in your browser...
📊 You can also run: make dashboard

✅ Quick Start Complete!

📍 Services Running:
  • Application: http://localhost:8080
  • Dashboard: Run 'make dashboard' to open
  • RabbitMQ UI: http://localhost:15672 (guest/guest)

🔧 Useful Commands:
  • make dashboard    - Open monitoring dashboard
  • make health       - Check application health
  • make metrics      - Show RabbitMQ metrics
  • make docker-logs  - View application logs
  • make docker-down  - Stop all services

🎉 Ready to monitor! Run 'make dashboard' now!
```

---

## 🎯 Next Step After `make start`

### **Open Dashboard:**

```bash
make dashboard
```

Dashboard akan:
- Start HTTP server di port 3000
- Open browser ke `http://localhost:3000/dashboard.html`
- Show real-time monitoring

---

## 📚 All Available Make Commands

### **Starting & Stopping**

```bash
make start           # ⭐ Full quick start (services + test)
make docker-up       # Start services only
make docker-down     # Stop all services
make docker-restart  # Restart services
make docker-rebuild  # Rebuild and restart
```

### **Monitoring & Dashboard**

```bash
make dashboard       # Open monitoring dashboard
make health          # Check application health
make metrics         # Show RabbitMQ metrics
make rabbitmq-ui     # Open RabbitMQ Management UI
```

### **Testing**

```bash
make test-api        # Test API endpoints
make test            # Run Go tests
make test-coverage   # Run tests with coverage
```

### **Development**

```bash
make run             # Run locally (without Docker)
make run-watch       # Run with live reload
make build           # Build binary
make build-local     # Build for local platform
```

### **Logs & Status**

```bash
make docker-ps       # Show container status
make docker-logs     # View app logs
make docker-logs-all # View all logs
make docker-shell    # Access container shell
```

### **Cleanup**

```bash
make clean           # Clean build artifacts
make clean-logs      # Clean log files
make docker-clean    # Clean Docker resources
make db-reset        # Reset database
make rabbitmq-reset  # Reset RabbitMQ
```

### **Info & Help**

```bash
make help            # Show all commands
make info            # Show project info
```

---

## 🎬 Complete Workflow Examples

### **Example 1: First Time Setup**

```bash
# Step 1: Start everything
make start

# Step 2: Open dashboard (in new terminal)
make dashboard

# Step 3: Generate more traffic (in new terminal)
for i in {1..20}; do
  curl -X POST http://localhost:8080/api/v1/messages \
    -H "Content-Type: application/json" \
    -d "{\"user_id\":1,\"content\":\"Test $i\",\"type\":\"test\"}"
done

# Step 4: Watch dashboard update!
```

---

### **Example 2: Daily Development**

```bash
# Morning: Start services
make start

# Open dashboard
make dashboard

# During development: Watch logs
make docker-logs

# Test changes
make test-api

# Evening: Stop
make docker-down
```

---

### **Example 3: Debugging**

```bash
# Start with logs
make start

# Check status
make docker-ps

# Check health
make health

# View logs
make docker-logs

# Check dashboard
make dashboard
```

---

### **Example 4: Clean Restart**

```bash
# Stop and clean everything
make docker-down
make docker-clean

# Fresh start
make start

# Open dashboard
make dashboard
```

---

## 🔍 Command Details

### **`make start` - Full Quick Start**

**What it does:**
1. Starts Docker Compose services
2. Waits 30 seconds for initialization
3. Runs API tests
4. Shows instructions

**When to use:**
- First time setup
- After restart
- Demo/presentation
- Testing from scratch

**Dependencies:**
- Docker Desktop running
- No other services on ports 8080, 5432, 5672, 15672

---

### **`make dashboard` - Open Dashboard**

**What it does:**
1. Detects Python availability
2. Starts HTTP server (if Python available)
3. Opens browser to dashboard

**When to use:**
- After `make start`
- For monitoring
- During development
- During testing

**Alternative:**
```bash
# Direct open (may have CORS issues)
make dashboard-direct
```

---

### **`make test-api` - Test Endpoints**

**What it does:**
1. Health check
2. Create test user
3. Show message stats

**When to use:**
- Verify services working
- Generate initial traffic
- Test API functionality

**Output:**
```json
1. Health Check:
"healthy"

2. Create Test User:
{ "id": 1, "name": "Test User" }

3. Message Stats:
{ "total_messages": 0 }
```

---

### **`make docker-logs` - View Logs**

**What it does:**
- Shows real-time application logs
- Follows log output

**When to use:**
- Debugging issues
- Monitoring activity
- Checking errors

**Tip:** Press `Ctrl+C` to exit

---

### **`make health` - Health Check**

**What it does:**
- Checks `/api/v1/monitoring/health`
- Pretty prints JSON

**Expected output:**
```json
{
  "status": "healthy",
  "timestamp": "1h23m45s",
  "uptime_seconds": 5025.5,
  "rabbitmq": {
    "connected": true,
    "messages_published": 150,
    "messages_consumed": 148
  }
}
```

---

## ⚠️ Prerequisites

Before running `make start`:

### **1. Docker Desktop Must Be Running**

**Check:**
```bash
docker ps
```

**If error:** Start Docker Desktop first

**Windows:**
- Start Menu → Docker Desktop
- Wait for "Docker Desktop is running"

---

### **2. Ports Must Be Available**

Required ports:
- `8080` - Application API
- `5432` - PostgreSQL
- `5672` - RabbitMQ AMQP
- `15672` - RabbitMQ Management UI
- `3000` - Dashboard HTTP server (when using `make dashboard`)

**Check ports:**
```bash
# Windows
netstat -ano | findstr :8080
netstat -ano | findstr :5432
netstat -ano | findstr :5672

# Linux/Mac
lsof -i :8080
lsof -i :5432
lsof -i :5672
```

---

### **3. Make Command Available**

**Check:**
```bash
make --version
```

**If not available:**
- Windows: Install via Chocolatey `choco install make`
- Or use Git Bash (comes with Git for Windows)
- Or use WSL

---

## 🐛 Troubleshooting

### **Issue 1: `make: command not found`**

**Solution:**
```bash
# Use Git Bash instead of CMD
# Or install make:
choco install make
```

---

### **Issue 2: Docker not running**

**Error:**
```
Cannot connect to Docker daemon
```

**Solution:**
1. Start Docker Desktop
2. Wait for "running" status
3. Run `docker ps` to verify
4. Retry `make start`

---

### **Issue 3: Port already in use**

**Error:**
```
Port 8080 is already allocated
```

**Solution:**
```bash
# Find process
netstat -ano | findstr :8080

# Kill process
taskkill /PID <PID> /F

# Or stop other services first
```

---

### **Issue 4: Services not healthy**

**Symptoms:**
```
make health
❌ Health check failed
```

**Solution:**
```bash
# Check logs
make docker-logs

# Check status
make docker-ps

# Restart if needed
make docker-rebuild
```

---

### **Issue 5: Dashboard doesn't open**

**Symptoms:**
- `make dashboard` doesn't open browser
- Or shows errors

**Solution:**
```bash
# Option 1: Direct open
make dashboard-direct

# Option 2: Manual
cd web
start dashboard.html   # Windows
```

---

## 💡 Pro Tips

### **Tip 1: Use Multiple Terminals**

**Terminal 1:** Backend & Logs
```bash
make start
make docker-logs
```

**Terminal 2:** Dashboard
```bash
make dashboard
```

**Terminal 3:** Testing
```bash
make test-api
# or custom curl commands
```

---

### **Tip 2: Quick Status Check**

```bash
# One-liner status
make docker-ps && make health && echo "All good!"
```

---

### **Tip 3: Auto-restart on Changes**

```bash
# Use watch mode
make run-watch

# Requires: make install-tools
```

---

### **Tip 4: Quick Cleanup**

```bash
# Stop and clean in one go
make docker-down && make docker-clean
```

---

### **Tip 5: Bookmark Commands**

Create aliases (in `.bashrc` or `.zshrc`):
```bash
alias rmq-start='cd ~/golang-rabbitmq-v2 && make start'
alias rmq-dash='cd ~/golang-rabbitmq-v2 && make dashboard'
alias rmq-stop='cd ~/golang-rabbitmq-v2 && make docker-down'
```

---

## 🎓 Learning Path

### **Beginner** (First time)

1. Run `make start`
2. Wait for completion
3. Run `make dashboard`
4. Explore UI

---

### **Intermediate** (Daily use)

1. `make start` in morning
2. `make dashboard` for monitoring
3. `make docker-logs` for debugging
4. `make docker-down` in evening

---

### **Advanced** (Development)

1. `make run-watch` for live reload
2. `make test` before commit
3. `make docker-rebuild` after changes
4. `make docker-logs` for debugging

---

## 📊 Command Cheat Sheet

| Task | Command |
|------|---------|
| **Quick Start** | `make start` |
| **Open Dashboard** | `make dashboard` |
| **Check Health** | `make health` |
| **View Logs** | `make docker-logs` |
| **Test API** | `make test-api` |
| **Stop All** | `make docker-down` |
| **Restart** | `make docker-rebuild` |
| **Show Help** | `make help` |

---

## ✅ Success Checklist

After `make start`, verify:

- [ ] 3 containers running (`make docker-ps`)
- [ ] Health check passes (`make health`)
- [ ] API test successful (`make test-api`)
- [ ] Dashboard opens (`make dashboard`)
- [ ] Dashboard shows "Connected"
- [ ] Can create messages
- [ ] Events appear in dashboard

**All checked?** ✅ Perfect!

---

## 🎉 Ready to Go!

Now you can:

```bash
# Start everything
make start

# Monitor in real-time
make dashboard

# That's it! 🚀
```

**Happy Monitoring!**

---

**Quick Reference:**
- Start: `make start`
- Dashboard: `make dashboard`
- Stop: `make docker-down`
- Help: `make help`
