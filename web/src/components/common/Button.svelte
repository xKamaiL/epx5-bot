<script>
	export let color = 'primary';
	export let to = '';
	export let size = 'md';
	export const className = `transition-all btn btn-${color}`;
</script>

{#if to !== ''}
	<a href={to} class={className}>
		<slot />
	</a>
{:else}
	<button class={className}>
		<slot />
	</button>
{/if}
