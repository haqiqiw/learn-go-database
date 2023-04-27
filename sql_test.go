package learn_go_database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestExecSql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	query := "INSERT INTO customer(id, name) VALUES('mamank', 'Mamank')"
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert new customer")
}

func TestQuerySql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	query := "SELECT id, name FROM customer"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, name string
		err := rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}

		fmt.Println("Id:", id)
		fmt.Println("Name:", name)
	}
}

func TestQuerySqlComplex(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	query := "SELECT id, name, email, balance, rating, birth_date, married, created_at FROM customer"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, name string
		var email sql.NullString
		var balance int32
		var rating float64
		var birthDate sql.NullTime
		var createdAt time.Time
		var married bool
		err := rows.Scan(&id, &name, &email, &balance,
			&rating, &birthDate, &married, &createdAt)
		if err != nil {
			panic(err)
		}

		fmt.Println("===========")
		fmt.Println("Id:", id)
		fmt.Println("Name:", name)
		if email.Valid {
			fmt.Println("Email:", email.String)
		}
		fmt.Println("Balance:", balance)
		fmt.Println("Rating:", rating)
		if birthDate.Valid {
			fmt.Println("Birth Date:", birthDate.Time)
		}
		fmt.Println("Married:", married)
		fmt.Println("Created At:", createdAt)
	}
}

func TestQuerySqlInjection(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	// username := "admin'; #" // mysql injection
	username := "admin"
	password := "admin"

	ctx := context.Background()

	query := "SELECT username FROM user where username = '" + username +
		"' AND password = '" + password + "' LIMIT 1"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		panic(err)
	}

	if rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}

		fmt.Println("Username:", username)
	} else {
		fmt.Println("Username not found")
	}
}

func TestQuerySqlInjectionSafe(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	// username := "admin'; #"
	username := "admin"
	password := "admin"

	ctx := context.Background()

	query := "SELECT username FROM user where username = ? AND password = ? LIMIT 1"
	rows, err := db.QueryContext(ctx, query, username, password)
	if err != nil {
		panic(err)
	}

	if rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}

		fmt.Println("Username:", username)
	} else {
		fmt.Println("Username not found")
	}
}

func TestExecSqlParameter(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	// username := "wawan'; DROP TABLE user; #"
	username := "wawan"
	password := "wawan"

	ctx := context.Background()

	query := "INSERT INTO user(username, password) VALUES(?, ?)"
	_, err := db.ExecContext(ctx, query, username, password)
	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert new user")
}

func TestAutoIncrement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	email := "wawan@gmail.com"
	comment := "ini comment"

	ctx := context.Background()

	query := "INSERT INTO comments(email, comment) VALUES(?, ?)"
	result, err := db.ExecContext(ctx, query, email, comment)
	if err != nil {
		panic(err)
	}

	inserId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert new comment with id", inserId)
}

func TestPrepareStatement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	query := "INSERT INTO comments(email, comment) VALUES(?, ?)"
	statement, err := db.PrepareContext(ctx, query)
	if err != nil {
		panic(err)
	}
	defer statement.Close()

	for i := 0; i < 10; i++ {
		email := "wawan" + strconv.Itoa(i) + "@gmail.com"
		comment := "ini comment ke " + strconv.Itoa(i)

		result, err := statement.ExecContext(ctx, email, comment)
		if err != nil {
			panic(err)
		}

		lastInsertId, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}
		fmt.Println("Comment id:", lastInsertId)
	}
}

func TestTransaction(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	query := "INSERT INTO comments(email, comment) VALUES(?, ?)"
	for i := 0; i < 10; i++ {
		email := "wawan" + strconv.Itoa(i) + "@gmail.com"
		comment := "ini comment ke " + strconv.Itoa(i)

		result, err := tx.ExecContext(ctx, query, email, comment)
		if err != nil {
			panic(err)
		}

		lastInsertId, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}
		fmt.Println("Comment id:", lastInsertId)
	}

	err = tx.Commit()
	// err = tx.Rollback()
	if err != nil {
		panic(err)
	}
}
