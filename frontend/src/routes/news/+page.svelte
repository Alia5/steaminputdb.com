<script lang="ts">
import { resolve } from '$app/paths';
import Author from '$lib/components/author.svelte';
import type { Picture } from '@sveltejs/enhanced-img';

type ArticleMeta = {
	title: string;
	date: Date;
	lead?: string;
	excerpt?: string;
	hero?: Picture;
	author?: {
		name: string;
		img?: Picture;
	};
};

const articlePages = import.meta.glob('/src/routes/news/\\(post\\)/**/+page.svx', {
	eager: true,
	import: 'metadata'
});
const sortedByDate = Object.entries(articlePages).sort(([, a], [, b]) => {
	const dateA = new Date((a as ArticleMeta).date).getTime();
	const dateB = new Date((b as ArticleMeta).date).getTime();
	return dateB - dateA;
});
</script>

<main>
	{#each sortedByDate as [r, m] (r)}
		{@const route = r as string}
		{@const meta = m as ArticleMeta}
		<!-- <enhanced:img src={articleMeta[0].hero}></enhanced:img> -->
		<a
			href={resolve(
				route.replace(/\/src\/routes\/news\/\(post\)/g, '/news').replace(/\/\+page\.svx$/g, '') as '/'
			)}
			style={meta.hero ? '--img-w: auto;' : '--img-w: 0;'}>
			{#if meta.hero}
				<enhanced:img src={meta.hero}></enhanced:img>
			{/if}
			<div>
				<span>{meta.title}</span>
				{#if meta.excerpt || meta.lead}
					{@const excerpt = meta.excerpt ?? meta.lead}
					<p>{excerpt?.replaceAll('\\n', '\n').replaceAll('\n\n', '\n')}</p>
				{/if}
				<div style="scale: 0.9;">
					<Author date={meta.date} author={meta.author} />
				</div>
			</div></a>
	{/each}
</main>

<style lang="postcss">
main {
	display: flex;
	flex-flow: row wrap;
	padding: 2em;
	gap: 0;

	place-self: baseline center;
	gap: 1em;
	min-width: 50%;
	--max-width: 1440px;
	max-width: min(100%, var(--max-width));
	text-align: justify;
	height: fit-content;

	& a {
		position: relative;
		padding: 2em 1em 2em 1em;
		color: var(--text-color);
		background: rgb(128 128 128 / 0.1);
		&:hover,
		&:focus-visible {
			color: var(--color-primary);
		}

		display: grid;
		grid-template-columns: var(--img-w) auto;
		row-gap: 1em;
		column-gap: 2em;

		& :global(picture) {
		}
		& :global(> picture),
		& :global(img) {
			height: 100%;
			width: 100%;
			object-fit: cover;
			object-position: center;
		}

		& > :last-child {
			display: flex;
			flex-direction: column;
			gap: 1em;
			height: 100%;
			align-items: baseline;
			& > :last-child {
				margin-top: auto;
				width: fit-content;
			}
			& > span {
				opacity: 0.95;
				font-size: 1.2em;
				font-weight: bold;
			}
			& > p {
				opacity: 0.8;
				font-size: 0.9em;
				text-align: justify;
			}
		}
	}

	& > :first-child {
		& :global(picture),
		& :global(img) {
			grid-row: 1 / span 1;
			grid-column: 1 / span 2;
		}
		flex: 1 1 calc(60% - 1em);
	}
	& > :nth-child(2) {
		& :global(picture),
		& :global(img) {
			grid-row: 1 / span 1;
			grid-column: 1 / span 2;
		}
		flex: 1 1 calc(45% - 1em);
	}
	& > :nth-child(1n + 3) {
		flex: 1 1 calc(20% - 1em);
		& :global(picture),
		& :global(img) {
			grid-row: 1 / span 1;
			grid-column: 1 / span 2;
		}
	}
	& > :nth-child(1n + 7) {
		flex: 1 1 calc(50% - 1em);
	}
	& > * {
		flex: 1;
		min-width: min(100%, 300px);
	}
	& > :last-child {
		grid-template-columns: auto;
		&:before {
			opacity: 0;
		}
	}
	@media (orientation: portrait) {
		& > * {
			grid-template-columns: auto;
			& :global(> picture) {
				grid-row: 1 / span 1;
				grid-column: 1 / span 2;
			}
		}
	}
}
</style>
