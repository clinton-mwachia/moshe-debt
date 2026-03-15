package main

import (
	"database/sql"
	"fmt"
	"moshe-debt/models"
	"moshe-debt/tables"
	"moshe-debt/utils"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Global slices for UI
var debts []models.Debt
var payments []models.Payment
var debtTable, paymentTable *widget.Table

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./debt_manager.db")
	if err != nil {
		panic(err)
	}
	utils.InitDB(db)

	a := app.NewWithID("moshe.debt")
	w := a.NewWindow("Debt & Payments Management - Moshe Crafts")
	w.Resize(fyne.NewSize(900, 600))

	// --- Tabs ---
	tabs := container.NewAppTabs(
		container.NewTabItem("Debts", buildDebtTab()),
		container.NewTabItem("Payments", buildPaymentTab()),
		container.NewTabItem("Contact", buildContactTab()),
	)

	w.SetContent(tabs)
	w.ShowAndRun()
}

// ---------------------- Debt Tab ----------------------
func buildDebtTab() fyne.CanvasObject {
	utils.LoadDebts(db, debts)

	addDebtBtn := widget.NewButton("Add Debt", func() { showDebtDialog(nil) })

	debtTable = tables.BuildDebtTable(db, debts, debtTable)
	debtContainer := container.NewVScroll(debtTable)
	debtContainer.SetMinSize(fyne.NewSize(600, 300))

	return container.NewVBox(
		container.NewHBox(addDebtBtn),
		debtContainer,
	)
}

func showDebtDialog(d *models.Debt) {
	a := fyne.CurrentApp()
	w := a.Driver().AllWindows()[0]

	customerEntry := widget.NewEntry()
	phoneEntry := widget.NewEntry()
	amountEntry := widget.NewEntry()

	if d != nil {
		customerEntry.SetText(d.Customer)
		phoneEntry.SetText(d.Phone)
		amountEntry.SetText(fmt.Sprintf("%.2f", d.Amount))
	}

	form := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("Customer", customerEntry),
			widget.NewFormItem("Phone", phoneEntry),
			widget.NewFormItem("Amount", amountEntry),
		},
		OnSubmit: func() {
			amt, err := strconv.ParseFloat(amountEntry.Text, 64)
			if err != nil {
				dialog.ShowError(fmt.Errorf("invalid amount"), w)
				return
			}

			if d == nil {
				res, _ := db.Exec("INSERT INTO debts(customer,phone,amount,balance) VALUES(?,?,?,?)",
					customerEntry.Text, phoneEntry.Text, amt, amt)
				id, _ := res.LastInsertId()
				debts = append(debts, models.Debt{ID: int(id), Customer: customerEntry.Text, Phone: phoneEntry.Text, Amount: amt, Balance: amt})
			} else {
				_, _ = db.Exec("UPDATE debts SET customer=?,phone=?,amount=? WHERE id=?",
					customerEntry.Text, phoneEntry.Text, amt, d.ID)
				loadDebts()
			}
			debtTable.Refresh()
		},
	}

	dlg := dialog.NewCustom("Debt Entry", "Cancel", form, w)
	dlg.Resize(fyne.NewSize(400, 250))
	dlg.Show()
}

// ---------------------- Payments Tab ----------------------
func buildPaymentTab() fyne.CanvasObject {
	utils.LoadPayments(db, payments)

	addPaymentBtn := widget.NewButton("Add Payment", func() { showPaymentDialog(nil) })

	paymentTable = tables.BuildPaymentTable(db, payments, debts, paymentTable)
	paymentContainer := container.NewVScroll(paymentTable)
	paymentContainer.SetMinSize(fyne.NewSize(600, 300))

	return container.NewVBox(
		container.NewHBox(addPaymentBtn),
		paymentContainer,
	)
}

func loadDebts() {
	utils.LoadDebts(db, debts)
}

func loadPayments() {
	utils.LoadPayments(db, payments)
}

func showPaymentDialog(p *models.Payment) {
	a := fyne.CurrentApp()
	w := a.Driver().AllWindows()[0]

	customerOptions := []string{}
	for _, d := range debts {
		customerOptions = append(customerOptions, d.Customer)
	}
	customerEntry := widget.NewSelect(customerOptions, nil)
	customerEntry.PlaceHolder = "Select customer"

	amountEntry := widget.NewEntry()
	balanceEntry := widget.NewEntry()
	balanceEntry.Disable()

	if p != nil {
		customerEntry.SetSelected(p.Customer)
		amountEntry.SetText(fmt.Sprintf("%.2f", p.Amount))
		balanceEntry.SetText(fmt.Sprintf("%.2f", p.Balance))
	}

	customerEntry.OnChanged = func(sel string) {
		for _, d := range debts {
			if d.Customer == sel {
				balanceEntry.SetText(fmt.Sprintf("%.2f", d.Balance))
			}
		}
	}

	form := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("Customer", customerEntry),
			widget.NewFormItem("Amount", amountEntry),
			widget.NewFormItem("Balance", balanceEntry),
		},
		OnSubmit: func() {
			amt, err := strconv.ParseFloat(amountEntry.Text, 64)
			if err != nil {
				dialog.ShowError(fmt.Errorf("invalid amount"), w)
				return
			}
			bal, _ := strconv.ParseFloat(balanceEntry.Text, 64)
			newBalance := bal - amt

			if p == nil {
				// Insert new payment
				now := time.Now().Format("02-01-2006 15:04:05")
				res, err := db.Exec("INSERT INTO payments(customer,amount,balance,created_at) VALUES(?,?,?,?)",
					customerEntry.Selected, amt, newBalance, now)
				if err != nil {
					dialog.ShowError(fmt.Errorf("failed to add payment: %v", err), w)
					return
				}

				id, _ := res.LastInsertId()
				payments = append(payments, models.Payment{ID: int(id), Customer: customerEntry.Selected, Amount: amt, Balance: newBalance})
			} else {
				// Update existing payment
				_, err := db.Exec("UPDATE payments SET customer=?,amount=?,balance=? WHERE id=?",
					customerEntry.Selected, amt, newBalance, p.ID)
				if err != nil {
					dialog.ShowError(fmt.Errorf("failed to update payment: %v", err), w)
					return
				}
				loadPayments()
				loadDebts()
			}

			// ✅ Only update debt after successful payment insert/update
			_, err = db.Exec("UPDATE debts SET balance=? WHERE customer=?", newBalance, customerEntry.Selected)
			if err != nil {
				dialog.ShowError(fmt.Errorf("failed to update debt balance: %v", err), w)
				return
			}

			// Refresh UI
			loadDebts()
			debtTable.Refresh()
			paymentTable.Refresh()

		},
	}

	dlg := dialog.NewCustom("Payment Entry", "Cancel", form, w)
	dlg.Resize(fyne.NewSize(400, 250))
	dlg.Show()
}

// ---------------------- Contact Tab ----------------------
func buildContactTab() fyne.CanvasObject {
	return container.NewVBox(
		widget.NewLabelWithStyle("Contact & Enquiries", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("For enquiries regarding this application or other custom applications, please reach out using the details below:"),
		widget.NewLabel(""),
		widget.NewLabel("> I design and develop all kinds of Applications."),
		widget.NewLabelWithStyle("Name:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Clinton Moshe"),
		widget.NewLabelWithStyle("Phone:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("+254-746-646-331"),
		widget.NewLabelWithStyle("Email:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("clintonmwachia9@gmail.com"),
	)
}
