# Testing Metrics - Step by Step

## Test 1: Prometheus Queries (Instant Test)

Go to: **http://localhost:9090/graph**

### Query 1: Container CPU Usage
```promql
rate(container_cpu_usage_seconds_total[5m]) * 100
```
**Expected**: CPU usage for all containers (blog-service, api-gateway, etc.)

### Query 2: Container Memory
```promql
container_memory_usage_bytes / 1024 / 1024
```
**Expected**: Memory usage in MB for each container

### Query 3: Network Traffic
```promql
rate(container_network_receive_bytes_total[5m])
```
**Expected**: Network receive bytes per second

**Click "Execute"** - if you see data, it's working! ✅

---

## Test 2: cAdvisor Web UI

Go to: **http://localhost:9091**

- You'll see all your running containers
- Click on any container (e.g., "blog-service")
- See real-time:
  - CPU usage
  - Memory usage
  - Network stats
  - File system stats

**If you can see this data, metrics are working! ✅**

---

## Test 3: Grafana Visualization

### Step 1: Login to Grafana
1. Go to: **http://localhost:3000**
2. Login: `admin` / `admin`

### Step 2: Add Prometheus Data Source
1. Click **Configuration** (gear icon) → **Data Sources**
2. Click **Add data source**
3. Select **Prometheus**
4. **URL**: `http://prometheus:9090`
5. Click **Save & Test**
6. You should see: ✅ "Data source is working"

### Step 3: Create a Test Dashboard
1. Click **Dashboards** → **New** → **New Dashboard**
2. Click **Add visualization**
3. In the query editor, paste:
   ```
   container_memory_usage_bytes
   ```
4. Click **Run query**
5. You should see memory data for containers! ✅

### Step 4: Create More Panels
**Panel 2 - CPU:**
```promql
rate(container_cpu_usage_seconds_total[5m]) * 100
```

**Panel 3 - Network:**
```promql
rate(container_network_receive_bytes_total[5m])
```

**Panel 4 - Disk:**
```promql
container_fs_usage_bytes
```

---

## Quick Verification Commands

Open PowerShell and run:

```powershell
# Test Prometheus is running
Invoke-WebRequest http://localhost:9090/api/v1/query?query=up

# Test cAdvisor is running  
Invoke-WebRequest http://localhost:9091/metrics

# Test Grafana is running
Invoke-WebRequest http://localhost:3000
```

---

## What Success Looks Like

✅ Prometheus shows: "Successfully queried data"
✅ cAdvisor shows: Real-time container stats
✅ Grafana shows: Connected data source
✅ Queries return: Actual numbers (not errors)

**If all 4 are ✅, your metrics are working perfectly!**

