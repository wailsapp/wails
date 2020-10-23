import { System, Browser } from '@wails/runtime';

function showPlatformHelp() {
  // Do things
  Browser.Open("https://wails.app/gettingstarted/" + System.Platform());
}