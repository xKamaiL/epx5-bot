<script context="module">
	/** @type {import('@sveltejs/kit').Load} */
	export function load({ error, status }) {
		return {
			props: {
				message: `${status}: ${error.message}`
			}
		};
	}
</script>

<script>
	import Button from '../components/common/Button.svelte';

	export let message;
</script>

<h1 class="leading-10 text-xl">เกิดข้อผิดพลาด {message}</h1>

<Button to="/">Home</Button>
