package index

type IndexBusiness struct {
	indexSpace int
}

func NewIndexBusiness() *IndexBusiness {
	return &IndexBusiness{
		indexSpace: 1000,
	}
}

func (ib *IndexBusiness) GetIndexSpace() int {
	return ib.indexSpace
}

func (ib *IndexBusiness) ConvertSentenceIdex(idx int) int {
	return (idx + 1) * ib.indexSpace
}
