package models

type Line struct {
	Model
	EnteredBy   string
	Description string
	LineAlias   string
	Sequence    string
}

func (l *Line) OwnerFieldName() string { return "entered_by" }
func (l *Line) ShortDesc() string      { return l.LineAlias }
func (l *Line) Desc() string           { return l.Description }
func (l *Line) GetSequence() string    { return l.Sequence }
