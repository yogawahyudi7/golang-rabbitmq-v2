# RabbitMQ Monitoring API Documentation

Dokumentasi lengkap untuk monitoring pub/sub RabbitMQ yang sedang berlangsung beserta riwayatnya.

## 📊 Overview

Sistem monitoring ini menyediakan:
- **Real-time Activity Monitoring** - Aktivitas pub/sub yang sedang berlangsung
- **Event History Tracking** - Riwayat lengkap setiap event publish/consume
- **Detailed Metrics** - Statistik performa dan kesehatan sistem
- **Event Filtering** - Filter event berdasarkan tipe dan limit

## 🔌 API Endpoints

### 1. Health Check

Mengecek kesehatan aplikasi dan koneksi RabbitMQ.

**Endpoint:** `GET /api/v1/monitoring/health`

**Response:**
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

### 2. Basic Metrics

Mendapatkan summary metrics RabbitMQ.

**Endpoint:** `GET /api/v1/monitoring/metrics`

**Response:**
```json
{
  "success": true,
  "data": {
    "messages": {
      "published": 150,
      "consumed": 148,
      "failed": 2,
      "requeued": 1
    },
    "confirmations": {
      "acked": 148,
      "nacked": 2,
      "timeout": 0
    },
    "connections": {
      "created": 1,
      "reconnect_attempts": 0,
      "errors": 0
    },
    "performance": {
      "avg_processing_time_ms": 45,
      "publish_rate_per_sec": 5.2,
      "consume_rate_per_sec": 5.0,
      "uptime_seconds": 5025.5
    }
  }
}
```

---

### 3. Detailed Metrics

Mendapatkan metrics detail dengan success rate dan ack rate.

**Endpoint:** `GET /api/v1/monitoring/metrics/detailed`

**Response:**
```json
{
  "success": true,
  "data": {
    "timestamp": "1h23m45s",
    "uptime_seconds": 5025.5,
    "messages": {
      "published": 150,
      "consumed": 148,
      "failed": 2,
      "requeued": 1,
      "success_rate": 98.67
    },
    "confirmations": {
      "acked": 148,
      "nacked": 2,
      "timeout": 0,
      "ack_rate": 98.67
    },
    "connections": {
      "created": 1,
      "reconnect_attempts": 0,
      "errors": 0
    },
    "performance": {
      "avg_processing_time_ms": 45,
      "publish_rate_per_sec": 5.2,
      "consume_rate_per_sec": 5.0
    },
    "activity": {
      "active_publishers": 0,
      "active_consumers": 0,
      "publishing_in_progress": 0,
      "consuming_in_progress": 0
    }
  }
}
```

---

### 4. Real-time Activity Snapshot 🔥

Mendapatkan snapshot aktivitas yang sedang berlangsung dengan recent events.

**Endpoint:** `GET /api/v1/monitoring/activity?limit=50`

**Query Parameters:**
- `limit` (optional) - Jumlah recent events yang ditampilkan (default: 50)

**Response:**
```json
{
  "success": true,
  "data": {
    "timestamp": "2025-07-29T10:15:30Z",
    "active_publishers": 2,
    "active_consumers": 5,
    "publishing_in_progress": 3,
    "consuming_in_progress": 8,
    "current_publish_rate": 5.2,
    "current_consume_rate": 5.0,
    "recent_events": [
      {
        "id": "evt_1722252930_145",
        "type": "consume",
        "timestamp": "2025-07-29T10:15:29Z",
        "routing_key": "app.message",
        "delivery_tag": 145,
        "duration_ms": 42,
        "success": true,
        "metadata": {
          "duration_ms": 42
        }
      },
      {
        "id": "evt_1722252929_144",
        "type": "publish",
        "timestamp": "2025-07-29T10:15:28Z",
        "routing_key": "app.message",
        "duration_ms": 5,
        "success": true,
        "metadata": {
          "duration_ms": 5
        }
      }
    ]
  }
}
```

**Use Case:**
- Monitoring real-time untuk melihat aktivitas yang sedang terjadi
- Dashboard monitoring dengan auto-refresh
- Debugging masalah performa

---

### 5. Event History 📜

Mendapatkan riwayat event dengan filtering.

**Endpoint:** `GET /api/v1/monitoring/events?type=publish&limit=100`

**Query Parameters:**
- `type` (optional) - Filter by event type: `publish`, `consume`, `confirm_ack`, `confirm_nack`, `confirm_timeout`, `failed`, `requeued`
- `limit` (optional) - Jumlah events yang ditampilkan (default: 100)

**Response:**
```json
{
  "success": true,
  "data": {
    "events": [
      {
        "id": "evt_1722252930_145",
        "type": "publish",
        "timestamp": "2025-07-29T10:15:29Z",
        "routing_key": "app.message",
        "duration_ms": 5,
        "success": true,
        "metadata": {
          "duration_ms": 5
        }
      }
    ],
    "count": 1,
    "filters": {
      "type": "publish",
      "limit": 100
    }
  }
}
```

**Use Case:**
- Audit trail untuk semua operasi pub/sub
- Investigasi error atau masalah
- Analisis performa historis

---

### 6. Publish Events History

Mendapatkan riwayat khusus publish events.

**Endpoint:** `GET /api/v1/monitoring/events/publish?limit=50`

**Query Parameters:**
- `limit` (optional) - Jumlah events (default: 50)

**Response:**
```json
{
  "success": true,
  "data": {
    "events": [
      {
        "id": "evt_1722252930_145",
        "type": "publish",
        "timestamp": "2025-07-29T10:15:29Z",
        "routing_key": "app.message",
        "duration_ms": 5,
        "success": true
      }
    ],
    "count": 1
  }
}
```

---

### 7. Consume Events History

Mendapatkan riwayat khusus consume events.

**Endpoint:** `GET /api/v1/monitoring/events/consume?limit=50`

**Query Parameters:**
- `limit` (optional) - Jumlah events (default: 50)

**Response:**
```json
{
  "success": true,
  "data": {
    "events": [
      {
        "id": "evt_1722252929_144",
        "type": "consume",
        "timestamp": "2025-07-29T10:15:28Z",
        "routing_key": "app.message",
        "delivery_tag": 144,
        "duration_ms": 42,
        "success": true
      }
    ],
    "count": 1
  }
}
```

---

### 8. Failed Events History

Mendapatkan riwayat khusus failed events untuk debugging.

**Endpoint:** `GET /api/v1/monitoring/events/failed?limit=50`

**Query Parameters:**
- `limit` (optional) - Jumlah events (default: 50)

**Response:**
```json
{
  "success": true,
  "data": {
    "events": [
      {
        "id": "evt_1722252925_140",
        "type": "failed",
        "timestamp": "2025-07-29T10:15:24Z",
        "routing_key": "app.message",
        "delivery_tag": 140,
        "duration_ms": 30,
        "success": false,
        "error": "database connection timeout",
        "metadata": {
          "duration_ms": 30
        }
      }
    ],
    "count": 1
  }
}
```

**Use Case:**
- Debugging error secara detail
- Monitoring error rate
- Root cause analysis

---

### 9. Clear Event History

Menghapus semua event history (untuk maintenance).

**Endpoint:** `DELETE /api/v1/monitoring/events`

**Response:**
```json
{
  "success": true,
  "message": "Event history cleared successfully"
}
```

**Use Case:**
- Maintenance dan cleanup
- Reset monitoring data
- Testing

---

## 📈 Event Types

| Type | Description |
|------|-------------|
| `publish` | Event saat message dipublish |
| `consume` | Event saat message dikonsumsi |
| `confirm_ack` | Publisher confirmation (ack) |
| `confirm_nack` | Publisher confirmation (nack) |
| `confirm_timeout` | Publisher confirmation timeout |
| `failed` | Event saat consume gagal |
| `requeued` | Event saat message di-requeue |

## 🎯 Use Cases

### 1. Real-time Dashboard Monitoring

Gunakan endpoint `/activity` dengan auto-refresh untuk monitoring real-time:

```bash
# Update setiap 5 detik
watch -n 5 curl -s http://localhost:8080/api/v1/monitoring/activity?limit=20
```

### 2. Debugging Failed Messages

Lihat detail error dari failed events:

```bash
curl http://localhost:8080/api/v1/monitoring/events/failed?limit=10
```

### 3. Performance Analysis

Analisis performa dengan detailed metrics:

```bash
curl http://localhost:8080/api/v1/monitoring/metrics/detailed
```

### 4. Audit Trail

Track semua operasi pub/sub:

```bash
# Semua events
curl http://localhost:8080/api/v1/monitoring/events?limit=100

# Hanya publish events
curl http://localhost:8080/api/v1/monitoring/events/publish?limit=50

# Hanya consume events
curl http://localhost:8080/api/v1/monitoring/events/consume?limit=50
```

## 🔧 Configuration

Event history disimpan dalam memory dengan limit:
- **Max Events**: 1000 events (FIFO)
- **Thread-safe**: Menggunakan RWMutex untuk concurrent access
- **Automatic Cleanup**: Old events otomatis terhapus saat limit tercapai

## 📊 Metrics Explanation

### Success Rate
```
Success Rate = (Consumed / (Consumed + Failed)) * 100
```

### Ack Rate
```
Ack Rate = (Acked / (Acked + Nacked)) * 100
```

### Processing Time
- Dihitung dari start hingga completion (success/fail)
- Average dihitung dari total processing time / total messages

### Rates
- Dihitung per detik berdasarkan interval sampling
- Updated setiap detik

## 🚀 Quick Start Examples

### Check System Health
```bash
curl http://localhost:8080/api/v1/monitoring/health
```

### View Current Activity
```bash
curl http://localhost:8080/api/v1/monitoring/activity
```

### View Last 50 Publish Events
```bash
curl http://localhost:8080/api/v1/monitoring/events/publish?limit=50
```

### View Last 50 Failed Events
```bash
curl http://localhost:8080/api/v1/monitoring/events/failed?limit=50
```

### Get Detailed Metrics
```bash
curl http://localhost:8080/api/v1/monitoring/metrics/detailed
```

## 🎨 Integration Examples

### JavaScript/React Dashboard
```javascript
// Real-time activity monitoring
const fetchActivity = async () => {
  const response = await fetch('http://localhost:8080/api/v1/monitoring/activity?limit=20');
  const data = await response.json();
  return data.data;
};

// Auto-refresh every 5 seconds
setInterval(async () => {
  const activity = await fetchActivity();
  updateDashboard(activity);
}, 5000);
```

### Python Monitoring Script
```python
import requests
import time

def monitor_activity():
    while True:
        response = requests.get('http://localhost:8080/api/v1/monitoring/activity?limit=20')
        data = response.json()['data']

        print(f"Active Publishers: {data['active_publishers']}")
        print(f"Active Consumers: {data['active_consumers']}")
        print(f"Publish Rate: {data['current_publish_rate']}/sec")
        print(f"Consume Rate: {data['current_consume_rate']}/sec")
        print(f"Recent Events: {len(data['recent_events'])}")
        print("-" * 50)

        time.sleep(5)

monitor_activity()
```

## 🔐 Security Notes

- Endpoint monitoring sebaiknya di-protect dengan authentication
- Rate limit untuk mencegah abuse
- Consider adding RBAC untuk production

## 📝 Notes

- Event history disimpan di memory, akan hilang saat restart
- Untuk persistent history, consider menggunakan database atau external monitoring tools
- Max 1000 events disimpan (configurable di `metrics.go`)

---

**Last Updated:** 2025-07-29
**Version:** 2.0
