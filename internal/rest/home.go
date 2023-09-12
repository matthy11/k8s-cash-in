package rest

import (
	"fmt"
	"heypay-cash-in-server/settings"
	"net/http"
)

func home(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, fmt.Sprintf("HeyPay Cash-In Server %s v%s", settings.V.Env, settings.V.Version))
}
