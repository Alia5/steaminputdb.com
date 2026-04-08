<script lang="ts">
import { resolve } from '$app/paths';
import { page } from '$app/state';
import Logo from '$lib/assets/logo.svg?component';
import { intersectionObserver } from '$lib/attachments/intersectionObserver.svelte';
import { setStylePropertyCallback } from '$lib/attachments/setStylePropertyCallback.svelte';
import Themetoggle from '$lib/components/theme/toggle.svelte';
import { swipeable } from '@svelte-put/swipeable';
import { cubicIn, cubicOut } from 'svelte/easing';
import { fly } from 'svelte/transition';
import Modal from '../Modal.svelte';
import Search from './Search.svelte';
import User from './User.svelte';

let modal = $state<Modal>()!;
let headerVisible = $state(false);
let header = $state<HTMLElement>()!;
let lastScrollY = $state(0);
let upScrollDistance = $state(0);
let downScrollDistance = $state(0);

const SHOW_THRESHOLD_PX = 72;
const HIDE_THRESHOLD_PX = 72;

let searchShowsResults = $state(false);

import IconNews from '~icons/material-symbols/news-rounded';
import IconAccount from '~icons/mdi/account';
import IconLogout from '~icons/mdi/logout';
import IconSearch from '~icons/mdi/magnify';
import IconMenu from '~icons/mdi/menu';

let windowHeight = $state(0);
let windowWidth = $state(0);

$effect(() => {
	if (windowHeight <= 742 && windowWidth >= windowHeight) {
		header.style.position = 'relative';
		header.style.removeProperty('transform');
	} else {
		header.style.position = 'sticky';
	}
});
</script>

<svelte:window
	onscroll={() => {
		if (windowHeight <= 742 && windowWidth >= windowHeight) {
			return;
		}
		if (searchShowsResults) {
			return;
		}
		const y = window.scrollY;
		const delta = y - lastScrollY;
		const direction = delta > 0 ? -1 : delta < 0 ? 1 : 0;
		lastScrollY = y;
		if (direction === 0) {
			return;
		}

		if (direction === 1) {
			upScrollDistance += Math.abs(delta);
			downScrollDistance = 0;
		} else {
			downScrollDistance += Math.abs(delta);
			upScrollDistance = 0;
		}

		if (direction === 1 && !headerVisible && upScrollDistance >= SHOW_THRESHOLD_PX) {
			header.style.removeProperty('transform');
			upScrollDistance = 0;
			return;
		}
		if (direction === -1 && headerVisible && downScrollDistance >= HIDE_THRESHOLD_PX) {
			header.style.transform = 'translateY(-100%)';
			downScrollDistance = 0;
			return;
		}
	}}
	bind:innerHeight={windowHeight}
	bind:innerWidth={windowWidth} />

<header
	bind:this={header}
	{@attach intersectionObserver(
		(isInterSecting) => {
			headerVisible = isInterSecting;
		},
		0.5,
		true
	)}>
	<a class="home" href={resolve('/')}>
		<Logo height="1.6em" />
		<span>SteamInputDB</span>
	</a>
	<button
		class="plain"
		onclick={() => {
			modal.setScrimOpacityMulti(1);
			modal.toggle();
		}}
		aria-label="Menu">
		<IconMenu style="width: 100%; height: 100%;" />
		{#if page.url.pathname.endsWith('/search')}
			<span>SteamInputDB</span>
		{/if}
	</button>
	{#if !page.url.pathname.endsWith('/search')}
		<Search ondisplayresultschange={(s) => (searchShowsResults = s)} />
	{:else}
		<nav>
			{@render navcontent()}
		</nav>
	{/if}
	<div>
		<User />
		<div>
			<Themetoggle />
		</div>
	</div>
	{#if !page.url.pathname.endsWith('/search')}
		<nav>
			{@render navcontent()}
		</nav>
	{:else}
		<div style="height: 0;"></div>
	{/if}
</header>
<Modal bind:this={modal} --backdrop-filter="blur(2px)">
	<aside
		in:fly={{ duration: 196, x: '-100%', y: 0, easing: cubicOut, opacity: 1 }}
		out:fly={{ duration: 196, x: '-100%', y: 0, easing: cubicIn, opacity: 1 }}
		use:swipeable={{
			direction: 'left'
		}}
		use:setStylePropertyCallback
		onsetstyleproperty={(e) => {
			const { name, value } = e.detail;
			if (name === '--swipe-distance-x') {
				let swipeDist = Math.min(Number(value.replace('px', '').trim()), 0);
				modal.setScrimOpacityMulti(1 - Math.abs(swipeDist) / (window.innerWidth / 0.75));
			}
		}}
		onswipeend={modal.swipeend}>
		<div>
			<a class="home" href={resolve('/')}>
				<Logo height="1.6em" width="1.6em" />
				<span>SteamInputDB</span>
			</a>
			<Themetoggle />
		</div>
		<nav>
			{@render navcontent()}
			{#if page.data.steamId}
				<div class="separator"></div>
				<a href={resolve(`/user/${page.data.steamId}`)}>
					<IconAccount style="width: 1.4em; height: 1.4em;" />
					<span>My Profile</span>
				</a>
				<a href={resolve(`/logout`)}>
					<IconLogout style="width: 1.4em; height: 1.4em;" />
					<span>Logout</span>
				</a>
			{/if}
		</nav>
	</aside>
</Modal>

{#snippet navcontent()}
	<a href={resolve('/news')} onclick={() => modal.close()}>
		<IconNews style="font-size: 1.4em;" />
		<span>News</span>
	</a>
	<a href={resolve('/config/search')} onclick={() => modal.close()}>
		<IconSearch style="font-size: 1.4em;" />
		<span>Advanced Search</span>
	</a>
{/snippet}

<style lang="postcss">
header {
	padding: 1em;
	position: sticky;
	top: 0;
	isolation: isolate;

	min-height: 6em;
	z-index: 2;

	&::after {
		content: '';
		position: absolute;
		inset: 0;
		z-index: -1;
		background: var(--card-background-noise);
		background-color: var(--card-color);
	}

	width: 100%;
	overflow-x: clip;

	display: grid;
	grid-template-columns: minmax(3.2em, 1fr) auto minmax(2em, 1fr);
	grid-template-rows: auto auto;
	@media (orientation: portrait) {
		grid-template-columns: minmax(2em, 1fr) auto minmax(2em, 1fr);
		grid-template-rows: auto;
	}
	align-items: center;
	box-shadow: 0 0px 4px var(--shadow-color);
	column-gap: 1em;
	row-gap: 0;
	transition-property: all;

	& > :first-child {
		margin-right: auto;
		padding: 0.25em;
		grid-row: 1 / -1;
	}

	& > :nth-child(2) {
		display: none;
		grid-row: 1 / -1;
	}
	& :global(> :nth-child(3)) {
		z-index: 2;
	}

	& > :nth-last-child(2) {
		display: grid;
		grid-auto-flow: column;
		gap: 1em;
		margin-left: auto;
		place-items: center;
		overflow: clip;
		max-width: 100%;
		padding: 0.25em;
		overflow-clip-margin: 2em;
		grid-row: 1 / -1;
		grid-column: 3;
	}
	& > nav {
		display: flex;
		flex-flow: row wrap;
		justify-content: center;
		align-items: center;
		gap: 2em;
		margin-top: 1em;
		z-index: 0;

		@media (orientation: portrait) {
			display: none;
		}

		> a {
			font-size: 1.16em;
			font-weight: bold;
			padding: 0;
			display: grid;
			grid-auto-flow: column;
			place-items: center;
			width: fit-content;
			& :global(> :first-child) {
				display: none;
			}
		}
	}
}
.home {
	grid-template-columns: min-content auto;
	font-weight: bold;
	height: fit-content;
	& span {
		overflow: hidden;
		text-overflow: ellipsis;
		width: 100%;
	}

	&:hover,
	&:focus-visible {
		:global(path) {
			fill: var(--color-primary) !important;
		}
	}
}

@media (orientation: portrait) {
	header {
		& > :first-child {
			display: none;
		}
		& > :nth-child(2) {
			display: grid;
		}
		& > :last-child {
			& > :last-child {
				display: none;
			}
		}
	}
}

a {
	font-size: 1.4em;
	white-space: nowrap;
	display: grid;
	place-items: center;
	grid-auto-flow: column;
	width: fit-content;
	gap: 0.5em;
	color: var(--text-color);
	&:hover,
	&:focus,
	&:focus-within {
		color: var(--color-primary);
	}
}

aside {
	nav {
		display: grid;
		justify-items: auto;
		align-items: start;
		padding: 2em;
		height: fit-content;
		gap: 1em;
	}
}

button {
	display: grid;
	grid-template-columns: 2em auto;
	gap: 1em;
	place-items: center;
	padding: 0;
	border: none;
	box-shadow: none;
	max-height: 2em;
	overflow: hidden;
	& > span {
		font-size: 1.4em;
		font-weight: bold;
		width: 100%;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}
}

aside {
	display: grid;
	grid-template-rows: min-content;

	min-width: 75%;
	width: fit-content;
	height: 100%;

	background: var(--card-background-noise);
	backdrop-filter: blur(1px);
	background-color: light-dark(
		color-mix(in srgb, var(--card-color), transparent 80%),
		color-mix(in srgb, var(--card-color), transparent 20%)
	);

	/* &::before {
		content: '';
		position: absolute;
		inset: 0;
		background: var(--card-background-noise);
		z-index: -1;
	} */

	transition: background var(--transition-duration) var(--default-ease);
	transform: translateX(clamp(-100%, var(--swipe-distance-x), 0px));
	& > :first-child {
		display: grid;
		grid-template-columns: auto min-content;
		min-width: 20ch;
		padding: 1em;
		gap: 1em;
		width: 100%;
		& > :first-child {
			margin-right: auto;
		}
	}
}

.separator {
	height: 1px;
	width: 100%;
	background: var(--text-color);
	opacity: 0.5;
	margin: 0.5em 0;
}
</style>
