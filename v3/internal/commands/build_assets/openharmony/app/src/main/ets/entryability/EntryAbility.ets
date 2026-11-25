/*
 * EntryAbility.ets — hosts the Wails page.
 *
 * The ability owns the ArkWeb view and the process-wide main window. The Go
 * backend runs in the same process and is reached through the native bridge.
 */

import { UIAbility, Want, EnvironmentCallback, Configuration, AbilityConstant } from '@kit.AbilityKit';
import { window } from '@kit.ArkUI';
import { hilog } from '@kit.PerformanceAnalysisKit';
import { BusinessError, batteryInfo, commonEventManager } from '@kit.BasicServicesKit';
import { connection } from '@kit.NetworkKit';
import * as native from 'libentry.so';

import { setContext, setMainWindow } from '../common/WailsContext';

const DOMAIN: number = 0x0000;
const TAG: string = 'Wails';

/** Common-event names used for battery and screen-lock notifications. */
const BATTERY_CHANGED: string = 'usual.event.BATTERY_CHANGED';
const SCREEN_OFF: string = 'usual.event.SCREEN_OFF';
const SCREEN_ON: string = 'usual.event.SCREEN_ON';

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
  private environmentId: number = -1;
  private systemSubscriber: commonEventManager.CommonEventSubscriber | null = null;
  private netConnection: connection.NetConnection | null = null;

  onCreate(want: Want): void {
    setContext(this.context);
  }

  onWindowStageCreate(stage: window.WindowStage): void {
    hilog.info(DOMAIN, TAG, 'Stage create');
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
    hilog.info(DOMAIN, TAG, 'Stage destroy');
  }

  onDestroy(): void {
    hilog.info(DOMAIN, TAG, 'Ability destroy');
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
        // Low memory is forwarded from Ability.onMemoryLevel by the Go side.
      }
    };
    try {
      this.environmentId = this.context.getApplicationContext().on('environment', environmentCallback);
    } catch (error) {
      hilog.error(DOMAIN, TAG, 'theme observer failed: %{public}s', JSON.stringify(error));
    }

    const subscribeInfo: commonEventManager.CommonEventSubscribeInfo = {
      events: [BATTERY_CHANGED, SCREEN_OFF, SCREEN_ON]
    };
    commonEventManager.createSubscriber(subscribeInfo).then((subscriber: commonEventManager.CommonEventSubscriber): void => {
      this.systemSubscriber = subscriber;
      commonEventManager.subscribe(subscriber,
        (error: BusinessError, data: commonEventManager.CommonEventData): void => {
          if (error !== undefined) {
            hilog.error(DOMAIN, TAG, 'system event subscribe failed: %{public}s', error.message);
            return;
          }
          if (data.event === BATTERY_CHANGED) {
            const charging: boolean = batteryInfo.chargingStatus === batteryInfo.BatteryChargeState.ENABLE;
            const payload: BatteryPayload = {
              level: batteryInfo.batterySOC / 100,
              state: charging ? 'charging' : 'discharging',
              lowPowerMode: false,
              charging: charging
            };
            native.emitSystemEvent('batteryChanged', JSON.stringify(payload));
          } else if (data.event === SCREEN_OFF) {
            native.emitSystemEvent('screenLocked', '{}');
          } else if (data.event === SCREEN_ON) {
            native.emitSystemEvent('screenUnlocked', '{}');
          }
        });
    }).catch((error: BusinessError): void => {
      hilog.error(DOMAIN, TAG, 'system subscriber failed: %{public}s', error.message);
    });

    const net: connection.NetConnection = connection.createNetConnection();
    net.on('netAvailable', (_handle: connection.NetHandle): void => this.emitNetwork(true));
    net.on('netCapabilitiesChange', (_info: connection.NetCapabilityInfo): void => this.emitNetwork(true));
    net.on('netLost', (_handle: connection.NetHandle): void => this.emitNetwork(false));
    net.on('netUnavailable', (): void => this.emitNetwork(false));
    net.register((error: BusinessError): void => {
      if (error !== undefined) {
        hilog.error(DOMAIN, TAG, 'network observer failed: %{public}s', error.message);
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
        hilog.error(DOMAIN, TAG, 'theme observer cleanup failed: %{public}s', JSON.stringify(error));
      }
      this.environmentId = -1;
    }
    if (this.systemSubscriber !== null) {
      commonEventManager.unsubscribe(this.systemSubscriber, (error: BusinessError): void => {
        if (error !== undefined) {
          hilog.error(DOMAIN, TAG, 'system unsubscribe failed: %{public}s', error.message);
        }
      });
      this.systemSubscriber = null;
    }
    if (this.netConnection !== null) {
      this.netConnection.unregister((error: BusinessError): void => {
        if (error !== undefined) {
          hilog.error(DOMAIN, TAG, 'network unregister failed: %{public}s', error.message);
        }
      });
      this.netConnection = null;
    }
  }
}
