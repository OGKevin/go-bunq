package bunq

type installation struct {
	ID              bunqID          `json:"Id"`
	Token           token           `json:"Token"`
	ServerPublicKey serverPublicKey `json:"ServerPublicKey"`
}

type bunqID struct {
	ID int `json:"id"`
}

type wrappedBunqID struct {
	ID bunqID `json:"Id"`
}

type token struct {
	common
	Token string `json:"token"`
}

type serverPublicKey struct {
	ServerPublicKey string `json:"server_public_key"`
}

type common struct {
	ID      int    `json:"id"`
	Created string `json:"created"`
	Updated string `json:"updated"`
}

type sessionServer struct {
	ID          bunqID      `json:"Id"`
	Token       token       `json:"Token"`
	UserCompany userCompany `json:"UserCompany"`
	UserPerson  userPerson  `json:"UserPerson"`
	UserAPIKey  userAPIKey  `json:"UserApiKey"`
}

// userType holds the flag of which user is currently authenticated.
type userType struct {
	isUserPerson  bool
	isUserCompany bool
	isUserAPIkey  bool
}

type user struct {
	common
	PublicUUID                         string                             `json:"public_uuid"`
	AddressMain                        address                            `json:"address_main"`
	Alias                              []alias                            `json:"alias"`
	AddressPostal                      address                            `json:"address_postal"`
	Avatar                             avatar                             `json:"avatar"`
	Status                             string                             `json:"status"`
	SubStatus                          string                             `json:"sub_status"`
	Region                             string                             `json:"region"`
	Language                           string                             `json:"language"`
	DailyLimitWithoutConfirmationLogin dailyLimitWithoutConfirmationLogin `json:"daily_limit_without_confirmation_login"`
	NotificationFilters                []notificationFilter               `json:"notification_filters"`
	VersionTermsOfService              string                             `json:"version_terms_of_service"`
	SessionTimeout                     int64                              `json:"session_timeout"`
	DisplayName                        string                             `json:"display_name"`
	PublicNickName                     string                             `json:"public_nick_name"`
}

type userCompany struct {
	user
	Name                    string            `json:"name"`
	Country                 string            `json:"country"`
	Ubo                     []ubo             `json:"ubo"`
	ChamberOfCommerceNumber string            `json:"chamber_of_commerce_number"`
	TypeOfBusinessEntity    string            `json:"type_of_business_entity"`
	SectorOfIndustry        string            `json:"sector_of_industry"`
	CounterBankIban         string            `json:"counter_bank_iban"`
	DirectorAlias           directorAlias     `json:"director_alias"`
	CardIds                 []bunqID          `json:"card_ids"`
	CardLimits              []cardLimits      `json:"card_limits"`
	Customer                customer          `json:"customer"`
	CustomerLimit           customer          `json:"customer_limit"`
	BillingContract         []billingContract `json:"billing_contract"`
}

type userPerson struct {
	user
	FirstName                 string        `json:"first_name"`
	MiddleName                string        `json:"middle_name"`
	LastName                  string        `json:"last_name"`
	TaxResident               []taxResident `json:"tax_resident"`
	DocumentType              string        `json:"document_type"`
	DocumentNumber            string        `json:"document_number"`
	DocumentCountryOfIssuance string        `json:"document_country_of_issuance"`
	DateOfBirth               string        `json:"date_of_birth"`
	PlaceOfBirth              string        `json:"place_of_birth"`
	CountryOfBirth            string        `json:"country_of_birth"`
	Nationality               string        `json:"nationality"`
	Gender                    string        `json:"gender"`
	LegalGuardianAlias        alias         `json:"legal_guardian_alias"`
	LegalName                 string        `json:"legal_name"`
}

type address struct {
	Street      string `json:"street"`
	HouseNumber string `json:"house_number"`
	PoBox       string `json:"po_box"`
	PostalCode  string `json:"postal_code"`
	City        string `json:"city"`
	Country     string `json:"country"`
	Province    string `json:"province"`
}

type ubo struct {
	Name        string `json:"name"`
	DateOfBirth string `json:"date_of_birth"`
	Nationality string `json:"nationality"`
}

type dailyLimitWithoutConfirmationLogin Amount

type notificationFilter struct {
	NotificationDeliveryMethod string `json:"notification_delivery_method"`
	NotificationTarget         string `json:"notification_target"`
	Category                   string `json:"category"`
}

type alias struct {
	Type  string `json:"type"`
	Value string `json:"value"`
	Name  string `json:"name"`
}

type avatar struct {
	UUID       string  `json:"uuid"`
	AnchorUUID string  `json:"anchor_uuid"`
	Image      []image `json:"image"`
}

type image struct {
	AttachmentPublicUUID string `json:"attachment_public_uuid"`
	ContentType          string `json:"content_type"`
	Height               int    `json:"height"`
	Width                int    `json:"width"`
}

type directorAlias struct {
	UUID           string `json:"uuid"`
	DisplayName    string `json:"display_name"`
	Country        string `json:"country"`
	Avatar         avatar `json:"avatar"`
	PublicNickName string `json:"public_nick_name"`
}

type cardLimits struct {
	DailyLimit string `json:"daily_limit"`
	Currency   string `json:"currency"`
	Type       string `json:"type"`
	ID         int    `json:"id"`
}

type customer struct {
	BillingAccountID              int    `json:"billing_account_id"`
	InvoiceNotificationPreference string `json:"invoice_notification_preference"`
	ID                            int    `json:"id"`
	Created                       string `json:"created"`
	Updated                       string `json:"updated"`
}

type customerLimit struct {
	LimitMonetaryAccount          int    `json:"limit_monetary_account"`
	LimitCardDebitMaestro         int    `json:"limit_card_debit_maestro"`
	LimitCardDebitMastercard      int    `json:"limit_card_debit_mastercard"`
	LimitCardDebitWildcard        int    `json:"limit_card_debit_wildcard"`
	LimitCardDebitReplacement     int    `json:"limit_card_debit_replacement"`
	LimitInviteUserPremiumLimited int    `json:"limit_invite_user_premium_limited"`
	LimitAmountMonthly            Amount `json:"limit_amount_monthly"`
	SpentAmountMonthly            Amount `json:"spent_amount_monthly"`
}

type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type billingContract struct {
	SubscriptionType          string `json:"subscription_type"`
	ID                        int    `json:"id"`
	Created                   string `json:"created"`
	Updated                   string `json:"updated"`
	ContractDateStart         string `json:"contract_date_start"`
	ContractDateEnd           string `json:"contract_date_end"`
	ContractVersion           int    `json:"contract_version"`
	SubscriptionTypeDowngrade string `json:"subscription_type_downgrade"`
	Status                    string `json:"status"`
	SubStatus                 string `json:"sub_status"`
}

type taxResident struct {
	Country   string `json:"country"`
	TaxNumber string `json:"tax_number"`
	Status    string `json:"status"`
}

type userAPIKey struct {
	common
	RequestedByUser requestedByUser `json:"requested_by_user"`
	GrantedByUser   grantedByUser   `json:"granted_by_user"`
}

type grantedByUser struct {
	UserPerson  userPerson  `json:"UserPerson"`
	UserCompany userCompany `json:"UserCompany"`
}

type requestedByUser struct {
	UserPerson  userPerson  `json:"UserPerson"`
	UserCompany userCompany `json:"UserCompany"`
}

// Pointer The pointer alias of a monetary account
type Pointer struct {
	PType string  `json:"type,omitempty"`
	Value string  `json:"value,omitempty"`
	Name  *string `json:"name,omitempty"`
}

// ClientContext holds the data that can be used to later on
// recreate the bunq client.
type ClientContext struct {
	PrivateKey           []byte         `json:"private_key"`
	InstallationContext  *installation  `json:"installation_context"`
	SessionServerContext *sessionServer `json:"session_server_context"`
	APIKey               string         `json:"api_key"`
	BaseURL              string         `json:"base_url"`
	UserID               uint           `json:"user_id"`
}

// MonetaryAccountBank The monetary account bank.
type MonetaryAccountBank struct {
	common
	Alias                  []Pointer              `json:"alias"`
	Avatar                 avatar                 `json:"avatar"`
	Balance                Amount                 `json:"balance"`
	Country                string                 `json:"country"`
	Currency               string                 `json:"currency"`
	DailyLimit             Amount                 `json:"daily_limit"`
	DailySpent             Amount                 `json:"daily_spent"`
	Description            string                 `json:"description"`
	PublicUUID             string                 `json:"public_uuid"`
	Status                 string                 `json:"status"`
	SubStatus              string                 `json:"sub_status"`
	Timezone               string                 `json:"timezone"`
	UserID                 int                    `json:"user_id"`
	MonetaryAccountProfile monetaryAccountProfile `json:"monetary_account_profile"`
	NotificationFilters    []notificationFilter   `json:"notification_filters"`
	Setting                monetaryAccountSetting `json:"setting"`
	OverdraftLimit         Amount                 `json:"overdraft_limit"`
}

// GetIBANPointer returns the IBAN Pointer for the given MA.
func (m *MonetaryAccountBank) GetIBANPointer() *Pointer {
	return getIBANPointer(m.Alias)
}

func getIBANPointer(allP []Pointer) *Pointer {
	for _, p := range allP {
		if p.PType == "IBAN" {
			return &p
		}
	}

	return nil
}

type monetaryAccountProfile struct {
	ProfileFill           interface{} `json:"profile_fill"`
	ProfileDrain          interface{} `json:"profile_drain"`
	ProfileActionRequired string      `json:"profile_action_required"`
	ProfileAmountRequired Amount      `json:"profile_amount_required"`
}

type monetaryAccountSetting struct {
	Color               string `json:"color"`
	DefaultAvatarStatus string `json:"default_avatar_status"`
	RestrictionChat     string `json:"restriction_chat"`
}

type draftPayment struct {
	common
	MonetaryAccountID            int                 `json:"monetary_account_id"`
	Status                       string              `json:"status"`
	Type                         string              `json:"type"`
	UserAliasCreated             labelUser           `json:"user_alias_created"`
	Responses                    interface{}         `json:"responses"`
	Entries                      []draftPaymentEntry `json:"entries"`
	Object                       interface{}         `json:"object"`
	RequestReferenceSplitTheBill []interface{}       `json:"request_reference_split_the_bill"`
}

type labelUser struct {
	UUID           string `json:"uuid"`
	DisplayName    string `json:"display_name"`
	Country        string `json:"country"`
	Avatar         avatar `json:"avatar"`
	PublicNickName string `json:"public_nick_name"`
}

type draftPaymentEntry struct {
	Amount            Amount                      `json:"Amount"`
	Alias             LabelMonetaryAccount        `json:"alias"`
	CounterpartyAlias LabelMonetaryAccount        `json:"counterparty_alias"`
	Description       string                      `json:"description"`
	Type              string                      `json:"type"`
	Attachment        []monetaryAccountAttachment `json:"monetaryAccountAttachment"`
	MerchantReference string                      `json:"merchant_reference"`
}

type LabelMonetaryAccount struct {
	IBAN                      string    `json:"iban"`
	IsLight                   bool      `json:"is_light"`
	DisplayName               string    `json:"display_name"`
	Avatar                    avatar    `json:"avatar"`
	LabelUser                 labelUser `json:"label_user"`
	Country                   string    `json:"country"`
	SwiftBic                  string    `json:"swift_bic"`
	SwiftAccountNumber        string    `json:"swift_account_number"`
	TransferwiseAccountNumber string    `json:"transferwise_account_number"`
	TransferwiseBankCode      string    `json:"transferwise_bank_code"`
	BunqMe                    bunqMe    `json:"bunq_me"`
}

type Pagination struct {
	FutureURL string `json:"future_url"`
	NewerURL  string `json:"newer_url"`
	OlderURL  string `json:"older_url"`
}

type masterCardAction struct {
	common
	MonetaryAccountID             int           `json:"monetary_account_id"`
	CardID                        int           `json:"card_id"`
	CardAuthorisationIDResponse   string        `json:"card_authorisation_id_response"`
	AmountLocal                   Amount        `json:"amount_local"`
	AmountConverted               Amount        `json:"amount_converted"`
	AmountBilling                 Amount        `json:"amount_billing"`
	AmountOriginalLocal           Amount        `json:"amount_original_local"`
	AmountOriginalBilling         Amount        `json:"amount_original_billing"`
	AmountFee                     Amount        `json:"amount_fee"`
	Decision                      string        `json:"decision"`
	DecisionDescription           string        `json:"decision_description"`
	DecisionDescriptionTranslated string        `json:"decision_description_translated"`
	Description                   string        `json:"description"`
	AuthorisationStatus           string        `json:"authorisation_status"`
	AuthorisationType             string        `json:"authorisation_type"`
	SettlementStatus              string        `json:"settlement_status"`
	City                          string        `json:"city"`
	Alias                         labelUser     `json:"alias"`
	CounterpartyAlias             labelUser     `json:"counterparty_alias"`
	LabelCard                     labelCard     `json:"label_card"`
	TokenStatus                   string        `json:"token_status"`
	ReservationExpiryTime         string        `json:"reservation_expiry_time"`
	AllowChat                     bool          `json:"allow_chat"`
	PanEntryModeUser              string        `json:"pan_entry_mode_user"`
	EligibleWhitelistID           int           `json:"eligible_whitelist_id"`
	SecureCodeID                  int           `json:"secure_code_id"`
	WalletProviderID              string        `json:"wallet_provider_id"`
	RequestReferenceSplitTheBill  []interface{} `json:"request_reference_split_the_bill"`
	AppliedLimit                  string        `json:"applied_limit"`
}

type labelCard struct {
	UUID       string    `json:"uuid"`
	Type       string    `json:"type"`
	SecondLine string    `json:"second_line"`
	ExpiryDate string    `json:"expiry_date"`
	Status     string    `json:"status"`
	LabelUser  labelUser `json:"label_user"`
}

// MonetaryAccountSaving The monetary account saving.
type MonetaryAccountSaving struct {
	common
	Alias                  []Pointer              `json:"alias"`
	Avatar                 avatar                 `json:"avatar"`
	Balance                Amount                 `json:"balance"`
	Country                string                 `json:"country"`
	Currency               string                 `json:"currency"`
	DailyLimit             Amount                 `json:"daily_limit"`
	DailySpent             Amount                 `json:"daily_spent"`
	Description            string                 `json:"description"`
	PublicUUID             string                 `json:"public_uuid"`
	Status                 string                 `json:"status"`
	SubStatus              string                 `json:"sub_status"`
	Timezone               string                 `json:"timezone"`
	UserID                 int                    `json:"user_id"`
	MonetaryAccountProfile monetaryAccountProfile `json:"monetary_account_profile"`
	NotificationFilters    []notificationFilter   `json:"notification_filters"`
	Setting                monetaryAccountSetting `json:"setting"`
	OverdraftLimit         Amount                 `json:"overdraft_limit"`
	SavingsGoal            Amount                 `json:"savings_goal"`
	SavingsGoalProgress    string                 `json:"savings_goal_progress"`
}

// GetIBANPointer returns the Pointer iban of this ma.
func (s *MonetaryAccountSaving) GetIBANPointer() *Pointer {
	return getIBANPointer(s.Alias)
}

type Payment struct {
	common
	MonetaryAccountID            int                            `json:"monetary_account_id"`
	Amount                       Amount                         `json:"Amount"`
	Alias                        LabelMonetaryAccount           `json:"alias"`
	CounterpartyAlias            LabelMonetaryAccount           `json:"counterparty_alias"`
	Description                  string                         `json:"description"`
	Type                         string                         `json:"type"`
	SubType                      string                         `json:"sub_type"`
	BunqtoStatus                 string                         `json:"bunqto_status"`
	BunqtoSubStatus              string                         `json:"bunqto_sub_status"`
	BunqtoShareURL               string                         `json:"bunqto_share_url"`
	BunqtoExpiry                 string                         `json:"bunqto_expiry"`
	BunqtoTimeResponded          string                         `json:"bunqto_time_responded"`
	Attachment                   []monetaryAccountAttachment    `json:"monetaryAccountAttachment"`
	MerchantReference            string                         `json:"merchant_reference"`
	BatchID                      int                            `json:"batch_id"`
	ScheduledID                  int                            `json:"scheduled_id"`
	AddressShipping              address                        `json:"address_shipping"`
	AddressBilling               address                        `json:"address_billing"`
	Geolocation                  geolocation                    `json:"geolocation"`
	AllowChat                    bool                           `json:"allow_chat"`
	RequestReferenceSplitTheBill []requestReferenceSplitTheBill `json:"request_reference_split_the_bill"`
	BalanceAfterMutation         Amount                         `json:"balance_after_mutation"`
}

// PaymentBatch a batch of payments
type PaymentBatch struct {
	Payments []Payment `json:"payments"`
}

// ScheduledPayment The scheduled payment
type ScheduledPayment struct {
	common
	MonetaryAccountID int                   `json:"monetary_account_id"`
	Payment           scheduledPaymentEntry `json:"payment"`
	Schedule          schedule              `json:"schedule"`
	Status            string                `json:"status"`
}

type scheduledPaymentEntry struct {
	Amount            Amount               `json:"amount"`
	Alias             LabelMonetaryAccount `json:"alias"`
	CounterpartyAlias LabelMonetaryAccount `json:"counterparty_alias"`
	Description       string               `json:"description"`
	MerchantReference string               `json:"merchant_reference"`
	AllowBunqTo       bool                 `json:"allow_bunqto"`
}

type schedule struct {
	TimeStart      string               `json:"time_start"`
	TimeEnd        string               `json:"time_end"`
	RecurrenceUnit string               `json:"recurrence_unit"`
	RecurrenceSize int                  `json:"recurrence_size"`
	Status         string               `json:"status"`
	Object         scheduleAnchorObject `json:"object"`
}

type scheduleAnchorObject struct {
	Payment      Payment      `json:"payment"`
	PaymentBatch PaymentBatch `json:"paymentBatch"`
}

type bunqMe struct {
	Type  string `json:"type"`
	Value string `json:"value"`
	Name  string `json:"name"`
}

type monetaryAccountAttachment struct {
	ID                int `json:"id"`
	MonetaryAccountID int `json:"monetary_account_id"`
}

type geolocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
	Radius    float64 `json:"radius"`
}

type requestReferenceSplitTheBill struct {
	Type string `json:"type"`
	ID   int    `json:"id"`
}

type RequestResponse struct {
	common
	SubType           string               `json:"sub_type"`
	MonetaryAccountID int                  `json:"monetary_account_id"`
	Amount            Amount               `json:"Amount"`
	AmountResponded   Amount               `json:"amount_responded"`
	AmountInquired    Amount               `json:"amount_inquired"`
	Alias             LabelMonetaryAccount `json:"alias"`
	CounterpartyAlias LabelMonetaryAccount `json:"counterparty_alias"`
	Description       string               `json:"description"`
	CreditSchemeID    string               `json:"credit_scheme_identifier"`
	MandateID         string               `json:"mandate_identifier"`
	Responded         string               `json:"time_responded"`
}
