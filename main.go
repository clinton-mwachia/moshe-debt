package main

import (
	"database/sql"
	"fmt"
	"moshe-debt/models"
	"moshe-debt/utils"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
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
	loadDebts()

	addDebtBtn := widget.NewButton("Add Debt", func() { showDebtDialog(nil) })

	debtTable = buildDebtTable()
	debtContainer := container.NewVScroll(debtTable)
	debtContainer.SetMinSize(fyne.NewSize(600, 300))

	return container.NewVBox(
		container.NewHBox(addDebtBtn),
		debtContainer,
	)
}

func loadDebts() {
	debts = nil
	rows, _ := db.Query("SELECT id,customer,phone,amount,balance FROM debts")
	defer rows.Close()
	for rows.Next() {
		var d models.Debt
		rows.Scan(&d.ID, &d.Customer, &d.Phone, &d.Amount, &d.Balance)
		debts = append(debts, d)
	}
}

func buildDebtTable() *widget.Table {
	table := widget.NewTable(
		func() (int, int) { return len(debts) + 1, 5 },
		func() fyne.CanvasObject {
			label := widget.NewLabel("template")
			delBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)
			actions := container.NewHBox(delBtn)
			return container.NewStack(label, actions)
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			cell := obj.(*fyne.Container)
			label := cell.Objects[0].(*widget.Label)
			if id.Row == 0 {
				headers := []string{"Customer", "Phone", "Amount", "Balance", "Actions"}
				label.SetText(headers[id.Col])
				label.Show()
				cell.Objects[1].Hide()
			} else {
				debt := debts[id.Row-1]
				switch id.Col {
				case 0:
					label.SetText(debt.Customer)
					label.Show()
					cell.Objects[1].Hide()
				case 1:
					label.SetText(debt.Phone)
					label.Show()
					cell.Objects[1].Hide()
				case 2:
					label.SetText(fmt.Sprintf("%.2f", debt.Amount))
					label.Show()
					cell.Objects[1].Hide()
				case 3:
					label.SetText(fmt.Sprintf("%.2f", debt.Balance))
					label.Show()
					cell.Objects[1].Hide()
				case 4:
					actions := cell.Objects[1].(*fyne.Container)
					delBtn := actions.Objects[0].(*widget.Button)
					delBtn.OnTapped = func() {
						_, _ = db.Exec("DELETE FROM debts WHERE id=?", debt.ID)
						loadDebts()
						debtTable.Refresh()
					}
					actions.Show()
					label.Hide()
				}
			}
		},
	)
	table.SetColumnWidth(0, 200)
	table.SetColumnWidth(1, 150)
	table.SetColumnWidth(2, 100)
	table.SetColumnWidth(3, 100)
	table.SetColumnWidth(4, 150)
	return table
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
	loadPayments()

	addPaymentBtn := widget.NewButton("Add Payment", func() { showPaymentDialog(nil) })

	paymentTable = buildPaymentTable()
	paymentContainer := container.NewVScroll(paymentTable)
	paymentContainer.SetMinSize(fyne.NewSize(600, 300))

	return container.NewVBox(
		container.NewHBox(addPaymentBtn),
		paymentContainer,
	)
}

func loadPayments() {
	payments = nil
	rows, _ := db.Query("SELECT id,customer,amount,balance FROM payments")
	defer rows.Close()
	for rows.Next() {
		var p models.Payment
		rows.Scan(&p.ID, &p.Customer, &p.Amount, &p.Balance)
		payments = append(payments, p)
	}
}

func buildPaymentTable() *widget.Table {
	table := widget.NewTable(
		func() (int, int) { return len(payments) + 1, 4 },
		func() fyne.CanvasObject {
			label := widget.NewLabel("template")
			delBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)
			actions := container.NewHBox(delBtn)
			return container.NewStack(label, actions)
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			cell := obj.(*fyne.Container)
			label := cell.Objects[0].(*widget.Label)
			if id.Row == 0 {
				headers := []string{"Customer", "Amount", "Balance", "Actions"}
				label.SetText(headers[id.Col])
				label.Show()
				cell.Objects[1].Hide()
			} else {
				p := payments[id.Row-1]
				switch id.Col {
				case 0:
					label.SetText(p.Customer)
					label.Show()
					cell.Objects[1].Hide()
				case 1:
					label.SetText(fmt.Sprintf("%.2f", p.Amount))
					label.Show()
					cell.Objects[1].Hide()
				case 2:
					label.SetText(fmt.Sprintf("%.2f", p.Balance))
					label.Show()
					cell.Objects[1].Hide()
				case 3:
					actions := cell.Objects[1].(*fyne.Container)
					delBtn := actions.Objects[0].(*widget.Button)
					delBtn.OnTapped = func() {
						// Delete the payment
						res, err := db.Exec("DELETE FROM payments WHERE id=?", p.ID)
						if err != nil {
							dialog.ShowError(fmt.Errorf("failed to delete payment: %v", err), fyne.CurrentApp().Driver().AllWindows()[0])
							return
						}

						rowsAffected, _ := res.RowsAffected()
						if rowsAffected > 0 {
							// Add back the payment balance to the debt
							_, err := db.Exec("UPDATE debts SET balance = balance + ? WHERE customer=?", p.Amount, p.Customer)
							if err != nil {
								dialog.ShowError(fmt.Errorf("failed to update debt balance: %v", err), fyne.CurrentApp().Driver().AllWindows()[0])
								return
							}
						}

						// Reload data and refresh tables
						loadPayments()
						loadDebts()
						paymentTable.Refresh()
						debtTable.Refresh()
					}

					actions.Show()
					label.Hide()
				}
			}
		},
	)
	table.SetColumnWidth(0, 200)
	table.SetColumnWidth(1, 100)
	table.SetColumnWidth(2, 100)
	table.SetColumnWidth(3, 150)
	return table
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
