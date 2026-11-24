/*
 * EntryAbility.ets — hosts the Wails page.
 *
 * The ability owns the ArkWeb view; the Go backend runs in the same process
 * and is reached through the native bridge.
 */

import { UIAbility, Want } from '@kit.AbilityKit';
import { window } from '@kit.ArkUI';
import { hilog } from '@kit.PerformanceAnalysisKit';

import { setContext } from '../common/WailsContext';

const DOMAIN: number = 0x0000;
const TAG: string = 'Wails';

export default class EntryAbility extends UIAbility {
  onCreate(want: Want): void {
    setContext(this.context);
  }

  onWindowStageCreate(stage: window.WindowStage): void {
    hilog.info(DOMAIN, TAG, 'Stage create');
    stage.loadContent('pages/Index');
  }

  onWindowStageDestroy(): void {
    hilog.info(DOMAIN, TAG, 'Stage destroy');
  }

  onDestroy(): void {
    hilog.info(DOMAIN, TAG, 'Ability destroy');
  }
}
