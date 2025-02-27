/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

const TYPED_EVENTS_MODULE = "\0wailsio_runtime_events_typed";

/**
 * A plugin that extends the wails runtime with locally generated code
 * to provide support for typed custom events.
 * With the plugin installed, vite will fail to build the project
 * unless wails bindings have been generated first.
 *
 * @param {string} [bindingsRoot] - The root import path for generated bindings
 */
export default function WailsTypedEvents(bindingsRoot) {
    let bindingsId = null,
        runtimeId = null,
        eventsId = null;

    return {
        name: "wails-typed-events",
        async buildStart() {
            const bindingsPath = `${bindingsRoot}/github.com/wailsapp/wails/v3/internal/eventcreate`;
            let resolution = await this.resolve(bindingsPath);
            if (!resolution || resolution.external) {
                this.error(`Event bindings module not found at import specifier '${bindingsPath}'. Please verify that the wails tool is up to date and the binding generator runs successfully. If you moved the bindings to a custom location, ensure you supplied the correct root path as the first argument to \`wailsTypedEventsPlugin\``);
                return;
            }
            bindingsId = resolution.id;

            resolution = await this.resolve("@wailsio/runtime");
            if (!resolution || resolution.external) { return; }
            runtimeId = resolution.id;

            resolution = await this.resolve("./events.js", runtimeId);
            if (!resolution || resolution.external) {
                this.error("Could not resolve events module within @wailsio/runtime package. Please verify that the module is correctly installed and up to date.");
                return;
            }

            eventsId = resolution.id;
        },
        resolveId: {
            order: 'pre',
            handler(id, importer) {
                if (
                    bindingsId !== null
                    && runtimeId !== null
                    && eventsId !== null
                    && importer === runtimeId
                    && id === "./events.js"
                ) {
                    return TYPED_EVENTS_MODULE;
                }
            }
        },
        load(id) {
            if (id === TYPED_EVENTS_MODULE) {
                return (
                    `import ${JSON.stringify(bindingsId)};\n`
                    + `export * from ${JSON.stringify(eventsId)};`
                );
            }
        }
    }
}
