import { afterEach, describe, expect, it, vi } from "vitest";

import * as IOS from "./ios.js";
import { objectNames, setTransport } from "./runtime.js";

describe("IOS.NativeTabs", () => {
    afterEach(() => setTransport(null));

    it("sends native tab controls through the iOS runtime object", async () => {
        const call = vi.fn().mockResolvedValueOnce(undefined).mockResolvedValueOnce(undefined).mockResolvedValueOnce(true);
        setTransport({ call });

        await IOS.NativeTabs.SetEnabled(false);
        await IOS.NativeTabs.Select(2);
        await expect(IOS.NativeTabs.IsEnabled()).resolves.toBe(true);

        expect(call).toHaveBeenNthCalledWith(1, objectNames.IOS, 9, "", { enabled: false });
        expect(call).toHaveBeenNthCalledWith(2, objectNames.IOS, 11, "", { index: 2 });
        expect(call).toHaveBeenNthCalledWith(3, objectNames.IOS, 10, "", null);
    });
});
