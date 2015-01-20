package common

type Phase interface {
	Id() string
}

type commonPhase struct {
	id string
}

func NewPhase(id string) Phase {
	return &commonPhase{
		id: id,
	}
}

func (p *commonPhase) Id() string {
	return p.id
}
