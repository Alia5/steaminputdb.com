<script lang="ts">
import { browser } from '$app/environment';
import { client, type ResponseType } from '$lib/api/client';
import { Debounced } from '$lib/debounce.svelte';
import { log } from '$lib/log';
import { cubicIn, cubicInOut, cubicOut } from 'svelte/easing';
import { fly, slide } from 'svelte/transition';
import Searchbar from '../search/searchbar.svelte';
import Spinner from '../Spinner.svelte';
import SearchResults from './SearchResults.svelte';

let { ondisplayresultschange }: { ondisplayresultschange?: (showsResults: boolean) => void } = $props();

let form = $state<HTMLFormElement>()!;
let debounced = new Debounced<string>(223).eager();
let focusState = new Debounced(96, false).eager();

let previewResults = $state<ResponseType<'POST', '/v1/search/'>['data']>({
	games: [],
	configs: []
});

let forceShowLoading = $state(false);
let showError = $state(false);

let isSteamBigPicture = $derived(browser ? navigator?.userAgent?.includes('Steam Gamepad') : false);

let shouldShowWhat = $derived.by(() => {
	if (isSteamBigPicture) {
		return false;
	}
	if (!focusState.value) {
		return false;
	}

	if (forceShowLoading) {
		return 'loading';
	}

	if (showError) {
		return 'error';
	}

	if ((previewResults?.games?.length || 0) === 0 && (previewResults?.configs?.length || 0) === 0) {
		if (debounced.isLocked()) {
			return 'loading';
		}
		return 'no_results';
	}
	return 'results';
});

$effect(() => {
	if (debounced.value == undefined) {
		return;
	}
	fetchLivePreview(debounced.value);
});

$effect(() => {
	ondisplayresultschange?.(shouldShowWhat === 'results');
});

const fetchLivePreview = (search_term = '') => {
	debounced.lock();

	const timeout = setTimeout(() => {
		forceShowLoading = true;
	}, 1250);

	showError = false;
	client
		.POST('/v1/search/', {
			body: {
				search_term,
				limit_games: 3,
				limit_configs: 5,
				limit: 10
			},
			signal: AbortSignal.timeout(12000)
		})
		.then((res) => {
			forceShowLoading = false;
			if (res.error) {
				showError = true;
				log.error('Error fetching search preview', 'err', res.error);
				return;
			}
			if (res.data) {
				previewResults = res.data;
				return;
			}
			showError = true;
		})
		.catch((err) => {
			forceShowLoading = false;
			showError = true;
			log.error('Error fetching search preview', 'err', err);
		})
		.finally(() => {
			debounced.unlock(true);
			clearTimeout(timeout);
		});
};
</script>

<search
	onfocusin={() => {
		if (debounced.value == undefined) {
			fetchLivePreview();
		}
		focusState.input = true;
	}}
	onfocusout={() => (focusState.input = false)}>
	<!-- TODO: change endpoint -->
	<form bind:this={form} method="POST" action="/?/search">
		<Searchbar
			name="searchtext"
			placeholder="Search configurations or games..."
			--box-shadow="inset 0 0.4em 0.4em 0 var(--shadow-color-dark)"
			bind:value={debounced.input} />
	</form>
	{#if shouldShowWhat}
		<dialog
			open
			data-no-focus-trap
			in:fly={{
				y: '-100%',
				x: 0,
				delay: 196 /* wait for child layout */,
				opacity: 0,
				duration: 196,
				easing: cubicOut
			}}
			out:fly|global={{ y: '-100%', x: 0, delay: 0, duration: 196, easing: cubicIn, opacity: 0 }}>
			{#if shouldShowWhat === 'loading'}
				<div class="ctr" transition:slide|global={{ duration: 196, easing: cubicInOut }}>
					<Spinner size="min(100dvw, 12em)" thickness="0.3em" />
				</div>
			{:else if shouldShowWhat === 'no_results'}
				<div class="ctr" transition:slide|global={{ duration: 196, easing: cubicInOut }}>
					<span>No results found</span>
				</div>
			{:else if shouldShowWhat === 'error'}
				<div class="ctr" transition:slide|global={{ duration: 196, easing: cubicInOut }}>
					<strong>Error fetching results</strong>
				</div>
			{:else}
				<div transition:slide|global={{ duration: 196, delay: 0, easing: cubicInOut }}>
					<SearchResults results={previewResults} />
				</div>
			{/if}
		</dialog>
	{/if}
</search>

<style lang="postcss">
search {
	position: relative;
}

form {
	place-self: center;
	display: grid;
	grid-template-columns: minmax(16ch, 48ch);
}

.ctr {
	display: grid;
	place-items: center;
	padding: 2em;
}

dialog[open] {
	--oversize: 50%;
	--corner-radius: 1em;
	position: absolute;
	top: 100%;
	left: 50%;
	min-height: 2em;
	translate: -50% 0;
	width: calc(100% + var(--oversize));
	max-width: 100dvw;
	z-index: -10;

	padding-top: 0.5em;

	border: none;
	border-radius: 0 0 var(--corner-radius) var(--corner-radius);
	margin-top: 0.5em;
	background: var(--card-background-noise);
	background-color: var(--card-color);
	box-shadow: 0 0.5em 1.5em -0.5em var(--shadow-color);

	display: grid;
	place-items: center;

	& > * {
		grid-area: 1 / 1;
		width: 100%;
	}

	&::before {
		@supports (corner-shape: scoop) {
			content: '';
		}
		position: absolute;
		top: 0.5em;
		left: calc(var(--corner-radius) * -1);
		corner-shape: scoop;
		border-radius: 0 0 0 var(--corner-radius);
		width: var(--corner-radius);
		height: var(--corner-radius);
		background: var(--card-background-noise);
		background-color: var(--card-color);
		z-index: -1;
	}
	&::after {
		@supports (corner-shape: scoop) {
			content: '';
		}
		position: absolute;
		top: 0.5em;
		right: calc(var(--corner-radius) * -1);
		corner-shape: scoop;
		border-radius: 0 0 var(--corner-radius) 0;
		width: var(--corner-radius);
		height: var(--corner-radius);
		background: var(--card-background-noise);
		background-color: var(--card-color);

		z-index: -1;
	}
}
</style>
