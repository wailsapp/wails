import { Component } from '@angular/core';

@Component({
  selector: '[id="app"]',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'my-app';

  clickMessage = '';

  onClickMe() {
    // @ts-ignore
    window.backend.basic().then(result =>
      this.clickMessage = result
    );
  }
}
