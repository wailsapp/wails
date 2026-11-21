import { UIAbility, Want } from '@kit.AbilityKit';
import { window } from '@kit.ArkUI';

export default class EntryAbility extends UIAbility {
  onCreate(want: Want): void {
    const mode = want.parameters?.['wailsLifecycle'];
    AppStorage.setOrCreate('wailsLifecycle', mode === 'refresh' || mode === 'rebuild' ? mode : '');
  }

  onWindowStageCreate(stage: window.WindowStage): void {
    stage.loadContent('pages/Index');
  }
}
