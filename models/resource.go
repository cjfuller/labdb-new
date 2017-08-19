package models

import "time"

type FieldDef struct {
	Name   string
	Lookup string
	Type   string
}

type InfoSection struct {
	Name         string
	Preformatted bool
	Lookup       string
	Single       bool
	Fields       []FieldDef
	InlineValue  *string
}

type SequenceInfo struct {
	Sequence FieldDef
	Verified FieldDef
}

type CoreLinks struct {
	Lookup string
	Name   string
	Links  []string
}

type CoreInfoSections struct {
}

// TODO(colin): json labels
type ResourceDef struct {
	Type               string
	ID                 int
	Timestamp          time.Time
	FieldData          Entity
	ResourcePath       string
	Name               string
	ShortDesc          string
	CoreLinks          *CoreLinks    // TODO(colin): type
	CoreInfoSections   []InfoSection // TODO(colin): type
	SequenceInfo       *SequenceInfo // TODO(colin): type
	SupplementalFields []FieldDef    // TODO(colin): type
}

func AsResourceDef(e Entity) ResourceDef {
	return ResourceDef{}
}
