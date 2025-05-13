package gsdb

func (db *Database) StartCounter(counterName string) {
	db.Counters.Lock.Lock()
	defer db.Counters.Lock.Unlock()
	db.Counters.Count[counterName] = 0
}

func (db *Database) GetCounter(counterName string) int64 {
	db.Counters.Lock.Lock()
	defer db.Counters.Lock.Unlock()
	return db.Counters.Count[counterName]
}

func (db *Database) IncCounter(counterName string) {
	db.Counters.Lock.Lock()
	defer db.Counters.Lock.Unlock()

	if _, ok := db.Counters.Count[counterName]; ok {
		db.Counters.Count[counterName]++
	}
}
