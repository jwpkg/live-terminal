package living_terminal

type LivingComponent interface {
	Init(requestReRender chan bool)
	Render() string
	Finish()
}
