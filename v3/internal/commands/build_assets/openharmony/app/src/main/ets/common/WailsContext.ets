/*
 * WailsContext.ets — access to the running UIAbility context.
 *
 * The ability registers its context during onCreate, so the bridge can reach
 * the system APIs that need an ability context (clipboard, preferences, file
 * picker, window).
 */

import { common } from '@kit.AbilityKit';
import { window } from '@kit.ArkUI';

let abilityContext: common.UIAbilityContext | undefined = undefined;

export function setContext(context: common.UIAbilityContext): void {
  abilityContext = context;
}

/**
 * Return the ability context.
 *
 * Throws when the shell has not registered a context yet: the ability
 * registers it during onCreate, so reaching this point means the shell is
 * misconfigured.
 */
export function getContext(): common.UIAbilityContext {
  if (abilityContext === undefined) {
    throw new Error('Wails ability context is not available');
  }
  return abilityContext;
}

/** Process-wide main window, set from EntryAbility.onWindowStageCreate. */
let mainWindow: window.Window | null = null;

export function setMainWindow(win: window.Window | null): void {
  mainWindow = win;
}

/** The main window, or null before the window stage is created. */
export function getMainWindow(): window.Window | null {
  return mainWindow;
}
