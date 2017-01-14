package profiles

type Profile struct {
	ID uint
	Username string
	Email  string
	Password string
}

type ProfileManager interface {
	Get(id uint) (*Profile, error)
	// etc
}

type Profiles struct {
	ProfileManager
}

func PageController(m ProfileManager) *Profiles {
	return &Profiles{ProfileManager: m}
}