package kernel

type Config struct {
	App struct {
		Host string `default:"0.0.0.0"`
		Port string `default:"8080"`
	}
	Redis struct {
		Host string `default:"redis"`
		Port string `default:"6379"`
	}
	MySql struct {
		User     string `default:"root"`
		Password string `default:"my_secret"`
		database string `default:"user_manager"`
	}
}
