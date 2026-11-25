import { UIAbility, Want } from '@kit.AbilityKit';
import { window } from '@kit.ArkUI';
import { hilog } from '@kit.PerformanceAnalysisKit';
import * as native from 'libentry.so';
import { setContext } from '../common/WailsContext';

// EntryAbility hosts the Wails page. KeeperAbility renders nothing; it exists
// so lifecycle tests can destroy a stage without taking the process with it.
export default class EntryAbility extends UIAbility {
  private name: string = '';
  private closeStage: boolean = false;

  onCreate(want: Want): void {
    setContext(this.context);
    this.name = want.abilityName ?? '';
    this.closeStage = want.parameters?.['closeStage'] === true;
    const mode: string | undefined = want.parameters?.['wailsLifecycle'] as string | undefined;
    const phase: string | undefined = want.parameters?.['wailsPhase'] as string | undefined;
    AppStorage.setOrCreate('wailsLifecycle',
      mode === 'refresh' || mode === 'rebuild' || mode === 'stage' ? mode : '');
    AppStorage.setOrCreate('wailsPhase', phase === 'new' ? 'new' : 'old');
  }

  onWindowStageCreate(stage: window.WindowStage): void {
    hilog.info(0x0, 'WailsM3', 'Stage create %{public}s', this.name);
    if (this.name === 'KeeperAbility') {
      return;
    }
    if (this.closeStage) {
      setTimeout((): void => {
        hilog.info(0x0, 'WailsM3', 'Terminating %{public}s', this.name);
        this.context.terminateSelf().catch((error: Error): void => {
          hilog.error(0x0, 'WailsM3', 'Terminate failed %{public}s', error.message);
        });
      }, 4000);
    }
    stage.loadContent('pages/Index');
    native.abilityStarted();
  }

  onForeground(): void {
    native.abilityResumed();
  }

  onBackground(): void {
    native.abilityPaused();
  }

  onWindowStageDestroy(): void {
    native.abilityStopped();
    hilog.info(0x0, 'WailsM3', 'Stage destroy %{public}s', this.name);
  }

  onDestroy(): void {
    hilog.info(0x0, 'WailsM3', 'Ability destroy %{public}s', this.name);
    native.shutdown();
  }
}
