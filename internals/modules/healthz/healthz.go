package healthz

type Usecase interface {
	CheckHealth() (map[string]string, error)
}

type Repository interface {
	PingDatabase() (string, error)
}
