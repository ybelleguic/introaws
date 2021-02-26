package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"

	"os/exec"
)

type Pet struct {
	Name string
	Owner string
	Species string 
	Sex string 
	BirthDate string 
	DeathDate sql.NullString
}

func getConnByPassword(user, password, endpoint string) *sql.DB {
	uri := fmt.Sprintf("%s:%s@tcp(%s:3306)/menagerie", user, password, endpoint)
	db, err := sql.Open("mysql", uri)

    // if there is an error opening the connection, handle it
    if err != nil {
        panic(err.Error())
	}
	
	return db
}

func ReadDatabase(user, password, authMethod, endpoint string) []Pet {
	endpoint = "database-2.cluster-c2oaznzpz4np.us-east-1.rds.amazonaws.com"
	args := []string{"rds", "generate-db-auth-token", "--hostname", endpoint, "--port", "3306", "--username", "tutu", "--region", "us-east-1"}
	cmd := exec.Command("aws",args...)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("error!")
	}
	password = string(stdoutStderr)
	fmt.Println(password)

    dsn := fmt.Sprintf("%s:%s@tcp(%s)/menagerie?allowCleartextPasswords=true&tls=true",
        "tutu", password, endpoint,
	)
	db, err := sql.Open("mysql", dsn)
    if err != nil {
        panic(err)
    }




	ret := []Pet{}
	// var db *sql.DB
	// if (authMethod == "password") {
	// 	db = getConnByPassword(user, password, endpoint)
	// }

	results, err := db.Query("SELECT * FROM pet")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }

    for results.Next() {
        var pet Pet
        // for each row, scan the result into our tag composite object
        err = results.Scan(&pet.Name, &pet.Owner, &pet.Species, &pet.Sex, &pet.BirthDate, &pet.DeathDate)
        if err != nil {
            panic(err.Error()) // proper error handling instead of panic in your app
        }
				// and then print out the tag's Name attribute
		ret = append(ret, pet)
    }

	db.Close()
	return ret
}

