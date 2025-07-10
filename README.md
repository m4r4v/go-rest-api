# Go REST API Framework v2.0 🚀

Una API REST moderna y escalable construida en Go, optimizada para Google Cloud Run con contenedores Docker.

## 🌟 Características

- ✅ **Compatible con Google Cloud Run**: Cumple con todos los estándares oficiales
- 🐳 **Containerizado**: Docker multi-stage para imágenes optimizadas
- 🔧 **Configuración por variables de entorno**: Especialmente variable `PORT` requerida por Cloud Run
- 🛡️ **Seguro**: Usuario no-root, CORS configurado, middleware de seguridad
- 📊 **Health checks**: Endpoints `/health` y `/v1/status` para monitoreo
- 🚀 **Graceful shutdown**: Manejo correcto de señales SIGTERM (requerido por Cloud Run)
- 📝 **Logging estructurado**: Logs en formato JSON para Cloud Logging
- ⚡ **Alto rendimiento**: Binario estático compilado con CGO_ENABLED=0

## 🏗️ Arquitectura

```
go-rest-api/
├── cmd/server/
│   └── main.go              # Aplicación principal optimizada para Cloud Run
├── Dockerfile               # Multi-stage build optimizado
├── .dockerignore           # Optimización del contexto de build
├── cloudbuild.yaml         # Configuración de Google Cloud Build
├── deploy.sh               # Script de despliegue automatizado
└── README.md               # Esta documentación
```

## 🚀 Despliegue en Google Cloud Run

### Opción 1: Despliegue Automático (Recomendado)

```bash
# 1. Asegúrate de tener gcloud CLI instalado y configurado
gcloud auth login
gcloud config set project YOUR_PROJECT_ID

# 2. Ejecuta el script de despliegue
chmod +x deploy.sh
./deploy.sh
```

### Opción 2: Despliegue Manual

```bash
# 1. Habilitar APIs necesarias
gcloud services enable cloudbuild.googleapis.com
gcloud services enable run.googleapis.com
gcloud services enable containerregistry.googleapis.com

# 2. Build y deploy con Cloud Build
gcloud builds submit --config cloudbuild.yaml

# 3. Verificar el despliegue
gcloud run services describe go-rest-api --region=us-central1
```

## 🐳 Desarrollo Local con Docker

### Construir la imagen

```bash
docker build -t go-rest-api .
```

### Ejecutar localmente

```bash
# Ejecutar en puerto 8080
docker run -p 8080:8080 go-rest-api

# Ejecutar en puerto personalizado
docker run -p 3000:8080 -e PORT=8080 go-rest-api
```

### Probar la aplicación

```bash
# Health check
curl http://localhost:8080/health

# Endpoint principal
curl http://localhost:8080/

# Status endpoint
curl http://localhost:8080/v1/status

# Ping endpoint
curl http://localhost:8080/v1/ping
```

## 💻 Desarrollo Local sin Docker

### Requisitos

- Go 1.23 o superior
- Git

### Instalación

```bash
# 1. Clonar el repositorio
git clone https://github.com/m4r4v/go-rest-api.git
cd go-rest-api

# 2. Descargar dependencias
go mod tidy

# 3. Compilar
go build -o main ./cmd/server

# 4. Ejecutar
PORT=8080 ./main
```

## 📋 Variables de Entorno

| Variable | Descripción | Valor por Defecto | Requerido |
|----------|-------------|-------------------|-----------|
| `PORT` | Puerto donde escucha el servidor | `8080` | ✅ (Cloud Run) |
| `LOG_LEVEL` | Nivel de logging | `info` | ❌ |

**Nota**: Google Cloud Run proporciona automáticamente la variable `PORT`. La aplicación está configurada para leerla correctamente.

## 🔍 Endpoints Disponibles

### Health Check
```http
GET /health
```
**Respuesta:**
```json
{
  "success": true,
  "status_code": 200,
  "status": "OK",
  "data": {
    "service": "go-rest-api",
    "version": "2.0.0",
    "status": "healthy",
    "timestamp": "2025-07-09T20:21:20Z"
  },
  "timestamp": "2025-07-09T20:21:20Z"
}
```

### Información General
```http
GET /
```

### Status del Sistema
```http
GET /v1/status
```

### Ping
```http
GET /v1/ping
```

## 🛠️ Configuración de Google Cloud Run

La aplicación está optimizada para Cloud Run con las siguientes características:

### ✅ Cumplimiento de Estándares

- **Puerto**: Escucha en `0.0.0.0:$PORT` (variable proporcionada por Cloud Run)
- **Graceful Shutdown**: Maneja SIGTERM con timeout de 10 segundos
- **Health Checks**: Endpoint `/health` para verificación de estado
- **Usuario no-root**: Contenedor ejecuta como usuario `appuser` (UID 1001)
- **Stateless**: No mantiene estado local, ideal para escalado automático

### ⚙️ Configuración Recomendada

```yaml
# En cloudbuild.yaml
args: [
  'run', 'deploy', 'go-rest-api',
  '--region', 'us-central1',
  '--platform', 'managed',
  '--allow-unauthenticated',
  '--port', '8080',
  '--memory', '512Mi',
  '--cpu', '1',
  '--concurrency', '80',
  '--max-instances', '100',
  '--timeout', '300'
]
```

## 🔧 Solución de Problemas

### Error: "Container failed to start"

**Causa**: El contenedor no está escuchando en el puerto correcto.

**Solución**: 
- Verificar que la aplicación lee la variable `PORT`
- Asegurar que escucha en `0.0.0.0:$PORT`, no en `localhost`

### Error: "Health check failed"

**Causa**: El endpoint `/health` no responde correctamente.

**Solución**:
- Verificar que `/health` retorna status 200
- Comprobar que el contenedor está corriendo

### Error de Build

**Causa**: Dependencias o configuración incorrecta.

**Solución**:
```bash
# Limpiar y reconstruir
go mod tidy
go clean -cache
docker build --no-cache -t go-rest-api .
```

## 📊 Monitoreo y Logs

### Ver logs en Cloud Run

```bash
# Logs en tiempo real
gcloud run services logs tail go-rest-api --region=us-central1

# Logs históricos
gcloud run services logs read go-rest-api --region=us-central1 --limit=50
```

### Métricas disponibles

- Latencia de requests
- Número de instancias activas
- CPU y memoria utilizadas
- Errores HTTP

## 🚀 Próximos Pasos

1. **Autenticación**: Implementar JWT o OAuth2
2. **Base de datos**: Conectar a Cloud SQL o Firestore
3. **Caching**: Implementar Redis para cache
4. **Monitoring**: Configurar alertas en Cloud Monitoring
5. **CI/CD**: Automatizar despliegues con GitHub Actions

## 📄 Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo [LICENSE](LICENSE) para más detalles.

## 🤝 Contribuir

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📞 Soporte

Si tienes problemas con el despliegue en Google Cloud Run:

1. Revisa los logs: `gcloud run services logs tail go-rest-api --region=us-central1`
2. Verifica la configuración: `gcloud run services describe go-rest-api --region=us-central1`
3. Prueba localmente primero con Docker
4. Consulta la [documentación oficial de Cloud Run](https://cloud.google.com/run/docs)

---

**¡Tu aplicación está lista para Google Cloud Run! 🎉**
