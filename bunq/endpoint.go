package bunq

const (
	endpointInstallationCreate string = "installation"

	endpointDeviceServerCreate string = "device-server"

	endpointSessionServerCreate string = "session-server"

	endpointUserPersonGet string = "user-person/%d"

	endpointPaymentBatchCreate string = "user/%d/monetary-account/%d/payment-batch"

	endpointDraftPaymentCreate string = "user/%d/monetary-account/%d/draft-payment"
	endpointDraftPaymentWithID string = "user/%d/monetary-account/%d/draft-payment/%d"

	endpointPaymentGet       string = "user/%d/monetary-account/%d/payment?count=200"
	endpointPaymentGetWithID string = "user/%d/monetary-account/%d/payment/%d"

	endpointSchedulePaymentGet string = "user/%d/monetary-account/%d/schedule-payment?count=200"

	endpointMonetaryAccountBankListing string = "user/%d/monetary-account-bank"
	endpointMonetaryAccountBankGet     string = "user/%d/monetary-account-bank/%d"

	endpointMonetaryAccountSavingsListing string = "user/%d/monetary-account-savings"
	endpointMonetaryAccountSavingsGet     string = "user/%d/monetary-account-savings/%d"

	endpointMasterCardActionGet string = "user/%d/monetary-account/%d/mastercard-action/%d"
)
