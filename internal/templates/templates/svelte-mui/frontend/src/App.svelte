<h1>Hello {name}!</h1>

<Snackbar bind:visible bg="#f44336">
    {response}
    <span slot="action">
        <Button color="#ff0" on:click={() => (visible = false)}>Close</Button>
    </span>
</Snackbar>

<Textfield
    autocomplete="off"
    label="Message"
    required
    bind:value
    message=""
/>

<Button
  outlined
  shaped
  color="Red"
  on:click={() => {
    window.backend.main.Basic.Greet(value).then((result) => {
      response = result;
      visible = true;
    });
  }}
>
    Hello
</Button>


<script>
  export let name;
  let value = '';
  let response = '';

  // optional import focus-visible polyfill only once
  import 'focus-visible';
  // import any components
  import { Button, Checkbox, Snackbar, Textfield } from 'svelte-mui';

  let checked = true;
  let visible = false;
</script>

<style>
    h1 {
        color: purple;
    }
</style>
