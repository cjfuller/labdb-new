package models

type Sample struct {
	Model
	EnteredBy   string
	Description string
	SampleAlias string
}

func (s *Sample) OwnerFieldName() string { return "entered_by" }
func (s *Sample) ShortDesc() string      { return s.SampleAlias }
func (s *Sample) Desc() string           { return s.Description }
