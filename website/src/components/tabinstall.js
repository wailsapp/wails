import React from 'react';
import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";

export default function TabsInstaller() {
    return (
//@formatter:off
<Tabs
defaultValue="Windows"
values={[
    {label: "Windows", value: "Windows"},
    {label: "MacOS", value: "MacOS"},
    {label: "Linux", value: "Linux"},
]}
>
<TabItem value="MacOS">
    Wails requires that the xcode command line tools are installed. This can be done by running:<br/>

    <code>xcode-select --install</code>
</TabItem>
<TabItem value="Windows">
    Wails requires that the <a
    href="https://developer.microsoft.com/en-us/microsoft-edge/webview2/">WebView2</a>{" "}
    runtime is installed. Some Windows installations will already have this installed. You can check using
    the{" "}
    <code>wails doctor</code> command (see below).
</TabItem>
<TabItem value="Linux">
    Linux required the standard <code>gcc</code> build tools
    plus <code>libgtk3</code> and <code>libwebkit</code>.
    Rather than list a ton of commands for different distros, Wails can try to determine
    what the installation commands are for your specific distribution. Run <code>wails doctor</code> after
    installation
    to be shown how to install the dependencies.
    If your distro/package manager is not supported, please consult the <a
    href="/docs/guides/linux-distro-support"> Add Linux Distro</a> guide.

</TabItem>

</Tabs>
//@formatter:off
    )
}
