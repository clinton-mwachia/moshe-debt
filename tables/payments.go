package tables

import (
	"database/sql"
	"fmt"
	"moshe-debt/models"
	"moshe-debt/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func BuildPaymentTable(db *sql.DB, payments []models.Payment, debts []models.Debt, paymentTable *widget.Table) *widget.Table {
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
						utils.LoadPayments(db, payments)
						utils.LoadDebts(db, debts)
						paymentTable.Refresh()
						paymentTable.Refresh()
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
