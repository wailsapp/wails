<script>
    import Link from '../../components/Link.svelte';
    import Browser from './Browser/Browser.svelte';
    
</script>
<div>
    Sometimes it is necessary to redirect users to external websites. The Wails runtime offers this ability through the <code>runtime.Browser</code> component. 
    
    By using <code>runtime.Browser.Open</code>, it is possible to open the user's default browser to either navigate to a url (pass in a URL), or open a local file (provide a local filepath).

    <Link href="https://wails.app">Here</Link> is an example.
    
    <br/><br/>
    
    <Browser></Browser>
</div>

