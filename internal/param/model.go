package param

type Param struct {
	Key   string `db:"key"`
	Value string `db:"value"`
}
