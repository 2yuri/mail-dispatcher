package domain

import "strings"

type MessageAttach struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type Message struct {
	Sender   string `json:"sender,omitempty"`
	Subject  string `json:"subject,omitempty"`
	HTMLBody string `json:"html_body,omitempty"`
	TextBody string `json:"text_body,omitempty"`

	//destinations
	Recipient    string `json:"recipient,omitempty"`
	CCAddresses  string `json:"cc,omitempty"`
	BCCAddresses string `json:"bcc,omitempty"`
	Attachments  []MessageAttach
}

func (dto *Message) CharSet() string {
	return "UTF-8"
}

func (dto *Message) CCRecipient() []string {
	return strings.Split(dto.CCAddresses, ";")
}

func (dto *Message) BCCRecipient() []string {
	return strings.Split(dto.BCCAddresses, ";")
}

func (dto *Message) ToRecipient() []string {
	return strings.Split(dto.Recipient, ";")
}
