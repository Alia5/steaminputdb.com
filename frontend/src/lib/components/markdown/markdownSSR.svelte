<script lang="ts" module>
export type MarkdownProps = {
	content: string;
};
</script>

<script lang="ts">
import { processMarkdown } from './mdprocess';

const {
	content
}: {
	content: string;
} = $props();

let html = $derived(
	await processMarkdown()
		.process(content)
		.catch(() => undefined)
);
</script>

{#if html !== undefined}
	<!-- eslint-disable-next-line svelte/no-at-html-tags -->
	{@html html}
{:else}
	{content}
{/if}
