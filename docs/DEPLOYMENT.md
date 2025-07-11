# Deployment Guide

## Overview

This guide provides comprehensive instructions for deploying the Go REST API Framework v2.0 in various environments, from local development to production cloud deployments.

## Local Development Deployment

### Prerequisites

- Go 1.23 or higher
- Docker and Docker Compose (optional)
- Git

### Quick Start

1. **Clone and Setup**
   ```bash
   git clone https://github.com/m4r4v/go-rest-api.git
   cd go-rest-api
   go mod download
   ```

2. **Environment Configuration**
   ```bash
   cp .env.example .env
   # Edit .env with development settings
   ```

3. **Build and Run**
   ```bash
   go build -o server ./cmd/server
   ./server
   ```

### Docker Development

1. **Using Docker Compose**
   ```bash
   docker-compose up --build
   ```

2. **Manual Docker Build**
   ```bash
   docker build -t go-rest-api:dev .
   docker run -p 8080:8080 --env-file .env go-rest-api:dev
   ```

## Production Deployment

### Environment Configuration

#### Required Environment Variables

```bash
# Security (REQUIRED)
JWT_SECRET=your-256-bit-secret-key-change-this-in-production
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://app.yourdomain.com

# Server Configuration
PORT=8080
SERVER_HOST=0.0.0.0
SERVER_READ_TIMEOUT=60s
SERVER_WRITE_TIMEOUT=60s
SERVER_IDLE_TIMEOUT=60s

# Authentication
JWT_EXPIRATION=24h
BCRYPT_COST=14

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

#### Optional Environment Variables

```bash
# Database (for future external DB integration)
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=your_username
DB_PASSWORD=your_password
DB_DATABASE=your_database
DB_SSLMODE=require

# Monitoring
METRICS_ENABLED=true
HEALTH_CHECK_INTERVAL=30s
```

### Security Configuration

#### JWT Secret Generation

Generate a strong JWT secret:

```bash
# Using OpenSSL
openssl rand -base64 32

# Using Go
go run -c 'package main; import ("crypto/rand"; "encoding/base64"; "fmt"); func main() { b := make([]byte, 32); rand.Read(b); fmt.Println(base64.StdEncoding.EncodeToString(b)) }'

# Using Python
python3 -c "import secrets; print(secrets.token_urlsafe(32))"
```

#### CORS Configuration

**Development**:
```bash
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080,http://127.0.0.1:3000
```

**Production**:
```bash
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://app.yourdomain.com,https://api.yourdomain.com
```

## Docker Production Deployment

### Dockerfile Optimization

The included Dockerfile is production-ready with:

- Multi-stage build for minimal image size
- Non-root user for security
- Alpine Linux base for reduced attack surface
- Health check support

### Building Production Image

```bash
# Build production image
docker build -t go-rest-api:production .

# Tag for registry
docker tag go-rest-api:production your-registry.com/go-rest-api:latest
docker tag go-rest-api:production your-registry.com/go-rest-api:v2.0.0

# Push to registry
docker push your-registry.com/go-rest-api:latest
docker push your-registry.com/go-rest-api:v2.0.0
```

### Docker Run Configuration

```bash
docker run -d \
  --name go-rest-api \
  --restart unless-stopped \
  -p 8080:8080 \
  -e JWT_SECRET="your-production-secret" \
  -e CORS_ALLOWED_ORIGINS="https://yourdomain.com" \
  -e LOG_LEVEL="info" \
  --memory="128m" \
  --cpus="0.5" \
  your-registry.com/go-rest-api:latest
```

### Docker Compose Production

```yaml
version: '3.8'

services:
  api:
    image: your-registry.com/go-rest-api:latest
    container_name: go-rest-api
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      - JWT_SECRET=${JWT_SECRET}
      - CORS_ALLOWED_ORIGINS=${CORS_ALLOWED_ORIGINS}
      - LOG_LEVEL=info
      - LOG_FORMAT=json
    deploy:
      resources:
        limits:
          memory: 128M
          cpus: '0.5'
        reservations:
          memory: 64M
          cpus: '0.25'
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  # Optional: Add reverse proxy
  nginx:
    image: nginx:alpine
    container_name: nginx-proxy
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/nginx/ssl:ro
    depends_on:
      - api
```

## Google Cloud Run Deployment

### Automated Deployment with Cloud Build

1. **Setup Cloud Build Trigger**
   ```bash
   # Enable required APIs
   gcloud services enable cloudbuild.googleapis.com
   gcloud services enable run.googleapis.com
   gcloud services enable containerregistry.googleapis.com
   ```

2. **Deploy using Cloud Build**
   ```bash
   # Submit build using included configuration
   gcloud builds submit --config cloudbuild.yaml
   ```

3. **Set Environment Variables**
   ```bash
   gcloud run services update go-rest-api \
     --region=us-central1 \
     --set-env-vars="JWT_SECRET=your-secret,CORS_ALLOWED_ORIGINS=https://yourdomain.com"
   ```

### Manual Cloud Run Deployment

1. **Build and Push to Container Registry**
   ```bash
   # Set project ID
   export PROJECT_ID=your-project-id
   
   # Build and tag
   docker build -t gcr.io/$PROJECT_ID/go-rest-api:latest .
   
   # Push to registry
   docker push gcr.io/$PROJECT_ID/go-rest-api:latest
   ```

2. **Deploy to Cloud Run**
   ```bash
   gcloud run deploy go-rest-api \
     --image=gcr.io/$PROJECT_ID/go-rest-api:latest \
     --platform=managed \
     --region=us-central1 \
     --allow-unauthenticated \
     --port=8080 \
     --memory=512Mi \
     --cpu=1 \
     --concurrency=80 \
     --min-instances=0 \
     --max-instances=10 \
     --timeout=300 \
     --execution-environment=gen2 \
     --set-env-vars="JWT_SECRET=your-secret,CORS_ALLOWED_ORIGINS=https://yourdomain.com"
   ```

### Cloud Run Configuration

#### Service Configuration

```yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: go-rest-api
  annotations:
    run.googleapis.com/ingress: all
    run.googleapis.com/execution-environment: gen2
spec:
  template:
    metadata:
      annotations:
        autoscaling.knative.dev/minScale: "0"
        autoscaling.knative.dev/maxScale: "10"
        run.googleapis.com/cpu-throttling: "false"
        run.googleapis.com/memory: "512Mi"
        run.googleapis.com/cpu: "1"
    spec:
      containerConcurrency: 80
      timeoutSeconds: 300
      containers:
      - image: gcr.io/PROJECT_ID/go-rest-api:latest
        ports:
        - containerPort: 8080
        env:
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: api-secrets
              key: jwt-secret
        - name: CORS_ALLOWED_ORIGINS
          value: "https://yourdomain.com"
        resources:
          limits:
            memory: "512Mi"
            cpu: "1"
```

#### Secret Management

```bash
# Create secret for JWT
echo -n "your-jwt-secret" | gcloud secrets create jwt-secret --data-file=-

# Grant access to Cloud Run service account
gcloud secrets add-iam-policy-binding jwt-secret \
  --member="serviceAccount:PROJECT_NUMBER-compute@developer.gserviceaccount.com" \
  --role="roles/secretmanager.secretAccessor"
```

## Kubernetes Deployment

### Namespace and Resources

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: go-rest-api
---
apiVersion: v1
kind: Secret
metadata:
  name: api-secrets
  namespace: go-rest-api
type: Opaque
data:
  jwt-secret: <base64-encoded-secret>
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: api-config
  namespace: go-rest-api
data:
  CORS_ALLOWED_ORIGINS: "https://yourdomain.com"
  LOG_LEVEL: "info"
  LOG_FORMAT: "json"
```

### Deployment Configuration

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-rest-api
  namespace: go-rest-api
  labels:
    app: go-rest-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-rest-api
  template:
    metadata:
      labels:
        app: go-rest-api
    spec:
      securityContext:
        runAsNonRoot: true
        runAsUser: 1001
        fsGroup: 1001
      containers:
      - name: api
        image: your-registry.com/go-rest-api:latest
        ports:
        - containerPort: 8080
          name: http
        env:
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: api-secrets
              key: jwt-secret
        envFrom:
        - configMapRef:
            name: api-config
        resources:
          requests:
            memory: "64Mi"
            cpu: "250m"
          limits:
            memory: "128Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /status
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        securityContext:
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: true
          capabilities:
            drop:
            - ALL
```

### Service and Ingress

```yaml
apiVersion: v1
kind: Service
metadata:
  name: go-rest-api-service
  namespace: go-rest-api
spec:
  selector:
    app: go-rest-api
  ports:
  - port: 80
    targetPort: 8080
    name: http
  type: ClusterIP
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: go-rest-api-ingress
  namespace: go-rest-api
  annotations:
    kubernetes.io/ingress.class: "nginx"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  tls:
  - hosts:
    - api.yourdomain.com
    secretName: api-tls
  rules:
  - host: api.yourdomain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: go-rest-api-service
            port:
              number: 80
```

## AWS ECS Deployment

### Task Definition

```json
{
  "family": "go-rest-api",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512",
  "executionRoleArn": "arn:aws:iam::ACCOUNT:role/ecsTaskExecutionRole",
  "taskRoleArn": "arn:aws:iam::ACCOUNT:role/ecsTaskRole",
  "containerDefinitions": [
    {
      "name": "go-rest-api",
      "image": "your-registry.com/go-rest-api:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "CORS_ALLOWED_ORIGINS",
          "value": "https://yourdomain.com"
        },
        {
          "name": "LOG_LEVEL",
          "value": "info"
        }
      ],
      "secrets": [
        {
          "name": "JWT_SECRET",
          "valueFrom": "arn:aws:secretsmanager:region:account:secret:jwt-secret"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/go-rest-api",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "ecs"
        }
      },
      "healthCheck": {
        "command": ["CMD-SHELL", "curl -f http://localhost:8080/health || exit 1"],
        "interval": 30,
        "timeout": 5,
        "retries": 3,
        "startPeriod": 60
      }
    }
  ]
}
```

### Service Configuration

```json
{
  "serviceName": "go-rest-api",
  "cluster": "production",
  "taskDefinition": "go-rest-api:1",
  "desiredCount": 2,
  "launchType": "FARGATE",
  "networkConfiguration": {
    "awsvpcConfiguration": {
      "subnets": ["subnet-12345", "subnet-67890"],
      "securityGroups": ["sg-abcdef"],
      "assignPublicIp": "ENABLED"
    }
  },
  "loadBalancers": [
    {
      "targetGroupArn": "arn:aws:elasticloadbalancing:region:account:targetgroup/api-tg",
      "containerName": "go-rest-api",
      "containerPort": 8080
    }
  ]
}
```

## Monitoring and Observability

### Health Checks

The framework provides built-in health check endpoints:

- `/health`: Container health check
- `/status`: Service readiness check

### Logging Configuration

```bash
# Production logging
LOG_LEVEL=info
LOG_FORMAT=json

# Development logging
LOG_LEVEL=debug
LOG_FORMAT=text
```

### Metrics Integration

#### Prometheus Metrics (Future Enhancement)

```go
// Example metrics endpoint
func metricsHandler(w http.ResponseWriter, r *http.Request) {
    // Expose Prometheus metrics
    promhttp.Handler().ServeHTTP(w, r)
}
```

#### Application Performance Monitoring

Integration points for APM tools:

- **New Relic**: Add agent to main.go
- **DataDog**: Include DD agent as sidecar
- **Jaeger**: Add tracing middleware

### Log Aggregation

The framework outputs structured JSON logs that can be easily integrated with log aggregation systems. All logs include timestamps, request IDs, and structured data for efficient parsing and analysis.

## Scaling Considerations

### Horizontal Scaling

The framework is designed for horizontal scaling:

1. **Stateless Design**: No server-side sessions
2. **Load Balancer Ready**: Standard HTTP interface
3. **Database Ready**: Easy external database integration

### Auto-scaling Configuration

#### Kubernetes HPA

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: go-rest-api-hpa
  namespace: go-rest-api
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: go-rest-api
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

#### Cloud Run Auto-scaling

```bash
gcloud run services update go-rest-api \
  --min-instances=1 \
  --max-instances=100 \
  --concurrency=80 \
  --cpu-throttling=false
```

## Backup and Disaster Recovery

### Data Backup Strategy

Since the framework uses in-memory storage by default:

1. **Configuration Backup**: Environment variables and secrets
2. **Application Backup**: Container images and deployment configs
3. **External Database**: When integrated, implement regular backups

### Disaster Recovery Plan

1. **Infrastructure as Code**: All deployment configs in version control
2. **Multi-Region Deployment**: Deploy across multiple regions
3. **Automated Recovery**: Use CI/CD for rapid redeployment
4. **Monitoring and Alerting**: Early detection of issues

## Troubleshooting

### Common Issues

#### Container Won't Start

```bash
# Check logs
docker logs go-rest-api

# Common issues:
# - Missing JWT_SECRET environment variable
# - Port already in use
# - Insufficient memory allocation
```

#### Health Check Failures

```bash
# Test health endpoint
curl -f http://localhost:8080/health

# Check container resources
docker stats go-rest-api
```

#### Authentication Issues

```bash
# Verify JWT secret is set
echo $JWT_SECRET

# Check token generation
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}'
```

### Performance Tuning

#### Memory Optimization

```bash
# Set Go garbage collection target
GOGC=100

# Limit memory usage
docker run --memory="128m" go-rest-api
```

#### CPU Optimization

```bash
# Set CPU limits
docker run --cpus="0.5" go-rest-api

# Use multiple CPU cores
GOMAXPROCS=2
```

This deployment guide provides comprehensive instructions for deploying the Go REST API Framework in various environments while maintaining security and performance best practices.
