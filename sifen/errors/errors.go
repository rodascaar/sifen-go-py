package errors

import (
	"fmt"
	"strings"
)

// ============================================================================
// Tipos de Error SIFEN
// ============================================================================

// ErrorType categoriza el tipo de error
type ErrorType string

const (
	// ErrorTypeValidation indica error de validación de datos
	ErrorTypeValidation ErrorType = "VALIDATION"
	// ErrorTypeCryptography indica error de firma/certificado
	ErrorTypeCryptography ErrorType = "CRYPTOGRAPHY"
	// ErrorTypeNetwork indica error de red/conexión
	ErrorTypeNetwork ErrorType = "NETWORK"
	// ErrorTypeSIFEN indica error retornado por SIFEN
	ErrorTypeSIFEN ErrorType = "SIFEN"
	// ErrorTypeBusiness indica error de regla de negocio
	ErrorTypeBusiness ErrorType = "BUSINESS"
	// ErrorTypeInternal indica error interno del cliente
	ErrorTypeInternal ErrorType = "INTERNAL"
)

// ============================================================================
// Error Principal
// ============================================================================

// SifenError representa un error del sistema SIFEN
type SifenError struct {
	// Tipo de error
	Type ErrorType
	// Código de error (puede ser código SIFEN o interno)
	Code string
	// Mensaje descriptivo
	Message string
	// Error original (si existe)
	Cause error
	// Datos adicionales del contexto
	Context map[string]interface{}
	// Indica si el error es recuperable (puede reintentarse)
	Recoverable bool
}

// Error implementa la interfaz error
func (e *SifenError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %s (cause: %v)", e.Type, e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s: %s", e.Type, e.Code, e.Message)
}

// Unwrap permite acceder al error original
func (e *SifenError) Unwrap() error {
	return e.Cause
}

// Is permite comparar errores
func (e *SifenError) Is(target error) bool {
	if t, ok := target.(*SifenError); ok {
		return e.Code == t.Code
	}
	return false
}

// WithContext agrega contexto al error
func (e *SifenError) WithContext(key string, value interface{}) *SifenError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// ============================================================================
// Constructores de Errores
// ============================================================================

// NewValidationError crea un error de validación
func NewValidationError(code, message string) *SifenError {
	return &SifenError{
		Type:        ErrorTypeValidation,
		Code:        code,
		Message:     message,
		Recoverable: false,
	}
}

// NewCryptoError crea un error de criptografía
func NewCryptoError(message string, cause error) *SifenError {
	return &SifenError{
		Type:        ErrorTypeCryptography,
		Code:        "CRYPTO_ERROR",
		Message:     message,
		Cause:       cause,
		Recoverable: false,
	}
}

// NewNetworkError crea un error de red
func NewNetworkError(message string, cause error) *SifenError {
	return &SifenError{
		Type:        ErrorTypeNetwork,
		Code:        "NETWORK_ERROR",
		Message:     message,
		Cause:       cause,
		Recoverable: true, // Los errores de red suelen ser recuperables
	}
}

// NewSifenResponseError crea un error desde respuesta SIFEN
func NewSifenResponseError(code, message string) *SifenError {
	return &SifenError{
		Type:        ErrorTypeSIFEN,
		Code:        code,
		Message:     message,
		Recoverable: isRecoverableSifenCode(code),
	}
}

// NewBusinessError crea un error de negocio
func NewBusinessError(code, message string) *SifenError {
	return &SifenError{
		Type:        ErrorTypeBusiness,
		Code:        code,
		Message:     message,
		Recoverable: false,
	}
}

// NewInternalError crea un error interno
func NewInternalError(message string, cause error) *SifenError {
	return &SifenError{
		Type:        ErrorTypeInternal,
		Code:        "INTERNAL_ERROR",
		Message:     message,
		Cause:       cause,
		Recoverable: false,
	}
}

// ============================================================================
// Errores Predefinidos - Validación
// ============================================================================

var (
	// ErrCDCInvalido indica que el CDC no tiene 44 caracteres
	ErrCDCInvalido = NewValidationError("VAL_001", "CDC debe tener 44 caracteres")

	// ErrRUCInvalido indica formato RUC inválido
	ErrRUCInvalido = NewValidationError("VAL_002", "Formato de RUC inválido")

	// ErrRUCDigitoVerificador indica dígito verificador incorrecto
	ErrRUCDigitoVerificador = NewValidationError("VAL_003", "Dígito verificador de RUC incorrecto")

	// ErrFechaInvalida indica fecha fuera de rango
	ErrFechaInvalida = NewValidationError("VAL_004", "Fecha inválida o fuera de rango")

	// ErrMontoInvalido indica monto negativo o inválido
	ErrMontoInvalido = NewValidationError("VAL_005", "Monto inválido")

	// ErrDocumentoVacio indica documento sin items
	ErrDocumentoVacio = NewValidationError("VAL_006", "Documento sin items")

	// ErrEstablecimientoInvalido indica establecimiento con formato incorrecto
	ErrEstablecimientoInvalido = NewValidationError("VAL_007", "Establecimiento debe tener 3 dígitos")

	// ErrPuntoExpedicionInvalido indica punto de expedición con formato incorrecto
	ErrPuntoExpedicionInvalido = NewValidationError("VAL_008", "Punto de expedición debe tener 3 dígitos")

	// ErrTimbradoInvalido indica número de timbrado inválido
	ErrTimbradoInvalido = NewValidationError("VAL_009", "Número de timbrado inválido")

	// ErrMotivoCancelacionRequerido indica que falta el motivo de cancelación
	ErrMotivoCancelacionRequerido = NewValidationError("VAL_010", "Motivo de cancelación es requerido")
)

// ============================================================================
// Errores Predefinidos - Criptografía
// ============================================================================

var (
	// ErrCertificadoNoEncontrado indica que no se encontró el certificado
	ErrCertificadoNoEncontrado = &SifenError{
		Type:    ErrorTypeCryptography,
		Code:    "CRYPTO_001",
		Message: "Certificado no encontrado",
	}

	// ErrCertificadoExpirado indica certificado vencido
	ErrCertificadoExpirado = &SifenError{
		Type:    ErrorTypeCryptography,
		Code:    "CRYPTO_002",
		Message: "Certificado expirado",
	}

	// ErrContrasenaCertificado indica contraseña incorrecta
	ErrContrasenaCertificado = &SifenError{
		Type:    ErrorTypeCryptography,
		Code:    "CRYPTO_003",
		Message: "Contraseña del certificado incorrecta",
	}

	// ErrFirmaInvalida indica error al firmar
	ErrFirmaInvalida = &SifenError{
		Type:    ErrorTypeCryptography,
		Code:    "CRYPTO_004",
		Message: "Error al generar firma digital",
	}

	// ErrClavePrivadaNoRSA indica que la clave no es RSA
	ErrClavePrivadaNoRSA = &SifenError{
		Type:    ErrorTypeCryptography,
		Code:    "CRYPTO_005",
		Message: "Clave privada no es RSA",
	}
)

// ============================================================================
// Errores Predefinidos - Red
// ============================================================================

var (
	// ErrTimeoutConexion indica timeout de conexión
	ErrTimeoutConexion = &SifenError{
		Type:        ErrorTypeNetwork,
		Code:        "NET_001",
		Message:     "Timeout de conexión con SIFEN",
		Recoverable: true,
	}

	// ErrTimeoutLectura indica timeout de lectura
	ErrTimeoutLectura = &SifenError{
		Type:        ErrorTypeNetwork,
		Code:        "NET_002",
		Message:     "Timeout de lectura de respuesta SIFEN",
		Recoverable: true,
	}

	// ErrConexionRechazada indica conexión rechazada
	ErrConexionRechazada = &SifenError{
		Type:        ErrorTypeNetwork,
		Code:        "NET_003",
		Message:     "Conexión rechazada por servidor SIFEN",
		Recoverable: true,
	}

	// ErrTLSHandshake indica error en handshake TLS
	ErrTLSHandshake = &SifenError{
		Type:        ErrorTypeNetwork,
		Code:        "NET_004",
		Message:     "Error en handshake TLS con SIFEN",
		Recoverable: false,
	}
)

// ============================================================================
// Errores Predefinidos - SIFEN
// ============================================================================

var (
	// ErrSifenCDCDuplicado indica CDC ya existe
	ErrSifenCDCDuplicado = NewSifenResponseError("0160", "CDC duplicado - documento ya existe")

	// ErrSifenEmisorNoAutorizado indica emisor no autorizado
	ErrSifenEmisorNoAutorizado = NewSifenResponseError("0161", "Emisor no autorizado para facturación electrónica")

	// ErrSifenFirmaInvalida indica firma digital inválida
	ErrSifenFirmaInvalida = NewSifenResponseError("0102", "Firma digital inválida")

	// ErrSifenXMLInvalido indica estructura XML inválida
	ErrSifenXMLInvalido = NewSifenResponseError("0101", "Estructura XML inválida")

	// ErrSifenServicioNoDisponible indica servicio no disponible
	ErrSifenServicioNoDisponible = NewSifenResponseError("0500", "Servicio SIFEN no disponible")
)

// ============================================================================
// Errores Predefinidos - Negocio
// ============================================================================

var (
	// ErrLoteVacio indica lote sin documentos
	ErrLoteVacio = NewBusinessError("BUS_001", "El lote debe contener al menos un documento")

	// ErrLoteExcedeMaximo indica demasiados documentos en lote
	ErrLoteExcedeMaximo = NewBusinessError("BUS_002", "El lote excede el máximo de 50 documentos")

	// ErrLoteTipoMixto indica mezcla de tipos en lote
	ErrLoteTipoMixto = NewBusinessError("BUS_003", "Todos los documentos del lote deben ser del mismo tipo")

	// ErrEventoCDCRequerido indica que falta CDC para evento
	ErrEventoCDCRequerido = NewBusinessError("BUS_004", "CDC es requerido para registrar evento")

	// ErrCancelacionFueraPlazo indica cancelación fuera del plazo permitido
	ErrCancelacionFueraPlazo = NewBusinessError("BUS_005", "Cancelación fuera del plazo permitido")
)

// ============================================================================
// Mapa de Códigos SIFEN
// ============================================================================

// SifenErrorCodes mapea códigos SIFEN a mensajes descriptivos
var SifenErrorCodes = map[string]string{
	// Éxito
	"0260": "Documento recibido correctamente",
	"0261": "Documento aprobado con observaciones",
	"0300": "Documento procesado correctamente",

	// Errores generales
	"0100": "Error interno del servidor SIFEN",
	"0101": "Estructura XML inválida",
	"0102": "Firma digital inválida",
	"0103": "CDC inválido",

	// Errores de documento
	"0160": "CDC duplicado",
	"0161": "Emisor no autorizado",
	"0162": "Receptor inválido",
	"0163": "Timbrado no vigente",
	"0164": "Número de documento fuera de rango",

	// Errores de lote
	"0320": "Tamaño de mensaje excede límite permitido",
	"0321": "Formato de lote inválido",
	"0322": "Lote contiene tipos de documento mixtos",

	// Errores de evento
	"0510": "Evento procesado correctamente",
	"0520": "Evento aceptado",
	"0530": "Evento rechazado",
	"0531": "Evento fuera de plazo",

	// Servicio no disponible
	"0500": "Servicio temporalmente no disponible",
	"0501": "Mantenimiento programado",
}

// GetSifenErrorMessage obtiene el mensaje para un código SIFEN
func GetSifenErrorMessage(code string) string {
	if msg, ok := SifenErrorCodes[code]; ok {
		return msg
	}
	return fmt.Sprintf("Código de error SIFEN: %s", code)
}

// isRecoverableSifenCode determina si un error SIFEN es recuperable
func isRecoverableSifenCode(code string) bool {
	// Errores de servicio no disponible son recuperables
	recoverableCodes := map[string]bool{
		"0500": true, // Servicio no disponible
		"0501": true, // Mantenimiento
		"0100": true, // Error interno (puede ser temporal)
	}
	return recoverableCodes[code]
}

// ============================================================================
// Helper Functions
// ============================================================================

// IsSifenError verifica si un error es de tipo SifenError
func IsSifenError(err error) bool {
	_, ok := err.(*SifenError)
	return ok
}

// AsSifenError convierte un error a SifenError si es posible
func AsSifenError(err error) (*SifenError, bool) {
	se, ok := err.(*SifenError)
	return se, ok
}

// IsRecoverable verifica si un error es recuperable
func IsRecoverable(err error) bool {
	if se, ok := err.(*SifenError); ok {
		return se.Recoverable
	}
	return false
}

// IsValidationError verifica si es error de validación
func IsValidationError(err error) bool {
	if se, ok := err.(*SifenError); ok {
		return se.Type == ErrorTypeValidation
	}
	return false
}

// IsNetworkError verifica si es error de red
func IsNetworkError(err error) bool {
	if se, ok := err.(*SifenError); ok {
		return se.Type == ErrorTypeNetwork
	}
	return false
}

// IsCryptoError verifica si es error de criptografía
func IsCryptoError(err error) bool {
	if se, ok := err.(*SifenError); ok {
		return se.Type == ErrorTypeCryptography
	}
	return false
}

// Wrap envuelve un error en SifenError
func Wrap(err error, message string) *SifenError {
	if err == nil {
		return nil
	}

	// Si ya es SifenError, agregar contexto
	if se, ok := err.(*SifenError); ok {
		se.Message = message + ": " + se.Message
		return se
	}

	// Detectar tipo de error por contenido
	errStr := strings.ToLower(err.Error())

	if strings.Contains(errStr, "timeout") {
		return NewNetworkError(message, err)
	}
	if strings.Contains(errStr, "connection refused") {
		return NewNetworkError(message, err)
	}
	if strings.Contains(errStr, "certificate") || strings.Contains(errStr, "tls") {
		return NewCryptoError(message, err)
	}

	return NewInternalError(message, err)
}
