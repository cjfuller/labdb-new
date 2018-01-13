package models

type Bacterium struct {
	Model
	EnteredBy   string
	Comments    string
	Strainalias string
	Sequence    string
}

func (b *Bacterium) OwnerFieldName() string { return "entered_by" }
func (b *Bacterium) ShortDesc() string      { return b.Strainalias }
func (b *Bacterium) Desc() string           { return b.Comments }
func (b *Bacterium) GetSequence() string    { return b.Sequence }
