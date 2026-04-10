<script lang="ts" module>
export { sectionInfo };
</script>

<script lang="ts">
import { resolve } from '$app/paths';
import type { components } from '$lib/api/openapi';
import { selectAllHandler } from '$lib/attachments/selectAllHandler.svelte';
import { configRating } from '$lib/snippets/configRating.svelte';
import { configurationFeatureList } from '$lib/snippets/configurationfeaturelist.svelte';
import { controllertype } from '$lib/snippets/controllertype.svelte';
import { format, formatDistance, formatDistanceToNow, formatDuration, intervalToDuration } from 'date-fns';
</script>

{#snippet sectionInfo({
	fileInfo,
	appInfo,
	creatorInfo
}: {
	fileInfo: components['schemas']['ConfigDetailResponse'];
	appInfo?: components['schemas']['AppItem'];

	creatorInfo?: components['schemas']['PlayerInfo'];
})}
	<section id="config-info">
		<dl class="card glass">
			<dt>Controller</dt>
			<dd>
				{@render controllertype({ item: fileInfo })}
			</dd>
			{#if creatorInfo}
				<dt>Author</dt>
				<dd>
					<a href={resolve(`/user/${fileInfo.creator_id}`)}>{creatorInfo.personaname}</a>
				</dd>
			{:else}
				<dt>Author ID</dt>
				<dd>
					<a href={resolve(`/user/${fileInfo.creator_id}`)}>{fileInfo.creator_id}</a>
				</dd>
			{/if}
			<dt>Created</dt>
			<dd>
				{format(new Date(fileInfo.time_created), 'PPpp')}
				<i>({formatDistanceToNow(new Date(fileInfo.time_created))} ago)</i>
			</dd>
			{#if fileInfo.file_url}
				<dt>File URL</dt>
				<dd>
					<code
						{@attach selectAllHandler(
							`outline: 1px solid transparent;
                                        background: rgb(128 128 128 / 0.10);`
						)}>{fileInfo.file_url}</code>
				</dd>
			{/if}
			{#if fileInfo.file_name}
				<dt>File Name</dt>
				<dd>
					<code
						{@attach selectAllHandler(
							`outline: 1px solid transparent;
                                        background: rgb(128 128 128 / 0.10);`
						)}>{fileInfo.file_name}</code>
				</dd>
			{/if}
			{#if fileInfo.file_size}
				<dt>File Size</dt>
				<dd>{(fileInfo.file_size / 1000).toFixed(0)} kB</dd>
			{/if}

			{#if fileInfo.playtime_seconds}
				{@const duration = intervalToDuration({
					start: new Date(0),
					end: new Date(fileInfo.playtime_seconds * 1000)
				})}
				<dt>Playtime</dt>
				<dd>
					{formatDuration(duration, {
						zero: false,
						format: ['years', 'months', 'days', 'hours']
					})}
					{#if (duration.hours || 0) > 0}
						and {formatDuration(duration, { format: ['minutes'] })}
					{/if}
					<i>(all users combined - last 14 days)</i>
				</dd>
			{/if}
			{#if fileInfo.playtime_sessions}
				<dt>Sessions</dt>
				<dd>
					{fileInfo.playtime_sessions.toLocaleString()}
					<i>(all users combined - last 14 days)</i>
				</dd>
			{/if}
			{#if fileInfo.lifetime_playtime_seconds}
				{@const duration = intervalToDuration({
					start: new Date(0),
					end: new Date(fileInfo.lifetime_playtime_seconds * 1000)
				})}
				<dt>Total</dt>
				<dd>
					<dl>
						<dt>Playtime</dt>
						<dd>
							{formatDuration(duration, {
								zero: false,
								format: ['years', 'months', 'days', 'hours']
							})}
							{#if (duration.hours || 0) > 0}
								and {formatDuration(duration, { format: ['minutes'] })}
							{/if}
						</dd>
						{#if fileInfo.lifetime_playtime_sessions}
							<dt>Sessions</dt>
							<dd>
								{fileInfo.lifetime_playtime_sessions.toLocaleString()}
							</dd>
						{/if}
					</dl>
				</dd>
			{/if}

			<dt>{appInfo ? 'Game' : 'Shortcut'}</dt>
			<dd>
				{#if appInfo}
					<a href={resolve(`/app/${appInfo.app_id}`)}> {appInfo.name} <i>({appInfo.app_id})</i></a>
				{:else}
					{fileInfo.app_id_string}
				{/if}
			</dd>
			{#if fileInfo.tags}
				<dt>Features</dt>
				<dd class="featurelist">
					{@render configurationFeatureList({ fileInfo })}
				</dd>
			{/if}
		</dl>
		<aside>
			<section>
				{#if fileInfo.votes}
					<div>
						{@render configRating({ item: fileInfo })}
					</div>
				{/if}
				{#if fileInfo.playtime_seconds || fileInfo.lifetime_playtime_seconds}
					<div class="playtime">
						{#if fileInfo.playtime_seconds}
							<span
								>{formatDistance(
									new Date(fileInfo.playtime_seconds * 1000),
									new Date(0)
								)}</span>
							<span>combined playtime in</span>
							<span>{fileInfo.playtime_sessions ?? 0}</span> <span>sessions</span>
							<i>(all users - last 14 days)</i>
						{:else}
							<span
								>{formatDistance(
									new Date((fileInfo.lifetime_playtime_seconds ?? 0) * 1000),
									new Date(0)
								)}</span>
							<span>combined playtime in</span>
							<span>{fileInfo.lifetime_playtime_sessions ?? 0}</span> <span>sessions</span>
							<i>(all users - since upload)</i>
						{/if}
					</div>
				{/if}
			</section>
			<p class="card glass">
				<span>Creator Description</span>
				{fileInfo.description?.replace(/\s\s/g, '\n')}
			</p>
		</aside>
	</section>
{/snippet}

<style lang="postcss">
:global(#config-info) {
	display: flex;
	flex-flow: row wrap-reverse;
	overflow: hidden;

	--gap: 1em;
	--info-min-width: 58ch;
	width: 100%;

	gap: 1em;
	padding: 0 1em;

	& > :first-child {
		padding: 1em;
		overflow-x: hidden;
		position: relative;
		backdrop-filter: blur(8px);

		& > dt {
			font-size: 1.2em;
		}

		flex-basis: calc(66% - var(--gap));
		min-width: min(100%, var(--info-min-width));
	}
	& > :last-child {
		display: grid;
		height: fit-content;
		align-self: flex-end;
		width: 100%;
		flex-shrink: 1;
		gap: 1em;

		--width: 20ch;

		flex-basis: calc(34%);
		min-width: min(100%, var(--width));

		& > :first-child {
			display: grid;
			place-items: center;
			margin: auto;
			gap: 1em;
			width: 100%;
			grid-template-columns: repeat(auto-fit, minmax(calc(var(--width) -1em), auto));
			:global(> *) {
				width: 100%;
				padding: 1em;
				position: relative;
				isolation: isolate;
				border-radius: var(--border-radius);
				box-shadow: var(--card-shadow);
				backdrop-filter: blur(8px);

				&::before {
					content: '';
					position: absolute;
					inset: 0;
					background: var(--card-glass);
					opacity: 0.5;
					border-radius: var(--border-radius);
					z-index: -1;
				}
				&::after {
					content: '';
					position: absolute;
					inset: 0;
					border-radius: inherit;
					border: 1px solid transparent;
					background: var(--card-border-pseudo-gradient) border-box;
					mask:
						linear-gradient(black, black) border-box,
						linear-gradient(black, black) padding-box;
					mask-composite: subtract;
					z-index: -1;
				}
			}
		}
	}

	@media (max-width: 200ch) {
		& > :first-child {
			flex-basis: 0;
			flex: 1;
			flex-grow: 1;
		}
		& > :last-child {
			flex-basis: 100%;
		}
	}

	dl {
		display: grid;
		place-items: center;
		grid-template-columns: min-content auto;
		grid-column-gap: 1em;
		grid-row-gap: 0.5em;

		& > :nth-child(2) {
			font-weight: bold;
			:global(svg) {
				translate: 0 0.2em;
			}
		}
		& > :nth-child(2n-1) {
			justify-self: end;
			white-space: nowrap;
		}
		& > :nth-child(2n) {
			justify-self: start;
			white-space: no-wrap;
		}
	}

	dd {
		max-width: 100%;
		overflow: auto;
		overflow-x: hidden;
		& a {
			font-weight: bold;
		}
	}

	dt {
		font-weight: bold;
	}

	dd > dl {
		grid-column-gap: 0.5em;
		& > dt {
			color: var(--text-muted);
		}
	}

	dl {
		& i {
			color: var(--text-muted);
			white-space: nowrap;
		}
	}

	aside {
		& > p {
			overflow: auto;
			overflow-x: hidden;
			contain: inline-size;

			padding: 1em;
			& > span {
				font-weight: bold;
				display: block;
				margin-bottom: 0.5em;
				font-size: 1.2em;
			}
		}

		& > section {
			& :global(> :first-child) {
				height: 100%;
			}
			& > .playtime {
				display: grid;
				place-items: center;
				height: 100%;

				& span {
					font-size: 1.2em;
					font-weight: 500;
					color: var(--text-color-dark);
					filter: drop-shadow(1px 1px 1px black);
				}

				& > :nth-child(2n) {
					font-size: 1em;
					font-weight: normal;
					text-align: center;
				}

				& i {
					color: var(--text-color-dark);
					opacity: 0.8;
					font-size: 0.8em;
					filter: drop-shadow(1px 1px 1px black) drop-shadow(0px 0px 2px black);
				}
				& :global(> div) {
					& > :first-child {
						font-size: 1.8em;
						filter: drop-shadow(1px 2px 3px black);
						transform: translate(0, 0.5em);
					}
					& > :last-child {
						font-size: 1.6em;
					}
				}
			}
		}
	}

	code {
		display: block;
	}

	.featurelist {
		display: flex;
		flex-wrap: wrap;
		gap: 0.5ch;
		overflow: clip;
		overflow-clip-margin: 1em;
	}
}
</style>
