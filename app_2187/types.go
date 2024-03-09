package main

import (
	"database/sql/driver"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// time.Time format for parsing/formatting input/output
const dateFmt = "2006-01-02Z0700"

// Date defines type to store a date.
type Date struct{ time.Time }

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Date) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*d, err = parseDate(s)
	if err != nil {
		return err
	}

	return nil
}

func parseDate(s string) (Date, error) {
	t, err := time.Parse(dateFmt, s)
	if err != nil {
		t, err = time.Parse(time.RFC3339Nano, s)
	}
	return Date{t}, err
}

// MarshalJSON implements the json.Marshaler interface.
func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Format(dateFmt))
}

// Scan implements the sql.Scanner interface.
func (d *Date) Scan(value interface{}) error {
	var ok bool
	d.Time, ok = value.(time.Time)
	if !ok {
		return fmt.Errorf("value not type time.Time")
	}

	return nil
}

// Value implements the driver Valuer interface.
func (d Date) Value() (driver.Value, error) {
	return d.Time, nil
}

// EventError defines type to store a error.
type EventError struct {
	StatusCode int
	Data       interface{}
	Message    string
}

// Error implements the error interface.
func (e *EventError) Error() string {
	return e.Message
}

// File defines type to store a file.
type File struct {
	FilePath  string
	CSVReader *csv.Reader
	// copy of file
	*os.File
}

// MarshalJSON implements the json.Marshaler interface.
func (f File) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.FilePath)
}

// CSVFileWrapper defines a wrapper to allow reading a file in CSV format.
type CSVFileWrapper File

func (f *CSVFileWrapper) Read() (record []string, err error) {
	if f.CSVReader == nil {
		f.CSVReader = csv.NewReader(f.File)
	}

	return f.CSVReader.Read()
}

// NonStrictInt defines a type to parse non-int into an int value.
type NonStrictInt struct{ int }

// UnmarshalJSON implements the json.Unmarshaler interface.
func (ri *NonStrictInt) UnmarshalJSON(buf []byte) (err error) {
	var i int

	// Attempt to convert string to int.
	if buf[0] == '"' {
		var s string
		if err = json.Unmarshal(buf, &s); err != nil {
			return err
		}

		if s == "" {
			s = "0"
		}

		if i, err = strconv.Atoi(s); err != nil {
			return err
		}
	} else {
		if err = json.Unmarshal(buf, &i); err != nil {
			return err
		}
	}

	*ri = NonStrictInt{i}
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (ri NonStrictInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(ri.int)
}

// NonStrictFloat defines a type to parse non-float into a float value.
type NonStrictFloat struct{ float64 }

// UnmarshalJSON implements the json.Unmarshaler interface.
func (ri *NonStrictFloat) UnmarshalJSON(buf []byte) (err error) {
	var f float64

	// Attempt to convert string to int.
	if buf[0] == '"' {
		var s string
		if err = json.Unmarshal(buf, &s); err != nil {
			return err
		}

		if s == "" {
			s = "0"
		}

		if f, err = strconv.ParseFloat(s, 64); err != nil {
			return err
		}
	} else {
		if err = json.Unmarshal(buf, &f); err != nil {
			return err
		}
	}

	*ri = NonStrictFloat{f}
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (ri NonStrictFloat) MarshalJSON() ([]byte, error) {
	return json.Marshal(ri.float64)
}

// NonStrictBool defines a type to parse non-bool into a bool value.
type NonStrictBool struct{ bool }

// UnmarshalJSON implements the json.Unmarshaler interface.
func (ri *NonStrictBool) UnmarshalJSON(buf []byte) (err error) {
	var b bool
	if b, err = strconv.ParseBool(strings.ToLower(strings.Trim(string(buf), "\""))); err != nil {
		return err
	}

	*ri = NonStrictBool{b}
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (ri NonStrictBool) MarshalJSON() ([]byte, error) {
	return json.Marshal(ri.bool)
}

// NonStrictString defines a type to parse non-string into a string value.
type NonStrictString struct{ string }

// UnmarshalJSON implements the json.Unmarshaler interface.
func (ri *NonStrictString) UnmarshalJSON(buf []byte) (err error) {
	var s string
	if buf[0] == '"' {
		if err = json.Unmarshal(buf, &s); err != nil {
			return err
		}
	} else {
		s = string(buf)
	}

	*ri = NonStrictString{s}
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (ri NonStrictString) MarshalJSON() ([]byte, error) {
	return json.Marshal(ri.string)
}

// Thumbnail defines type to store the "thumbnail" system construct.
type Thumbnail struct {
	Width  *int   `json:"width" db:"width" csv:"width"`
	Height *int   `json:"height" db:"height" csv:"height"`
	Url    string `json:"url" db:"url" csv:"url"`
}

// Scan implements sql.Scanner interface.
func (a *Thumbnail) Scan(src interface{}) error {
	switch t := src.(type) {
	case []byte:
		return json.Unmarshal(t, a)
	case nil:
		return nil
	default:
		return fmt.Errorf("cannot convert %T to Thumbnail", src)
	}
}

// Value implements driver.Valuer interface.
func (a Thumbnail) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Attachment defines type to store the "attachment" system construct.
type Attachment struct {
	Id       string     `json:"id" db:"id" csv:"id"`
	Url      string     `json:"url" db:"url" csv:"url"`
	Filename string     `json:"filename" db:"filename" csv:"filename"`
	Width    *int       `json:"width" db:"width" csv:"width"`
	Height   *int       `json:"height" db:"height" csv:"height"`
	Type     string     `json:"type" db:"type" csv:"type"`
	Size     int        `json:"size" db:"size" csv:"size"`
	Small    *Thumbnail `json:"small" db:"small" csv:"small,inline"`
	Medium   *Thumbnail `json:"medium" db:"medium" csv:"medium,inline"`
	Large    *Thumbnail `json:"large" db:"large" csv:"large,inline"`
	Filepath string     `json:"filepath" db:"filepath" csv:"filepath"`
}

// Scan implements sql.Scanner interface.
func (a *Attachment) Scan(src interface{}) error {
	switch t := src.(type) {
	case []byte:
		return json.Unmarshal(t, a)
	case nil:
		return nil
	default:
		return fmt.Errorf("cannot convert %T to Attachment", src)
	}
}

// Value implements driver.Valuer interface.
func (a Attachment) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// AttachmentArray defines a type to store a slice of type "Attachment".
type AttachmentArray []Attachment

// Scan implements sql.Scanner interface.
func (a *AttachmentArray) Scan(src interface{}) error {
	switch t := src.(type) {
	case []byte:
		return json.Unmarshal(t, a)
	case nil:
		return nil
	default:
		return fmt.Errorf("cannot convert %T to AttachmentArray", src)
	}
}

// Value implements driver.Valuer interface.
func (a AttachmentArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// RequestPasswordResetEPIConstruct defines type to store the "requestPasswordResetEPI" construct.
type RequestPasswordResetEPIConstruct struct {
	Email NonStrictString `json:"email" yaml:"email" schema:"email"  db:"email" csv:"email""`
}

// LoginEPIConstruct defines type to store the "loginEPI" construct.
type LoginEPIConstruct struct {
	Email    string `json:"email" yaml:"email" schema:"email"  db:"email" csv:"email""`
	Password string `json:"password" yaml:"password" schema:"password"  db:"password" csv:"password""`
}

// SignUpEPIConstruct defines type to store the "signUpEPI" construct.
type SignUpEPIConstruct struct {
	Password  string `json:"password" yaml:"password" schema:"password"  db:"password" csv:"password""`
	Email     string `json:"email" yaml:"email" schema:"email"  db:"email" csv:"email""`
	FirstName string `json:"firstName" yaml:"firstName" schema:"firstName"  db:"firstName" csv:"firstName""`
	LastName  string `json:"lastName" yaml:"lastName" schema:"lastName"  db:"lastName" csv:"lastName""`
}

// CreateListingEPIConstruct defines type to store the "createListingEPI" construct.
type CreateListingEPIConstruct struct {
	GuestLimit     int  `json:"guest_limit" yaml:"guest_limit" schema:"guest_limit"  db:"guest_limit" csv:"guest_limit""`
	Price          *int `json:"price" yaml:"price" schema:"price"  db:"price" csv:"price""`
	OrderTypeId    *int `json:"order_type_id" yaml:"order_type_id" schema:"order_type_id"  db:"order_type_id" csv:"order_type_id""`
	PartnerSiteId  int  `json:"partnerSiteId" yaml:"partnerSiteId" schema:"partnerSiteId"  db:"partnerSiteId" csv:"partnerSiteId""`
	CheckInTimeId  *int `json:"check_in_time_id" yaml:"check_in_time_id" schema:"check_in_time_id"  db:"check_in_time_id" csv:"check_in_time_id""`
	CheckOutTimeId *int `json:"check_out_time_id" yaml:"check_out_time_id" schema:"check_out_time_id"  db:"check_out_time_id" csv:"check_out_time_id""`
}

// GetPropertiesEPRConstruct defines type to store the "getPropertiesEPR" construct.
type GetPropertiesEPRConstruct struct {
	Id          int              `json:"id" yaml:"id" schema:"id"  db:"id" csv:"id""`
	Description *string          `json:"description" yaml:"description" schema:"description"  db:"description" csv:"description""`
	LandType    *string          `json:"land_type" yaml:"land_type" schema:"land_type"  db:"land_type" csv:"land_type""`
	Address     string           `json:"address" yaml:"address" schema:"address"  db:"address" csv:"address""`
	City        string           `json:"city" yaml:"city" schema:"city"  db:"city" csv:"city""`
	State       string           `json:"state" yaml:"state" schema:"state"  db:"state" csv:"state""`
	Zip         int              `json:"zip" yaml:"zip" schema:"zip"  db:"zip" csv:"zip""`
	Status      *string          `json:"status" yaml:"status" schema:"status"  db:"status" csv:"status""`
	Images      *AttachmentArray `json:"images" yaml:"images" schema:"images"  db:"images" csv:"-""`
	Name        string           `json:"name" yaml:"name" schema:"name"  db:"name" csv:"name""`
	CreatorId   int              `json:"creator_id" yaml:"creator_id" schema:"creator_id"  db:"creator_id" csv:"creator_id""`
	Amenity     *string          `json:"amenity" yaml:"amenity" schema:"amenity"  db:"amenity" csv:"amenity""`
	Acres       *int             `json:"acres" yaml:"acres" schema:"acres"  db:"acres" csv:"acres""`
	File        *AttachmentArray `json:"file" yaml:"file" schema:"file"  db:"file" csv:"-""`
}

// CreatePropertiesEPIConstruct defines type to store the "createPropertiesEPI" construct.
type CreatePropertiesEPIConstruct struct {
	Address     string           `json:"address" yaml:"address" schema:"address"  db:"address" csv:"address""`
	City        string           `json:"city" yaml:"city" schema:"city"  db:"city" csv:"city""`
	Name        string           `json:"name" yaml:"name" schema:"name"  db:"name" csv:"name""`
	Description string           `json:"description" yaml:"description" schema:"description"  db:"description" csv:"description""`
	Zip         int              `json:"zip" yaml:"zip" schema:"zip"  db:"zip" csv:"zip""`
	Amenity     *NonStrictString `json:"amenity" yaml:"amenity" schema:"amenity"  db:"amenity" csv:"amenity""`
	LandType    *NonStrictString `json:"land_type" yaml:"land_type" schema:"land_type"  db:"land_type" csv:"land_type""`
	State       *NonStrictString `json:"state" yaml:"state" schema:"state"  db:"state" csv:"state""`
	Status      *NonStrictString `json:"status" yaml:"status" schema:"status"  db:"status" csv:"status""`
	Acres       int              `json:"acres" yaml:"acres" schema:"acres"  db:"acres" csv:"acres""`
	Images      *AttachmentArray `json:"images" yaml:"images" schema:"images"  db:"images" csv:"-""`
}

// UpdatePropertiesEPIConstruct defines type to store the "updatePropertiesEPI" construct.
type UpdatePropertiesEPIConstruct struct {
	Status      *string          `json:"status" yaml:"status" schema:"status"  db:"status" csv:"status""`
	CreatorId   *int             `json:"creator_id" yaml:"creator_id" schema:"creator_id"  db:"creator_id" csv:"creator_id""`
	City        *string          `json:"city" yaml:"city" schema:"city"  db:"city" csv:"city""`
	State       *string          `json:"state" yaml:"state" schema:"state"  db:"state" csv:"state""`
	Zip         *int             `json:"zip" yaml:"zip" schema:"zip"  db:"zip" csv:"zip""`
	Name        *string          `json:"name" yaml:"name" schema:"name"  db:"name" csv:"name""`
	Description *string          `json:"description" yaml:"description" schema:"description"  db:"description" csv:"description""`
	Address     *string          `json:"address" yaml:"address" schema:"address"  db:"address" csv:"address""`
	Amenity     *NonStrictString `json:"amenity" yaml:"amenity" schema:"amenity"  db:"amenity" csv:"amenity""`
	LandType    *string          `json:"land_type" yaml:"land_type" schema:"land_type"  db:"land_type" csv:"land_type""`
}

// GetListingsEPRConstruct defines type to store the "getListingsEPR" construct.
type GetListingsEPRConstruct struct {
	Amenity        *string          `json:"amenity" yaml:"amenity" schema:"amenity"  db:"amenity" csv:"amenity""`
	Id             int              `json:"id" yaml:"id" schema:"id"  db:"id" csv:"id""`
	GuestLimit     int              `json:"guest_limit" yaml:"guest_limit" schema:"guest_limit"  db:"guest_limit" csv:"guest_limit""`
	PropertyId     int              `json:"property_id" yaml:"property_id" schema:"property_id"  db:"property_id" csv:"property_id""`
	Price          *float64         `json:"price" yaml:"price" schema:"price"  db:"price" csv:"price""`
	NicheAmenity   *string          `json:"niche_amenity" yaml:"niche_amenity" schema:"niche_amenity"  db:"niche_amenity" csv:"niche_amenity""`
	OrderTypeId    *int             `json:"order_type_id" yaml:"order_type_id" schema:"order_type_id"  db:"order_type_id" csv:"order_type_id""`
	Description    string           `json:"description" yaml:"description" schema:"description"  db:"description" csv:"description""`
	Address        string           `json:"address" yaml:"address" schema:"address"  db:"address" csv:"address""`
	City           string           `json:"city" yaml:"city" schema:"city"  db:"city" csv:"city""`
	State          string           `json:"state" yaml:"state" schema:"state"  db:"state" csv:"state""`
	Zip            int              `json:"zip" yaml:"zip" schema:"zip"  db:"zip" csv:"zip""`
	Images         AttachmentArray  `json:"images" yaml:"images" schema:"images"  db:"images" csv:"-""`
	Status         string           `json:"status" yaml:"status" schema:"status"  db:"status" csv:"status""`
	LandType       *string          `json:"land_type" yaml:"land_type" schema:"land_type"  db:"land_type" csv:"land_type""`
	CreatorId      int              `json:"creator_id" yaml:"creator_id" schema:"creator_id"  db:"creator_id" csv:"creator_id""`
	Email          string           `json:"email" yaml:"email" schema:"email"  db:"email" csv:"email""`
	FirstName      *string          `json:"first_name" yaml:"first_name" schema:"first_name"  db:"first_name" csv:"first_name""`
	LastName       *string          `json:"last_name" yaml:"last_name" schema:"last_name"  db:"last_name" csv:"last_name""`
	PartnerSiteId  *int             `json:"partner_site_id" yaml:"partner_site_id" schema:"partner_site_id"  db:"partner_site_id" csv:"partner_site_id""`
	Acres          int              `json:"acres" yaml:"acres" schema:"acres"  db:"acres" csv:"acres""`
	CheckOutTimeId *int             `json:"check_out_time_id" yaml:"check_out_time_id" schema:"check_out_time_id"  db:"check_out_time_id" csv:"check_out_time_id""`
	CheckInTimeId  *int             `json:"check_in_time_id" yaml:"check_in_time_id" schema:"check_in_time_id"  db:"check_in_time_id" csv:"check_in_time_id""`
	OrderTypeName  *string          `json:"order_type_name" yaml:"order_type_name" schema:"order_type_name"  db:"order_type_name" csv:"order_type_name""`
	OutTimeHour    *string          `json:"out_time_hour" yaml:"out_time_hour" schema:"out_time_hour"  db:"out_time_hour" csv:"out_time_hour""`
	InTimeHour     *string          `json:"in_time_hour" yaml:"in_time_hour" schema:"in_time_hour"  db:"in_time_hour" csv:"in_time_hour""`
	ProfilePicture *AttachmentArray `json:"profile_picture" yaml:"profile_picture" schema:"profile_picture"  db:"profile_picture" csv:"-""`
	Name           string           `json:"name" yaml:"name" schema:"name"  db:"name" csv:"name""`
	File           *AttachmentArray `json:"file" yaml:"file" schema:"file"  db:"file" csv:"-""`
}

// CalendarsConstruct defines type to store the "calendars" construct.
type CalendarsConstruct struct {
	Id         int     `json:"Id" yaml:"Id" schema:"Id"  db:"Id" csv:"Id""`
	PropertyId int     `json:"property_id" yaml:"property_id" schema:"property_id"  db:"property_id" csv:"property_id""`
	Month      *string `json:"month" yaml:"month" schema:"month"  db:"month" csv:"month""`
	Year       int     `json:"year" yaml:"year" schema:"year"  db:"year" csv:"year""`
	Date       *Date   `json:"date" yaml:"date" schema:"date"  db:"date" csv:"date""`
}

// FormGetOneCalendarsEndpointConstruct defines type to store the "formGetOneCalendarsEndpoint" construct.
type FormGetOneCalendarsEndpointConstruct struct {
	Id NonStrictInt `json:"id" yaml:"id" schema:"id"  db:"id" csv:"id""`
}

// FormCreateCalendarsEndpointConstruct defines type to store the "FormCreateCalendarsEndpoint" construct.
type FormCreateCalendarsEndpointConstruct struct {
	Id         int    `json:"Id" yaml:"Id" schema:"Id"  db:"Id" csv:"Id""`
	PropertyId int    `json:"property_id" yaml:"property_id" schema:"property_id"  db:"property_id" csv:"property_id""`
	Month      string `json:"month" yaml:"month" schema:"month"  db:"month" csv:"month""`
	Year       int    `json:"year" yaml:"year" schema:"year"  db:"year" csv:"year""`
	Date       Date   `json:"date" yaml:"date" schema:"date"  db:"date" csv:"date""`
}

// FormUpdateCalendarsEndpointConstruct defines type to store the "FormUpdateCalendarsEndpoint" construct.
type FormUpdateCalendarsEndpointConstruct struct {
	PropertyId int    `json:"property_id" yaml:"property_id" schema:"property_id"  db:"property_id" csv:"property_id""`
	Month      string `json:"month" yaml:"month" schema:"month"  db:"month" csv:"month""`
	Year       int    `json:"year" yaml:"year" schema:"year"  db:"year" csv:"year""`
	Date       Date   `json:"date" yaml:"date" schema:"date"  db:"date" csv:"date""`
	Id         int    `json:"Id" yaml:"Id" schema:"Id"  db:"Id" csv:"Id""`
}

// FormDeleteCalendarsEndpointConstruct defines type to store the "FormDeleteCalendarsEndpoint" construct.
type FormDeleteCalendarsEndpointConstruct struct {
	Id NonStrictInt `json:"id" yaml:"id" schema:"id"  db:"id" csv:"id""`
}

// GetUsStatesConstruct defines type to store the "getUsStates" construct.
type GetUsStatesConstruct struct {
	Id           int     `json:"Id" yaml:"Id" schema:"Id"  db:"Id" csv:"Id""`
	State        *string `json:"State" yaml:"State" schema:"State"  db:"State" csv:"State""`
	Abbreviation *string `json:"Abbreviation" yaml:"Abbreviation" schema:"Abbreviation"  db:"Abbreviation" csv:"Abbreviation""`
}

// GetPropertyListingEPRConstruct defines type to store the "getPropertyListingEPR" construct.
type GetPropertyListingEPRConstruct struct {
	Id             *int             `json:"id" yaml:"id" schema:"id"  db:"id" csv:"id""`
	Name           string           `json:"name" yaml:"name" schema:"name"  db:"name" csv:"name""`
	Avatar         *AttachmentArray `json:"avatar" yaml:"avatar" schema:"avatar"  db:"avatar" csv:"-""`
	PropertyId     int              `json:"property_id" yaml:"property_id" schema:"property_id"  db:"property_id" csv:"property_id""`
	Address        string           `json:"address" yaml:"address" schema:"address"  db:"address" csv:"address""`
	City           string           `json:"city" yaml:"city" schema:"city"  db:"city" csv:"city""`
	State          string           `json:"state" yaml:"state" schema:"state"  db:"state" csv:"state""`
	CheckInTimeId  *int             `json:"check_in_time_id" yaml:"check_in_time_id" schema:"check_in_time_id"  db:"check_in_time_id" csv:"check_in_time_id""`
	CheckOutTimeId *int             `json:"check_out_time_id" yaml:"check_out_time_id" schema:"check_out_time_id"  db:"check_out_time_id" csv:"check_out_time_id""`
	Hour           *string          `json:"hour" yaml:"hour" schema:"hour"  db:"hour" csv:"hour""`
	Amenity        string           `json:"amenity" yaml:"amenity" schema:"amenity"  db:"amenity" csv:"amenity""`
	Price          int              `json:"price" yaml:"price" schema:"price"  db:"price" csv:"price""`
	GuestLimit     int              `json:"guest_limit" yaml:"guest_limit" schema:"guest_limit"  db:"guest_limit" csv:"guest_limit""`
	OrderTypeId    int              `json:"order_type_id" yaml:"order_type_id" schema:"order_type_id"  db:"order_type_id" csv:"order_type_id""`
	CreatorId      int              `json:"creator_id" yaml:"creator_id" schema:"creator_id"  db:"creator_id" csv:"creator_id""`
	PartnerSiteId  int              `json:"partner_site_id" yaml:"partner_site_id" schema:"partner_site_id"  db:"partner_site_id" csv:"partner_site_id""`
}

// CalendarModelConstruct defines type to store the "calendarModel" construct.
type CalendarModelConstruct struct {
	CurrentMonth int `json:"currentMonth" yaml:"currentMonth" schema:"currentMonth"  db:"currentMonth" csv:"currentMonth""`
	CurrentYear  int `json:"currentYear" yaml:"currentYear" schema:"currentYear"  db:"currentYear" csv:"currentYear""`
	DaysInMonth  int `json:"daysInMonth" yaml:"daysInMonth" schema:"daysInMonth"  db:"daysInMonth" csv:"daysInMonth""`
}

// FormGetOneBookingsEndpointConstruct defines type to store the "FormGetOneBookingsEndpoint" construct.
type FormGetOneBookingsEndpointConstruct struct {
	Id NonStrictInt `json:"id" yaml:"id" schema:"id"  db:"id" csv:"id""`
}

// FormCreateBookingsEndpointConstruct defines type to store the "FormCreateBookingsEndpoint" construct.
type FormCreateBookingsEndpointConstruct struct {
	Id                int             `json:"Id" yaml:"Id" schema:"Id"  db:"Id" csv:"Id""`
	ListingId         int             `json:"listing_id" yaml:"listing_id" schema:"listing_id"  db:"listing_id" csv:"listing_id""`
	BookingDate       time.Time       `json:"booking_date" yaml:"booking_date" schema:"booking_date"  db:"booking_date" csv:"booking_date""`
	CheckoutSessionId string          `json:"checkout_session_id" yaml:"checkout_session_id" schema:"checkout_session_id"  db:"checkout_session_id" csv:"checkout_session_id""`
	GuestId           int             `json:"guest_id" yaml:"guest_id" schema:"guest_id"  db:"guest_id" csv:"guest_id""`
	Status            string          `json:"status" yaml:"status" schema:"status"  db:"status" csv:"status""`
	Request           json.RawMessage `json:"request" yaml:"request" schema:"request"  db:"request" csv:"request""`
}

// UpdateBookingsEPIConstruct defines type to store the "updateBookingsEPI" construct.
type UpdateBookingsEPIConstruct struct {
	Status            string          `json:"status" yaml:"status" schema:"status"  db:"status" csv:"status""`
	Request           json.RawMessage `json:"request" yaml:"request" schema:"request"  db:"request" csv:"request""`
	Id                int             `json:"Id" yaml:"Id" schema:"Id"  db:"Id" csv:"Id""`
	ListingId         int             `json:"listing_id" yaml:"listing_id" schema:"listing_id"  db:"listing_id" csv:"listing_id""`
	BookingDate       time.Time       `json:"booking_date" yaml:"booking_date" schema:"booking_date"  db:"booking_date" csv:"booking_date""`
	CheckoutSessionId string          `json:"checkout_session_id" yaml:"checkout_session_id" schema:"checkout_session_id"  db:"checkout_session_id" csv:"checkout_session_id""`
	GuestId           int             `json:"guest_id" yaml:"guest_id" schema:"guest_id"  db:"guest_id" csv:"guest_id""`
}

// FormDeleteBookingsEndpointConstruct defines type to store the "FormDeleteBookingsEndpoint" construct.
type FormDeleteBookingsEndpointConstruct struct {
	Id NonStrictInt `json:"id" yaml:"id" schema:"id"  db:"id" csv:"id""`
}

// GetBookingConstruct defines type to store the "getBooking" construct.
type GetBookingConstruct struct {
	Id                int       `json:"Id" yaml:"Id" schema:"Id"  db:"Id" csv:"Id""`
	ListingId         int       `json:"listing_id" yaml:"listing_id" schema:"listing_id"  db:"listing_id" csv:"listing_id""`
	BookingDate       time.Time `json:"booking_date" yaml:"booking_date" schema:"booking_date"  db:"booking_date" csv:"booking_date""`
	CheckoutSessionId string    `json:"checkout_session_id" yaml:"checkout_session_id" schema:"checkout_session_id"  db:"checkout_session_id" csv:"checkout_session_id""`
	GuestId           int       `json:"guest_id" yaml:"guest_id" schema:"guest_id"  db:"guest_id" csv:"guest_id""`
	Status            *string   `json:"status" yaml:"status" schema:"status"  db:"status" csv:"status""`
	Request           *string   `json:"request" yaml:"request" schema:"request"  db:"request" csv:"request""`
}

// ApproveBookingRequestConstruct defines type to store the "approveBookingRequest" construct.
type ApproveBookingRequestConstruct struct {
	Status    string `json:"status" yaml:"status" schema:"status"  db:"status" csv:"status""`
	BookingId int    `json:"booking_id" yaml:"booking_id" schema:"booking_id"  db:"booking_id" csv:"booking_id""`
}

// GetPropertyPhotosConstruct defines type to store the "getPropertyPhotos" construct.
type GetPropertyPhotosConstruct struct {
	PropertyPhotos *AttachmentArray `json:"propertyPhotos" yaml:"propertyPhotos" schema:"propertyPhotos"  db:"propertyPhotos" csv:"-""`
}

// FormGetOneHelpPostsEndpointConstruct defines type to store the "FormGetOneHelpPostsEndpoint" construct.
type FormGetOneHelpPostsEndpointConstruct struct {
	Id NonStrictInt `json:"id" yaml:"id" schema:"id"  db:"id" csv:"id""`
}

// GetListingsFromPartnerConstruct defines type to store the "getListingsFromPartner" construct.
type GetListingsFromPartnerConstruct struct {
	Description   string           `json:"description" yaml:"description" schema:"description"  db:"description" csv:"description""`
	OrderTypeId   int              `json:"order_type_id" yaml:"order_type_id" schema:"order_type_id"  db:"order_type_id" csv:"order_type_id""`
	GuestLimit    int              `json:"guest_limit" yaml:"guest_limit" schema:"guest_limit"  db:"guest_limit" csv:"guest_limit""`
	ListingId     int              `json:"listing_id" yaml:"listing_id" schema:"listing_id"  db:"listing_id" csv:"listing_id""`
	Zip           int              `json:"zip" yaml:"zip" schema:"zip"  db:"zip" csv:"zip""`
	Images        *AttachmentArray `json:"images" yaml:"images" schema:"images"  db:"images" csv:"-""`
	LandType      string           `json:"land_type" yaml:"land_type" schema:"land_type"  db:"land_type" csv:"land_type""`
	Status        string           `json:"status" yaml:"status" schema:"status"  db:"status" csv:"status""`
	Acres         int              `json:"acres" yaml:"acres" schema:"acres"  db:"acres" csv:"acres""`
	Amenity       *string          `json:"amenity" yaml:"amenity" schema:"amenity"  db:"amenity" csv:"amenity""`
	City          string           `json:"city" yaml:"city" schema:"city"  db:"city" csv:"city""`
	PropertyId    int              `json:"property_id" yaml:"property_id" schema:"property_id"  db:"property_id" csv:"property_id""`
	CreatorId     int              `json:"creator_id" yaml:"creator_id" schema:"creator_id"  db:"creator_id" csv:"creator_id""`
	PartnerSiteId int              `json:"partner_site_id" yaml:"partner_site_id" schema:"partner_site_id"  db:"partner_site_id" csv:"partner_site_id""`
	PropertyName  string           `json:"property_name" yaml:"property_name" schema:"property_name"  db:"property_name" csv:"property_name""`
	Address       string           `json:"address" yaml:"address" schema:"address"  db:"address" csv:"address""`
	State         string           `json:"state" yaml:"state" schema:"state"  db:"state" csv:"state""`
	Price         int              `json:"price" yaml:"price" schema:"price"  db:"price" csv:"price""`
}

// GetListingsFromPartnerEPRConstruct defines type to store the "getListingsFromPartnerEPR" construct.
type GetListingsFromPartnerEPRConstruct struct {
	Data []GetListingsFromPartnerConstruct `json:"Data" yaml:"Data" schema:"Data"  db:"Data" csv:"-""`
}

// GetPropertyListingDataEPRConstruct defines type to store the "getPropertyListingDataEPR" construct.
type GetPropertyListingDataEPRConstruct struct {
	Id            int  `json:"Id" yaml:"Id" schema:"Id"  db:"Id" csv:"Id""`
	GuestLimit    int  `json:"guest_limit" yaml:"guest_limit" schema:"guest_limit"  db:"guest_limit" csv:"guest_limit""`
	PropertyId    int  `json:"property_id" yaml:"property_id" schema:"property_id"  db:"property_id" csv:"property_id""`
	OrderTypeId   *int `json:"order_type_id" yaml:"order_type_id" schema:"order_type_id"  db:"order_type_id" csv:"order_type_id""`
	CreatorId     *int `json:"creator_id" yaml:"creator_id" schema:"creator_id"  db:"creator_id" csv:"creator_id""`
	Price         int  `json:"price" yaml:"price" schema:"price"  db:"price" csv:"price""`
	PartnerSiteId int  `json:"partner_site_id" yaml:"partner_site_id" schema:"partner_site_id"  db:"partner_site_id" csv:"partner_site_id""`
}

// GetHostBookingsConstruct defines type to store the "getHostBookings" construct.
type GetHostBookingsConstruct struct {
	Id                int              `json:"Id" yaml:"Id" schema:"Id"  db:"Id" csv:"Id""`
	BookingDate       time.Time        `json:"booking_date" yaml:"booking_date" schema:"booking_date"  db:"booking_date" csv:"booking_date""`
	CheckoutSessionId string           `json:"checkout_session_id" yaml:"checkout_session_id" schema:"checkout_session_id"  db:"checkout_session_id" csv:"checkout_session_id""`
	GuestId           int              `json:"guest_id" yaml:"guest_id" schema:"guest_id"  db:"guest_id" csv:"guest_id""`
	Status            *string          `json:"status" yaml:"status" schema:"status"  db:"status" csv:"status""`
	Request           *json.RawMessage `json:"request" yaml:"request" schema:"request"  db:"request" csv:"request""`
	ListingId         *int             `json:"listing_id" yaml:"listing_id" schema:"listing_id"  db:"listing_id" csv:"listing_id""`
}

// GetHostBookingConstruct defines type to store the "getHostBooking" construct.
type GetHostBookingConstruct struct {
	BookingId     int    `json:"booking_id" yaml:"booking_id" schema:"booking_id"  db:"booking_id" csv:"booking_id""`
	BookingDate   Date   `json:"booking_date" yaml:"booking_date" schema:"booking_date"  db:"booking_date" csv:"booking_date""`
	GuestId       int    `json:"guest_id" yaml:"guest_id" schema:"guest_id"  db:"guest_id" csv:"guest_id""`
	Request       string `json:"request" yaml:"request" schema:"request"  db:"request" csv:"request""`
	Id            int    `json:"id" yaml:"id" schema:"id"  db:"id" csv:"id""`
	Name          string `json:"name" yaml:"name" schema:"name"  db:"name" csv:"name""`
	Email         string `json:"email" yaml:"email" schema:"email"  db:"email" csv:"email""`
	Acres         int    `json:"acres" yaml:"acres" schema:"acres"  db:"acres" csv:"acres""`
	GuestLimit    int    `json:"guest_limit" yaml:"guest_limit" schema:"guest_limit"  db:"guest_limit" csv:"guest_limit""`
	Price         int    `json:"price" yaml:"price" schema:"price"  db:"price" csv:"price""`
	OrderTypeId   int    `json:"order_type_id" yaml:"order_type_id" schema:"order_type_id"  db:"order_type_id" csv:"order_type_id""`
	State         string `json:"state" yaml:"state" schema:"state"  db:"state" csv:"state""`
	Status        string `json:"status" yaml:"status" schema:"status"  db:"status" csv:"status""`
	PropertyId    int    `json:"property_id" yaml:"property_id" schema:"property_id"  db:"property_id" csv:"property_id""`
	Description   string `json:"description" yaml:"description" schema:"description"  db:"description" csv:"description""`
	Zip           int    `json:"zip" yaml:"zip" schema:"zip"  db:"zip" csv:"zip""`
	LandType      string `json:"land_type" yaml:"land_type" schema:"land_type"  db:"land_type" csv:"land_type""`
	PartnerSiteId int    `json:"partner_site_id" yaml:"partner_site_id" schema:"partner_site_id"  db:"partner_site_id" csv:"partner_site_id""`
	Amenity       string `json:"amenity" yaml:"amenity" schema:"amenity"  db:"amenity" csv:"amenity""`
	Address       string `json:"address" yaml:"address" schema:"address"  db:"address" csv:"address""`
	City          string `json:"city" yaml:"city" schema:"city"  db:"city" csv:"city""`
	CreatorId     int    `json:"creator_id" yaml:"creator_id" schema:"creator_id"  db:"creator_id" csv:"creator_id""`
	FirstName     string `json:"first_name" yaml:"first_name" schema:"first_name"  db:"first_name" csv:"first_name""`
	LastName      string `json:"last_name" yaml:"last_name" schema:"last_name"  db:"last_name" csv:"last_name""`
}

// FormUpdateUsersEndpointConstruct defines type to store the "FormUpdateUsersEndpoint" construct.
type FormUpdateUsersEndpointConstruct struct {
	LastLogInDateTime *time.Time       `json:"last_log_in_date_time" yaml:"last_log_in_date_time" schema:"last_log_in_date_time"  db:"last_log_in_date_time" csv:"last_log_in_date_time""`
	Bio               *NonStrictString `json:"bio" yaml:"bio" schema:"bio"  db:"bio" csv:"bio""`
	PartnerSiteId     *int             `json:"partner_site_id" yaml:"partner_site_id" schema:"partner_site_id"  db:"partner_site_id" csv:"partner_site_id""`
	Id                *NonStrictInt    `json:"Id" yaml:"Id" schema:"Id"  db:"Id" csv:"Id""`
	LastName          *NonStrictString `json:"last_name" yaml:"last_name" schema:"last_name"  db:"last_name" csv:"last_name""`
	Password          *string          `json:"password" yaml:"password" schema:"password"  db:"password" csv:"password""`
	FirstName         *string          `json:"first_name" yaml:"first_name" schema:"first_name"  db:"first_name" csv:"first_name""`
	StripeCustomerId  *NonStrictString `json:"stripe_customer_id" yaml:"stripe_customer_id" schema:"stripe_customer_id"  db:"stripe_customer_id" csv:"stripe_customer_id""`
	TypeOfUser        *string          `json:"type_of_user" yaml:"type_of_user" schema:"type_of_user"  db:"type_of_user" csv:"type_of_user""`
	Email             *NonStrictString `json:"email" yaml:"email" schema:"email"  db:"email" csv:"email""`
	PhoneNumber       *NonStrictString `json:"phone number" yaml:"phone number" schema:"phone number"  db:"phone number" csv:"phone number""`
	ProfilePicture    *AttachmentArray `json:"profilePicture" yaml:"profilePicture" schema:"profilePicture"  db:"profilePicture" csv:"-""`
}

// GetBookingsEPRConstruct defines type to store the "getBookingsEPR" construct.
type GetBookingsEPRConstruct struct {
	BookingDate       Date    `json:"booking_date" yaml:"booking_date" schema:"booking_date"  db:"booking_date" csv:"booking_date""`
	Id                int     `json:"id" yaml:"id" schema:"id"  db:"id" csv:"id""`
	CheckoutSessionId string  `json:"checkout_session_id" yaml:"checkout_session_id" schema:"checkout_session_id"  db:"checkout_session_id" csv:"checkout_session_id""`
	GuestId           int     `json:"guest_id" yaml:"guest_id" schema:"guest_id"  db:"guest_id" csv:"guest_id""`
	Status            *string `json:"status" yaml:"status" schema:"status"  db:"status" csv:"status""`
	ListingId         *int    `json:"listing_id" yaml:"listing_id" schema:"listing_id"  db:"listing_id" csv:"listing_id""`
	FirstName         string  `json:"first_name" yaml:"first_name" schema:"first_name"  db:"first_name" csv:"first_name""`
	LastName          string  `json:"last_name" yaml:"last_name" schema:"last_name"  db:"last_name" csv:"last_name""`
	Name              string  `json:"name" yaml:"name" schema:"name"  db:"name" csv:"name""`
	RequestDetail     *string `json:"request_detail" yaml:"request_detail" schema:"request_detail"  db:"request_detail" csv:"request_detail""`
	CreatorId         int     `json:"creator_id" yaml:"creator_id" schema:"creator_id"  db:"creator_id" csv:"creator_id""`
}

// BookingStatusConstruct defines type to store the "bookingStatus" construct.
type BookingStatusConstruct struct {
	Id     int     `json:"Id" yaml:"Id" schema:"Id"  db:"Id" csv:"Id""`
	Status *string `json:"status" yaml:"status" schema:"status"  db:"status" csv:"status""`
}

// UpdateListingsEPIConstruct defines type to store the "updateListingsEPI" construct.
type UpdateListingsEPIConstruct struct {
	CheckInId   int           `json:"check_in_id" yaml:"check_in_id" schema:"check_in_id"  db:"check_in_id" csv:"check_in_id""`
	CheckOutId  int           `json:"check_out_id" yaml:"check_out_id" schema:"check_out_id"  db:"check_out_id" csv:"check_out_id""`
	GuestLimit  *NonStrictInt `json:"guest_limit" yaml:"guest_limit" schema:"guest_limit"  db:"guest_limit" csv:"guest_limit""`
	OrderTypeId *NonStrictInt `json:"order_type_id" yaml:"order_type_id" schema:"order_type_id"  db:"order_type_id" csv:"order_type_id""`
	Price       *NonStrictInt `json:"price" yaml:"price" schema:"price"  db:"price" csv:"price""`
}

// MessageEPIConstruct defines type to store the "messageEPI" construct.
type MessageEPIConstruct struct {
	SenderId        int        `json:"sender_id" yaml:"sender_id" schema:"sender_id"  db:"sender_id" csv:"sender_id""`
	ReciverId       int        `json:"reciver_id" yaml:"reciver_id" schema:"reciver_id"  db:"reciver_id" csv:"reciver_id""`
	Message         string     `json:"message" yaml:"message" schema:"message"  db:"message" csv:"message""`
	Id              int        `json:"id" yaml:"id" schema:"id"  db:"id" csv:"id""`
	SenderFirstName string     `json:"sender_first_name" yaml:"sender_first_name" schema:"sender_first_name"  db:"sender_first_name" csv:"sender_first_name""`
	SenderLastName  string     `json:"sender_last_name" yaml:"sender_last_name" schema:"sender_last_name"  db:"sender_last_name" csv:"sender_last_name""`
	SentAt          *time.Time `json:"sent_at" yaml:"sent_at" schema:"sent_at"  db:"sent_at" csv:"sent_at""`
}

// AmenitiesArrayConstruct defines type to store the "amenitiesArray" construct.
type AmenitiesArrayConstruct struct {
	Amenity *[]string `json:"amenity" yaml:"amenity" schema:"amenity"  db:"amenity" csv:"amenity""`
}

// FormGetOnePropertyImagesEndpointConstruct defines type to store the "FormGetOnePropertyImagesEndpoint" construct.
type FormGetOnePropertyImagesEndpointConstruct struct {
	Id NonStrictInt `json:"id" yaml:"id" schema:"id"  db:"id" csv:"id""`
}

// PropertyImagesConstruct defines type to store the "property_images" construct.
type PropertyImagesConstruct struct {
	Id         int              `json:"Id" yaml:"Id" schema:"Id"  db:"Id" csv:"Id""`
	FileName   *string          `json:"file_name" yaml:"file_name" schema:"file_name"  db:"file_name" csv:"file_name""`
	File       *AttachmentArray `json:"file" yaml:"file" schema:"file"  db:"file" csv:"-""`
	PropertyId *int             `json:"property_id" yaml:"property_id" schema:"property_id"  db:"property_id" csv:"property_id""`
}

// UpdateUserEmailConstruct defines type to store the "updateUserEmail" construct.
type UpdateUserEmailConstruct struct {
	Email string `json:"email" yaml:"email" schema:"email"  db:"email" csv:"email""`
}

// GetMessageEPRConstruct defines type to store the "getMessageEPR" construct.
type GetMessageEPRConstruct struct {
	SenderId               int              `json:"sender_id" yaml:"sender_id" schema:"sender_id"  db:"sender_id" csv:"sender_id""`
	ReceiverId             int              `json:"receiver_id" yaml:"receiver_id" schema:"receiver_id"  db:"receiver_id" csv:"receiver_id""`
	Message                *string          `json:"message" yaml:"message" schema:"message"  db:"message" csv:"message""`
	SentDate               *time.Time       `json:"sent_date" yaml:"sent_date" schema:"sent_date"  db:"sent_date" csv:"sent_date""`
	Id                     *int             `json:"id" yaml:"id" schema:"id"  db:"id" csv:"id""`
	Email                  string           `json:"email" yaml:"email" schema:"email"  db:"email" csv:"email""`
	SenderFirstName        string           `json:"sender_first_name" yaml:"sender_first_name" schema:"sender_first_name"  db:"sender_first_name" csv:"sender_first_name""`
	ReceiverFirstName      string           `json:"receiver_first_name" yaml:"receiver_first_name" schema:"receiver_first_name"  db:"receiver_first_name" csv:"receiver_first_name""`
	ReceiverLastName       string           `json:"receiver_last_name" yaml:"receiver_last_name" schema:"receiver_last_name"  db:"receiver_last_name" csv:"receiver_last_name""`
	SenderLastName         string           `json:"sender_last_name" yaml:"sender_last_name" schema:"sender_last_name"  db:"sender_last_name" csv:"sender_last_name""`
	ReceiverProfilePicture *AttachmentArray `json:"receiver_profile_picture" yaml:"receiver_profile_picture" schema:"receiver_profile_picture"  db:"receiver_profile_picture" csv:"-""`
	SenderProfilePicture   *AttachmentArray `json:"sender_profile_picture" yaml:"sender_profile_picture" schema:"sender_profile_picture"  db:"sender_profile_picture" csv:"-""`
	ReceiverMessage        *string          `json:"receiver_message" yaml:"receiver_message" schema:"receiver_message"  db:"receiver_message" csv:"receiver_message""`
	SenderMessage          *string          `json:"sender_message" yaml:"sender_message" schema:"sender_message"  db:"sender_message" csv:"sender_message""`
}

// GetConversationEPRConstruct defines type to store the "getConversationEPR" construct.
type GetConversationEPRConstruct struct {
	ConversationData GetMessageEPRConstruct `json:"conversationData" yaml:"conversationData" schema:"conversationData"  db:"conversationData" csv:"conversationData,inline""`
}

// EndpointInputLoginConstruct defines type to store the "endpointInputLogin" construct.
type EndpointInputLoginConstruct struct {
}

// UsersModel defines type to store the "users" model.
type UsersModel struct {
	Id                int              `db:"id" json:"Id" schema:"Id" yaml:"Id"`
	Email             string           `db:"email" json:"email" schema:"email" yaml:"email"`
	Password          string           `db:"password" json:"password" schema:"password" yaml:"password"`
	ProfilePicture    *AttachmentArray `db:"profile_picture" json:"profile_picture" schema:"profile_picture" yaml:"profile_picture"`
	Bio               *string          `db:"bio" json:"bio" schema:"bio" yaml:"bio"`
	FirstName         *string          `db:"first_name" json:"first_name" schema:"first_name" yaml:"first_name"`
	LastName          *string          `db:"last_name" json:"last_name" schema:"last_name" yaml:"last_name"`
	PhoneNumber       *string          `db:"phone_number" json:"phone number" schema:"phone number" yaml:"phone number"`
	LastLogInDateTime time.Time        `db:"last_log_in_date_time" json:"last_log_in_date_time" schema:"last_log_in_date_time" yaml:"last_log_in_date_time"`
	TypeOfUser        *string          `db:"type_of_user" json:"type_of_user" schema:"type_of_user" yaml:"type_of_user"`
	StripeCustomerId  *string          `db:"stripe_customer_id" json:"stripe_customer_id" schema:"stripe_customer_id" yaml:"stripe_customer_id"`
	PartnerSiteId     *int             `db:"partner_site_id" json:"partner_site_id" schema:"partner_site_id" yaml:"partner_site_id"`
}

// PropertiesModel defines type to store the "properties" model.
type PropertiesModel struct {
	Id          int              `db:"id" json:"Id" schema:"Id" yaml:"Id"`
	Name        string           `db:"name" json:"name" schema:"name" yaml:"name"`
	Description string           `db:"description" json:"description" schema:"description" yaml:"description"`
	CreatorId   int              `db:"creator_id" json:"creator_id" schema:"creator_id" yaml:"creator_id"`
	Address     string           `db:"address" json:"address" schema:"address" yaml:"address"`
	City        string           `db:"city" json:"city" schema:"city" yaml:"city"`
	State       string           `db:"state" json:"state" schema:"state" yaml:"state"`
	Zip         int              `db:"zip" json:"zip" schema:"zip" yaml:"zip"`
	Images      *AttachmentArray `db:"images" json:"images" schema:"images" yaml:"images"`
	LandType    string           `db:"land_type" json:"land_type" schema:"land_type" yaml:"land_type"`
	Amenity     *string          `db:"amenity" json:"amenity" schema:"amenity" yaml:"amenity"`
	Status      *string          `db:"status" json:"status" schema:"status" yaml:"status"`
	Acres       int              `db:"acres" json:"acres" schema:"acres" yaml:"acres"`
}

// OrderTypesModel defines type to store the "order_types" model.
type OrderTypesModel struct {
	Id   int    `db:"id" json:"Id" schema:"Id" yaml:"Id"`
	Name string `db:"name" json:"name" schema:"name" yaml:"name"`
}

// AmenitiesModel defines type to store the "amenities" model.
type AmenitiesModel struct {
	Id          int              `db:"id" json:"Id" schema:"Id" yaml:"Id"`
	Name        string           `db:"name" json:"name" schema:"name" yaml:"name"`
	Description *string          `db:"description" json:"description" schema:"description" yaml:"description"`
	Atta        *AttachmentArray `db:"atta" json:"atta" schema:"atta" yaml:"atta"`
}

// ListingsModel defines type to store the "listings" model.
type ListingsModel struct {
	Id             int  `db:"id" json:"Id" schema:"Id" yaml:"Id"`
	GuestLimit     int  `db:"guest_limit" json:"guest_limit" schema:"guest_limit" yaml:"guest_limit"`
	PropertyId     int  `db:"property_id" json:"property_id" schema:"property_id" yaml:"property_id"`
	OrderTypeId    *int `db:"order_type_id" json:"order_type_id" schema:"order_type_id" yaml:"order_type_id"`
	CreatorId      *int `db:"creator_id" json:"creator_id" schema:"creator_id" yaml:"creator_id"`
	Price          int  `db:"price" json:"price" schema:"price" yaml:"price"`
	PartnerSiteId  int  `db:"partner_site_id" json:"partner_site_id" schema:"partner_site_id" yaml:"partner_site_id"`
	CheckInTimeId  *int `db:"check_in_time_id" json:"check_in_time_id" schema:"check_in_time_id" yaml:"check_in_time_id"`
	CheckOutTimeId *int `db:"check_out_time_id" json:"check_out_time_id" schema:"check_out_time_id" yaml:"check_out_time_id"`
}

// ProfitRangesModel defines type to store the "profit_ranges" model.
type ProfitRangesModel struct {
	Id          int     `db:"id" json:"Id" schema:"Id" yaml:"Id"`
	ProfitRange *string `db:"profit_range" json:"profit_range" schema:"profit_range" yaml:"profit_range"`
}

// PartnerSitesModel defines type to store the "partner_sites" model.
type PartnerSitesModel struct {
	Id     int              `db:"id" json:"Id" schema:"Id" yaml:"Id"`
	Name   string           `db:"name" json:"name" schema:"name" yaml:"name"`
	Avatar *AttachmentArray `db:"avatar" json:"avatar" schema:"avatar" yaml:"avatar"`
}

// UsStatesModel defines type to store the "us_states" model.
type UsStatesModel struct {
	Id           int     `db:"id" json:"Id" schema:"Id" yaml:"Id"`
	State        *string `db:"state" json:"State" schema:"State" yaml:"State"`
	Abbreviation *string `db:"abbreviation" json:"Abbreviation" schema:"Abbreviation" yaml:"Abbreviation"`
}

// HelpPostsModel defines type to store the "help_posts" model.
type HelpPostsModel struct {
	Id          int              `db:"id" json:"Id" schema:"Id" yaml:"Id"`
	Title       string           `db:"title" json:"title" schema:"title" yaml:"title"`
	Contents    *json.RawMessage `db:"contents" json:"contents" schema:"contents" yaml:"contents"`
	AuthorId    *int             `db:"author_id" json:"author_id" schema:"author_id" yaml:"author_id"`
	PostImage   *AttachmentArray `db:"post_image" json:"post_image" schema:"post_image" yaml:"post_image"`
	CreatedDate *Date            `db:"created_date" json:"created_date" schema:"created_date" yaml:"created_date"`
	Topic       *string          `db:"topic" json:"topic" schema:"topic" yaml:"topic"`
}

// BookingsModel defines type to store the "bookings" model.
type BookingsModel struct {
	Id                int       `db:"id" json:"Id" schema:"Id" yaml:"Id"`
	ListingId         *int      `db:"listing_id" json:"listing_id" schema:"listing_id" yaml:"listing_id"`
	BookingDate       time.Time `db:"booking_date" json:"booking_date" schema:"booking_date" yaml:"booking_date"`
	CheckoutSessionId string    `db:"checkout_session_id" json:"checkout_session_id" schema:"checkout_session_id" yaml:"checkout_session_id"`
	GuestId           int       `db:"guest_id" json:"guest_id" schema:"guest_id" yaml:"guest_id"`
	Status            *string   `db:"status" json:"status" schema:"status" yaml:"status"`
	RequestDetail     *string   `db:"request_detail" json:"request_detail" schema:"request_detail" yaml:"request_detail"`
}

// ConversationModel defines type to store the "conversation" model.
type ConversationModel struct {
	Id         int        `db:"id" json:"Id" schema:"Id" yaml:"Id"`
	SenderId   *int       `db:"sender_id" json:"sender_id" schema:"sender_id" yaml:"sender_id"`
	ReceiverId *int       `db:"receiver_id" json:"receiver_id" schema:"receiver_id" yaml:"receiver_id"`
	Message    *string    `db:"message" json:"message" schema:"message" yaml:"message"`
	SentDate   *time.Time `db:"sent_date" json:"sent_date" schema:"sent_date" yaml:"sent_date"`
}

// ReservedPropertyDatesModel defines type to store the "reserved_property_dates" model.
type ReservedPropertyDatesModel struct {
	Id          int        `db:"id" json:"Id" schema:"Id" yaml:"Id"`
	PropertyId  int        `db:"property_id" json:"property_id" schema:"property_id" yaml:"property_id"`
	ReservedDay *time.Time `db:"reserved_day" json:"reserved_day" schema:"reserved_day" yaml:"reserved_day"`
}

// StatusTypesModel defines type to store the "status_types" model.
type StatusTypesModel struct {
	Id         int     `db:"id" json:"Id" schema:"Id" yaml:"Id"`
	StatusName *string `db:"status_name" json:"status_name" schema:"status_name" yaml:"status_name"`
}

// CheckInOutTimeModel defines type to store the "check_in_out_time" model.
type CheckInOutTimeModel struct {
	Id   int     `db:"id" json:"Id" schema:"Id" yaml:"Id"`
	Hour *string `db:"hour" json:"Hour" schema:"Hour" yaml:"Hour"`
}

// PropertyImagesModel defines type to store the "property_images" model.
type PropertyImagesModel struct {
	Id         int              `db:"id" json:"Id" schema:"Id" yaml:"Id"`
	FileName   *string          `db:"file_name" json:"file_name" schema:"file_name" yaml:"file_name"`
	File       *AttachmentArray `db:"file" json:"file" schema:"file" yaml:"file"`
	PropertyId *int             `db:"property_id" json:"property_id" schema:"property_id" yaml:"property_id"`
	UserId     *int             `db:"user_id" json:"user_id" schema:"user_id" yaml:"user_id"`
}
