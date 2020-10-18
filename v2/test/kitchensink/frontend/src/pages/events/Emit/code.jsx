import { Events } from '@wails/runtime';

function processButtonPress(name, address) {
  Events.Emit("new user", name, address);
}

