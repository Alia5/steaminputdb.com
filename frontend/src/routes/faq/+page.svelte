<script lang="ts">
import Content from './content.svx';
</script>

<svelte:head>
	<title>SteamInputDB - FAQ</title>
	<meta property="og:title" content="SteamInputDB - FAQ" />
	<meta name="description" content="Frequently Asked Questions about SteamInputDB" />
	<meta property="og:description" content="Frequently Asked Questions about SteamInputDB" />
	<meta name="twitter:card" content="summary_large_image" />
	<meta name="twitter:title" content="SteamInputDB - FAQ" />
	<meta name="twitter:description" content="Frequently Asked Questions about SteamInputDB" />
</svelte:head>

<main>
	<div>
		<Content />
	</div>
</main>

<style lang="postcss">
main {
	padding: 2em;
	display: grid;
	gap: 0.75em;
	height: 100%;
	min-width: 50%;
	--max-width: 1440px;
	max-width: min(100%, var(--max-width));
	place-self: center;
}
div {
	display: grid;
	text-align: justify;
	gap: 0.5em;
	height: fit-content;
	:global(h1),
	:global(h2),
	:global(h3),
	:global(h4),
	:global(h5),
	:global(h6) {
		margin-top: 0.5em;
		margin-bottom: 0.2em;
	}
	:global(li) {
		margin-bottom: 0.33em;
	}
	:global(ol),
	:global(ul) {
		overflow: hidden;
		padding-left: 1.5em;
		margin-left: 0;
		margin-bottom: 1em;
	}
	:global(p) {
		margin-bottom: 0.5em;
	}
}
</style>
