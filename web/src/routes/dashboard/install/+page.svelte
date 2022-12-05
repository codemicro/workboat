<script lang="ts">
    import { checkLogin, makeAPIURL } from "$lib/util";
    import { onMount } from "svelte";
    import Page from "$lib/components/Page.svelte";

    checkLogin();

    let repositories = undefined;

    onMount(() => {
        fetch(makeAPIURL("/api/install/getRepositories")).then((res) => {
            if (!res.ok) {
                return;
            }
            return res.json();
        }).then((jsonData) => {
            repositories = jsonData;
        });
    });
</script>

<Page>
    <h2>Install into new repository</h2>

    {#if repositories !== undefined}
        {#if repositories.length === 0}
            <p class="text-secondary">No repositories available.</p>
        {:else}
            <select class="form-select" aria-label="Repository selection">
                {#each repositories as repository}
                    <option value={repository.id}>{repository.name}</option>
                {/each}
            </select>
        {/if}
    {/if}
</Page>