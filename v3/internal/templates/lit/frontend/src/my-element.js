import {css, html, LitElement} from 'lit'
import {GreetService} from "../bindings/changeme";
import {Events} from "@wailsio/runtime";

export class MyElement extends LitElement {
    static properties = {
        name: {type: String},
        result: {type: String},
        time: {type: String},
    };

    constructor() {
        super();
        this.name = '';
        this.result = 'Please enter your name below 👇';
        this.time = 'Listening for Time event...';
        Events.On('time', (timeValue) => {
            this.time = timeValue.data;
        });
    }

    doGreet() {
        let name = this.name;
        if (!name) {
            name = 'anonymous';
        }
        GreetService.Greet(name).then((resultValue) => {
            this.result = resultValue;
        }).catch((err) => {
            console.log(err);
        });
    }

    render() {
        return html`
            <div class="container">
                <div>
                    <a wml-openURL="https://wails.io">
                        <img src="/wails.png" class="logo" alt="Wails logo"/>
                    </a>
                    <a wml-openURL="https://lit.dev">
                        <img src="/lit.svg" class="logo lit" alt="Lit logo"/>
                    </a>
                </div>
                <slot></slot>
                <div class="result">${this.result}</div>
                <div class="card">
                    <div class="input-box">
                        <input class="input" .value=${this.name} @input=${e => this.name = e.target.value} type="text"
                               autocomplete="off"/>
                        <button class="btn" @click=${this.doGreet}>Greet</button>
                    </div>
                </div>
                <div class="footer">
                    <div><p>Click on the Wails logo to learn more</p></div>
                    <div><p>${this.time}</p></div>
                </div>
            </div>
        `;
    }

    static styles = css`
        :host {
            max-width: 1280px;
            margin: 0 auto;
            padding: 2rem;
            text-align: center;
        }

        h3 {
            font-size: 3em;
            line-height: 1.1;
        }

        a {
            font-weight: 500;
            color: #646cff;
            text-decoration: inherit;
        }

        a:hover {
            color: #535bf2;
        }

        button {
            width: 60px;
            height: 30px;
            line-height: 30px;
            border-radius: 3px;
            border: none;
            margin: 0 0 0 20px;
            padding: 0 8px;
            cursor: pointer;
        }

        .container {
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
        }

        .logo {
            height: 6em;
            padding: 1.5em;
            will-change: filter;
        }

        .logo:hover {
            filter: drop-shadow(0 0 2em #e80000aa);
        }

        .logo.lit:hover {
            filter: drop-shadow(0 0 2em #325cffaa);
        }

        .result {
            height: 20px;
            line-height: 20px;
            margin: 1.5rem auto;
            text-align: center;
        }

        .footer {
            margin-top: 1rem;
            align-content: center;
            text-align: center;
            color: rgba(255, 255, 255, 0.67);
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
            color: black;
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
  `;
}

window.customElements.define('my-element', MyElement);
