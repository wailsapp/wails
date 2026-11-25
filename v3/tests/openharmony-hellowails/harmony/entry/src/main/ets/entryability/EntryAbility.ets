import { UIAbility, Want, EnvironmentCallback, Configuration, AbilityConstant } from '@kit.AbilityKit';
import { window } from '@kit.ArkUI';
import { hilog } from '@kit.PerformanceAnalysisKit';
import { BusinessError, batteryInfo, commonEventManager } from '@kit.BasicServicesKit';
import { connection } from '@kit.NetworkKit';
import * as native from 'libentry.so';
import { setContext, setMainWindow } from '../common/WailsContext';

// EntryAbility hosts the Wails page. KeeperAbility renders nothing; it exists
// so lifecycle tests can destroy a stage without taking the process with it.
interface ThemePayload {
  isDarkMode: boolean;
}

interface BatteryPayload {
  level: number;
  state: string;
  lowPowerMode: boolean;
  charging: boolean;
}

interface NetworkPayload {
  connected: boolean;
  type: string;
  metered: boolean;
}

export default class EntryAbility extends UIAbility {
  private name: string = '';
  private closeStage: boolean = false;
  private environmentId: number = -1;
  private systemSubscriber: commonEventManager.CommonEventSubscriber | null = null;
  private netConnection: connection.NetConnection | null = null;

  onCreate(want: Want): void {
    setContext(this.context);
    this.name = want.abilityName ?? '';
    this.closeStage = want.parameters?.['closeStage'] === true;
    const mode: string | undefined = want.parameters?.['wailsLifecycle'] as string | undefined;
    const phase: string | undefined = want.parameters?.['wailsPhase'] as string | undefined;
    AppStorage.setOrCreate('wailsLifecycle',
      mode === 'all' || mode === 'api' || mode === 'refresh' || mode === 'rebuild' || mode === 'stage' ? mode : '');
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
    setMainWindow(stage.getMainWindowSync());
    stage.loadContent('pages/Index');
    native.abilityStarted();
    this.registerSystemObservers();
  }

  onForeground(): void {
    native.abilityResumed();
  }

  onBackground(): void {
    native.abilityPaused();
  }

  onWindowStageDestroy(): void {
    this.unregisterSystemObservers();
    setMainWindow(null);
    native.abilityStopped();
    hilog.info(0x0, 'WailsM3', 'Stage destroy %{public}s', this.name);
  }

  onDestroy(): void {
    hilog.info(0x0, 'WailsM3', 'Ability destroy %{public}s', this.name);
    native.shutdown();
  }

  /** Forward OS environment changes to the Go event bridge. */
  private registerSystemObservers(): void {
    const environmentCallback: EnvironmentCallback = {
      onConfigurationUpdated: (config: Configuration): void => {
        const payload: ThemePayload = { isDarkMode: config.colorMode === 0 };
        native.emitSystemEvent('themeChanged', JSON.stringify(payload));
      },
      onMemoryLevel: (_level: AbilityConstant.MemoryLevel): void => {
        // The Go side handles low-memory lifecycle events separately.
      }
    };
    try {
      this.environmentId = this.context.getApplicationContext().on('environment', environmentCallback);
    } catch (error) {
      hilog.error(0x0, 'WailsM3', 'theme observer failed: %{public}s', JSON.stringify(error));
    }

    const subscribeInfo: commonEventManager.CommonEventSubscribeInfo = {
      events: ['usual.event.BATTERY_CHANGED', 'usual.event.SCREEN_OFF', 'usual.event.SCREEN_ON']
    };
    commonEventManager.createSubscriber(subscribeInfo).then((subscriber: commonEventManager.CommonEventSubscriber): void => {
      this.systemSubscriber = subscriber;
      commonEventManager.subscribe(subscriber,
        (error: BusinessError, data: commonEventManager.CommonEventData): void => {
          if (error !== undefined) {
            hilog.error(0x0, 'WailsM3', 'system event subscribe failed: %{public}s', error.message);
            return;
          }
          if (data.event === 'usual.event.BATTERY_CHANGED') {
            const charging: boolean = batteryInfo.chargingStatus === batteryInfo.BatteryChargeState.ENABLE;
            const payload: BatteryPayload = {
              level: batteryInfo.batterySOC / 100,
              state: charging ? 'charging' : 'discharging',
              lowPowerMode: false,
              charging: charging
            };
            native.emitSystemEvent('batteryChanged', JSON.stringify(payload));
          } else if (data.event === 'usual.event.SCREEN_OFF') {
            native.emitSystemEvent('screenLocked', '{}');
          } else if (data.event === 'usual.event.SCREEN_ON') {
            native.emitSystemEvent('screenUnlocked', '{}');
          }
        });
    }).catch((error: BusinessError): void => {
      hilog.error(0x0, 'WailsM3', 'system subscriber failed: %{public}s', error.message);
    });

    const net: connection.NetConnection = connection.createNetConnection();
    net.on('netAvailable', (_handle: connection.NetHandle): void => this.emitNetwork(true));
    net.on('netCapabilitiesChange', (_info: connection.NetCapabilityInfo): void => this.emitNetwork(true));
    net.on('netLost', (_handle: connection.NetHandle): void => this.emitNetwork(false));
    net.on('netUnavailable', (): void => this.emitNetwork(false));
    net.register((error: BusinessError): void => {
      if (error !== undefined) {
        hilog.error(0x0, 'WailsM3', 'network observer failed: %{public}s', error.message);
      }
    });
    this.netConnection = net;
  }

  private emitNetwork(connected: boolean): void {
    let type: string = 'none';
    if (connected) {
      const handle: connection.NetHandle = connection.getDefaultNetSync();
      if (handle.netId !== 0) {
        const capabilities: connection.NetCapabilities = connection.getNetCapabilitiesSync(handle);
        const bearers: connection.NetBearType[] = capabilities.bearerTypes;
        for (const bearer of bearers) {
          if (bearer === connection.NetBearType.BEARER_WIFI) {
            type = 'wifi';
            break;
          }
          if (bearer === connection.NetBearType.BEARER_CELLULAR) {
            type = 'cellular';
            break;
          }
          if (bearer === connection.NetBearType.BEARER_ETHERNET) {
            type = 'ethernet';
            break;
          }
        }
      }
    }
    const payload: NetworkPayload = {
      connected: connected,
      type: type,
      metered: type === 'cellular'
    };
    native.emitSystemEvent('networkChanged', JSON.stringify(payload));
  }

  private unregisterSystemObservers(): void {
    if (this.environmentId >= 0) {
      try {
        this.context.getApplicationContext().off('environment', this.environmentId);
      } catch (error) {
        hilog.error(0x0, 'WailsM3', 'theme observer cleanup failed: %{public}s', JSON.stringify(error));
      }
      this.environmentId = -1;
    }
    if (this.systemSubscriber !== null) {
      commonEventManager.unsubscribe(this.systemSubscriber, (error: BusinessError): void => {
        if (error !== undefined) {
          hilog.error(0x0, 'WailsM3', 'system unsubscribe failed: %{public}s', error.message);
        }
      });
      this.systemSubscriber = null;
    }
    if (this.netConnection !== null) {
      this.netConnection.unregister((error: BusinessError): void => {
        if (error !== undefined) {
          hilog.error(0x0, 'WailsM3', 'network unregister failed: %{public}s', error.message);
        }
      });
      this.netConnection = null;
    }
  }
}
