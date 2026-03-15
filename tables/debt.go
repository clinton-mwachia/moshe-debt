package tables

import (
	"database/sql"
	"fmt"
	"moshe-debt/models"
	"moshe-debt/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func BuildDebtTable(db *sql.DB, debts []models.Debt, debtTable *widget.Table) *widget.Table {
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
						utils.LoadDebts(db, debts)
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
