package info

type Contributor interface {
	GetInfo() (info interface{}, err error)
}

type ContributorAdapter func() (interface{}, error)

func (ca ContributorAdapter) GetInfo() (interface{}, error) {
	return ca()
}
