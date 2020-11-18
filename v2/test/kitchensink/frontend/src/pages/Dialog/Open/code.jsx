import { Dialog } from '@wails/runtime';

let notes = "";

function loadNotes() {
  // Prompt the user to select a single file
  let filename = Dialog.Open({
    "DefaultFilename": "notes.md",
    "Filters": "*.md",
    "AllowFiles": true,
  });

  // Do something with the file
  backend.main.LoadNotes(filename).then( (result) => {
    if (result.length == 0) {
      // Cancelled
      return
    }
    // We only prompted for a single file
    notes = result[0];
  }).catch( (err) => {
    // Show an alert
    showAlert(err);
  })
}
