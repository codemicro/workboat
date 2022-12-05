<script lang="ts">
    import { checkLogin } from "$lib/util";
    import Page from "$lib/components/Page.svelte";

    checkLogin();
</script>

<Page>
    <h2>Dashboard</h2>

    <a class="btn btn-success" href="/dashboard/install" role="button">Install into new repository</a>
</Page>