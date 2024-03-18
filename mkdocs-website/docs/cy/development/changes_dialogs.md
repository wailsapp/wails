## Sgyrsiau

Mae sgyrsiau bellach ar gael yn JavaScript!

### Ffenestri

Nid yw botymau sgyrsiau yn Windows yn ffurfweddadwy ac yn gyson yn dibynnu ar y
math o sgwrs. I drîgeru galwad enw pan fo botwm yn cael ei wasgu, crëwch botwm
â'r un enw â'r botwm yr ydych am gael y galwad enw i'w gysylltu ag ef.
Enghraifft: Crëwch fotwm â'r label `Ok` a defnyddio `OnClick()` i osod y
dull galwad enw:

```go
        dialog := app.QuestionDialog().
			SetTitle("Diweddaru").
			SetMessage("Mae'r botwm canslo yn cael ei ddewis pan fo'r bys dianc yn cael ei wasgu")
		ok := dialog.AddButton("Ok")
		ok.OnClick(func() {
			// Gwneud rhywbeth
		})
		no := dialog.AddButton("Canslo")
		dialog.SetDefaultButton(ok)
		dialog.SetCancelButton(no)
		dialog.Show()
```