package lending

type InsertLendingRequest struct {
	AccountID uint    `json:"account_id"`
	Amount    float64 `json:"amount"`
	Tenor     int     `json:"tenor"`
}

type SimulateLendingResponse struct {
	Fee          float64 `json:"fee"`
	FeeStampDuty float64 `json:"fee_stamp_duty"`
	Interest     float64 `json:"interest"`
	TotalPayment float64 `json:"total_payment"`
}

type SimulateRepaymentResponse struct {
	Month            int     `json:"month"`
	Year             int     `json:"year"`
	Fee              float64 `json:"fee"`
	FeeStampDuty     float64 `json:"fee_stamp_duty"`
	Interest         float64 `json:"interest"`
	PokokYangDibayar float64 `json:"pokok_yang_dibayar"`
	Tagihan          float64 `json:"tagihan"`
	Outstanding      float64 `json:"outstanding"`
}
