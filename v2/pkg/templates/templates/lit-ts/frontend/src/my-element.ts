import {css, html, LitElement} from 'lit'
import logo from './assets/images/logo-universal.png'
import {Greet} from "../wailsjs/go/main/App";
import {customElement, property} from 'lit/decorators.js'
import './style.css';

/**
 * An example element.
 *
 * @slot - This element has a slot
 * @csspart button - The button
 */
@customElement('my-element')
export class MyElement extends LitElement {
    static styles = css`
  #logo {
    display: block;
    width: 50%;
    height: 50%;
    margin: auto;
    padding: 10% 0 0;
    background-position: center;
    background-repeat: no-repeat;
    background-size: 100% 100%;
    background-origin: content-box;
  }

  .result {
    height: 20px;
    line-height: 20px;
    margin: 1.5rem auto;
  }

  .input-box .btn {
    width: 60px;
    height: 30px;
    line-height: 30px;
    border-radius: 3px;
    border: none;
    margin: 0 0 0 20px;
    padding: 0 8px;
    cursor: pointer;
  }

  .input-box .btn:hover {
    background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
    color: #333333;
  }

  .input-box .input {
    border: none;
    border-radius: 3px;
    outline: none;
    height: 30px;
    line-height: 30px;
    padding: 0 10px;
    background-color: rgba(240, 240, 240, 1);
    -webkit-font-smoothing: antialiased;
  }

  .input-box .input:hover {
    border: none;
    background-color: rgba(255, 255, 255, 1);
  }

  .input-box .input:focus {
    border: none;
    background-color: rgba(255, 255, 255, 1);
  }

    `

    @property()
    resultText = "Please enter your name below 👇"

    greet() {
        let thisName = (this.shadowRoot?.getElementById('name') as HTMLInputElement)?.value;
        if (thisName) {
            Greet(thisName).then(result => {
                this.resultText = result
            });
        }
    }

    render() {
        return html`
            <main>
                <img id="logo" src=${logo} alt="Wails logo">
                <div class="result" id="result">${this.resultText}</div>
                <div class="input-box" id="input">
                    <input class="input" id="name" type="text" autocomplete="off"/>
                    <button @click=${this.greet} class="btn">Greet</button>
                </div>
            </main>
        `
    }
}

declare global {
    interface HTMLElementTagNameMap {
        'my-element': MyElement
    }
}