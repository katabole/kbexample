package actions

type Environment string

const (
	DevelopmentEnvironment Environment = "development"
	TestEnvironment                    = "test"
	ProductionEnvironment              = "production"
)

func (e Environment) IsProduction() bool {
	return e != DevelopmentEnvironment && e != TestEnvironment
}
