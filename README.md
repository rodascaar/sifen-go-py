# SIFEN Go - Cliente de Facturación Electrónica para Paraguay

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

Cliente Go para integración con **SIFEN** (Sistema Integrado de Facturación Electrónica Nacional) de la SET (Subsecretaría de Estado de Tributación) de Paraguay.

**Versión del Manual Técnico soportado: 150**

## Características

- ✅ Generación de Documentos Electrónicos (DE)
- ✅ Firma digital XML (RSA-SHA256)
- ✅ Envío individual y por lotes
- ✅ Consulta de documentos y RUC
- ✅ Eventos (Cancelación, Inutilización, Conformidad, etc.)
- ✅ Generación automática de CDC
- ✅ Validación de RUC
- ✅ Cálculo automático de totales e IVA
- ✅ **NUEVO:** Generación de KuDE (Representación Gráfica) en HTML/PDF
- ✅ **NUEVO:** Sistema de Caché optimizado (RUC y Consultas)
- ✅ **NUEVO:** Manejo de errores tipados y específicos

## Instalación

```bash
go get github.com/rodascaar/sifen-go-py
```

## Documentación

- [Guía de Inicio Rápido](README.md)
- [Auditoría de Cumplimiento](auditoria_cumplimiento_sifen.md)
- [Guía de Características Avanzadas (KuDE, Cache, Errores)](manual_usuario_avanzado.md)

## Inicio Rápido

### Configuración del Cliente

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/rodascaar/sifen-go-py/sifen"
    "github.com/rodascaar/sifen-go-py/sifen/cache"
)

func main() {
    // Crear configuración
    config := sifen.NewSifenConfig()
    config.SetAmbiente(sifen.TipoAmbienteDev) // Usar TipoAmbienteProd para producción
    config.CertificadoCliente = "/ruta/al/certificado.pfx"
    config.ContrasenaCertificadoCliente = "password"
    config.IdCSC = "0001"
    config.CSC = "TU_CSC_SECRETO"
    
    // Configurar Cache (Opcional, habilitado por defecto)
    config.CacheConfig = cache.DefaultCacheConfig()

    // Crear cliente
    client, err := sifen.NewSifenClient(config)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close() // Importante para limpiar recursos de caché

    // Consultar RUC (usa caché automáticamente)
    resp, err := client.ConsultaRUC("80069563-1")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Razón Social: %s\n", resp.XContRUC.DRazCons)
}
```

### Crear una Factura Electrónica y Generar KuDE

```go
package main

import (
    "log"
    "time"
    
    "github.com/rodascaar/sifen-go-py/sifen"
    "github.com/rodascaar/sifen-go-py/sifen/models"
    "github.com/rodascaar/sifen-go-py/sifen/types"
    "github.com/rodascaar/sifen-go-py/sifen/kude"
)

func main() {
    // ... configuración cliente ...
    
    // Crear DE (ver ejemplos completos en manuales)
    de := models.NewDE("01800695631001001000000612024123017595714694")
    // ... rellenar datos del DE ...

    // Generar KuDE
    kudeConfig := kude.KuDEConfig{
        BaseURLQR: sifen.GetURLConsultaQR(sifen.TipoAmbienteDev),
        CSC:       "ABCD...", 
        IdCSC:     "0001",
    }
    generator := kude.NewKuDEGenerator(kudeConfig)
    
    // Convertir DE a datos de visualización
    kudeData := generator.GenerateFromDE(de)
    
    // Generar HTML
    htmlContent, err := generator.GenerateHTML(kudeData)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println("KuDE generado exitosamente")
}
```

## Estructura del Proyecto

```
go-implementation/
├── cmd/demo/           # Ejemplo de uso
├── internal/
│   ├── signature/      # Firma digital XML
│   ├── soap/           # Cliente SOAP
│   └── util/           # Utilidades (CDC, RUC, cálculos)
└── sifen/
    ├── client.go       # Cliente principal
    ├── config.go       # Configuración
    ├── events/         # Eventos SIFEN
    ├── models/         # Modelos de datos XML
    ├── kude/           # Generador de Representación Gráfica (NUEVO)
    ├── cache/          # Sistema de Caché (NUEVO)
    ├── errors/         # Errores Tipados (NUEVO)
    ├── request/        # Tipos de solicitud
    ├── response/       # Tipos de respuesta
    └── types/          # Enums y tipos base
```

## Configuración y Features Avanzados

### Caché
El sistema incluye un caché en memoria para reducir llamadas redundantes a la SET:
- **Consulta RUC:** TTL 30 minutos
- **Consulta DE:** TTL 10 minutos (solo estados finales)

### Errores Tipados
Manejo robusto de errores con el paquete `sifen/errors`:
```go
if sifenErr, ok := err.(*errors.SifenError); ok {
   fmt.Printf("Código SIFEN: %s, Categoría: %s\n", sifenErr.Code, sifenErr.Category)
}
```

## Testing

```bash
cd go-implementation
go test ./... -v
```

## Licencia

MIT License - ver [LICENSE](LICENSE) para más detalles.

## Referencias

- [Manual Técnico SIFEN v150](https://ekuatia.set.gov.py/)
- [Portal eKuatia](https://ekuatia.set.gov.py/)
