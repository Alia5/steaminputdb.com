<script lang="ts">
import { resolve } from '$app/paths';
import SC2 from '$lib/assets/SC2_Googley.svg.svelte';
import { createHomeSchemaJsonLd } from '$lib/schema/home';
import { onMount } from 'svelte';
import type { PageProps } from './$types';
import HotTodayBar from './HotTodayBar.svelte';

const { data }: PageProps = $props();

let eyes = $state<{
	left: HTMLElement;
	right: HTMLElement;
}>({
	left: undefined!,
	right: undefined!
})!;

onMount(() => {
	const group = document.querySelector('#sc2>*>svg>g>:last-child')!;
	eyes.left = group.children[0] as HTMLElement;
	eyes.right = group.children[1] as HTMLElement;
});

$inspect(data.hotToday);
</script>

<svelte:head>
	<title>SteamInputDB | Database of every Steam Input configuration</title>
	<meta
		name="description"
		content="Community-driven database of Steam Input configurations using the Steam API." />
	<meta
		name="keywords"
		content="Steam Input DB, Steam DB, DB, Steam Deck, Steam Input, Steam controller configs, controller layouts, community database, Steam API, gamepad configurations, controller presets" />
	<meta
		name="robots"
		content="index, follow, max-image-preview:large, max-snippet:-1, max-video-preview:-1" />

	<meta property="og:site_name" content="SteamInputDB" />
	<meta property="og:type" content="website" />
	<meta property="og:title" content="SteamInputDB | Database of every Steam Input configuration" />
	<meta
		property="og:description"
		content="SteamInputDB is a database of every Steam Input configuration using the Steam API." />
	<meta property="og:url" content="https://www.steaminputdb.com/" />
	<meta property="og:image" content="https://www.steaminputdb.com/ogimage.png" />
	<meta property="og:image:alt" content="SteamInputDB preview image" />

	<meta name="twitter:card" content="summary_large_image" />
	<meta name="twitter:title" content="SteamInputDB | Database of every Steam Input configuration" />
	<meta
		name="twitter:description"
		content="SteamInputDB is a database of every Steam Input configuration using the Steam API." />
	<meta name="twitter:image" content="https://www.steaminputdb.com/ogimage.png" />
	<meta name="twitter:image:alt" content="SteamInputDB preview image" />

	<svelte:element this={'script'} type="application/ld+json">{createHomeSchemaJsonLd()}</svelte:element>
</svelte:head>

<svelte:window
	onmousemove={(e) => {
		if (!eyes || !eyes.left || !eyes.right) {
			return;
		}
		Object.values(eyes).forEach((eye) => {
			const rect = eye.getBoundingClientRect();
			const eyeX = rect.left + rect.width / 2;
			const eyeY = rect.top + rect.height / 2;
			const deltaX = e.clientX - eyeX;
			const deltaY = e.clientY - eyeY;
			const angle = Math.atan2(deltaY, deltaX);
			const distance = Math.min(16, Math.hypot(deltaX, deltaY) / 2);
			const translateX = Math.cos(angle) * distance;
			const translateY = Math.sin(angle) * distance;
			eye.style.transform = `translate(${translateX}px, ${translateY}px)`;
		});
	}} />

<main>
	<div>
		<div>
			<section id="sc2">
				<SC2
					height="100%"
					--eyes-color="black"
					--eyes-white-color="var(--text-color-dark)"
					--eyes-border-color="light-dark(var(--text-color-light), transparent)" />
			</section>
			<section>
				<h1>
					Community driven database of <span>SteamInput</span>
					configurations
					<strong>utilizing Steam API</strong>
				</h1>
				<p>
					Providing an easy way to <strong>browse</strong>, <strong>filter</strong>,
					<strong>share</strong>, and <strong>apply</strong> community-shared SteamInput controller layouts/configurations.
				</p>
				<br />
				<p>SteamInputDB uses the same APIs as Steam itself</p>
				<p>That means <strong>every</strong> configuration on Steam is also available here!</p>
				<p>And yes, it works with <strong>non-Steam</strong> games, too! 😎</p>
			</section>
		</div>

		<!-- <section id="search">
			<Search />
		</section> -->

		<section id="hot">
			<h2>Hot Today</h2>
			<HotTodayBar hotConfigs={data.hotToday} />
		</section>

		<section id="buddy">
			<h2><span>SteamInputDB-<strong>Buddy App</strong></span> available now! <em>(Beta)</em></h2>
			<p>SteamInputDB-Buddy directly integrates this website with your Steam client.</p>
			<p>
				Browse this website directly in Steam or apply configurations to any game with a single click
			</p>
			<a href={resolve('/buddy-app/install')} class="button">Learn more | Install</a>
		</section>
	</div>
</main>

<style lang="postcss">
main {
	padding: 2em;
	display: grid;
	place-items: center;
	position: relative;
	& > div {
		display: grid;
		place-items: center;
		gap: 4em;
		width: 100%;
		height: 100%;
	}
}

#search {
	margin: 2em 0;
	@media (orientation: landscape) {
		font-size: 1.5em;
	}
}

main > div > :first-child {
	display: grid;
	place-items: center;
	gap: 1em;
}

h1 {
	& strong {
		color: var(--color-primary);
	}
	& span {
		color: var(--highlight-color);
	}
	margin-bottom: 0.5em;
}
p {
	font-size: 1.2em;
	font-weight: bold;
	& strong {
		color: var(--highlight-color);
	}
}

#sc2 {
	filter: drop-shadow(0px 0.25em 0.2em var(--shadow-color));
	max-height: 16em;
	height: 100%;
	width: 100%;
	display: grid;
	place-items: center;
	padding: 0.5em;
	@media (any-pointer: coarse) {
		:global(ellipse) {
			transition: transform calc(var(--transition-duration) * 2) var(--default-ease);
		}
	}
}

#buddy {
	display: grid;

	& :global(strong) {
		color: var(--color-primary);
	}

	justify-self: center;
	& strong {
		color: var(--color-primary);
	}
	& span {
		color: var(--highlight-color);
	}
	& em {
		font-size: 0.8em;
		opacity: 0.75;
	}

	& > :first-child {
		font-size: 1.8em;
		justify-self: center;
		margin-bottom: 0.5em;
		@media screen and (max-width: 400px) {
			font-size: 1.5em;
		}
	}

	a.button {
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 1.6em;
		width: fit-content;
		justify-self: center;
		margin-top: 1em;
		gap: 0.5em;
		padding: 0.5em 1.5em;
		color: var(--text-color);
		border-radius: 99vw;
		font-weight: bold;
		& em {
			font-size: 0.8em;
			opacity: 0.75;
		}
	}
}

#hot {
	max-width: 100%;
	overflow: hidden;
	display: grid;
	place-items: center;
	& > :first-child {
		font-size: 1.8em;
		justify-self: center;
		margin-bottom: 0.5em;
		@media screen and (max-width: 400px) {
			font-size: 1.5em;
		}
	}
}
</style>
