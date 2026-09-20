<script lang="ts" module>
export { configRating };
const calculateScoreColor = (score: number, up: number, down: number) => {
	if (up === down && up === 0) {
		return 'currentColor';
	}
	if (!score) {
		return 'currentColor';
	}
	if (down > up) {
		return 'hsl(0deg 100% 60%)';
	}
	if (score > 0.8) {
		return 'hsl(125deg 100% 50%)';
	}
	if (score > 0.7) {
		return 'hsl(80deg 100% 50%)';
	}
	if (score >= 0.51) {
		return 'hsl(60deg 100% 50%)';
	}
	if (score > 0.4) {
		return 'hsl(30deg 100% 50%)';
	}
	return 'hsl(0deg 100% 60%)';
};
</script>

<script lang="ts">
import type { components } from '$lib/api/openapi';
import { tooltip } from '$lib/attachments/tooltip.svelte';
import Icon from '@iconify/svelte';
</script>

{#snippet configRating({ item }: { item: components['schemas']['ConfigItem' | 'ConfigDetailResponse'] })}
	{#if item.votes}
		{@const scoreColor = calculateScoreColor(
			item.votes.score ?? 0,
			item.votes.up ?? 0,
			item.votes.down ?? 0
		)}
		<div
			class="rating"
			style="--rating-color: {scoreColor};"
			{@attach tooltip({
				snippet: tooltipContent,
				snippetInDefaultBackground: true,
				outDelay: 200,
				arrow: true,
				arrowFollowCursor: true
			})}
		>
			<div>
				<span>
					{#if (item.votes?.down || 0) > (item.votes?.up || 0)}
						😣
					{:else if (item.votes?.score || 0) > 0.8}
						😍
					{:else if (item.votes?.score || 0) > 0.7}
						🤩
					{:else if (item.votes?.score || 0) >= 0.51}
						😎
					{:else if (item.votes?.score || 0) > 0.4}
						🙁
					{:else}
						🤔
					{/if}
				</span>
				<span>
					{item.votes.up == item.votes.down && item.votes.up === 0
						? '???'
						: (item.votes?.score ?? 0).toLocaleString(undefined, {
								style: 'percent',
								minimumFractionDigits: 0,
								maximumFractionDigits: 1
							})}
				</span>
			</div>
			<div>
				<span
					>{item.votes.up?.toLocaleString(undefined, {
						notation: 'compact',
						minimumFractionDigits: 0,
						maximumFractionDigits: 1
					})}
					<Icon icon="mdi:thumb-up" />
				</span><span>/</span><span>
					{item.votes.down?.toLocaleString(undefined, {
						notation: 'compact',
						minimumFractionDigits: 0,
						maximumFractionDigits: 1
					})}
					<Icon icon="mdi:thumb-down" /></span
				>
			</div>
			<i>(Rating from Steam)</i>
		</div>
	{/if}
{/snippet}

{#snippet tooltipContent()}
	<div style="display: grid; place-items: center; max-width: 80svw;">
		<p style="text-align: center;">The ranking system is provided by Steam</p>
		<p style="text-align: center;">I do not know and can only guess on how it rates</p>
	</div>
{/snippet}

<style lang="postcss">
:global(.rating) {
	display: grid;
	place-items: center;
	& span {
		filter: drop-shadow(1px 1px 1px black);
		color: var(--rating-color);
		font-size: 1.1em;
		font-weight: 500;
	}
	& i {
		color: var(--text-color-dark);
		opacity: 0.8;
		font-size: 0.8em;
		filter: drop-shadow(1px 1px 1px black) drop-shadow(0px 0px 2px black);
	}
	& > :first-child {
		& > :first-child {
			font-size: 1.8em;
			filter: drop-shadow(1px 2px 3px black);
			transform: translate(0, 0.5em);
		}
		& > :last-child {
			font-size: 1.7em;
			filter: drop-shadow(1px 1px 1px black) drop-shadow(1px 1px 2px var(--shadow-color));
		}
	}
	& > :nth-child(2) {
		display: flex;
		gap: 0.5ch;
		& > :first-child {
			color: hsl(108, 100%, 50%);
		}
		& > :nth-child(2) {
			font-size: 1.1em;
			color: var(--text-color-dark);
		}
		& > :last-child {
			color: hsl(0, 100%, 60%);
		}
	}
	height: 100%;
}
</style>
