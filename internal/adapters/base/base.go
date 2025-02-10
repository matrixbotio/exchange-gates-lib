package baseadp

type AdapterBase struct {
	ExchangeID int
	Name       string
	Tag        string
}

func NewAdapterBase(id int, name, tag string) AdapterBase {
	return AdapterBase{
		ExchangeID: id,
		Name:       name,
		Tag:        tag,
	}
}

func (a *AdapterBase) GetTag() string {
	return a.Tag
}

func (a *AdapterBase) GetID() int {
	return a.ExchangeID
}

func (a *AdapterBase) GetName() string {
	return a.Name
}
