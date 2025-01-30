package scripts

import (
  "database/sql"
  "fmt"
	_ "github.com/lib/pq"
)

func RunPrecomputation() {
	var config Config = *GetConfig()
	psqlconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName)

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	defer db.Close()

	err = db.Ping()
	CheckError(err)
	fmt.Println("Connected to db.")

	fmt.Println("Creating route_edges table...")
	_, err = db.Exec(st1)
	CheckError(err)

	fmt.Println("Filling route_edges with data...")
	_, err = db.Exec(st2)
	CheckError(err)

	fmt.Println("Creating precomputed_travel_times table...")
	_, err = db.Exec(st3)
	CheckError(err)

  fmt.Println("Done.")
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
