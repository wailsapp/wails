import { Context, Converter, ReflectionKind } from "typedoc";
import ts from "typescript";

const projects = new WeakMap();
const createId = Context.prototype.createSymbolId;
const typeDeclarations = ts.SymbolFlags.Class | ts.SymbolFlags.Interface |
    ts.SymbolFlags.TypeAlias | ts.SymbolFlags.Enum;

function remember(context, symbol) {
    const state = projects.get(context.converter);
    if (!state?.active) return;
    let target = context.resolveAliasedSymbol(symbol);
    const declaration = target.declarations?.find(node => node.name);
    if (declaration) target = context.getSymbolAtLocation(declaration.name) ?? target;
    if (!(target.flags & typeDeclarations) || state.pending.has(target)) return;
    state.pending.set(target, context.withScope(context.scope));
}

// Standard helper pages already linked by the runtime documentation.
const libraryHelpers = new Set([
    "ArrayBufferView", "CallSite", "ErrorOptions", "Iterable", "PromiseFulfilledResult",
    "PromiseLike", "PromiseRejectedResult", "PromiseWithResolvers", "ArrayBufferLike",
    "Awaited", "EventListenerOrEventListenerObject", "Exclude", "Omit", "Partial",
    "Pick", "PromiseSettledResult", "PropertyKey", "Readonly", "Record",
]);

// Keep the source program while TypeDoc converts public signatures.
Context.prototype.createSymbolId = function (symbol, ...args) {
    const id = createId.call(this, symbol, ...args);
    remember(this, symbol);
    return id;
};

export function load(app) {
    if (projects.has(app.converter)) return;
    const state = { active: false, pending: new Map() };
    projects.set(app.converter, state);
    app.converter.on(Converter.EVENT_BEGIN, () => {
        state.pending.clear();
        state.active = true;
    });
    // Create helpers before TypeDoc resolves source links and removes private members.
    app.converter.on(Converter.EVENT_RESOLVE_BEGIN, context => {
        let helpers;
        for (const [symbol, source] of state.pending) {
            if (source.getReflectionFromSymbol(symbol) || source.shouldIgnore(symbol)) continue;
            const file = symbol.declarations?.[0]?.getSourceFile();
            const external = file && (source.program.isSourceFileDefaultLibrary(file) ||
                source.program.isSourceFileFromExternalLibrary(file));
            if (external && !libraryHelpers.has(symbol.getName())) continue;
            if (!helpers) {
                helpers = context.createDeclarationReflection(ReflectionKind.Module, undefined, undefined, "<internal>");
                context.finalizeDeclarationReflection(helpers);
            }
            const destination = context.withScope(helpers);
            destination.setActiveProgram(source.program);
            destination.preventInline.add(symbol.getName());
            app.converter.convertSymbol(destination, symbol);
        }
        state.active = false;
    }, 10000);
    app.converter.on(Converter.EVENT_END, () => {
        state.active = false;
        state.pending.clear();
    });
}
