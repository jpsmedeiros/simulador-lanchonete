package simulation

type System struct {
	FC  int //Fila de cliente
	FH  int //Fila de hamburguer
	FR  int //Fila normal de refrigerante
	FPR int //Fila de prioridades de refrigerante
}

func (s *System) AddClient() {
	s.FC = s.FC + 1
}

func (s *System) AddHamburguer() {
	s.FH = s.FH + 1
}
func (s *System) AddSoda() {
	s.FR = s.FR + 1
}

func (s *System) AddSodaP() {
	s.FPR = s.FPR + 1
}

func (s *System) RemoveClient() {
	s.FC = s.FC - 1
}

func (s *System) RemoveHamburguer() {
	s.FH = s.FH - 1
}

func (s *System) RemoveSoda() {
	s.FR = s.FR - 1
}

func (s *System) RemoveSodaP() {
	s.FPR = s.FPR - 1
}
