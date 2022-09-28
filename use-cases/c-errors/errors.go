package cErrors

type Custom struct {
	Property     string
	MessageID    string
	TemplateData map[string]interface{}
	message      string
}

func (c *Custom) SetLocalesErrMsg(msg string) {
	c.message = msg
}

func (c *Custom) Error() string {
	return c.message
}
