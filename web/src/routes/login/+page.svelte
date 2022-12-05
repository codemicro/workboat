<script lang="ts">
    import { doLogin } from "$lib/util";
    import Page from "$lib/components/Page.svelte";
</script>

<Page>
    <a href="#" on:click={doLogin()}>Click here to login with Gitea.</a>
</Page>