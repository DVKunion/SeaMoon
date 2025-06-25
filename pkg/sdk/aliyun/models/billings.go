package models

type BillingResponse struct {
	StatusCode int                    `json:"statusCode"`
	Headers    map[string]interface{} `json:"headers"`
	Body       Billing                `json:"body"`
}

type Billing struct {
	Code      string                 `json:"Code"`
	Message   string                 `json:"Message"`
	RequestId string                 `json:"RequestId"`
	Success   bool                   `json:"Success"`
	Data      map[string]interface{} `json:"Data"`
}
