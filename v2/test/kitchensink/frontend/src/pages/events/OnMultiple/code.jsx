import { Events } from '@wails/runtime';

let notes = [];

// Do some things
Events.On("notes loaded", (newNotes) => {
  notes = newNotes;
});
