<script context="module">
	/** @type {import('./[...path]').Load} */
	export async function load({ params, fetch, session, stuff }) {
		console.log(params.path);
		return {
			status: 200,
			props: {
				path: params.path
			},
			error: new Error(`asdas`)
		};
	}
</script>

<script>
	export let path;
</script>

<svelte:head>
	<title>About</title>
</svelte:head>

<h1>About this site</h1>
<p>TODO...</p>

<p>Path: {path}</p>
