package domain

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"restAPI/errs"
	"time"
)

type CustomerRepositoryDb struct {
	client *sql.DB
}

func (c CustomerRepositoryDb) GetAll(status string) ([]Customer, *errs.AppError) {
	var findAllSQL string
	var rows *sql.Rows
	var err error
	if status == "" {
		findAllSQL = "SELECT customer_id , name , city , zipcode ,date_of_birth , status FROM customers"
		rows, err = c.client.Query(findAllSQL)
	} else {
		findAllSQL = "select * from customers where status = ?"
		rows, err = c.client.Query(findAllSQL, status)
	}

	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, errs.NewUnExpectedError("Error while fetching data from database")
	}
	customers := make([]Customer, 0)
	for rows.Next() {
		var customer Customer
		err = rows.Scan(&customer.Id, &customer.Name, &customer.City, &customer.Zipcode,
			&customer.DateOfBirth, &customer.Status)
		if err != nil {
			log.Printf("Error: %v\n", err)
			return nil, errs.NewUnExpectedError("Error while getting customers")
		}
		customers = append(customers, customer)
	}
	return customers, nil
}

func (c CustomerRepositoryDb) GetById(id int64) (*Customer, *errs.AppError) {
	findByIdSQL := "SELECT customer_id , name , city , zipcode ,date_of_birth , status FROM customers WHERE customer_id = ?"
	row := c.client.QueryRow(findByIdSQL, id)
	var customer Customer
	err := row.Scan(&customer.Id, &customer.Name, &customer.City, &customer.Zipcode,
		&customer.DateOfBirth, &customer.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Customer not found")
		} else {
			log.Printf("Error: %v\n", err)
			return nil, errs.NewUnExpectedError("Error while getting customer")
		}
	}
	return &customer, nil
}

func NewCustomerRepositoryDb() CustomerRepositoryDb {
	db, err := sql.Open("mysql", "root:1@tcp(localhost:3306)/banking")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return CustomerRepositoryDb{db}
}
