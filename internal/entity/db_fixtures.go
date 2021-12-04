package entity

// CreateDefaultFixtures inserts default fixtures for test and production.
func CreateDefaultFixtures() {
	CreateUnknownAddress()
	CreateDefaultUsers()
	CreateUnknownPlace()
	CreateUnknownLocation()
	CreateUnknownCountry()
	CreateUnknownCamera()
	CreateUnknownLens()
}

// ResetTestFixtures re-creates registered database tables and inserts test fixtures.
func ResetTestFixtures() {
	Entities.Migrate(Db(), false)
	Entities.WaitForMigration(Db())
	Entities.Truncate(Db())

	CreateDefaultFixtures()

	CreateTestFixtures()
}
