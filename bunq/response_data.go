package bunq

type responseInstallation struct {
	Response []installation
}

type responseError struct {
	Error []bunqError `json:"Error"`
}

type responseDeviceServer struct {
	Response []wrappedBunqID
}

type responseSessionServer struct {
	Response []sessionServer
}

type responseUserPerson struct {
	Response []struct {
		UserPerson userPerson
	}
}

type responseBunqID struct {
	Response []wrappedBunqID
}

// ResponseMonetaryAccountBankGet The monetary account bank response object.
type ResponseMonetaryAccountBankGet struct {
	Response []struct {
		MonetaryAccountBank MonetaryAccountBank `json:"MonetaryAccountBank"`
	} `json:"Response"`
	Pagination Pagination `json:"Pagination"`
}

// ResponseMonetaryAccountSavingGet The monetary account savings response object.
type ResponseMonetaryAccountSavingGet struct {
	Response []struct {
		MonetaryAccountSaving MonetaryAccountSaving `json:"MonetaryAccountSavings"`
	} `json:"Response"`
	Pagination Pagination `json:"Pagination"`
}

type responseDraftPaymentGet struct {
	Response []struct {
		DraftPayment draftPayment `json:"DraftPayment"`
	} `json:"Response"`
}

// ResponsePaymentGet The payment response data.
type ResponsePaymentGet struct {
	Response []struct {
		Payment Payment `json:"Payment"`
	} `json:"Response"`
	Pagination Pagination `json:"Pagination"`
}

type responseMasterCardActionGet struct {
	Response []struct {
		MasterCardAction masterCardAction `json:"MasterCardAction"`
	} `json:"Response"`
	Pagination Pagination `json:"Pagination"`
}

type bunqError struct {
	ErrorDescription           string `json:"error_description"`
	ErrorDescriptionTranslated string `json:"error_description_translated"`
}
