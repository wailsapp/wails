/*
 * WailsContext.ets — access to the running UIAbility context.
 *
 * The ability registers its context during onCreate, so the bridge can reach
 * the system APIs that need an ability context (clipboard, theme, window).
 */

import { common } from '@kit.AbilityKit';

let abilityContext: common.UIAbilityContext | undefined = undefined;

export function setContext(context: common.UIAbilityContext): void {
  abilityContext = context;
}

export function getContext(): common.UIAbilityContext {
  if (abilityContext === undefined) {
    throw new Error('Wails ability context is not available');
  }
  return abilityContext;
}
