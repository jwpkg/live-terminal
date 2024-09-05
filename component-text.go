package living_terminal

type LivingText struct {
	Text            string
	requestReRender chan bool
}

func NewLivingText(text string) *LivingText {
	return &LivingText{
		Text: text,
	}
}

func (livingText *LivingText) Update(text string) {
	livingText.Text = text
	if livingText.requestReRender != nil {
		livingText.requestReRender <- true
	}
}

func (livingText *LivingText) Init(requestReRender chan bool) {
	livingText.requestReRender = requestReRender
}

func (livingText *LivingText) Render() string {
	return string(livingText.Text)
}

func (livingText *LivingText) Finish() {
	livingText.requestReRender = nil
}
