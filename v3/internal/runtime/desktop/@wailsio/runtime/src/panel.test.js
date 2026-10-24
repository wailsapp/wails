import {readFileSync} from 'node:fs';
import {resolve} from 'node:path';
import {describe, it, expect, vi, beforeEach} from 'vitest';

const {caller, factory} = vi.hoisted(() => ({caller:vi.fn(), factory:vi.fn()}));
vi.mock('./runtime.js', () => ({
    objectNames: {Panel:13},
    newRuntimeCaller: factory.mockImplementation(() => caller),
}));
import {Panel} from './panel';

const go = readFileSync(resolve('../../../../../pkg/application/messageprocessor_panel.go'), 'utf8');
const backend = Object.fromEntries([...go.matchAll(/\bPanel(\w+)\s*=\s*(\d+)/g)].map(([,name,id]) => [name,Number(id)]));
const methods = [
    ['SetBounds', [{x:1,y:2,width:300,height:200}], {x:1,y:2,width:300,height:200}],
    ['GetBounds', [], {}], ['SetZIndex', [7], {zIndex:7}],
    ['Reload', [], {}], ['ForceReload', [], {}], ['Show', [], {}], ['Hide', [], {}],
    ['IsVisible', [], {}], ['SetZoom', [1.5], {zoom:1.5}], ['GetZoom', [], {}],
    ['Focus', [], {}], ['IsFocused', [], {}], ['OpenDevTools', [], {}],
    ['Destroy', [], {}], ['Name', [], {}],
];

beforeEach(() => { caller.mockReset(); factory.mockClear(); });

describe('Panel runtime contract', () => {
    it.each(methods)('%s dispatches to the corresponding Go method', async (name,args,fields) => {
        const result = {value:'backend result'};
        caller.mockResolvedValue(result);
        const panel = Panel.Get('browser','main');
        expect(factory).toHaveBeenCalledWith(13,'main');
        expect(await panel[name](...args)).toBe(result);
        expect(caller).toHaveBeenCalledWith(backend[name],{panel:'browser',...fields});
        expect(backend[name]).toBeTypeOf('number');
    });
    it('covers every exported instance method and every Go method', () => {
        const names = methods.map(([name])=>name).sort();
        expect(Object.getOwnPropertyNames(Panel.prototype).filter(n=>n!=='constructor').sort()).toEqual(names);
        expect(Object.keys(backend).sort()).toEqual(names);
    });
    it('reserves removed navigation and code execution IDs', () => {
        const panel = new Panel('browser');
        for (const method of ['SetURL','SetHTML','ExecJS']) expect(panel[method]).toBeUndefined();
        for (const id of [3,4,5]) expect(Object.values(backend)).not.toContain(id);
    });
    it('does not send Destroy when querying focus (PR 4880 regression)', async () => {
        const panel = new Panel('browser');
        const isFocused = panel.IsFocused;
        await isFocused();
        expect(caller).toHaveBeenCalledWith(14,{panel:'browser'});
        expect(backend.IsFocused).toBe(14);
        expect(backend.Destroy).toBe(16);
    });
    it('propagates backend errors', async () => {
        caller.mockRejectedValue(new Error('panel not found'));
        await expect(Panel.Get('missing').Show()).rejects.toThrow('panel not found');
    });
});
