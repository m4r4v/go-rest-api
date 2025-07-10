# Google Cloud Run Deployment Guide - Visual Portal Method

## Prerequisites Checklist

Before starting, ensure you have:

- [ ] Google Cloud Project created and billing enabled
- [ ] Cloud Run API enabled
- [ ] Cloud Build API enabled  
- [ ] Container Registry API enabled
- [ ] gcloud CLI installed and authenticated

## Step 1: Enable Required APIs

### Via Google Cloud Console:
1. Go to [Google Cloud Console](https://console.cloud.google.com)
2. Navigate to **APIs & Services** > **Library**
3. Search and enable these APIs:
   - Cloud Run API
   - Cloud Build API
   - Container Registry API
   - Artifact Registry API (recommended over Container Registry)

### Via CLI (Alternative):
```bash
gcloud services enable cloudbuild.googleapis.com run.googleapis.com containerregistry.googleapis.com
```

## Step 2: Build Container Image

### Option A: Using Cloud Build (Recommended)
```bash
# From your project directory
gcloud builds submit --config cloudbuild.yaml
```

### Option B: Local Build + Push
```bash
# Build locally
docker build -t gcr.io/YOUR_PROJECT_ID/go-rest-api .

# Configure Docker for GCR
gcloud auth configure-docker

# Push to Container Registry
docker push gcr.io/YOUR_PROJECT_ID/go-rest-api
```

## Step 3: Deploy via Google Cloud Console (Visual Method)

### 3.1 Navigate to Cloud Run
1. Go to [Google Cloud Console](https://console.cloud.google.com)
2. In the navigation menu, click **Cloud Run**
3. Click **CREATE SERVICE**

### 3.2 Configure Container
1. **Container Image URL**: 
   - Click **SELECT** next to "Container image URL"
   - Choose **Container Registry** or **Artifact Registry**
   - Select your project
   - Choose `go-rest-api` repository
   - Select the latest image tag
   - Click **SELECT**

### 3.3 Service Configuration
1. **Service name**: `go-rest-api`
2. **Region**: `us-central1` (or your preferred region)
3. **CPU allocation and pricing**: 
   - Select **CPU is only allocated during request processing**

### 3.4 Advanced Settings - Container Tab
Click **CONTAINER, VARIABLES & SECRETS, CONNECTIONS, SECURITY**

#### Container Settings:
- **Container port**: `8080`
- **Memory**: `128 MiB` (as requested)
- **CPU**: `1`
- **Request timeout**: `300` seconds
- **Maximum concurrent requests per instance**: `80`

#### Variables & Secrets:
- Add environment variable:
  - **Name**: `LOG_LEVEL`
  - **Value**: `info`

### 3.5 Advanced Settings - Autoscaling Tab
- **Minimum number of instances**: `0`
- **Maximum number of instances**: `1`

### 3.6 Advanced Settings - Security Tab
- **Service account**: Use default compute service account
- **Container security**: Keep defaults

### 3.7 Advanced Settings - Connections Tab
- **CPU throttling**: Uncheck "CPU throttling" (disable it)
- **Execution environment**: Select **Second generation**

### 3.8 Traffic Settings
- **Authentication**: Select **Allow unauthenticated invocations**

### 3.9 Deploy
1. Review all settings
2. Click **CREATE**
3. Wait for deployment to complete (usually 2-3 minutes)

## Step 4: Verify Deployment

### 4.1 Check Service Status
1. In Cloud Run console, you should see your service with a green checkmark
2. Note the service URL (something like `https://go-rest-api-xxx-uc.a.run.app`)

### 4.2 Test Endpoints
```bash
# Replace YOUR_SERVICE_URL with your actual Cloud Run URL
curl https://YOUR_SERVICE_URL/health
curl https://YOUR_SERVICE_URL/
curl https://YOUR_SERVICE_URL/v1/status
```

Expected response for `/health`:
```json
{
  "success": true,
  "status_code": 200,
  "status": "OK",
  "data": {
    "service": "go-rest-api",
    "version": "2.0.0",
    "status": "healthy",
    "timestamp": "2025-07-09T21:23:00Z"
  },
  "timestamp": "2025-07-09T21:23:00Z"
}
```

## Step 5: Monitor and Debug

### 5.1 View Logs
1. In Cloud Run console, click on your service
2. Go to **LOGS** tab
3. View real-time logs and any error messages

### 5.2 Check Metrics
1. Go to **METRICS** tab
2. Monitor:
   - Request count
   - Request latency
   - Instance count
   - Memory utilization

### 5.3 Common Issues and Solutions

#### Issue: "Container failed to start"
**Symptoms**: Service shows error status, logs show port binding issues

**Solutions**:
1. Verify container port is set to `8080`
2. Check that application listens on `0.0.0.0:$PORT`
3. Ensure health check endpoint `/health` responds correctly

#### Issue: "Memory limit exceeded"
**Symptoms**: Service restarts frequently, 503 errors

**Solutions**:
1. Increase memory allocation (try 256Mi)
2. Optimize application memory usage
3. Check for memory leaks in logs

#### Issue: "Request timeout"
**Symptoms**: 504 Gateway Timeout errors

**Solutions**:
1. Increase request timeout in service settings
2. Optimize application response time
3. Check for blocking operations in code

## Step 6: Update Deployment

### 6.1 Deploy New Revision
1. Build new container image with updated code
2. In Cloud Run console, click **EDIT & DEPLOY NEW REVISION**
3. Update container image URL to new version
4. Adjust settings if needed
5. Click **DEPLOY**

### 6.2 Traffic Management
1. Go to **MANAGE TRAFFIC** tab
2. Allocate traffic between revisions
3. Use for blue-green deployments or gradual rollouts

## Configuration Summary

Your final Cloud Run service configuration:

```yaml
Service: go-rest-api
Region: us-central1
Memory: 128Mi
CPU: 1
Min instances: 0
Max instances: 1
Concurrency: 80
Port: 8080
Timeout: 300s
Authentication: Allow unauthenticated
CPU throttling: Disabled
Execution environment: Second generation
```

## Troubleshooting Commands

```bash
# View service details
gcloud run services describe go-rest-api --region=us-central1

# View logs
gcloud run services logs tail go-rest-api --region=us-central1

# Update service settings
gcloud run services update go-rest-api \
  --region=us-central1 \
  --memory=256Mi \
  --max-instances=2

# Delete service (if needed)
gcloud run services delete go-rest-api --region=us-central1
```

## Next Steps

1. **Custom Domain**: Set up custom domain mapping
2. **Authentication**: Implement IAM-based authentication
3. **Monitoring**: Set up Cloud Monitoring alerts
4. **CI/CD**: Automate deployments with Cloud Build triggers
5. **Load Testing**: Test with realistic traffic patterns

---

**Note**: The visual deployment method through Google Cloud Console provides the same functionality as CLI deployment but with a user-friendly interface. Both methods result in identical Cloud Run services.
