package sip2

import (
	"bytes"
	"reflect"
)

type LANG int

const (
	Unknown LANG = iota
	English
	French
	German
	Italian
	Dutch
	Swedish
	Finnish
	Spanish
	Danish
	Portuguese
	CanadianFrench
	Norwegian
	Hebrew
	Japanese
	Russian
	Arabic
	Polish
	Greek
	Chinese
	Korean
	NorthAmericanSpanish
	Tamil
	Malay
	UnitedKingdom
	Icelandic
	Belgian
	Taiwanese
)

type CURRENCY string

const (
	USD CURRENCY = "USD"
	CAD CURRENCY = "CAD"
	GBP CURRENCY = "GBP"
	FRF CURRENCY = "FRF"
	DEM CURRENCY = "DEM"
	ITL CURRENCY = "ITL"
	ESP CURRENCY = "ESP"
	JPY CURRENCY = "JPY"
)

type PatronStatusRequest struct {
	CommandID        `json:"command_id"`
	Language         `json:"language"`
	TransactionDate  `json:"transaction_date"`
	InstitutionID    `json:"institution_id"`
	PatronID         `json:"patron_id"`
	TerminalPassword `json:"terminal_password"`
	PatronPassword   `json:"patron_password"`
}

func NewPatronStatusRequest() *PatronStatusRequest {
	req := &PatronStatusRequest{}
	InitRequest(req)
	*(req.CommandID.StrValue) = StrValue("23")
	return req
}

type PatronInformationRequest struct {
	CommandID        `json:"command_id"`
	Language         `json:"language"`
	TransactionDate  `json:"transaction_date"`
	Summary          `json:"summary"`
	InstitutionID    `json:"institution_id"`
	PatronID         `json:"patron_id"`
	TerminalPassword `json:"terminal_password"`
	PatronPassword   `json:"patron_password"`
	StartItem        `json:"start_item"`
	EndItem          `json:"end_item"`
}

func NewPatronInformationRequest() *PatronInformationRequest {
	req := &PatronInformationRequest{}
	InitRequest(req)
	*(req.CommandID.StrValue) = StrValue("63")
	return req
}

type ItemInformationRequest struct {
	CommandID        `json:"command_id"`
	TransactionDate  `json:"transaction_date"`
	InstitutionID    `json:"institution_id"`
	ItemID           `json:"item_id"`
	TerminalPassword `json:"terminal_password"`
}

func NewItemInformationRequest() *ItemInformationRequest {
	req := &ItemInformationRequest{}
	InitRequest(req)
	*(req.CommandID.StrValue) = StrValue("17")
	return req
}

type CheckoutRequest struct {
	CommandID        `json:"command_id"`
	SCRenewalPolicy  `json:"sc_renewal_policy"`
	NoBlock          `json:"no_block"`
	TransactionDate  `json:"transaction_date"`
	NBDueDate        `json:"nb_due_date"`
	InstitutionID    `json:"institution_id"`
	PatronID         `json:"patron_id"`
	ItemID           `json:"item_id"`
	TerminalPassword `json:"terminal_password"`
	ItemProperties   `json:"item_properites"`
	PatronPassword   `json:"patron_password"`
	FeeAcknowledged  `json:"fee_acknowledged"`
	Cancel           `json:"cancel"`
}

func NewCheckoutRequest() *CheckoutRequest {
	req := &CheckoutRequest{}
	InitRequest(req)
	*(req.CommandID.StrValue) = StrValue("11")
	return req
}

type CheckinRequest struct {
	CommandID        `json:"command_id"`
	NoBlock          `json:"no_block"`
	TransactionDate  `json:"transaction_date"`
	ReturnDate       `json:"return_date"`
	CurrentLocation  `json:"current_location"`
	InstitutionID    `json:"institution_id"`
	ItemID           `json:"item_id"`
	TerminalPassword `json:"terminal_password"`
	ItemProperties   `json:"item_properties"`
	Cancel           `json:"cancel"`
}

func NewCheckinRequest() *CheckinRequest {
	req := &CheckinRequest{}
	InitRequest(req)
	*(req.CommandID.StrValue) = StrValue("09")
	return req
}

type BlockPatronRequest struct {
	CommandID        `json:"command_id"`
	CardRetained     `json:"card_retained"`
	TransactionDate  `json:"transaction_date"`
	InstitutionID    `json:"institution_id"`
	BlockedCardMsg   `json:"blocked_card_msg"`
	PatronID         `json:"patron_id"`
	TerminalPassword `json:"terminal_password"`
}

func NewBlockPatronRequest() *BlockPatronRequest {
	req := &BlockPatronRequest{}
	InitRequest(req)
	*(req.CommandID.StrValue) = StrValue("01")
	return req
}

type SCStatusRequest struct {
	CommandID       `json:"command_id"`
	StatusCode      `json:"status_code"`
	MaxPrintWidth   `json:"max_print_width"`
	ProtocolVersion `json:"protocal_version"`
}

func NewSCStatusRequest() *SCStatusRequest {
	req := &SCStatusRequest{}
	InitRequest(req)
	*(req.CommandID.StrValue) = StrValue("99")
	return req
}

type LoginRequest struct {
	CommandID     `json:"command_id"`
	UIDAlgorithm  `json:"uid_algorithm"`
	PWDAlgorithm  `json:"pwd_algorithm"`
	LoginUserID   `json:"login_user_id"`
	LoginPassword `json:"login_password"`
	LocationCode  `json:"location_code"`
}

func NewLoginRequest() *LoginRequest {
	req := &LoginRequest{}
	InitRequest(req)
	*(req.CommandID.StrValue) = StrValue("93")
	return req
}

type ResendRequest struct {
	*CommandID `json:"command_id"`
}

func NewResendRequest() *ResendRequest {
	req := &ResendRequest{}
	InitRequest(req)
	*(req.CommandID.StrValue) = StrValue("97")
	return req
}

type EndPatronSessionRequest struct {
	CommandID        `json:"command_id"`
	TransactionDate  `json:"transaction_date"`
	InstitutionID    `json:"institution_id"`
	PatronID         `json:"patron_id"`
	TerminalPassword `json:"terminal_password"`
	PatronPassword   `json:"patron_password"`
}

func NewEndPatronSessionRequest() *EndPatronSessionRequest {
	req := &EndPatronSessionRequest{}
	InitRequest(req)
	*(req.CommandID.StrValue) = StrValue("35")
	return req
}

type FeePaidRequest struct {
	CommandID        `json:"command_id"`
	TransactionDate  `json:"transaction_date"`
	FeeType          `json:"fee_type"`
	PaymentType      `json:"payment_type"`
	CurrencyType     `json:"currency_type"`
	FeeAmount        `json:"fee_amount"`
	InstitutionID    `json:"institution_id"`
	PatronID         `json:"patron_id"`
	TerminalPassword `json:"terminal_password"`
	FeeID            `json:"fee_id"`
	TransactionID    `json:"transaction_id"`
}

func NewFeePaidRequest() *FeePaidRequest {
	req := &FeePaidRequest{}
	InitRequest(req)
	*(req.CommandID.StrValue) = StrValue("37")
	return req
}

type ItemStatusUpdateRequest struct {
	CommandID        `json:"command_id"`
	TransactionDate  `json:"transaction_date"`
	InstitutionID    `json:"institution_id"`
	ItemID           `json:"item_id"`
	TerminalPassword `json:"terminal_password"`
	ItemProperties   `json:"item_properties"`
}

func NewItemStatusUpdateRequest() *ItemStatusUpdateRequest {
	req := &ItemStatusUpdateRequest{}
	InitRequest(req)
	*(req.CommandID.StrValue) = StrValue("19")
	return req
}

type PatronEnableRequest struct {
	CommandID        `json:"command_id"`
	TransactionDate  `json:"transaction_date"`
	InstitutionID    `json:"institution_id"`
	PatronID         `json:"patron_id"`
	TerminalPassword `json:"terminal_password"`
	PatronPassword   `json:"patron_password"`
}

func NewPatronEnableRequest() *PatronEnableRequest {
	req := &PatronEnableRequest{}
	InitRequest(req)
	*(req.CommandID.StrValue) = StrValue("25")
	return req
}

type HoldRequest struct {
	CommandID        `json:"command_id"`
	TransactionDate  `json:"transaction_date"`
	ExpirationDate   `json:"expiration_date"`
	PickupLocation   `json:"pickup_location"`
	HoldType         `json:"hold_type"`
	InstitutionID    `json:"institution_id"`
	PatronID         `json:"patron_id"`
	PatronPassword   `json:"patron_password"`
	ItemID           `json:"item_id"`
	TitleID          `json:"title_id"`
	TerminalPassword `json:"terminal_password"`
	FeeAcknowledged  `json:"fee_acknowledged"`
}

func NewHoldRequest() *HoldRequest {
	req := &HoldRequest{}
	InitRequest(req)
	*(req.CommandID.StrValue) = StrValue("15")
	return req
}

type RenewRequest struct {
	CommandID         `json:"command_id"`
	ThirdPartyAllowed `json:"third_party_allowed"`
	NoBlock           `json:"no_block"`
	TransactionDate   `json:"transaction_date"`
	NBDueDate         `json:"nb_due_date"`
	InstitutionID     `json:"institution_id"`
	PatronID          `json:"patron_id"`
	PatronPassword    `json:"patron_password"`
	ItemID            `json:"item_id"`
	TitleID           `json:"title_id"`
	TerminalPassword  `json:"terminal_password"`
	ItemProperties    `json:"item_properties"`
	FeeAcknowledged   `json:"fee_acknowledged"`
}

func NewRenewRequest() *RenewRequest {
	req := &RenewRequest{}
	InitRequest(req)
	*(req.CommandID.StrValue) = StrValue("29")
	return req
}

type RenewAllRequest struct {
	CommandID        `json:"command_id"`
	TransactionDate  `json:"transaction_date"`
	InstitutionID    `json:"institution_id"`
	PatronID         `json:"patron_id"`
	PatronPassword   `json:"patron_password"`
	TerminalPassword `json:"terminal_password"`
	FeeAcknowledged  `json:"fee_acknowledged"`
}

func NewRenewAllRequest() *RenewAllRequest {
	req := &RenewAllRequest{}
	InitRequest(req)
	*(req.CommandID.StrValue) = StrValue("65")
	return req
}

func EncodeRequest(req interface{}) ([]byte, error) {
	val := reflect.ValueOf(req).Elem()
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i).Interface().(SipField)
		id, _, length := field.Info()
		buffer.Write(field.Encode(id, length))
	}
	buffer.WriteString("AY0AZ")
	buffer.WriteString(genChecksum(buffer.Bytes()))
	return buffer.Bytes(), nil
}

func InitRequest(req interface{}) {
	val := reflect.ValueOf(req).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i).Field(0)
		fieldType := field.Type().Elem()
		field.Set(reflect.New(fieldType))
	}
}
