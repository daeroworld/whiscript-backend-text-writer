package index

type IndexBusiness struct {
}

func NewIndexBusiness() *IndexBusiness {
	return &IndexBusiness{}
}

func (ib *IndexBusiness) ConvertSentenceIdex(idx int) int {
	return (idx + 1) * 1000
}
