import { Events } from '@wails/runtime';

// Respond to the unlock event 3 times
Events.OnMultiple("unlock", (password) => {
  // Check password
}, 3);
