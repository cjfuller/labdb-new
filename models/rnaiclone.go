package models

import (
	"strconv"
)

type RNAiClone struct {
	Model
	Number          int
	Alias           string
	Notebook        string
	Description     string
	EnteredBy       string
	SequenceName    string
	Library         string
	HostStrain      string // TODO: default is HT115
	PlasmidBackbone string // TODO: default is L4440
	Antibiotic      string // TODO: default is Amp
	Location        string
	Sequenced       bool
}

func (r *RNAiClone) GetCoreInfoSections() []InfoSection {
	return []InfoSection{
		InfoSection{
			Name:   "Description",
			Lookup: "Description",
			Single: true,
		},
		InfoSection{
			Name: "Library Information",
			Fields: []FieldDef{
				FieldDef{
					Name:   "Library",
					Lookup: "Library",
					Type:   "value",
				},
				FieldDef{
					Name:   "Host strain",
					Lookup: "HostStrain",
					Type:   "value",
				},
				FieldDef{
					Name:   "Plasmid backbone",
					Lookup: "PlasmidBackbone",
					Type:   "value",
				},
				FieldDef{
					Name:   "Antibiotic",
					Lookup: "Antibiotic",
					Type:   "value",
				},
				FieldDef{
					Name:   "Location",
					Lookup: "Location",
					Type:   "value",
				},
				FieldDef{
					Name:   "Sequenced?",
					Lookup: "Sequenced",
					Type:   "boolean",
				},
			},
		},
	}
}

func (r *RNAiClone) GetSequenceInfo() *SequenceInfo { return nil }
func (r *RNAiClone) GetSupplementalFields() []FieldDef {
	return []FieldDef{
		FieldDef{
			Name:   "Entered by",
			Lookup: "EnteredBy",
			Type:   "value",
		},
		FieldDef{
			Name:   "Date entered",
			Lookup: "CreatedAt",
			Type:   "value",
		},
	}
}

func (r *RNAiClone) AutoFill(userName string) {
	r.EnteredBy = userName
	r.HostStrain = "HT115"
	r.PlasmidBackbone = "L4440"
	r.Antibiotic = "Amp"
	r.Number = NextAvailableNumber(r)
}

func (r RNAiClone) TableName() string {
	return "rnai_clones"
}

func (r *RNAiClone) GetNumber() int {
	return r.Number
}

func (r *RNAiClone) Kind() string {
	return "rnai_clone"
}

func (r *RNAiClone) GetName() string {
	// TODO(colin): naming module
	return "RNAiC" + strconv.Itoa(r.GetNumber())
}

func (r *RNAiClone) OwnerFieldName() string { return "entered_by" }
func (r *RNAiClone) ShortDesc() string {
	return r.Alias
}

func (r *RNAiClone) Desc() string {
	return r.Description
}

func (r *RNAiClone) GetCoreLinks() *CoreLinks { return nil }
