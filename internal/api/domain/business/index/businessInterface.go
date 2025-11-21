package index

type IIndexBusiness interface {
	GetIndexSpace() int
	ConvertSentenceIdex(int) int
}
