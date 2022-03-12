package model

import "time"

const (
	EVENT_DATA   = "event-data"
	EVENT        = "event"
	IP           = "ip"
	MESSAGE      = "message"
	MESSAGE_ID   = "message-id"
	HEADERS      = "headers"
	MAILGUN      = "mailgun"
	RECIPIENT    = "recipient"
	GEO_LOCATION = "geolocation"
	COUNTRY      = "country"
	CAMPAIGNS    = "campaigns"
	TAGS         = "tags"
)

// model struct for keeping webhook information
type Callback struct {
	// invocation details
	MessageID  string
	Event      string
	Provider   string
	ReceivedAt time.Time
	Recipient  string
	// destination details
	IsQueued   bool
	IsMessaged bool
	IsFiled    bool
	// network details
	IP  string
	Geo string
	// promotion details
	Campaigns []string
	Tags      []string
}

// 'Domain Object Model pattern'
// for parsing information received on the event body.
func (log *Callback) ParseMailgun(eventBody map[string]interface{}) error {
	eventData := eventBody[EVENT_DATA].(map[string]interface{})
	message := eventData[MESSAGE].(map[string]interface{})
	headers := message[HEADERS].(map[string]interface{})
	geoLocation := eventData[GEO_LOCATION].(map[string]interface{})
	// invocation details
	log.MessageID = headers[MESSAGE_ID].(string)
	log.Event = eventData[EVENT].(string)
	log.Provider = MAILGUN
	// save in utc & apply time difference for localization in reports
	log.ReceivedAt = time.Now().UTC()
	log.Recipient = eventData[RECIPIENT].(string)
	// network details
	log.IP = eventData[IP].(string)
	log.Geo = geoLocation[COUNTRY].(string)
	// promotion details: campaign array mapping
	campaignIF := eventData[CAMPAIGNS].([]interface{})
	campaigns := []string{}
	for _, name := range campaignIF {
		campaigns = append(campaigns, name.(string))
	}
	log.Campaigns = campaigns
	// tags
	tagIF := eventData[TAGS].([]interface{})
	tags := []string{}
	for _, name := range tagIF {
		tags = append(tags, name.(string))
	}
	log.Tags = tags
	return nil
}

// TODO: implement twilio event and adhere to common callback model
func (log *Callback) ParseTwilio(event map[string]interface{}) error {
	return nil
}
