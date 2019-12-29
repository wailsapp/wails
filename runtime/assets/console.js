(function () {

    window.wailsconsole = {};

    var debugconsole = document.createElement("div");
    var header = document.createElement("div");
    var consoleOut = document.createElement("div");


    document.addEventListener('keyup', logKey);

    debugconsole.id = "wailsdebug";
    debugconsole.style.fontSize = "18px";
    debugconsole.style.width = "100%";
    debugconsole.style.height = "35%";
    debugconsole.style.maxHeight = "35%";
    debugconsole.style.position = "fixed";
    debugconsole.style.left = "0px";
    debugconsole.style.backgroundColor = "rgba(255,255,255,0.8)";
    debugconsole.style.borderTop = '1px solid black';
    debugconsole.style.color = "black";
    debugconsole.style.display = "none";

    header.style.width = "100%";
    header.style.height = "30px";
    header.style.display = "block";
    // header.style.paddingTop = "3px";
    header.style.verticalAlign = "middle";
    header.style.paddingLeft = "10px";
    header.style.background = "rgba(255,255,255,0.8)";
    header.innerHTML = "  <span style='vertical-align: middle'> Wails Console > <input id='conin' style='border: solid 1px black; width: 50%'></input><span style='padding-left: 5px; cursor:pointer;' onclick='window.wailsconsole.clearConsole()'>Clear</span></span>";

    consoleOut.style.position = "absolute";
    consoleOut.style.width = "100%";
    consoleOut.style.height = "auto";
    consoleOut.style.top = "30px";
    // consoleOut.style.paddingLeft = "10px";
    consoleOut.style.bottom = "0px";
    consoleOut.style.backgroundColor = "rgba(200,200,200,1)";
    consoleOut.style.overflowY = "scroll";
    consoleOut.style.msOverflowStyle = "-ms-autohiding-scrollbar";

    debugconsole.appendChild(header);
    debugconsole.appendChild(consoleOut);
    document.body.appendChild(debugconsole);
    console.log(debugconsole.style.display)

    function logKey(e) {
        var conin = document.getElementById('conin');
        if (e.which == 27 && e.shiftKey) {
            toggleConsole(conin);
        }
        if (e.which == 13 && consoleVisible()) {
            var command = conin.value.trim();
            if (command.length > 0) {
                console.log("> " + command)
                try {
                    evaluateInput(command);
                } catch (e) {
                    console.error(e.message);
                }
                conin.value = "";
            }
        }
    };


    function consoleVisible() {
        return debugconsole.style.display == "block";
    }

    function toggleConsole(conin) {
        var display = "none"
        if (debugconsole.style.display == "none") display = "block";
        debugconsole.style.display = display;
        if (display == "block") {
            conin.focus();
        }
    }

    function evaluateExpression(expression) {

        var pathSegments = [].concat(expression.split('.'));
        if (pathSegments[0] == 'window') {
            pathSegments.shift()
        }
        var currentObject = window;
        for (var i = 0; i < pathSegments.length; i++) {
            var pathSegment = pathSegments[i];
            if (currentObject[pathSegment] == undefined) {
                return false;
            }
            currentObject = currentObject[pathSegment];
        }
        console.log(JSON.stringify(currentObject));

        return true;
    }

    function evaluateInput(command) {
        try {
            if (evaluateExpression(command)) {
                return
            } else {
                eval(command);
            }
        } catch (e) {
            console.error(e.message)
        }
    }


    // Set us up as a listener
    function hookIntoIPC() {
        if (window.wails && window.wails._ && window.wails._.AddIPCListener) {
            window.wails._.AddIPCListener(processIPCMessage);
        } else {
            setTimeout(hookIntoIPC, 100);
        }
    }
    hookIntoIPC();

    function processIPCMessage(message) {
        console.log(message);
        var parsedMessage;
        try {
            parsedMessage = JSON.parse(message);
        } catch (e) {
            console.error("Error in parsing IPC message: " + e.message);
            return false;
        }
        var logmessage = "[IPC] "
        switch (parsedMessage.type) {
            case 'call':
                logmessage += " Call: " + parsedMessage.payload.bindingName;
                var params = "";
                var parsedParams = JSON.parse(parsedMessage.payload.data);
                if (parsedParams.length > 0) {
                    params = parsedParams;
                }
                logmessage += "(" + params + ")";
                break;
            case 'log':
                logmessage += "Log (" + parsedMessage.payload.level + "): " + parsedMessage.payload.message;
                break;
            default:
                logmessage = message;
        }
        console.log(logmessage);
    }

    window.wailsconsole.clearConsole = function () {
        consoleOut.innerHTML = "";
    }

    console.log = function (message) {
        consoleOut.innerHTML = consoleOut.innerHTML + "<span style='padding-left: 5px'>" + message + '</span><br/>';
        consoleOut.scrollTop = consoleOut.scrollHeight;

    };
    console.error = function (message) {
        consoleOut.innerHTML = consoleOut.innerHTML + "<span style='color:red; padding-left: 5px'> Error: " + message + '</span><br/>';
        consoleOut.scrollTop = consoleOut.scrollHeight;
    };
    // var h = document.getElementsByTagName("html")[0];
    // console.log("html margin: " + h.style.marginLeft);
    // console.log("html padding: " + h.style.paddingLeft);

    // setInterval(function() { console.log("test");}, 1000);
    // setInterval(function() { console.error("oops");}, 3000);
    // var script = document.createElement('script'); 
    // script.src = "https://cdnjs.cloudflare.com/ajax/libs/firebug-lite/1.4.0/firebug-lite.js#startOpened=true"; 
    // document.body.appendChild(script); 

})();