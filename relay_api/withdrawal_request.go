package relay_api

type WithdrawalRequest struct {
	ID                uint64
	Timestamp         uint64
	Nonce             uint64
	Blockchain        string
	Address           string
	BeneficialAddress string
	Amount            string
	Signature         string
}
