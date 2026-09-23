<script lang="ts">
import type { LayoutFile } from './config.types';
import { niceInputMap } from './const';
import { niceInputName } from './helper';

const {
	name,
	b,
	parsedVdf
}: {
	name: string;
	b?: string;
	parsedVdf?: LayoutFile;
} = $props();

const localizeGameAction = (input: string): string | undefined => {
	const parts = input.split(' ');
	const actionName = parts[2];
	if (!actionName) return undefined;

	const english = parsedVdf?.controller_mappings?.localization?.english;
	if (english) {
		const localized = english[actionName];
		if (typeof localized === 'string' && localized) {
			return localized;
		}
	}
	return undefined;
};

const presetTitle = (presetId: string): string | undefined => {
	if (!parsedVdf) {
		return;
	}
	const presets = Array.isArray(parsedVdf.controller_mappings.preset)
		? parsedVdf.controller_mappings.preset
		: [parsedVdf.controller_mappings.preset];
	const preset = presets.find((p) => p.id == presetId);
	if (!preset) {
		return;
	}
	return (
		parsedVdf.controller_mappings.action_layers?.[preset.name]?.title ??
		parsedVdf.controller_mappings.actions?.[preset.name]?.title ??
		undefined
	);
};
</script>

{#if b}
	{@const [input, description] = b.split(',').map((s: string) => s.trim())}
	{@const isEmpty = input?.includes('empty_binding')}
	{@const isModeShiftBinding = input?.startsWith('mode_shift ')}
	{@const isHoldLayerBinding = input?.includes('hold_layer')}
	{#if isEmpty}
		<span>
			{#if description}
				{description}
				(Cleared from Parent)
			{:else}
				--
			{/if}
		</span>
	{:else if isModeShiftBinding}
		{@const parts = input?.split(' ') ?? []}
		{@const target = (niceInputMap as Record<string, string>)[parts[1] ?? ''] ?? parts[1]}
		<span>Mode Shift ({target})</span>
	{:else if isHoldLayerBinding}
		{@const parts = input?.split(' ') ?? []}
		{@const layerNum = parts[2]}
		{@const layerTitle = layerNum ? presetTitle(String(Number(layerNum) - 1)) : undefined}
		<span
			>{niceInputName(input)}
			{#if layerTitle}
				({layerTitle})
			{/if}
		</span>
	{:else}
		{@const isGameAction = input?.toLowerCase().includes('game_action')}
		{@const gameActionLabel = isGameAction ? localizeGameAction(input!) : undefined}
		<span>
			{#if isGameAction}
				{gameActionLabel ?? description ?? input?.split(' ')[2] ?? niceInputName(input)}
			{:else if description}
				{description} ({niceInputName(input)})
			{:else if name == 'Full_Press'}
				{niceInputName(input)}
			{:else}
				({name.replace(/_/g, ' ')}) {niceInputName(input)}
			{/if}
		</span>
	{/if}
{/if}
