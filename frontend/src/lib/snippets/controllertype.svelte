<script lang="ts" module>
export { controllertype };
</script>

<script lang="ts">
import type { components } from '$lib/api/openapi';
import EightBitDo from '$lib/assets/steam_controller_type_svgs/8bitdo_ultimate.svg?component';
import Hori from '$lib/assets/steam_controller_type_svgs/hori.svg?component';
import PS4 from '$lib/assets/steam_controller_type_svgs/ps4.svg?component';
import PS5 from '$lib/assets/steam_controller_type_svgs/ps5.svg?component';
import Gordon from '$lib/assets/steam_controller_type_svgs/steam.svg?component';
import SwitchPro from '$lib/assets/steam_controller_type_svgs/switchpro.svg?component';
import Triton from '$lib/assets/steam_controller_type_svgs/triton.svg?component';
import XBox from '$lib/assets/steam_controller_type_svgs/xbox.svg?component';
import Icon from '@iconify/svelte';
</script>

{#snippet controllertype({ item }: { item: components['schemas']['ConfigItem' | 'ConfigDetailResponse'] })}
	{#if item.controller_type === 'controller_neptune'}
		<Icon icon="simple-icons:steamdeck" width="1.2em" />
	{:else if item.controller_type === 'controller_triton'}
		<Triton width="1.2em" />
	{:else if item.controller_type === 'controller_steamcontroller_gordon'}
		<Gordon width="1.2em" />
	{:else if item.controller_type === 'controller_ps5'}
		<PS5 width="1.2em" />
	{:else if item.controller_type === 'controller_ps5_edge'}
		<PS5 width="1.2em" />
	{:else if item.controller_type === 'controller_ps4'}
		<PS4 width="1.2em" />
	{:else if item.controller_type === 'controller_ps3'}
		<Icon icon="iconoir:gamepad" width="1.2em" />
	{:else if item.controller_type === 'controller_xbox360'}
		<XBox width="1.2em" />
	{:else if item.controller_type === 'controller_xboxone'}
		<Icon icon="fluent:xbox-controller-24-filled" width="1.2em" />
	{:else if item.controller_type === 'controller_xboxelite'}
		<Icon icon="fluent:xbox-controller-24-filled" width="1.2em" />
	{:else if item.controller_type === 'controller_switch_pro'}
		<SwitchPro width="1.2em" />
	{:else if item.controller_type === 'controller_switch2_pro'}
		<SwitchPro width="1.2em" />
	{:else if item.controller_type === 'controller_8bitdo'}
		<EightBitDo width="1.2em" />
	{:else if item.controller_type === 'controller_legion_go_s'}
		<Icon icon="game-icons:spartan-helmet" width="1.2em" />
	{:else if item.controller_type === 'controller_hori_steam'}
		<Hori width="1.2em" />
	{:else if item.controller_type === 'controller_rog_ally'}
		<Icon icon="simple-icons:republicofgamers" width="1.2em" />
	{:else if item.controller_type === 'controller_mobile_touch'}
		<Icon icon="mdi:cellphone" width="1.2em" />
	{:else if item.controller_type === 'controller_android'}
		<Icon icon="mdi:android" width="1.2em" />
	{:else}
		<Icon icon="mdi:gamepad" height="1.2em" />
	{/if}

	{item.controller_type_nice || item.controller_type || 'Generic Controller'}
{/snippet}
