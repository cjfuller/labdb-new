package models

// TODO(colin): other fields
type Plasmid struct {
	Model
	Alias       string
	Description string
	Sequence    string
	Creator     string
}

func (p *Plasmid) OwnerFieldName() string { return "creator" }
func (p *Plasmid) ShortDesc() string      { return p.Alias }
func (p *Plasmid) Desc() string           { return p.Description }
func (p *Plasmid) GetSequence() string    { return p.Sequence }
