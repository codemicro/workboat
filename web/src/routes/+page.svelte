<script lang="ts">
    import { onMount } from "svelte";
    import { isLoggedIn } from "$lib/util";
    import Page from "$lib/components/Page.svelte";

    let cookieExists = true;

    onMount(() => {
        cookieExists = isLoggedIn();
    })
</script>

<Page>
    {#if cookieExists}
        <p>You're logged in!</p>
    {:else}
        <a href="/login">Click here to login with Gitea</a>
    {/if}
</Page>