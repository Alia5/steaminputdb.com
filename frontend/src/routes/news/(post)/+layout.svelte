<script lang="ts">
const { children } = $props();
</script>

<svelte:head>
	<meta property="og:type" content="article" />
</svelte:head>

<main class="mdcontainer">
	{@render children?.()}
</main>

<style lang="postcss">
main {
	grid-row: 2 / span 1;
	grid-column: 1 / span 1;

	display: grid;
	place-self: baseline center;
	gap: 1em;
	min-width: 50%;
	--max-width: 1440px;
	max-width: min(100%, var(--max-width));

	display: grid;
	text-align: justify;
	height: fit-content;
	:global(h1),
	:global(h2),
	:global(h3),
	:global(h4),
	:global(h5),
	:global(h6) {
		margin-top: 1em;
		margin-bottom: 0.5em;
	}

	:global(> :first-child) {
		margin-bottom: 0;
	}
	:global(.lead) {
		opacity: 0.7;
	}

	:global(p) {
		text-align: justify;
	}

	padding: 0 2em;
	padding-bottom: 2em;

	& :global(.hljs) {
		background: light-dark(#282c34, transparent);
	}
	& :global(a) {
		color: var(--highlight-color);
	}
}
</style>
