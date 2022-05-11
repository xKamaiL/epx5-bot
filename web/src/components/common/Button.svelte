<script>
	export let color = 'primary';
	export let to = '';
	export const className = `btn btn-${color}`;
</script>

{#if to !== ''}
	<a href={to} class="{className}">
		<slot />
	</a>
{:else}
	<button class="{className}">
		<slot />
	</button>
{/if}
