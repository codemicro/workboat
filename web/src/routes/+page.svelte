<script lang="ts">
    import { onMount } from "svelte";
    import { isLoggedIn } from "$lib/util";
    import Page from "$lib/components/Page.svelte";

    let cookieExists = true;
    let nextURL;

    onMount(() => {
        cookieExists = isLoggedIn();

        if (cookieExists) {
            nextURL = "/dashboard";
        } else {
            nextURL = "/login";
        }
        window.location.replace(nextURL);
    })
</script>

<Page>
    <p>If you're not redirected shortly, click <a href="{nextURL}">here</a>.</p>
</Page>