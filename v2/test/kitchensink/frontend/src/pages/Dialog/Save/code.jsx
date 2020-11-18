import { Dialog } from '@wails/runtime';

let notes = "";

function saveNotes() {
  // Prompt the user to select a single file
  let filename = Dialog.Save({
    "DefaultFilename": "notes.md",
    "Filters": "*.md",
  });

  // Do something with the file
  backend.main.SaveNotes(filename, notes).then( (result) => {
    if ( !result ) {
      // Cancelled
      return
    }
    showMessage('Notes saved!');
  }).catch( (err) => {
    // Show an alert
    showAlert(err);
  })
}
