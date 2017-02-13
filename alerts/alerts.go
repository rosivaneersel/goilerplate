package alerts

type Alert struct {
	Title   string
	Class   string
	Message string
}

type Alerts []Alert

func (a *Alerts) New(title string, class string, message string) {
	*a = append(*a, Alert{title, class, message})
}

func (a *Alerts) Get() []Alert {
	c := *a
	*a = nil
	return c
}
