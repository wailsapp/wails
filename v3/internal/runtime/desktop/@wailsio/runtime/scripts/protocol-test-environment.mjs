// Protocol tests need browsing-context state, not DOM rendering or parsing.
class SessionValues {
    #values = new Map();

    getItem(key) {
        return this.#values.get(String(key)) ?? null;
    }

    setItem(key, value) {
        this.#values.set(String(key), String(value));
    }

    clear() {
        this.#values.clear();
    }
}

export default {
    name: "wails-protocol",
    viteEnvironment: "ssr",
    setup(global) {
        const values = {
            window: {
                name: "",
                location: new URL("http://localhost:3000"),
                sessionStorage: new SessionValues(),
            },
            document: Object.freeze({}),
            Storage: SessionValues,
        };
        const originals = new Map();
        for (const [key, value] of Object.entries(values)) {
            originals.set(key, Object.getOwnPropertyDescriptor(global, key));
            Object.defineProperty(global, key, { value, configurable: true, writable: true });
        }
        return {
            teardown() {
                for (const [key, descriptor] of originals) {
                    if (descriptor) Object.defineProperty(global, key, descriptor);
                    else delete global[key];
                }
            },
        };
    },
};
