<script lang="ts" module>
export { markdown };
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
</script>

<!-- enable use as component -->
{@render markdown({ content })}

{#snippet markdown(p: MarkdownProps)}
	<svelte:boundary>
		<!-- eslint-disable-next-line svelte/no-at-html-tags -->
		{@html await processMarkdown().process(p.content)}
		{#snippet pending()}
			<!-- fallback to plaintext -->
			{p.content}
		{/snippet}
		{#snippet failed()}
			<!-- fallback to plaintext -->
			{p.content}
		{/snippet}
	</svelte:boundary>
{/snippet}
