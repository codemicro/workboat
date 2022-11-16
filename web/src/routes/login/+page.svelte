<script lang="ts">
    import { onMount } from 'svelte';
    import { makeURL } from "$lib/urls";
    import Page from "$lib/components/Page.svelte";

    let loginURL :string;

    onMount(async () => {
        const res = await fetch(makeURL("/api/auth/newLogin"));
        if (!res.ok) {
            return;
        }

        loginURL = await res.json();
        window.location.replace(loginURL);
    });
</script>

<Page>
    <h2>You're being redirected...</h2>

    {#if loginURL }
        <p>If nothing happens, try clicking <a href="{loginURL}">here</a>.</p>
    {:else}
        <p>Please wait</p>
    {/if}
</Page>