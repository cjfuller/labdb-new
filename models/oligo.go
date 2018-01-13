package models

type Oligo struct {
	Model
	EnteredBy  string
	Oligoalias string
	Purpose    string
	Sequence   string
}

func (o *Oligo) OwnerFieldName() string { return "entered_by" }
func (o *Oligo) ShortDesc() string      { return o.Oligoalias }
func (o *Oligo) Desc() string           { return o.Purpose }
func (o *Oligo) GetSequence() string    { return o.Sequence }
