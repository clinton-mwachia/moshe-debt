package utils

import (
	"database/sql"
	"moshe-debt/models"
)

func LoadDebts(db *sql.DB, debts []models.Debt) {
	debts = nil
	rows, _ := db.Query("SELECT id,customer,phone,amount,balance FROM debts")
	defer rows.Close()
	for rows.Next() {
		var d models.Debt
		rows.Scan(&d.ID, &d.Customer, &d.Phone, &d.Amount, &d.Balance)
		debts = append(debts, d)
	}
}

func LoadPayments(db *sql.DB, payments []models.Payment) {
	payments = nil
	rows, _ := db.Query("SELECT id,customer,amount,balance FROM payments")
	defer rows.Close()
	for rows.Next() {
		var p models.Payment
		rows.Scan(&p.ID, &p.Customer, &p.Amount, &p.Balance)
		payments = append(payments, p)
	}
}
