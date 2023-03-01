function GreetService(method) {
    return {
        packageName: "main",
        serviceName: "GreetService",
        methodName: method,
        args: Array.prototype.slice.call(arguments, 1),
    };
}

/**
 * GreetService.Greet
 * Greet someone
 * @param name {string}
 * @returns {Promise<string>}
 */
function Greet(name) {
    return wails.Call(GreetService("Greet", name));
}

window.go = window.go || {};
Object.window.go.main = {
    GreetService: {
        Greet,
    }
};
