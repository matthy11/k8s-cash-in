package constants

var BankID = map[string]string{
	"0001": "bch-chile",
	"0009": "internacional-chile",
	"0012": "banco-estado-chile",
	"0014": "scotiabank-chile",
	"0016": "bci-chile",
	"0027": "corpbanca-chile",
	"0028": "bice-chile",
	"0031": "hsbc-chile",
	"0037": "santander-chile",
	"0039": "itau-chile",
	"0049": "security-chile",
	"0051": "falabella-chile",
	"0053": "ripley-chile",
	"0054": "rabobank-chile",
	"0055": "consorcio-chile",
	"0057": "paris-chile",
	"0504": "scotiabank-azul-chile",
	"0507": "desarrollo-chile",
}

var BankAccountTypeID = map[string]string{
	"10": "cuenta-ahorro",
	"20": "cuenta-corriente",
	"40": "cuenta-vista",
}

type ContextKey string

const (
	ContextKeyCurrentStepCode  ContextKey = "currentStepCode"
	ContextKeyCurrentStepLabel ContextKey = "currentStepLabel"
	ContextKeyCurrentStepTime  ContextKey = "currentStepTime"
	ContextKeyIPAddress        ContextKey = "ip"
	ContextKeyLoggerID         ContextKey = "loggerId"
	ContextKeyStartTime        ContextKey = "startTime"
)

const (
	ErrorLabel_Error_In_Server = "Error in server."
)

const (
	ErrorCode_Http_Error    = "http-error"
	ErrorCode_Unknown_Error = "unknown-error"
)
