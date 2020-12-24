package sequence

type Sequence interface {
	NextId(name string) (nextId int64, err error)
}
