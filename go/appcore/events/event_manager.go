package events

type EventManager struct {
	db *DB
}

func NewEventManager(dataDir string) (*EventManager, error) {
	db, err := NewDB(dataDir)
	if err != nil {
		return nil, err
	}

	return &EventManager{
		db: db,
	}, nil
}
