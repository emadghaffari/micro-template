package config

type SMS struct {
	Service string      // with service need this action
	Token   string      // correlationID for tracing
	Data    interface{} // data is attrs for message
	To      string      // phone
}

type EMAIL struct {
	Service  string      // with service need this action
	Token    string      // correlationID for tracing
	Data     interface{} // data is attrs for message
	To       []string    // emails
	BCC      []string    // bcc
	Title    string      // title for email
	SubTitle string      // subtitles
}
