package sip2

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type SipField interface {
	Info() (id, name string, length int)
	Encode(string, int) []byte
	Decode(*bytes.Reader, string, int) error
}

type StrValue string

func (sv *StrValue) Encode(id string, length int) []byte {
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	buffer.WriteString(id)
	if length == -1 {
		buffer.WriteString(string(*sv) + "|")
	} else {
		b := make([]byte, length)
		template := fmt.Sprintf("%%%dv", length)
		copy(b, fmt.Sprintf(template, *sv))
		buffer.Write(b)
	}
	return buffer.Bytes()
}

func (sv *StrValue) Decode(r *bytes.Reader, id string, length int) error {
	err := checkID(r, id)
	if err != nil {
		return err
	}
	content, err := readContent(r, length)
	if err != nil {
		return err
	}
	*sv = StrValue(content)
	return nil
}

func (sv *StrValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(*sv))
}

func (sv *StrValue) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*sv = StrValue(s)
	return nil
}

type BoolValue bool

func (bv *BoolValue) Encode(id string, length int) []byte {
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	buffer.WriteString(id)
	if *bv {
		buffer.WriteString("Y|")
	} else {
		buffer.WriteString("N|")
	}
	return buffer.Bytes()
}

func (bv *BoolValue) Decode(r *bytes.Reader, id string, length int) error {
	err := checkID(r, id)
	if err != nil {
		return err
	}
	content, err := readN(r, 2)
	if err != nil {
		return err
	}
	if content[0] == 'Y' {
		*bv = BoolValue(true)
	} else {
		*bv = BoolValue(false)
	}
	return nil
}

func (bv *BoolValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(bool(*bv))
}

func (bv *BoolValue) UnmarshalJSON(b []byte) error {
	var bo bool
	err := json.Unmarshal(b, &bo)
	if err != nil {
		return err
	}
	*bv = BoolValue(bo)
	return nil
}

type IntValue int

func (iv *IntValue) Encode(id string, length int) []byte {
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	buffer.WriteString(id)
	if length == -1 {
		buffer.WriteString(fmt.Sprintf("%d", *iv) + "|")
	} else {
		b := make([]byte, length)
		template := fmt.Sprintf("%%0%dd", length)
		copy(b, []byte(fmt.Sprintf(template, *iv)))
		buffer.Write(b)
	}
	return buffer.Bytes()
}

func (iv *IntValue) Decode(r *bytes.Reader, id string, length int) error {
	err := checkID(r, id)
	if err != nil {
		return err
	}
	content, err := readContent(r, length)
	if err != nil {
		return err
	}
	i64, err := strconv.ParseInt(string(content), 10, 64)
	if err != nil {
		return err
	}
	*iv = IntValue(i64)
	return nil
}

func (iv *IntValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(*iv))
}

func (iv *IntValue) UnmarshalJSON(b []byte) error {
	var i int
	err := json.Unmarshal(b, &i)
	if err != nil {
		return err
	}
	*iv = IntValue(i)
	return nil
}

type FloatValue float64

func (fv *FloatValue) Encode(id string, length int) []byte {
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	buffer.WriteString(id)
	buffer.WriteString(fmt.Sprintf("%f", *fv) + "|")
	return buffer.Bytes()
}

func (fv *FloatValue) Decode(r *bytes.Reader, id string, length int) error {
	err := checkID(r, id)
	if err != nil {
		return err
	}
	content, err := readContent(r, length)
	if err != nil {
		return err
	}
	v, err := strconv.ParseFloat(string(content), 10)
	if err != nil {
		return err
	}
	*fv = FloatValue(v)
	return nil
}

func (fv *FloatValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(*fv))
}

func (fv *FloatValue) UnmarshalJSON(b []byte) error {
	var f float64
	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}
	*fv = FloatValue(f)
	return nil
}

type TimeValue time.Time

func (tv *TimeValue) Encode(id string, length int) []byte {
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	buffer.WriteString(id)
	buffer.WriteString(time.Time(*tv).Format("20060102    150405"))
	return buffer.Bytes()
}

func (tv *TimeValue) Decode(r *bytes.Reader, id string, length int) error {
	err := checkID(r, id)
	if err != nil {
		return err
	}
	content, err := readContent(r, length)
	if err != nil {
		return err
	}
	timeVal, err := time.Parse("20060102    150405", string(content))
	if err != nil {
		return errors.New("TimeValue Decode: " + err.Error())
	}
	*tv = TimeValue(timeVal)
	return nil
}

func (tv *TimeValue) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, time.Time(*tv).Format("2006-01-02 15:04:05"))), nil
}

func (tv *TimeValue) UnmarshalJSON(b []byte) error {
	t, err := time.Parse("2006-01-02 15:04:05", string(b))
	if err != nil {
		return err
	}
	*tv = TimeValue(t)
	return nil
}

type StrSliceValue []string

func (ssv *StrSliceValue) Encode(id string, length int) []byte {
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	buffer.WriteString(id)
	buffer.WriteString(strings.Join(([]string)(*ssv), ",") + "|")
	return buffer.Bytes()
}

func (ssv *StrSliceValue) Decode(r *bytes.Reader, id string, length int) error {
	err := checkID(r, id)
	if err != nil {
		return err
	}
	content, err := readContent(r, length)
	if err != nil {
		return err
	}
	*ssv = strings.Split(string(content), ",")
	return nil
}

func (ssv *StrSliceValue) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string(*ssv))
}

func (ssv *StrSliceValue) UnmarshalJSON(b []byte) error {
	ss := make([]string, 0, 16)
	err := json.Unmarshal(b, ss)
	if err != nil {
		return err
	}
	*ssv = StrSliceValue(ss)
	return nil
}

type CommandID struct {
	*StrValue
}

func (cid CommandID) Info() (id, name string, length int) {
	return "", "command_id", 2
}

type Language struct {
	*IntValue
}

func (l Language) Info() (id, name string, length int) {
	return "", "language_id", 3
}

type TransactionDate struct {
	*TimeValue
}

func (td TransactionDate) Info() (id, name string, length int) {
	return "", "transaction_date", 18
}

type InstitutionID struct {
	*StrValue
}

func (iid InstitutionID) Info() (id, name string, length int) {
	return "AO", "institution_id", -1
}

type PatronID struct {
	*StrValue
}

func (pid PatronID) Info() (id, name string, length int) {
	return "AA", "patron_id", -1
}

type TerminalPassword struct {
	*StrValue
}

func (tp TerminalPassword) Info() (id, name string, length int) {
	return "AC", "terminal_password", -1
}

type PatronPassword struct {
	*StrValue
}

func (pp PatronPassword) Info() (id, name string, length int) {
	return "AD", "patron_password", -1
}

type FeeType struct {
	*IntValue
}

func (ft FeeType) Info() (id, name string, length int) {
	return "", "fee type", 2
}

type PaymentType struct {
	*IntValue
}

func (pt PaymentType) Info() (id, name string, length int) {
	return "", "payment_type", 2
}

type CurrencyType struct {
	*StrValue
}

func (ct CurrencyType) Info() (id, name string, length int) {
	return "BH", "currency_type", 3
}

type FeeAmount struct {
	*FloatValue
}

func (fa FeeAmount) Info() (id, name string, length int) {
	return "BV", "fee_amount", -1
}

type FeeID struct {
	*StrValue
}

func (fid FeeID) Info() (id, name string, length int) {
	return "CG", "fee_id", -1
}

type TransactionID struct {
	*StrValue
}

func (tid TransactionID) Info() (id, name string, length int) {
	return "BK", "transaction_id", -1
}

type ItemID struct {
	*StrValue
}

func (iid ItemID) Info() (id, name string, length int) {
	return "AB", "item_id", -1
}

type UIDAlgorithm struct {
	*IntValue
}

func (ua UIDAlgorithm) Info() (id, name string, length int) {
	return "", "uid_algorithm", 1
}

type PWDAlgorithm struct {
	*IntValue
}

func (pa PWDAlgorithm) Info() (id, name string, length int) {
	return "", "pwd_algorithm", 1
}

type LoginUserID struct {
	*StrValue
}

func (luid LoginUserID) Info() (id, name string, length int) {
	return "CN", "login_user_id", -1
}

type LocationCode struct {
	*StrValue
}

func (lc LocationCode) Info() (id, name string, length int) {
	return "CP", "location_code", -1
}

type ItemProperties struct {
	*StrSliceValue
}

func (ip ItemProperties) Info() (id, name string, length int) {
	return "CH", "item_properties", -1
}

type LoginPassword struct {
	*StrValue
}

func (lp LoginPassword) Info() (id, name string, length int) {
	return "CO", "login_password", -1
}

type StatusCode struct {
	*IntValue
}

func (sc StatusCode) Info() (id, name string, length int) {
	return "", "status_code", 1
}

type MaxPrintWidth struct {
	*IntValue
}

func (mpw MaxPrintWidth) Info() (id, name string, length int) {
	return "", "max_print_width", 3
}

type ProtocolVersion struct {
	*StrValue
}

func (pv ProtocolVersion) Info() (id, name string, length int) {
	return "", "protocal_version", 4
}

type CardRetained struct {
	*BoolValue
}

func (cr CardRetained) Info() (id, name string, length int) {
	return "", "card_retained", 1
}

type BlockedCardMsg struct {
	*StrValue
}

func (bcm BlockedCardMsg) Info() (id, name string, length int) {
	return "AL", "blocked_card_msg", -1
}

type NoBlock struct {
	*BoolValue
}

func (nb NoBlock) Info() (id, name string, length int) {
	return "", "no_block", 1
}

type ReturnDate struct {
	*TimeValue
}

func (rd ReturnDate) Info() (id, name string, length int) {
	return "", "return_date", 18
}

type CurrentLocation struct {
	*StrValue
}

func (cl CurrentLocation) Info() (id, name string, length int) {
	return "AP", "current_location", -1
}

type Summary struct {
	*StrValue
}

func (s Summary) Info() (id, name string, length int) {
	return "", "summary", 10
}

type StartItem struct {
	*IntValue
}

func (si StartItem) Info() (id, name string, length int) {
	return "BP", "start_item", -1
}

type EndItem struct {
	*IntValue
}

func (ei EndItem) Info() (id, name string, length int) {
	return "BQ", "end_item", -1
}

type Cancel struct {
	*BoolValue
}

func (c Cancel) Info() (id, name string, length int) {
	return "BI", "cancel", 1
}

type SCRenewalPolicy struct {
	*BoolValue
}

func (scrp SCRenewalPolicy) Info() (id, name string, length int) {
	return "", "sc_renewal_policy", 1
}

type NBDueDate struct {
	*TimeValue
}

func (nbdd NBDueDate) Info() (id, name string, length int) {
	return "", "nb_due_date", 18
}

type FeeAcknowledged struct {
	*BoolValue
}

func (fa FeeAcknowledged) Info() (id, name string, length int) {
	return "BO", "fee_acknowledged", 1
}

type ExpirationDate struct {
	*TimeValue
}

func (ed ExpirationDate) Info() (id, name string, length int) {
	return "BW", "expiration_date", 18
}

type PickupLocation struct {
	*StrValue
}

func (pl PickupLocation) Info() (id, name string, length int) {
	return "BS", "pickup_location", -1
}

type HoldType struct {
	*IntValue
}

func (ht HoldType) Info() (id, name string, length int) {
	return "BY", "hold_type", 1
}

type TitleID struct {
	*StrValue
}

func (tid TitleID) Info() (id, name string, length int) {
	return "AJ", "title_id", -1
}

type ThirdPartyAllowed struct {
	*BoolValue
}

func (tpa ThirdPartyAllowed) Info() (id, name string, length int) {
	return "", "third_party_allowed", 1
}

type PatronStatus struct {
	*StrValue
}

func (ps PatronStatus) Info() (id, name string, length int) {
	return "", "patron_status", 14
}

type PersonalName struct {
	*StrValue
}

func (pn PersonalName) Info() (id, name string, length int) {
	return "AE", "personal_name", -1
}

type ValidPatron struct {
	*BoolValue
}

func (vp ValidPatron) Info() (id, name string, length int) {
	return "BL", "valid_patron", 1
}

type ValidPatronPassword struct {
	*BoolValue
}

func (vpp ValidPatronPassword) Info() (id, name string, length int) {
	return "CQ", "valid_patron_password", 1
}

type ScreenMessage struct {
	*StrValue
}

func (sm ScreenMessage) Info() (id, name string, length int) {
	return "AF", "screen_message", -1
}

type PrintLine struct {
	*StrValue
}

func (pl PrintLine) Info() (id, name string, length int) {
	return "AG", "print_line", -1
}

type OK struct {
	*BoolValue
}

func (ok OK) Info() (id, name string, length int) {
	return "", "ok", 1
}

type RenewalOK struct {
	*BoolValue
}

func (rok RenewalOK) info() (id, name string, length int) {
	return "", "renewal_ok", 1
}

type MagneticMedia struct {
	*BoolValue
}

func (mm MagneticMedia) Info() (id, name string, length int) {
	return "", "megnetic_media", 1
}

type Desensitize struct {
	*BoolValue
}

func (d Desensitize) Info() (id, name string, length int) {
	return "", "desensitize", 1
}

type DueDate struct {
	*TimeValue
}

func (dd DueDate) Info() (id, name string, length int) {
	return "AH", "due_date", -1
}

type SecurityInhibit struct {
	*BoolValue
}

func (si SecurityInhibit) Info() (id, name string, length int) {
	return "CI", "security_inhibit", 1
}

type MediaType struct {
	*IntValue
}

func (mt MediaType) Info() (id, name string, length int) {
	return "CK", "media_type", 3
}

type Resensitize struct {
	*BoolValue
}

func (r Resensitize) Info() (id, name string, length int) {
	return "", "resensitize", 1
}

type Alert struct {
	*BoolValue
}

func (a Alert) Info() (id, name string, length int) {
	return "", "alert", 1
}

type PermanentLocation struct {
	*StrValue
}

func (pl PermanentLocation) Info() (id, name string, length int) {
	return "AQ", "permanent_location", -1
}

type SortBin struct {
	*StrValue
}

func (sb SortBin) Info() (id, name string, length int) {
	return "CL", "sort_bin", -1
}

type OnlineStatus struct {
	*BoolValue
}

func (os OnlineStatus) Info() (id, name string, length int) {
	return "", "online_status", 1
}

type CheckinOK struct {
	*BoolValue
}

func (cok CheckinOK) Info() (id, name string, length int) {
	return "", "checkin_ok", 1
}

type CheckoutOK struct {
	*BoolValue
}

func (cok CheckoutOK) Info() (id, name string, length int) {
	return "", "checkout_ok", 1
}

type ACSRenewalPolicy struct {
	*BoolValue
}

func (arp ACSRenewalPolicy) Info() (id, name string, length int) {
	return "", "acs_renewal_policy", 1
}

type StatusUpdateOK struct {
	*BoolValue
}

func (suo StatusUpdateOK) Info() (id, name string, length int) {
	return "", "status_update_ok", 1
}

type OfflineOK struct {
	*BoolValue
}

func (ook OfflineOK) Info() (id, name string, length int) {
	return "", "offline_ok", 1
}

type TimeoutPeriod struct {
	*IntValue
}

func (tp TimeoutPeriod) Info() (id, name string, length int) {
	return "", "timeout_period", 3
}

type RetriesAllowed struct {
	*BoolValue
}

func (ra RetriesAllowed) Info() (id, name string, length int) {
	return "", "retries_allowed", 3
}

type DateTimeSync struct {
	*TimeValue
}

func (dts DateTimeSync) Info() (id, name string, length int) {
	return "", "datetime_sync", 18
}

type LibraryName struct {
	*StrValue
}

func (ln LibraryName) Info() (id, name string, length int) {
	return "AM", "library_name", -1
}

type SupportedMessages struct {
	*StrValue
}

func (sm SupportedMessages) Info() (id, name string, length int) {
	return "BX", "supported_messages", -1
}

type TerminalLocation struct {
	*StrValue
}

func (tl TerminalLocation) Info() (id, name string, length int) {
	return "AN", "terminal_location", -1
}

type HoldItemsCount struct {
	*IntValue
}

func (hic HoldItemsCount) Info() (id, name string, length int) {
	return "", "hold_item_count", 4
}

type OverdueItemsCount struct {
	*IntValue
}

func (oic OverdueItemsCount) Info() (id, name string, length int) {
	return "", "overdue_items_count", 4
}

type ChargedItemsCount struct {
	*IntValue
}

func (cic ChargedItemsCount) Info() (id, name string, length int) {
	return "", "charged_items_count", 4
}

type FindItemsCount struct {
	*IntValue
}

func (fic FindItemsCount) Info() (id, name string, length int) {
	return "", "fine_items_count", 4
}

type RecallItemsCount struct {
	*IntValue
}

func (ric RecallItemsCount) Info() (id, name string, length int) {
	return "", "recall_items_count", 4
}

type UnavailableHoldsCount struct {
	*IntValue
}

func (uhc UnavailableHoldsCount) Info() (id, name string, length int) {
	return "", "unavailable_holds_count", 4
}

type HoldItemsLimit struct {
	*IntValue
}

func (hil HoldItemsLimit) Info() (id, name string, length int) {
	return "BZ", "hold_items_limit", 4
}

type OverdueItemsLimit struct {
	*IntValue
}

func (oil OverdueItemsLimit) Info() (id, name string, length int) {
	return "CA", "overdue_items_limit", 4
}

type ChargedItemsLimit struct {
	*IntValue
}

func (cil ChargedItemsLimit) Info() (id, name string, length int) {
	return "CB", "charged_items_limit", 4
}

type FeeLimit struct {
	*IntValue
}

func (fl FeeLimit) Info() (id, name string, length int) {
	return "CC", "fee_limit", -1
}

type HoldItems struct {
	*StrSliceValue
}

func (hi HoldItems) Info() (id, name string, length int) {
	return "AS", "hold_items", -1
}

type OverdueItems struct {
	*StrSliceValue
}

func (oi OverdueItems) Info() (id, name string, length int) {
	return "AT", "overdue_items", -1
}

type ChargedItems struct {
	*StrSliceValue
}

func (ci ChargedItems) Info() (id, name string, length int) {
	return "AU", "charged_items", -1
}

type FineItems struct {
	*StrSliceValue
}

func (fi FineItems) Info() (id, name string, length int) {
	return "AV", "fine_items", -1
}

type RecallItems struct {
	*StrSliceValue
}

func (ri RecallItems) Info() (id, name string, length int) {
	return "BU", "recall_items", -1
}

type UnavailableHoldItems struct {
	*StrSliceValue
}

func (uhi UnavailableHoldItems) Info() (id, name string, length int) {
	return "CD", "unavailable_hold_items", -1
}

type HomeAddress struct {
	*StrValue
}

func (ha HomeAddress) Info() (id, name string, length int) {
	return "BD", "home_address", -1
}

type EmailAddress struct {
	*StrValue
}

func (ea EmailAddress) Info() (id, name string, length int) {
	return "BE", "email_address", -1
}

type HomePhoneNumber struct {
	*StrValue
}

func (hpn HomePhoneNumber) Info() (id, name string, length int) {
	return "BF", "home_phone_number", -1
}

type EndSession struct {
	*BoolValue
}

func (es EndSession) Info() (id, name string, length int) {
	return "", "end_session", 1
}

type PaymentAccepted struct {
	*BoolValue
}

func (pa PaymentAccepted) Info() (id, name string, length int) {
	return "", "payment_accepted", 1
}

type CirculationStatus struct {
	*IntValue
}

func (cs CirculationStatus) Info() (id, name string, length int) {
	return "", "circulation_status", 2
}

type SecurityMarker struct {
	*IntValue
}

func (sm SecurityMarker) Info() (id, name string, length int) {
	return "", "security_maker", 2
}

type HoldQueueLength struct {
	*FloatValue
}

func (hq HoldQueueLength) Info() (id, name string, length int) {
	return "CF", "hold_queue_length", -1
}

type RecallDate struct {
	*TimeValue
}

func (rd RecallDate) Info() (id, name string, length int) {
	return "CJ", "recall_date", 18
}

type HoldPickupDate struct {
	*TimeValue
}

func (hpd HoldPickupDate) Info() (id, name string, length int) {
	return "CM", "hold_pickup_date", 18
}

type Owner struct {
	*StrValue
}

func (o Owner) Info() (id, name string, length int) {
	return "BG", "owner", -1
}

type ItemPropertiesOK struct {
	*BoolValue
}

func (ipt ItemPropertiesOK) Info() (id, name string, length int) {
	return "", "item_properties_ok", 1
}

type Available struct {
	*BoolValue
}

func (a Available) Info() (id, name string, length int) {
	return "", "available", 1
}

type QueuePosition struct {
	*IntValue
}

func (qp QueuePosition) Info() (id, name string, length int) {
	return "BR", "queue_position", -1
}

type RenewedCount struct {
	*IntValue
}

func (rc RenewedCount) Info() (id, name string, length int) {
	return "", "renewed_count", 4
}

type UnrenewedCount struct {
	*IntValue
}

func (uc UnrenewedCount) Info() (id, name string, length int) {
	return "", "unrenewed_count", 4
}

type RenewedItems struct {
	*StrSliceValue
}

func (ri RenewedItems) Info() (id, name string, length int) {
	return "BM", "renewed_items", -1
}

type UnrenewedItems struct {
	*StrSliceValue
}

func (ui UnrenewedItems) Info() (id, name string, length int) {
	return "BN", "unrenewed_items", -1
}

type FineItemsCount struct {
	*IntValue
}

func (fic FineItemsCount) Info() (id, name string, length int) {
	return "", "fine_items_count", 4
}

// type UnknowedJE struct {
// 	*StrValue
// }
//
// func (uje *UnknowedJE) Info() (id, name string, length int) {
// 	return "JE", "unknowed je", -1
// }
//
// type UnknowedJF struct {
// 	*StrValue
// }
//
// func (ujf *UnknowedJF) Info() (id, name string, length int) {
// 	return "JF", "unknowed jf", -1
// }
//
// type UnknowedXV struct {
// 	*StrValue
// }
//
// func (uxv *UnknowedXV) Info() (id, name string, length int) {
// 	return "XV", "unknowed xv", -1
// }
//
// type UnknowedXR struct {
// 	*StrValue
// }
//
// func (uxr *UnknowedXR) Info() (id, name string, length int) {
// 	return "XR", "unknowed xr", -1
// }
//
// type UnknowedXC struct {
// 	*StrValue
// }
//
// func (uxc *UnknowedXC) Info() (id, name string, length int) {
// 	return "XC", "unknowed xr", -1
// }
//
// type UnknowedLL struct {
// 	*StrValue
// }
//
// func (ull *UnknowedLL) Info() (id, name string, length int) {
// 	return "LL", "unknowed ll", -1
// }
//
// type UnknowedGL struct {
// 	*StrValue
// }
//
// func (ugl *UnknowedGL) Info() (id, name string, length int) {
// 	return "GL", "unknowed gl", -1
// }
//
// type UnknowedXT struct {
// 	*StrValue
// }
//
// func (uxt *UnknowedXT) Info() (id, name string, length int) {
// 	return "XT", "unknowed XT", -1
// }
//
// type UnknowedXD struct {
// 	*StrValue
// }
//
// func (uxd *UnknowedXD) Info() (id, name string, length int) {
// 	return "XD", "unknowed XD", -1
// }
//
// type UnknowedXO struct {
// 	*StrValue
// }
//
// func (uxo *UnknowedXO) Info() (id, name string, length int) {
// 	return "XO", "unknowed XO", -1
// }
//
// type UnknowedRU struct {
// 	*StrValue
// }
//
// func (uru *UnknowedRU) Info() (id, name string, length int) {
// 	return "RU", "unknowed RU", -1
// }
//
// type UnknowedXH struct {
// 	*StrValue
// }
//
// func (uxh *UnknowedXH) Info() (id, name string, length int) {
// 	return "XH", "unknowed XH", -1
// }
//
// type UnknowedRS struct {
// 	*StrValue
// }
//
// func (urs *UnknowedRS) Info() (id, name string, length int) {
// 	return "RS", "unknowed RS", -1
// }
//
// type UnknowedMC struct {
// 	*StrValue
// }
//
// func (umc *UnknowedMC) Info() (id, name string, length int) {
// 	return "MC", "unknowed MC", -1
// }
//
// type UnknowedSR struct {
// 	*StrValue
// }
//
// func (usr *UnknowedSR) Info() (id, name string, length int) {
// 	return "SR", "unknowed SR", -1
// }
//
// type UnknowedPN struct {
// 	*StrValue
// }
//
// func (upn *UnknowedPN) Info() (id, name string, length int) {
// 	return "PN", "unknowed PN", -1
// }
//
// type UnknowedXF struct {
// 	*StrValue
// }
//
// func (uxf *UnknowedXF) Info() (id, name string, length int) {
// 	return "XF", "unknowed XF", -1
// }
//
// type SequenceNumber struct {
// 	*IntValue
// }

// func (sn *SequenceNumber) Info() (id, name string, length int) {
// 	return "AY", "sequence number", 1
// }
