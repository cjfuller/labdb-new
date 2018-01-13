package models

type Yeaststrain struct {
	Model
	EnteredBy   string
	Comments    string
	Strainalias string
	Sequence    string
}

func (y *Yeaststrain) OwnerFieldName() string { return "entered_by" }
func (y *Yeaststrain) ShortDesc() string      { return y.Strainalias }
func (y *Yeaststrain) Desc() string           { return y.Comments }
func (y *Yeaststrain) GetSequence() string    { return y.Sequence }
