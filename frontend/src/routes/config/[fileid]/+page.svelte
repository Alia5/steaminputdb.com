<script lang="ts">
import { browser } from '$app/environment';
import { page } from '$app/state';
import type { components } from '$lib/api/openapi';
import LayoutPreview from '$lib/components/layout_preview/LayoutPreview.svelte';
import MarkdownSSR from '$lib/components/markdown/markdownSSR.svelte';
import { configurationFeatureList } from '$lib/snippets/configurationfeaturelist.svelte';
import { assetUrlBase, storePageBackgroundBase } from '$lib/steamapi/const';
import ConfigInfoHeader from './configInfoHeader.svelte';

const fileInfo: components['schemas']['ConfigDetailResponse'] = $derived(page.data.fileInfo);
const appInfo: components['schemas']['AppItem'] = $derived(page.data.appInfo);
const creatorInfo: components['schemas']['UserInfoResponse'] | undefined = $derived(page.data.creatorInfo);

const pageBGURL = $derived.by(() => {
	if (!appInfo?.assets) {
		return;
	}
	if (appInfo.assets.page_background) {
		return `${assetUrlBase}${appInfo.assets.asset_url_format?.replace(
			'${FILENAME}',
			appInfo.assets.page_background
		)}`;
	}
	if (appInfo.assets.raw_page_background) {
		// TODO: find out if is correct base url
		return `${storePageBackgroundBase}${appInfo.assets.asset_url_format?.replace(
			'${FILENAME}',
			appInfo.assets.raw_page_background
		)}`;
	}
	if (appInfo.assets.page_background_path) {
		return `${storePageBackgroundBase}${appInfo.assets.page_background_path}`;
	}
});

let isMobileBrowser = $state(false);
if (browser) {
	const uaDataMobile =
		navigator.userAgent.toLowerCase().includes('mobile') ||
		(navigator as unknown as { userAgentData?: { mobile?: boolean } }).userAgentData?.mobile;
	isMobileBrowser = !!uaDataMobile;
}
</script>

<svelte:head>
	<title>SteamInputDB - Config: {fileInfo?.title} ({appInfo?.name}) by {creatorInfo?.personaname}</title>
	<link rel="canonical" href={page.url.href} />
	<meta property="og:site_name" content="SteamInputDB" />
	<meta property="og:type" content="website" />
	<meta property="og:url" content={page.url.href} />
	<meta
		property="og:title"
		content="SteamInputDB - Config: {fileInfo?.title} ({appInfo?.name}) by {creatorInfo?.personaname}"
	/>
	<meta
		name="description"
		content={fileInfo?.description ??
			`Steam Input configuration ${fileInfo?.title} (${appInfo?.name}) by ${creatorInfo?.personaname}`}
	/>
	<meta
		property="og:description"
		content={fileInfo?.description ??
			`Steam Input configuration ${fileInfo?.title} (${appInfo?.name}) by ${creatorInfo?.personaname}`}
	/>
	{#if appInfo?.assets}
		{@const assets = appInfo?.assets}
		{@const assetChosen =
			assets.main_capsule ?? assets.header ?? assets.hero_capsule ?? assets.library_hero ?? 'none.svg'}
		{#if assetChosen}
			<meta
				property="og:image"
				content={`${assetUrlBase}${assets.asset_url_format?.replace('${FILENAME}', assetChosen)}`}
			/>
			<meta
				name="twitter:image"
				content={`${assetUrlBase}${assets.asset_url_format?.replace('${FILENAME}', assetChosen)}`}
			/>
		{/if}
	{/if}
	<meta name="twitter:card" content="summary_large_image" />
	<meta
		name="twitter:title"
		content="SteamInputDB - Config: {fileInfo?.title} ({appInfo?.name}) by {creatorInfo?.personaname}"
	/>
	<meta
		name="twitter:description"
		content={fileInfo?.description ??
			`Steam Input configuration ${fileInfo?.title} (${appInfo?.name}) by ${creatorInfo?.personaname}`}
	/>
</svelte:head>

<main style={pageBGURL ? `--bg: url('${pageBGURL}')` : ''}>
	<div>
		<ConfigInfoHeader
			fileInfo={fileInfo}
			appInfo={appInfo}
			isMobileBrowser={isMobileBrowser}
			creatorInfo={creatorInfo}
		/>

		<div class="description-container">
			<div class="card glass">
				{#if fileInfo.tags}
					<div>
						<span>Features</span>
						<div>
							<div class="featurelist">
								{@render configurationFeatureList({ fileInfo })}
							</div>
						</div>
					</div>
				{/if}
				{#if fileInfo.description}
					<div>
						<span>Creator Description</span>
						<div style="display: grid; gap: 0.25em; padding: 0 1em;">
							<MarkdownSSR content={fileInfo.description} />
						</div>
					</div>
				{/if}
			</div>
		</div>
		{#if fileInfo.file_url}
			<LayoutPreview vdfLink={fileInfo.file_url || undefined} />
		{:else}
			<div class="legacy">
				<div class="card glass">
					<p>Preview not possible for legacy configurations 🫤</p>
				</div>
			</div>
		{/if}
	</div>
</main>

<style lang="postcss">
main {
	position: relative;
	display: grid;
	padding: 1em 0;

	place-items: center;
	grid-template-rows: min-content;
	grid-template-columns: minmax(min(100%, auto), 50%);
	width: 100%;
	& > div {
		display: grid;
		place-self: center;
		gap: 1em;
		place-items: center;
		min-width: 60%;
		--max-width: 1440px;
		max-width: min(100%, var(--max-width));
		/* container: main / inline-size;*/
		:global(> :first-child) {
			width: 100%;
		}
	}

	&::before {
		content: '';
		position: absolute;
		left: 0;
		right: 0;
		top: 0;
		height: min(100%, 100dvh);
		background: var(--bg, transparent) top/cover no-repeat;
		mask: linear-gradient(0deg, transparent, white 100%);
		mask-type: alpha;
		z-index: -2;
	}
}

.featurelist {
	width: 100%;
	display: flex;
	flex-wrap: wrap;
	gap: 1em;
	overflow: clip;
	overflow-clip-margin: 1em;
	justify-content: center;
}

.description-container {
	padding: 0 1em;
	width: 100%;
	& > div {
		display: flex;
		flex-flow: row wrap;
		gap: 1em;
		align-items: center;
		justify-content: center;
		& > div {
			flex: 1 1 auto;
			display: grid;
			gap: 1em;
			place-items: center;
		}
		& > div > :first-child {
			font-weight: bold;
			font-size: 1.2em;
		}
	}
}

.legacy {
	width: 100%;
	padding: 0 1em;
	display: grid;
	place-items: center;
	text-align: center;
}
</style>
