# --------------------------------------------------------
# ETAPA 1: Builder (Compilación)
# Usamos Alpine para una base ligera de construcción
# --------------------------------------------------------
FROM golang:1.25-alpine AS builder

# Instalamos certificados y git (necesario para algunas dependencias)
# y tzdata para zonas horarias.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Creamos un usuario "appuser" para no correr como root al final
ENV USER=appuser
ENV UID=10001 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

WORKDIR /app

# 1. Copiamos solo los archivos de dependencias primero (Optimización de Caché)
COPY go.mod go.sum ./

# 2. Descargamos dependencias. Si go.mod no cambia, Docker usa la caché aquí.
RUN go mod download

# 3. Copiamos el código fuente
COPY . .

# 4. Compilación "Super Ultra"
# CGO_ENABLED=0: Deshabilita CGO para crear un binario 100% estático.
# -ldflags="-w -s": Elimina información de debug y tablas de símbolos (reduce tamaño ~30%).
# -trimpath: Quita rutas del sistema de archivos del binario (seguridad/reproducibilidad).
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -trimpath \
    -o main cmd/main.go

# --------------------------------------------------------
# ETAPA 2: Runtime (Ejecución)
# Usamos 'scratch', que es una imagen vacía (0 MB overhead)
# --------------------------------------------------------
FROM scratch

# Importamos cosas esenciales del builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copiamos el binario compilado
COPY --from=builder /app/main /main

# Usamos el usuario sin privilegios
USER appuser:appuser

# Exponemos el puerto (documentación, ajusta según tu app)
EXPOSE 3000

ENTRYPOINT ["/main"]