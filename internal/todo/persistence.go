package todo

// Repo persists a list of [Item] values. Implementations are typically file or DB
// adapters (Dependency Inversion: app depends on this port, not on JSON/OS details).
type Repo interface {
	Load() ([]Item, error)
	Save(items []Item) error
}
