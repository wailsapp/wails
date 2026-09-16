import { Call } from '@wailsio/runtime';
Call.ByName('main.AcceptanceService.Ping', 'native-dev').then(async (value: string) => {
 document.getElementById('result')!.textContent = value;
 await Call.ByName('main.AcceptanceService.Confirm', value);
}).catch(error => {document.getElementById('result')!.textContent = String(error)});
