<script lang="ts">
import { page } from '$app/state';
import { buildSteamLoginUrl } from '$lib/steam_login';
import Icon from '@iconify/svelte';
import Content from './content.svx';
</script>

<svelte:head>
	<title>SteamInputDB - Login</title>
	<meta property="og:title" content="SteamInputDB - Login" />
	<meta name="description" content="Login to SteamInputDB using your Steam account" />
	<meta name="twitter:card" content="summary_large_image" />
	<meta name="twitter:title" content="SteamInputDB - Login" />
</svelte:head>

<main>
	<div>
		<div>
			<Content />
		</div>
		<a class="button" href={buildSteamLoginUrl(page.url.origin)} rel="external"
			>Sign in through <Icon icon="mdi:steam" width="1.4em" height="1.4em" /> Steam
		</a>
	</div>
</main>

<style lang="postcss">
main {
	padding: 2em;
	display: grid;
	place-items: center;
	width: 100%;
	height: 100%;
	& > div {
		display: grid;
		place-items: center;
		gap: 2em;
		& > :first-child {
			display: grid;
			gap: 0.75em;
			text-align: justify;
			max-width: 80ch;
			:global(h1),
			:global(h2),
			:global(h3),
			:global(h4),
			:global(h5),
			:global(h6) {
				margin-top: 1em;
			}
		}
	}
}

a {
	border-radius: 100vh;
	font-size: 1.4em;
	& > :global(svg) {
		margin: 0 0.5ch;
	}
	width: fit-content;
	height: fit-content;
	background: var(--card-glass);
	background-color: color(from var(--color-primary) srgb r g b / 0.2);
	padding: 0.75em 1.25em;
	color: var(--text-color);
	font-weight: bold;

	&:hover,
	&:focus-visible {
		background-color: color(from var(--color-primary) srgb r g b / 0.4);
	}
}
</style>
