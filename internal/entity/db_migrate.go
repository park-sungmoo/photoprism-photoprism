package entity

// MigrateDb creates database tables and inserts default fixtures as needed.
func MigrateDb(dropDeprecated, runFailed bool) {
	if dropDeprecated {
		DeprecatedTables.Drop(Db())
	}

	Entities.Migrate(Db(), runFailed)
	Entities.WaitForMigration(Db())

	CreateDefaultFixtures()
}

// InitTestDb connects to and completely initializes the test database incl fixtures.
func InitTestDb(driver, dsn string) *Gorm {
	if HasDbProvider() {
		return nil
	}

	if driver == "test" || driver == "sqlite" || driver == "" || dsn == "" {
		driver = "sqlite3"
		dsn = ".test.db"
	}

	log.Infof("initializing %s test db in %s", driver, dsn)

	db := &Gorm{
		Driver: driver,
		Dsn:    dsn,
	}

	SetDbProvider(db)
	ResetTestFixtures()

	return db
}
