import React from 'react';
import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";

export default function TabsFrameworks() {
    return (
//@formatter:off
<Tabs
    defaultValue="Svelte"
    values={[
        {label: "Svelte", value: "Svelte"},
        {label: "React", value: "React"},
        {label: "Vue", value: "Vue"},
        {label: "Preact", value: "Preact"},
        {label: "Lit", value: "Lit"},
        {label: "Vanilla", value: "Vanilla"},
    ]}
>
<TabItem value="Svelte">
    Generate a <a href="https://svelte.dev/">Svelte</a> project using Javascript with:<br/>

    <code>wails init -n myproject -t svelte</code>

    If you would rather use Typescript:

    <code>wails init -n myproject -t svelte-ts</code>
</TabItem>
<TabItem value="React">
    Generate a <a href="https://reactjs.org/">React</a> project using Javascript with:<br/>

    <code>wails init -n myproject -t react</code>

    If you would rather use Typescript:

    <code>wails init -n myproject -t react-ts</code>
</TabItem>
<TabItem value="Vue">
    Generate a <a href="https://vuejs.org/">Vue</a> project using Javascript with:<br/>

    <code>wails init -n myproject -t vue</code>

    If you would rather use Typescript:

    <code>wails init -n myproject -t vue-ts</code>
</TabItem>
<TabItem value="Preact">
    Generate a <a href="https://preactjs.com/">Preact</a> project using Javascript with:<br/>

    <code>wails init -n myproject -t preact</code>

    If you would rather use Typescript:

    <code>wails init -n myproject -t preact-ts</code>
</TabItem>
<TabItem value="Lit">
    Generate a <a href="https://lit.dev/">Lit</a> project using Javascript with:<br/>

    <code>wails init -n myproject -t lit</code>

    If you would rather use Typescript:

    <code>wails init -n myproject -t lit-ts</code>
</TabItem>
<TabItem value="Vanilla">
    Generate a Vanilla project using Javascript with:<br/>

    <code>wails init -n myproject -t vanilla</code>

    If you would rather use Typescript:

    <code>wails init -n myproject -t vanilla-ts</code>
</TabItem>
</Tabs>
//@formatter:on
    )
}





