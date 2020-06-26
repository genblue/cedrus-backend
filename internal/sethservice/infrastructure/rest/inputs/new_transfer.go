package inputs

type NewTransfer struct {
	ClaimCode string `json:"claim-code"`
	Address   string `json:"address"`
}
