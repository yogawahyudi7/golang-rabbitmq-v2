# 📑 Documentation Index

Quick navigation untuk semua dokumentasi project ini.

---

## 🚀 Quick Start

**Baru mulai?** Baca dalam urutan ini:

1. **[GETTING_STARTED.md](GETTING_STARTED.md)** ⭐ START HERE
   - Step-by-step setup dari awal
   - Prerequisites & installation
   - Multiple setup scenarios
   - Troubleshooting startup issues

2. **[DASHBOARD_GUIDE.md](DASHBOARD_GUIDE.md)**
   - Cara menggunakan dashboard
   - Fitur-fitur dashboard
   - Customization guide
   - Tips & tricks

3. **[MONITORING_API.md](MONITORING_API.md)**
   - Complete API reference
   - 9 monitoring endpoints
   - Query parameters & examples
   - Integration examples

4. **[COMPLETE_SETUP_SUMMARY.md](COMPLETE_SETUP_SUMMARY.md)**
   - Summary lengkap semua yang ada
   - Quick command reference
   - Learning path
   - Next steps

---

## 📚 Documentation by Purpose

### **For New Users** 🆕
- Start: [GETTING_STARTED.md](GETTING_STARTED.md)
- Quick Start in README: [README.md](README.md#-quick-start)
- Summary: [COMPLETE_SETUP_SUMMARY.md](COMPLETE_SETUP_SUMMARY.md)

### **For Dashboard Users** 🎨
- Dashboard Guide: [DASHBOARD_GUIDE.md](DASHBOARD_GUIDE.md)
- Technical Docs: [web/README.md](web/README.md)
- Quick Open: `make dashboard`

### **For API Users** 🔌
- API Documentation: [MONITORING_API.md](MONITORING_API.md)
- Test Endpoints: `make test-api`
- Health Check: `make health`

### **For Developers** 👨‍💻
- Project Overview: [README.md](README.md)
- Project Summary: [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)
- Architecture: See README.md#architecture
- Code Structure: See README.md#project-structure

---

## 📖 Documentation by Topic

### **Setup & Installation**
| Document | Purpose |
|----------|---------|
| [GETTING_STARTED.md](GETTING_STARTED.md) | Complete setup guide |
| [README.md](README.md) | Project overview & quick start |
| [COMPLETE_SETUP_SUMMARY.md](COMPLETE_SETUP_SUMMARY.md) | What's been built & how to use |

### **Dashboard**
| Document | Purpose |
|----------|---------|
| [DASHBOARD_GUIDE.md](DASHBOARD_GUIDE.md) | Dashboard usage guide |
| [web/README.md](web/README.md) | Technical documentation |
| Dashboard HTML: `web/dashboard.html` | The actual dashboard |

### **API & Monitoring**
| Document | Purpose |
|----------|---------|
| [MONITORING_API.md](MONITORING_API.md) | Complete API reference |
| Monitoring Handler: `internal/handler/monitoring_handler.go` | Code implementation |

### **Project Info**
| Document | Purpose |
|----------|---------|
| [README.md](README.md) | Main project documentation |
| [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) | Discussion & decisions |
| [Makefile](Makefile) | All available commands |

---

## 🎯 Common Tasks - Quick Reference

### **Starting the System**

**Super Quick (1 command!):**
```bash
make start          # Starts everything!
make dashboard      # Open monitoring
```

**Or Manual Steps:**
```bash
# 1. Start backend
make docker-up

# 2. Open dashboard
make dashboard

# 3. Generate test traffic
make test-api
```

**Docs:**
- [RUNNING_WITH_MAKE.md](RUNNING_WITH_MAKE.md) - Make commands guide
- [GETTING_STARTED.md](GETTING_STARTED.md) - Detailed setup

---

### **Monitoring**
```bash
# Dashboard (visual)
make dashboard

# RabbitMQ Management UI
make rabbitmq-ui

# API (command line)
make health
make metrics
```

**Docs:** [DASHBOARD_GUIDE.md](DASHBOARD_GUIDE.md) | [MONITORING_API.md](MONITORING_API.md)

---

### **Development**
```bash
# Run locally
make run

# Run with live reload
make run-watch

# Build
make build

# Test
make test
```

**Docs:** [README.md](README.md#development-commands)

---

### **Troubleshooting**
```bash
# Check status
make docker-ps

# View logs
make docker-logs

# Health check
make health

# Restart
make docker-rebuild
```

**Docs:** [GETTING_STARTED.md](GETTING_STARTED.md#troubleshooting-startup-issues)

---

## 🗂️ File Structure

```
📁 golang-rabbitmq-v2/
│
├── 📄 INDEX.md                          # This file (navigation)
├── 📄 README.md                         # Main documentation
├── 📄 GETTING_STARTED.md                # Setup guide ⭐
├── 📄 DASHBOARD_GUIDE.md                # Dashboard usage
├── 📄 MONITORING_API.md                 # API reference
├── 📄 COMPLETE_SETUP_SUMMARY.md         # Complete summary
├── 📄 PROJECT_SUMMARY.md                # Project discussion
│
├── 📁 web/                              # Dashboard files
│   ├── 📄 dashboard.html                # Dashboard UI
│   ├── 📄 dashboard.js                  # Dashboard logic
│   ├── 📄 README.md                     # Technical docs
│   ├── 📄 open-dashboard.bat            # Windows launcher
│   └── 📄 open-dashboard.sh             # Linux/Mac launcher
│
├── 📁 cmd/                              # Application code
├── 📁 internal/                         # Internal packages
├── 📁 pkg/                              # Shared packages
├── 📄 Makefile                          # Commands
├── 📄 docker-compose.yml                # Docker config
└── 📄 Dockerfile                        # Container config
```

---

## 🎓 Learning Path

### **Beginner** (Never used this before)
1. Read [GETTING_STARTED.md](GETTING_STARTED.md) - Quick start section
2. Run `make docker-up`
3. Run `make dashboard`
4. Explore the dashboard!

### **Intermediate** (Want to understand features)
1. Read [DASHBOARD_GUIDE.md](DASHBOARD_GUIDE.md) - Full guide
2. Read [MONITORING_API.md](MONITORING_API.md) - API endpoints
3. Try filtering events
4. Test API with curl

### **Advanced** (Want to customize/integrate)
1. Read [web/README.md](web/README.md) - Technical details
2. Customize dashboard (colors, rates)
3. Build custom monitoring tools
4. Integrate with external systems

---

## 💡 Best Practices

### **Daily Development Workflow**
```bash
# Morning:
make docker-up          # Start services
make dashboard          # Open dashboard

# During development:
make test-api           # Test changes
make docker-logs        # Check logs

# Evening:
make docker-down        # Stop services
```

### **Debugging Issues**
```bash
# 1. Check services
make docker-ps

# 2. Check health
make health

# 3. Check dashboard
make dashboard
# Click "Failed" filter

# 4. Check logs
make docker-logs
```

### **Testing New Features**
```bash
# Terminal 1: Monitor
make dashboard

# Terminal 2: Generate traffic
make test-api
# or your custom tests

# Watch real-time updates!
```

---

## 🔗 External Resources

### **RabbitMQ**
- Management UI: http://localhost:15672 (guest/guest)
- Official Docs: https://www.rabbitmq.com/documentation.html

### **Services**
- Application API: http://localhost:8080
- Monitoring Dashboard: http://localhost:3000/dashboard.html
- PostgreSQL: localhost:5432

---

## 🆘 Need Help?

### **Quick Checks**
```bash
make help           # Show all commands
make info           # Project information
make health         # Check system health
make docker-logs    # View logs
```

### **Common Issues**
- **Can't start?** → [GETTING_STARTED.md](GETTING_STARTED.md#troubleshooting-startup-issues)
- **Dashboard not working?** → [DASHBOARD_GUIDE.md](DASHBOARD_GUIDE.md#troubleshooting)
- **API errors?** → [MONITORING_API.md](MONITORING_API.md#troubleshooting)

### **Documentation Issues**
If dokumentasi tidak jelas atau ada yang kurang:
1. Check semua 4 main docs
2. Run `make help` untuk commands
3. Check logs dengan `make docker-logs`

---

## 📊 What's Available

### **Web UI** 🎨
- ✅ Monitoring Dashboard (real-time)
- ✅ RabbitMQ Management UI (built-in)

### **API Endpoints** 🔌
- ✅ 9 monitoring endpoints
- ✅ User management endpoints
- ✅ Message management endpoints
- ✅ Health & metrics endpoints

### **Documentation** 📚
- ✅ 7 documentation files
- ✅ API examples
- ✅ Integration guides
- ✅ Troubleshooting guides

### **Tools** 🛠️
- ✅ Makefile with 40+ commands
- ✅ Docker Compose setup
- ✅ Launcher scripts
- ✅ Development tools

---

## ⭐ Quick Tips

1. **Always start with** `GETTING_STARTED.md` jika pertama kali
2. **Bookmark** `make dashboard` - paling sering dipakai
3. **Check** `COMPLETE_SETUP_SUMMARY.md` untuk quick reference
4. **Use** `make help` untuk lihat semua commands
5. **Monitor with** dashboard - lebih mudah dari curl

---

## 🚀 Start Here

**Belum tahu mau mulai dari mana?**

```bash
# Just run these 3 commands:
make docker-up      # 1. Start everything
make dashboard      # 2. Open monitoring
make test-api       # 3. Test it works

# Then read:
# - GETTING_STARTED.md for details
# - DASHBOARD_GUIDE.md for features
# - MONITORING_API.md for API
```

**That's it!** Enjoy! 🎉

---

**Last Updated:** 2025-07-29
**Version:** 2.0
