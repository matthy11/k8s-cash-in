package depositvalidations

type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
	Valid   bool   `json:"valid"`
}
