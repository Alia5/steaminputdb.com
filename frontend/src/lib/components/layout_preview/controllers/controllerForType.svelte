<script lang="ts" module>
export { controllerForType };
</script>

<script lang="ts">
import Generic from './Generic.svg?component';
import Neptune from './Neptune.svg?component';
import PS4 from './PS4.svg?component';
import PS5 from './PS5.svg?component';
import Gordon from './SC.svg?component';
import SwitchPro from './SwitchPro.svg?component';
import Triton from './Triton.svg?component';
import XBox from './Xbox.svg?component';
</script>

{#snippet controllerForType(controllerType: string)}
	{#if controllerType === 'controller_neptune'}
		<Neptune />
	{:else if controllerType === 'controller_triton'}
		<Triton />
	{:else if controllerType === 'controller_steamcontroller_gordon'}
		<Gordon />
	{:else if controllerType === 'controller_ps5'}
		<PS5 />
	{:else if controllerType === 'controller_ps5_edge'}
		<PS5 />
	{:else if controllerType === 'controller_ps4'}
		<PS4 />
	{:else if controllerType === 'controller_xbox360'}
		<XBox />
	{:else if controllerType === 'controller_xboxone'}
		<XBox />
	{:else if controllerType === 'controller_xboxelite'}
		<XBox />
	{:else if controllerType === 'controller_switch_pro'}
		<SwitchPro />
	{:else if controllerType === 'controller_8bitdo'}
		<SwitchPro />
	{:else}
		<Generic />
	{/if}
{/snippet}
