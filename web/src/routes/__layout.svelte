<script>
	import '../app.css';
	import Navbar from '../components/Navbar.svelte';
</script>

<svelte:head />
<svelte:body class="bg-darken-200" />

<div class="">
	<Navbar />
	<div class="bg-page md:mt-[105px] mt-[15px] min-h-screen relative">
		<div class="mx-auto container relative">
			<div>
				<slot />
			</div>
		</div>
	</div>
</div>
