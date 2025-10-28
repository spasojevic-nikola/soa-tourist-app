# Metrics Monitoring Setup

This project implements comprehensive metrics monitoring using Prometheus, node-exporter, cAdvisor, and Grafana.

## Components

### 1. Prometheus (`http://localhost:9090`)
- **Purpose**: Metrics collection and time-series database
- **Port**: 9090
- **Scrapes**:
  - node-exporter (OS metrics)
  - cAdvisor (container metrics)
  - Application services

### 2. node-exporter (Port 9100)
- **Purpose**: Operating system metrics
- **Metrics Collected**:
  - CPU usage
  - RAM/memory usage
  - File system usage
  - Network traffic flow
  - Disk I/O
  - System load

### 3. cAdvisor (`http://localhost:9091`)
- **Purpose**: Container metrics
- **Metrics Collected**:
  - Container CPU usage
  - Container RAM/memory usage
  - Container file system
  - Container network traffic
  - Container I/O statistics

### 4. Grafana (`http://localhost:3000`)
- **Purpose**: Metrics visualization
- **Login**: admin / admin
- **Data Source**: Prometheus (auto-configured)

## How to Use

### 1. Start All Services

```bash
docker-compose up -d
```

### 2. Access Dashboards

#### Prometheus UI
```
http://localhost:9090
```

Features:
- Query metrics using PromQL
- View targets being scraped
- Check metrics status

#### Grafana UI
```
http://localhost:3000
```

Login with:
- Username: `admin`
- Password: `admin`

### 3. Explore Metrics

#### Operating System Metrics (from node-exporter)

**CPU Usage**:
```promql
100 - (avg by (instance) (rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)
```

**Memory Usage**:
```promql
node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes
```

**Disk Usage**:
```promql
node_filesystem_avail_bytes / node_filesystem_size_bytes * 100
```

**Network Traffic**:
```promql
rate(node_network_receive_bytes_total[5m])
```

#### Container Metrics (from cAdvisor)

**Container CPU Usage**:
```promql
rate(container_cpu_usage_seconds_total{container!="POD"}[5m])
```

**Container Memory Usage**:
```promql
container_memory_usage_bytes
```

**Container Network**:
```promql
rate(container_network_receive_bytes_total[5m])
```

### 4. Import Pre-built Dashboards (Optional)

In Grafana:
1. Go to Dashboards → Import
2. Search for dashboards:
   - **Node Exporter Full**: 1860
   - **Docker Monitoring**: 179
   - **cAdvisor**: 14282

These provide ready-made visualizations for all the required metrics.

## Required Metrics Coverage

### ✅ Operating System Metrics
- ✅ **CPU**: node_exporter `node_cpu_seconds_total`
- ✅ **RAM**: node_exporter `node_memory_*`
- ✅ **File System**: node_exporter `node_filesystem_*`
- ✅ **Network**: node_exporter `node_network_*`

### ✅ Container Metrics
- ✅ **CPU**: cAdvisor `container_cpu_usage_seconds_total`
- ✅ **RAM**: cAdvisor `container_memory_usage_bytes`
- ✅ **File System**: cAdvisor `container_fs_*`
- ✅ **Network**: cAdvisor `container_network_*`

## Service Endpoints

| Service | Port | Endpoint |
|---------|------|----------|
| Prometheus | 9090 | http://localhost:9090 |
| Grafana | 3000 | http://localhost:3000 |
| cAdvisor | 9091 | http://localhost:9091/metrics |
| node-exporter | 9100 | http://localhost:9100/metrics |

## Troubleshooting

### Prometheus not scraping targets?

Check targets status:
```
http://localhost:9090/targets
```

All targets should show "UP" status.

### Node-exporter not working?

node-exporter uses `network_mode: host` and accesses system files. Make sure it has proper permissions.

### cAdvisor not showing container metrics?

cAdvisor needs privileged access. Check if it's running:
```bash
docker ps | grep cadvisor
```

### Missing metrics in Grafana?

1. Check Prometheus has data: http://localhost:9090/graph
2. Verify Prometheus is added as data source in Grafana
3. Check time range in Grafana (last 5-15 minutes)

## Advanced Configuration

### Custom Prometheus Rules

Create `alerts.yml` for alerting:
```yaml
groups:
  - name: instance
    rules:
      - alert: HighCPUUsage
        expr: 100 - avg(irate(node_cpu_seconds_total{mode="idle"}[5m])) by (instance) * 100 > 80
        for: 5m
```

### Grafana Dashboards

Create custom dashboards by:
1. Going to Dashboards → New Dashboard
2. Adding panels with PromQL queries
3. Configuring visualization (time series, gauges, etc.)

## Summary

This setup provides:
- ✅ OS metrics (CPU, RAM, filesystem, network)
- ✅ Container metrics (CPU, RAM, filesystem, network)
- ✅ Centralized metrics collection (Prometheus)
- ✅ Beautiful visualizations (Grafana)
- ✅ Historical data retention
- ✅ Query capabilities

All required metrics are now collected and visualized!
