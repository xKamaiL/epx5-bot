<script context="module">
	/** @type {import('./[...path]').Load} */
	export async function load({ params, fetch, session, stuff }) {
		console.log(params.path);
		return {
			status: 200,
			props: {
				path: params.path
			}
		};
	}
</script>

<script>
	import { FileManagerService } from '../../services';
	import { fileList } from '../../store.ts';

	export let path;
	export let data = [];
	// import { onMount } from 'svelte';

	// onMount(async () => {
	$: FileManagerService.list(path)
		.then((files) => {
			fileList.set(files);
		})
		.catch((err) => {
			return [];
		});
	// });

	fileList.subscribe((value) => {
		data = value;
	});

	// add prefix with slash
	function cleanPath(path) {
		console.log('path', path);
		if (path === '') return '';
		return !path.startsWith('/') ? '/' + path : path;
	}

	function upperPathCheck(originalName) {
		if (originalName === path + '/') {
			const spt = `/${originalName}`.split('/');
			return `/${spt[spt.length - 3]}`;
		}
		return originalName;
	}
</script>

<h1 class="text-4xl divide-gray-500 ">
	<a href="/f/" class="hover:text-primary-500 font-bold">Folder</a>
	<span class="">
		{#each path.split(`/`) as p, i}
			<span class="text-gray-500">/</span>
			<a
				href="/f{cleanPath(path.split(`/`).splice(0, i).join(`/`))}/{p}"
				class="hover:text-primary-500 font-bold"
			>
				{p}
			</a>
		{/each}
	</span>
</h1>
<div class=" h-2 w-full" />

<div class="px-1 py-5 shadow rounded">
	{#each data as item}
		<a href={`/f/${upperPathCheck(item.originalName)}`} class="">
			<div
				class="p-2 divide-amber-300  cursor-pointer  group {item.originalName === path + `/`
					? `mb-2 bg-gray-200 rounded hover:bg-gray-300`
					: `rounded bg-gray-100 hover:bg-gray-200`}"
			>
				<div class="grid grid-cols-12">
					<div class="col-span-9">
						<div class="flex-row flex">
							<div class="icon mr-2">
								{#if item.type === 0}
									<i class="fas fa-file" />
								{:else}
									<i class="fas fa-folder" />
								{/if}
							</div>
							<span class="text-primary-500 group-hover:text-primary-700">
								{#if item.name === '/' && item.type === 1}
									{item.originalName}
								{:else}
									{item.name}
								{/if}
							</span>
						</div>
					</div>
					<div class="col-span-3">
						{item.createdAt}
					</div>
				</div>
			</div>
		</a>
	{/each}
</div>
