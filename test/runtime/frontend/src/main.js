
require('file-loader?name=[name].[ext]!./index.html');
import 'bootstrap';
import 'bootstrap/dist/css/bootstrap.min.css';

// Import main css
import './assets/css/main.css';

// Import logo
const logo = require('./assets/images/wails-logo.png');
document.getElementById("logo").src = logo;

window.scripts = {
    injectCalls: () => {
        var calls = document.getElementById("calls-info");
        if (window.wailsbindings != "{}") {
            var bindings = JSON.parse(window.wailsbindings);

            for (const pkg in bindings) {
                for (const struct in bindings[pkg]) {
                    for (const func in bindings[pkg][struct]) {
                        var form = document.createElement("form");
                        calls.appendChild(form);

                        var details = bindings[pkg][struct][func];
                        if (!details.inputs) {
                            details.inputs = [];
                        }
                        if (!details.outputs) {
                            details.outputs = [];
                        }
                        var funcName = details.name;
                        var div = document.createElement("div");
                        form.appendChild(div);
                        div.classList.add("form-group");
                        div.classList.add("row");
                        var label = document.createElement("label");
                        div.appendChild(label);
                        label.classList.add("col-sm-2");
                        label.classList.add("col-form-label");
                        label.innerText = details.name;

                        var inputs = document.createElement("div");
                        inputs.classList.add("form-group");
                        inputs.classList.add("row");
                        form.appendChild(inputs);

                        label = document.createElement("label");
                        label.classList.add("col-sm-2");
                        label.classList.add("col-form-label");
                        label.innerText = "Inputs";
                        inputs.appendChild(label);

                        for (const input in details["inputs"]) {
                            var idiv = document.createElement("div");
                            idiv.classList.add("col");
                            var i = document.createElement("input");
                            idiv.appendChild(i);
                            i.classList.add("form-control");
                            i.setAttribute("type", "text");
                            i.setAttribute("placeholder", details["inputs"][input]["type"]);
                            i.setAttribute("id", details.name + ".input" + input);
                            inputs.appendChild(idiv);
                        }

                        var outputs = document.createElement("div");
                        outputs.classList.add("form-group");
                        outputs.classList.add("row");
                        form.appendChild(outputs);

                        label = document.createElement("label");
                        label.classList.add("col-sm-2");
                        label.classList.add("col-form-label");
                        label.innerText = "Outputs";
                        outputs.appendChild(label);

                        for (const output in details["outputs"]) {
                            var odiv = document.createElement("div");
                            outputs.appendChild(odiv);
                            odiv.classList.add("col");
                            var o = document.createElement("input");
                            o.disabled = true;
                            odiv.appendChild(o);
                            o.classList.add("form-control");
                            o.setAttribute("type", "text");
                            o.setAttribute("placeholder", details["outputs"][output]["type"]);
                            o.setAttribute("id", details.name + ".output" + output);
                        }

                        var button = document.createElement("button");
                        button.classList.add("btn");
                        button.classList.add("btn-primary");
                        button.classList.add("float-right");
                        button.innerText = "Call";
                        button.onclick = function (details) {

                            return function () {

                                // Clear result panels
                                var output = document.getElementById(details.name + ".output0");
                                if (output) {
                                    output.value = "";
                                }
                                output = document.getElementById(details.name + ".output1");
                                if (output) {
                                    output.value = "";
                                }

                                var f = eval("window.backend." + details.name);
                                var inputs = [];
                                for (var i = 0; i < details.inputs.length; i++) {
                                    var input = document.getElementById(details.name + ".input" + i);
                                    switch (details.inputs[i].type) {
                                        case "float":
                                            inputs.push(parseFloat(input.value));
                                            break;
                                        case "int":
                                            inputs.push(parseInt(input.value));
                                            break;
                                        case "string":
                                            inputs.push(input.value);
                                    }
                                }
                                console.log("Outputs length = " + details.outputs.length);
                                switch (details.outputs.length) {
                                    case 0:
                                        f(...inputs);
                                        break;
                                    case 1:
                                        f(...inputs).then(result => {
                                            document.getElementById(details.name + ".output0").value = result;
                                        });
                                        break;
                                    case 2:

                                        f(...inputs).then((result) => {
                                            document.getElementById(details.name + ".output0").value = result;
                                        }).catch(error => {
                                            document.getElementById(details.name + ".output0").value = "";
                                            document.getElementById(details.name + ".output1").value = error;
                                        });
                                        break;
                                };

                                return false;
                            };
                        }(details);
                        form.appendChild(button);
                    }
                }
            }
        } else {
            calls.appendChild(document.createTextNode("None defined"));
        }
    },

    setActive: (el) => {
        var activeElement = document.getElementsByClassName("active")[0];
        if (activeElement) {
            activeElement.classList.remove("active");
        }

        el.classList.add("active");
    },
    openURL: () => {
        var url = document.getElementById("url").value;
        window.wails.Browser.OpenURL(url);
    },
    sendEvent: () => {
        try {
            var eventName = document.getElementById("eventname").value;
            var eventParameters = document.getElementById("eventparameters").value;
            eventParameters = "[" + eventParameters + "]";
            var parsedParams = JSON.parse(eventParameters);
            console.log(eventName)
            console.log(parsedParams)
            window.wails.Events.Emit(eventName, ...parsedParams)
        } catch (e) {
            console.log(e)
        }
    },
    doLogging: () => {
        var radios = document.getElementsByName('logRadio');
        var selected = "";

        for (var i = 0, length = radios.length; i < length; i++) {
            if (radios[i].checked) {
                // do whatever you want with the checked radio
                selected = radios[i].value;

                // only one radio can be logically checked, don't check the rest
                break;
            }
        }

        var message = document.getElementById("message").value;

        switch (selected) {
            case "Debug":
                window.wails.Log.Debug(message);
                break;
            case "Info":
                window.wails.Log.Info(message);
                break;
            case "Warning":
                window.wails.Log.Warning(message);
                break;
            case "Error":
                window.wails.Log.Error(message);
                break;
            case "Fatal":
                window.wails.Log.Fatal(message);
                break;
            default:
                alert("Unknwon log type: " + selected);
        }
    },

    // Windowing
    hideWindow: () => {
        window.wails.Window.Hide();
        setTimeout( ()=> {
            window.wails.Window.Show();
        }, 3000);
    },
    center: () => {
        window.wails.Window.Center();
    },
    maximiseWindow: () => {
        window.wails.Window.Maximise();
    },
    unmaximiseWindow: () => {
        window.wails.Window.Unmaximise();
    },
    minimiseWindow: () => {
        window.wails.Window.Minimise();
        setTimeout( ()=> {
            window.wails.Window.Unminimise();
        }, 3000);
    },
    unminimiseWindow: () => {
        window.wails.Window.Unminimise();
    },
    setposition: () => {
        var windowx = parseInt(document.getElementById("windowx").value);
        var windowy = parseInt(document.getElementById("windowy").value);

        if( windowx == NaN || windowy == NaN) {
            return;
        }

        window.wails.Window.SetPosition(windowx, windowy);
    },
    setsize: () => {
        var width = parseInt(document.getElementById("windowwidth").value);
        var height = parseInt(document.getElementById("windowheight").value);

        if( width == NaN || height == NaN) {
            return;
        }

        window.wails.Window.SetSize(width, height);   
    }
};


window.scripts.injectCalls();
