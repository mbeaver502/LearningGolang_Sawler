module myapp

go 1.19

replace github.com/mbeaver502/LearningGolang_Sawler/go_laravel/celeritas => ../celeritas

require github.com/mbeaver502/LearningGolang_Sawler/go_laravel/celeritas v0.0.0-00010101000000-000000000000

require (
	github.com/go-chi/chi/v5 v5.0.7 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
)
