<script>
	import '../app.css';
	import Navbar from '../components/Navbar.svelte';
	import { SvelteToast } from '@zerodevx/svelte-toast'
	// Optionally set default options here
	const options = {
	}
</script>

<svelte:head />
<svelte:body class="bg-darken-200" />

<div class="">
	<Navbar />
	<div class="bg-page md:mt-[105px] mt-[15px] relative">
		<div class="mx-auto container relative">
			<div>
				<slot />
			</div>
			<footer class="text-darken mt-10 divide-amber-300 divide-y">
				<div class="font-light">Copyright (c) 2020</div>
			</footer>
		</div>
	</div>
</div>
<SvelteToast {options} />

