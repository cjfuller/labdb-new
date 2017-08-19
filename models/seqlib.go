package models

type SeqLib struct {
	Model
	Genome           string
	Method           string
	EnteredBy        string
	Project          string
	StorageLocation  string
	Concentration    float64
	SizeDistribution string
	IndexID          string
	IndexSeq         string `labdb_role:"Sequence"`
	Description      string
	LinkedItems      string
	Alias            string
	Notebook         string
	Number           int
}

func (r *SeqLib) AutoFill(userName string) {
	r.EnteredBy = userName
	r.Number = NextAvailableNumber(r)
}

func (r *SeqLib) GetNumber() int {
	return r.Number
}
