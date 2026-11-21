import { open, request, close } from 'libentry.so';
import { hilog } from '@kit.PerformanceAnalysisKit';

@Entry
@Component
struct Index {
  @State private status: string = 'Running';
  private replies: Map<number, number> = new Map<number, number>();
  private tick: number = 0;
  private timer: number = -1;
  private deadline: boolean = false;

  aboutToAppear(): void {
    const start = Date.now();
    this.timer = setInterval((): void => { this.tick++; }, 25);
    const opened: boolean = open((id: number, value: number): void => { this.replies.set(id, value); });
    const sent: boolean = opened && request(3, 30) && request(4, 80) && request(5, 40) && request(6, 1000);
    setTimeout((): void => { this.deadline = !this.replies.has(6); }, 180);
    setTimeout((): void => {
      const complete = this.replies.get(3) === 9 && this.replies.get(4) === 16 && this.replies.get(5) === 25;
      const pending: boolean = request(7, 600);
      const closed: boolean = close();
      const ticks = this.tick;
      const elapsed = Date.now() - start;
      setTimeout((): void => {
        const cancelled = !this.replies.has(6) && !this.replies.has(7);
        const pass: boolean = sent && complete && pending && closed && cancelled && this.deadline && ticks >= 3;
        this.status = pass ? 'M2 PASS' : 'M2 FAIL';
        hilog.info(0x0, 'WailsM2', '%{public}s sent=%{public}s complete=%{public}s pending=%{public}s closed=%{public}s cancelled=%{public}s deadline=%{public}s ticks=%{public}d elapsed=%{public}d', this.status, String(sent), String(complete), String(pending), String(closed), String(cancelled), String(this.deadline), ticks, elapsed);
        clearInterval(this.timer);
      }, 1250);
    }, 375);
  }

  aboutToDisappear(): void {
    close();
    clearInterval(this.timer);
  }

  build() {
    Column() { Text(this.status).fontSize(24) }
  }
}
