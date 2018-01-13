package models

type Antibody struct {
	Model
	EnteredBy string
	Alias     string
	Comments  string
}

func (a *Antibody) OwnerFieldName() string { return "entered_by" }
func (a *Antibody) ShortDesc() string      { return a.Alias }
func (a *Antibody) Desc() string           { return a.Comments }
