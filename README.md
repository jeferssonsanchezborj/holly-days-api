# Holidays API 🇨🇱

Servicio REST desarrollado en **Go + Gin Gonic** que expone los feriados
chilenos publicados por [`api.boostr.cl`](https://api.boostr.cl/holidays.json),
permitiendo filtrar por **tipo** y **rango de fechas**, con soporte de
respuesta en **JSON** y **XML**.

Construido siguiendo **Arquitectura Hexagonal (Ports & Adapters)**, con logs
estructurados, tests unitarios, documentación OpenAPI 3.x y despliegue vía
Docker.

---

## 📐 Arquitectura

```
cmd/api/main.go                          → Composition root (wiring de dependencias)
internal/
  domain/                                → Entidades y reglas de negocio puras
    holiday.go                           → Holiday, HolidayFilter
    ports/ports.go                       → Interfaces (puertos): HolidayRepository, HolidayService
  application/                           → Casos de uso (lógica de filtrado)
    holiday_service.go
    holiday_service_test.go
  infrastructure/                        → Adaptadores (detalles técnicos)
    repository/                          → Adaptador secundario: cliente HTTP a api.boostr.cl
      boostr_repository.go               → Fetch ÚNICO al iniciar + cache en memoria
      boostr_repository_test.go
    http/
      handler/                           → Adaptador primario: controladores Gin
      middleware/                        → Logging de requests
      router/                            → Definición de rutas
      dto/                               → Objetos de transferencia (JSON/XML)
pkg/logger/                              → Logger estructurado (logrus)
docs/openapi.yaml                        → Especificación OpenAPI 3.0.3
```

**Regla clave del desafío**: el servicio externo `https://api.boostr.cl/holidays.json`
se invoca **una sola vez**, durante el bootstrap de la aplicación
(`main.go` → `repository.NewBoostrHolidayRepository`). Todas las peticiones
posteriores del API son resueltas 100% en memoria, sin volver a llamar al
servicio externo. Esto queda demostrado en los logs de arranque
(`"fetching holidays from upstream service (one-time bootstrap call)"`)
que aparece una única vez por ciclo de vida del proceso.

### Puertos y adaptadores

| Tipo                | Puerto (interfaz)                  | Adaptador                              |
|---------------------|-------------------------------------|-----------------------------------------|
| Primario (driving)  | `ports.HolidayService`             | `handler.HolidayHandler` (HTTP/Gin)     |
| Secundario (driven) | `ports.HolidayRepository`          | `repository.BoostrHolidayRepository`    |

---

## 🚀 Levantar el servicio con Docker

### Opción 1: Docker Compose (recomendado)

```bash
docker compose up --build
```

### Opción 2: Docker manual

```bash
docker build -t holidays-api:latest .
docker run -d --name holidays-api -p 8080:8080 holidays-api:latest
```

El servicio quedará disponible en `http://localhost:8080`.

### Variables de entorno soportadas

| Variable                | Default                                  | Descripción                               |
|--------------------------|-------------------------------------------|--------------------------------------------|
| `PORT`                   | `8080`                                    | Puerto HTTP del servicio                   |
| `LOG_LEVEL`              | `info`                                    | Nivel de log (`debug`, `info`, `warn`, `error`) |
| `HOLIDAYS_UPSTREAM_URL`  | `https://api.boostr.cl/holidays.json`     | URL del servicio de feriados (fetch único) |

---

## 💻 Ejecutar en local (sin Docker)

Requiere Go 1.24+ (recomendado 1.26).

```bash
go mod tidy
go run ./cmd/api
```

---

## 🧪 Tests unitarios

```bash
go test ./... -v -cover
```

Cobertura obtenida:
- `internal/application`: 100%
- `internal/infrastructure/repository`: 85.2%
- `internal/infrastructure/http/handler`: 82.1%

---

## 📡 Endpoints

### `GET /api/v1/holidays`

Retorna la lista de feriados, con **fecha, título, teléfono, tipo,
inalienable y extra** por cada uno.

**Query params (opcionales):**

| Parámetro    | Formato      | Descripción                                             |
|--------------|--------------|----------------------------------------------------------|
| `type`       | texto        | Filtra por tipo (ej. `Civil`, `Religioso`). Case-insensitive |
| `start_date` | `YYYY-MM-DD` | Fecha inicial del rango (inclusive)                      |
| `end_date`   | `YYYY-MM-DD` | Fecha final del rango (inclusive)                        |

**Negociación de contenido:** enviar el header `Accept` para elegir el formato:
- `Accept: application/json` (default)
- `Accept: application/xml`

### `GET /health`
Health check simple (`{"status":"ok"}`).

### `GET /docs`
Swagger UI navegable, sirviendo `docs/openapi.yaml`.

### `GET /openapi.yaml`
Especificación OpenAPI 3.0.3 cruda.

---

## 📋 Ejemplos de uso (cURL)

### Todos los feriados (JSON, por defecto)
```bash
curl "http://localhost:8080/api/v1/holidays"
```

### Todos los feriados en XML
```bash
curl -H "Accept: application/xml" "http://localhost:8080/api/v1/holidays"
```

### Filtrar por tipo
```bash
curl "http://localhost:8080/api/v1/holidays?type=Civil"
curl "http://localhost:8080/api/v1/holidays?type=Religioso"
```

### Filtrar por rango de fechas
```bash
curl "http://localhost:8080/api/v1/holidays?start_date=2026-01-01&end_date=2026-06-30"
```

### Combinando tipo + rango de fechas, en XML
```bash
curl -H "Accept: application/xml" \
  "http://localhost:8080/api/v1/holidays?type=Civil&start_date=2026-01-01&end_date=2026-12-31"
```

### Ejemplo de respuesta (JSON)
```json
{
  "total": 1,
  "data": [
    {
      "date": "2026-01-01",
      "title": "Año Nuevo",
      "phone": "",
      "type": "Civil",
      "inalienable": true,
      "extra": "Civil e Irrenunciable"
    }
  ]
}
```

> **Nota sobre el campo `phone`**: el servicio origen (`api.boostr.cl`) no
> expone actualmente un teléfono de contacto asociado a cada feriado. El
> campo se mantiene en el contrato de datos (dominio, DTO y OpenAPI) tal
> como fue solicitado en el desafío, retornando cadena vacía mientras la
> fuente de datos no provea dicho valor.

### Health check
```bash
curl "http://localhost:8080/health"
```

### Caso de error (fecha inválida) → HTTP 400
```bash
curl "http://localhost:8080/api/v1/holidays?start_date=not-a-date"
```
```json
{ "message": "invalid value for 'start_date': 'not-a-date', expected format YYYY-MM-DD" }
```

---

## 📚 Documentación OpenAPI

La especificación completa está en [`docs/openapi.yaml`](docs/openapi.yaml)
(OpenAPI 3.0.3). Puede visualizarse de forma interactiva accediendo a
`http://localhost:8080/docs` una vez levantado el servicio, o importando el
archivo directamente en [Swagger Editor](https://editor.swagger.io/).

---

## 🛠️ Stack Tecnológico

- **Go 1.26**
- **Gin Gonic** — framework HTTP
- **Logrus** — logging estructurado (JSON)
- **Testify** — aserciones para tests unitarios
- **Docker / Docker Compose**
- **OpenAPI 3.0.3**

---

## 📦 Estructura de commits / repositorio

Este repositorio fue inicializado localmente con `git init`. Para publicarlo
en un repositorio remoto público (por ejemplo en GitHub):

```bash
git init
git add .
git commit -m "feat: Holidays API - Hexagonal Architecture, Gin, Docker"
git branch -M main
git remote add origin <URL_DE_TU_REPOSITORIO_REMOTO>
git push -u origin main
```

