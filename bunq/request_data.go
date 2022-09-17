package bunq

type requestInstallation struct {
	ClientPublicKey string `json:"client_public_key"`
}

type requestDeviceServer struct {
	Description  string   `json:"description"`
	Secret       string   `json:"secret"`
	PermittedIps []string `json:"permitted_ips"`
}

type requestSessionServer struct {
	Secret string `json:"secret"`
}

type requestUserPersonPut struct {
	NotificationFilters []notificationFilter `json:"notification_filters,omitempty"`
}

type requestCreateDraftPayment struct {
	Entries                 []draftPaymentEntryCreate `json:"entries"`
	NumberOfRequiredAccepts *int                      `json:"number_of_required_accepts,omitempty"`
}

type requestUpdateDraftPayment struct {
	requestCreateDraftPayment
	UpdatedTimestamp string `json:"previous_updated_timestamp"`
	status           *string
}

type draftPaymentEntryCreate struct {
	Amount            Amount  `json:"Amount,omitempty"`
	CounterpartyAlias Pointer `json:"counterparty_alias,omitempty"`
	Description       string  `json:"description,omitempty"`
	MerchantReference *string `json:"merchant_reference,omitempty"`
}

type PaymentBatchCreate struct {
	Payments []PaymentCreate `json:"payments"`
}

type PaymentCreate struct {
	Amount            Amount  `json:"amount"`
	CounterpartyAlias Pointer `json:"counterparty_alias"`
	Description       string  `json:"description"`
	AllowBunqto       bool    `json:"allow_bunqto"`
}
