package httpserver //nolint:typecheck

type DataBaser interface {
	SetCounter(key string, value int64)
	SetGauge(key string, value float64)
	GetCounter(name string) (value int64, err error)
	GetGauge(name string) (value float64, err error)
	GetCounters() map[string]int64
	GetGauges() map[string]float64
	Ping() error
}

type DataBase struct {
	DataBaser
}

func NewDataBase(db DataBaser) *DataBase {
	return &DataBase{db}
}
