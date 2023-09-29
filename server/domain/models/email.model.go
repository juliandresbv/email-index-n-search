package models

type EmailModel struct {
	Id        string `json:"id,omitempty"`
	Bcc       string `json:"bcc,omitempty"`
	Body      string `json:"body,omitempty"`
	Cc        string `json:"cc,omitempty"`
	Date      string `json:"date,omitempty"`
	From      string `json:"from,omitempty"`
	MessageId string `json:"messageId,omitempty"`
	Subject   string `json:"subject,omitempty"`
	To        string `json:"to,omitempty"`
	XBcc      string `json:"xBcc,omitempty"`
	XCc       string `json:"xCc,omitempty"`
	XFileName string `json:"xFileName,omitempty"`
	XFolder   string `json:"xFolder,omitempty"`
	XFrom     string `json:"xFrom,omitempty"`
	XOrigin   string `json:"xOrigin,omitempty"`
	XTo       string `json:"xTo,omitempty"`
}
