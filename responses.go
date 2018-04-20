package sip2

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"
)

var ResponseMap = map[string]reflect.Type{
	"24": reflect.TypeOf(PatronStatusResponse{}),
	"12": reflect.TypeOf(CheckoutResponse{}),
	"10": reflect.TypeOf(CheckinResponse{}),
	"98": reflect.TypeOf(ACSStatusResponse{}),
	"96": reflect.TypeOf(RequestSCResendResponse{}),
	"94": reflect.TypeOf(LoginResponse{}),
	"64": reflect.TypeOf(PatronInformationResponse{}),
	"36": reflect.TypeOf(EndSessionResponse{}),
	"38": reflect.TypeOf(FeePaidResponse{}),
	"18": reflect.TypeOf(ItemInformationResponse{}),
	"20": reflect.TypeOf(ItemStatusUpdateResponse{}),
	"26": reflect.TypeOf(PatronEnableResponse{}),
	"16": reflect.TypeOf(HoldResponse{}),
	"30": reflect.TypeOf(RenewResponse{}),
	"66": reflect.TypeOf(RenewAllResponse{}),
}

type PatronStatusResponse struct {
	PatronStatus        `json:"patron_status"`
	Language            `json:"language"`
	TransactionDate     `json:"transaction_date"`
	InstitutionID       `json:"institution_id"`
	PatronID            `json:"patron_id"`
	PersonalName        `json:"personal_name"`
	ValidPatron         `json:"valid_patron"`
	ValidPatronPassword `json:"valid_patron_password"`
	ScreenMessage       `json:"screen_message"`
	PrintLine           `json:"print_line"`
}

type CheckoutResponse struct {
	OK              `json:"ok"`
	RenewalOK       `json:"renewal_ok"`
	MagneticMedia   `json:"magnetic_media"`
	Desensitize     `json:"desensitize"`
	TransactionDate `json:"transaction_date"`
	InstitutionID   `json:"institution_id"`
	PatronID        `json:"patron_id"`
	ItemID          `json:"item_id"`
	TitleID         `json:"title_id"`
	DueDate         `json:"due_date"`
	FeeType         `json:"fee_type"`
	SecurityInhibit `json:"security_inhibit"`
	CurrencyType    `json:"currency_type"`
	FeeAmount       `json:"fee_amount"`
	MediaType       `json:"media_type"`
	ItemProperties  `json:"item_properties"`
	TransactionID   `json:"transaction_id"`
	ScreenMessage   `json:"screen_message"`
	PrintLine       `json:"print_line"`
}

type CheckinResponse struct {
	OK                `json:"ok"`
	Resensitize       `json:"resensitize"`
	MagneticMedia     `json:"magnetic_media"`
	Alert             `json:"alert"`
	TransactionDate   `json:"transaction_date"`
	InstitutionID     `json:"institution_id"`
	ItemID            `json:"item_id"`
	PermanentLocation `json:"permanent_location"`
	TitleID           `json:"title_id"`
	SortBin           `json:"sort_bin"`
	PatronID          `json:"patron_id"`
	MediaType         `json:"media_type"`
	ItemProperties    `json:"item_properties"`
	ScreenMessage     `json:"screen_message"`
	PrintLine         `json:"print_line"`
}

type ACSStatusResponse struct {
	OnlineStatus      `json:"online_status"`
	CheckinOK         `json:"checkin_ok"`
	CheckoutOK        `json:"checkout_ok"`
	ACSRenewalPolicy  `json:"ACS_renewal_policy"`
	StatusUpdateOK    `json:"status_update_ok"`
	OfflineOK         `json:"offline_ok"`
	TimeoutPeriod     `json:"timeout_period"`
	RetriesAllowed    `json:"retrise_allowed"`
	DateTimeSync      `json:"date_time_sync"`
	ProtocolVersion   `json:"protocal_version"`
	InstitutionID     `json:"institution_id"`
	LibraryName       `json:"library_name"`
	SupportedMessages `json:"supported_message"`
	TerminalLocation  `json:"terminal_location"`
	ScreenMessage     `json:"screen_message"`
	PrintLine         `json:"print_line"`
}

type RequestSCResendResponse struct{}

type LoginResponse struct {
	*OK `json:"ok"`
}

type EndSessionResponse struct {
	EndSession      `json:"end_session"`
	TransactionDate `json:"transaction_date"`
	InstitutionID   `json:"institution_id"`
	PatronID        `json:"patron_id"`
	ScreenMessage   `json:"screen_message"`
	PrintLine       `json:"print_line"`
}

type FeePaidResponse struct {
	PaymentAccepted `json:"payment_accepted"`
	TransactionDate `json:"transaction_date"`
	InstitutionID   `json:"institution_id"`
	PatronID        `json:"patron_id"`
	TransactionID   `json:"transaction_id"`
	ScreenMessage   `json:"screen_message"`
	PrintLine       `json:"print_line"`
}

type ItemInformationResponse struct {
	CirculationStatus `json:"circulation_status"`
	SecurityMarker    `json:"security_maker"`
	FeeType           `json:"fee_type"`
	TransactionDate   `json:"transaction_date"`
	HoldQueueLength   `json:"hold_queue_length"`
	DueDate           `json:"due_date"`
	RecallDate        `json:"recall_date"`
	HoldPickupDate    `json:"hold_picked_date"`
	ItemID            `json:"item_id"`
	TitleID           `json:"title_id"`
	Author            `json:"author"`
	ISBN              `json:"isbn"`
	Owner             `json:"owner"`
	CurrencyType      `json:"currency_type"`
	FeeAmount         `json:"fee_amount"`
	MediaType         `json:"media_type"`
	PermanentLocation `json:"permanent_location"`
	CurrentLocation   `json:"current_location"`
	ItemProperties    `json:"item_properties"`
	ScreenMessage     `json:"screen_message"`
	PrintLine         `json:"print_line"`
	Publisher         `json:"publisher"`
}

type ItemStatusUpdateResponse struct {
	ItemPropertiesOK `json:"item_properties_ok"`
	TransactionDate  `json:"transaction_date"`
	ItemID           `json:"item_id"`
	TitleID          `json:"title_id"`
	ItemProperties   `json:"item_properties"`
	ScreenMessage    `json:"screen_message"`
	PrintLine        `json:"print_line"`
}

type PatronEnableResponse struct {
	PatronStatus        `json:"patron_status"`
	Language            `json:"language"`
	TransactionDate     `json:"transaction_date"`
	InstitutionID       `json:"institution_id"`
	PatronID            `json:"patron_id"`
	PersonalName        `json:"personal_name"`
	ValidPatron         `json:"valid_patron"`
	ValidPatronPassword `json:"valid_patron_password"`
	ScreenMessage       `json:"screen_message"`
	PrintLine           `json:"print_line"`
}

type HoldResponse struct {
	OK              `json:"ok"`
	TransactionDate `json:"transaction_date"`
	ExpirationDate  `json:"expiration_date"`
	QueuePosition   `json:"queue_position"`
	PickupLocation  `json:"pickup_location"`
	InstitutionID   `json:"institution_id"`
	PatronID        `json:"patron_id"`
	ItemID          `json:"item_id"`
	TitleID         `json:"title_id"`
	ScreenMessage   `json:"screen_message"`
	PrintLine       `json:"print_line"`
}

type RenewResponse struct {
	OK              `json:"ok"`
	RenewalOK       `json:"renewal_ok"`
	MagneticMedia   `json:"magnetic_media"`
	Desensitize     `json:"desensitize"`
	TransactionDate `json:"transaction_date"`
	InstitutionID   `json:"institution_id"`
	PatronID        `json:"patron_id"`
	TitleID         `json:"title_id"`
	DueDate         `json:"due_date"`
	FeeType         `json:"fee_type"`
	SecurityInhibit `json:"security_inhibit"`
	CurrencyType    `json:"currency_type"`
	FeeAmount       `json:"fee_amount"`
	MediaType       `json:"media_type"`
	ItemProperties  `json:"item_properties"`
	TransactionID   `json:"transaction_id"`
	ScreenMessage   `json:"screen_message"`
	PrintLine       `json:"print_line"`
}

type RenewAllResponse struct {
	OK              `json:"ok"`
	RenewedCount    `json:"renewed_count"`
	UnrenewedCount  `json:"unrenewed_count"`
	TransactionDate `json:"transaction_date"`
	InstitutionID   `json:"institution_id"`
	RenewedItems    `json:"renewed_items"`
	UnrenewedItems  `json:"unrenewed_items"`
	ScreenMessage   `json:"screen_message"`
	PrintLine       `json:"print_line"`
}

type PatronInformationResponse struct {
	PatronStatus          `json:"patron_status"`
	Language              `json:"language"`
	TransactionDate       `json:"transaction_date"`
	HoldItemsCount        `json:"hold_item_count"`
	OverdueItemsCount     `json:"overdue_item_count"`
	ChargedItemsCount     `json:"charged_item_count"`
	FineItemsCount        `json:"fine_items_count"`
	RecallItemsCount      `json:"recall_items_count"`
	UnavailableHoldsCount `json:"unavailable_holds_count"`
	InstitutionID         `json:"institution_id"`
	PatronID              `json:"patron_id"`
	PersonalName          `json:"personal_name"`
	HoldQueueLength       `json:"hold_queue_length"`
	OverdueItemsLimit     `json:"overdue_items_limit"`
	ChargedItemsLimit     `json:"charged_items_limit"`
	ValidPatron           `json:"valid_patron"`
	ValidPatronPassword   `json:"valid_patron_password"`
	CurrencyType          `json:"currency_type"`
	FeeAmount             `json:"fee_amount"`
	FeeLimit              `json:"fee_limit"`
	HoldItemsLimit        `json:"hold_items_limit"`
	HoldItems             `json:"hold_items"`
	StartItem             `json:"start_item"`
	RenewedItems          `json:"renewed_items"`
	EmailAddress          `json:"email_address"`
	HomeAddress           `json:"home_address"`
	ScreenMessage         `json:"screen_message"`
	PrintLine             `json:"print_line"`
}

// func (p *ClientPool) DecodeResponse(b []byte) (interface{}, error) {
// 	reader := bytes.NewReader(b)
// 	commandID := make([]byte, 2)
// 	_, err := reader.Read(commandID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp, err := GenResponse(string(commandID))
// 	if err != nil {
// 		return nil, err
// 	}
// 	respVal := reflect.ValueOf(resp).Elem()
// 	for i := 0; i < respVal.NumField(); i++ {
// 		field := respVal.Field(i).Interface().(SipField)
// 		id, _, length := field.Info()
// 		field.Decode(reader, id, length)
// 	}
// 	return resp, nil
// }

func classifyFields(resp interface{}) ([]SipField, map[string]SipField) {
	fixedFields := make([]SipField, 0, 16)
	variableFields := make(map[string]SipField)
	val := reflect.ValueOf(resp).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i).Interface().(SipField)
		id, _, length := field.Info()
		if length != -1 && id == "" {
			fixedFields = append(fixedFields, field)
		} else {
			variableFields[id] = field
		}
	}
	return fixedFields, variableFields
}

func decodeVarFields(r *bytes.Reader, varFieldsMap map[string]SipField) error {
	bs, _ := ioutil.ReadAll(r)
	fieldBytes := bytes.Split(bs, []byte("|"))
	for _, fb := range fieldBytes {
		if field, ok := varFieldsMap[string(fb[:2])]; ok {
			reader := bytes.NewReader(append(fb, '|'))
			id, _, length := field.Info()
			err := field.Decode(reader, id, length)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *ClientPool) DecodeResponse(b []byte) (interface{}, error) {
	reader := bytes.NewReader(b)
	commandID := make([]byte, 2)
	_, err := reader.Read(commandID)
	if err != nil {
		return nil, err
	}
	resp, err := GenResponse(string(commandID))
	if err != nil {
		return nil, err
	}
	fixed, variable := classifyFields(resp)
	for _, field := range fixed {
		id, _, length := field.Info()
		err = field.Decode(reader, id, length)
		if err != nil {
			return nil, err
		}
	}
	err = decodeVarFields(reader, variable)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GenResponse(commandID string) (interface{}, error) {
	respType, ok := ResponseMap[commandID]
	if !ok {
		return nil, fmt.Errorf("DecodeResponse: %s response not exist", respType)
	}
	resp := reflect.New(respType).Interface()
	InitResponse(resp)
	return resp, nil
}

func InitResponse(resp interface{}) {
	val := reflect.ValueOf(resp).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i).Field(0)
		fieldType := field.Type().Elem()
		field.Set(reflect.New(fieldType))
	}
}
