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
	import { goto } from '$app/navigation';

	import { FileManagerService } from '../../services';
	import { fileList } from '../../store.ts';
	import { toast } from '@zerodevx/svelte-toast';
	import Button from '../../components/common/Button.svelte';
	import dayjs from 'dayjs';

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

	function routeToPathOrFile(route, obj) {
		return async (e) => {
			e.preventDefault();
			if (obj.type === 1) {
				goto(`${route}`);
				return;
			}
			let id = toast.push(`Loading...`);
			await FileManagerService.view(obj.originalName)
				.catch((err) => {
					console.log('view error: ' + err);
				})
				.finally(() => {
					toast.pop(id);
				});
		};
	}
</script>

<div class="grid grid-cols-12">
	<div class="col-span-6">
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
	</div>
	<div class="col-span-6 flex md:justify-end space-y-4">
		<Button>Upload</Button>
		<input placeholder="Search..." />
	</div>
</div>
<div class=" h-2 w-full" />

<div class="px-1 py-5 shadow rounded">
	{#each data as item}
		<a>
			<div
				class="p-2 divide-amber-300  cursor-pointer  group {item.originalName === path + `/`
					? `mb-2 bg-gray-200 rounded hover:bg-gray-300`
					: `rounded bg-gray-100 hover:bg-gray-200`}"
			>
				<div class="grid grid-cols-12">
					<div class="col-span-9">
						<div class="flex-row flex">
							<div class="icon ">
								{#if item.type === 0}
									<i class="fas fa-file" />
								{:else}
									<i class="fas fa-folder" />
								{/if}
							</div>
							<a
								on:click={routeToPathOrFile(`/f/${upperPathCheck(item.originalName)}`, item)}
								href={`/f/${upperPathCheck(item.originalName)}`}
								class="text-primary-500 group-hover:text-primary-700 w-1/2"
							>
								{#if item.name === '/' && item.type === 1}
									{item.originalName}
								{:else}
									{item.name}
								{/if}
							</a>
							<i class="fa fa-user" />
							{item.owner.split('user-')[1].toLowerCase()}
						</div>
					</div>
					<div class="col-span-3 text-right">
						<Button>
							<i class="fa fa-download" />
						</Button>
						{dayjs(item.createdAt).format(`DD/MM/YYYY HH:mm`)}
					</div>
				</div>
			</div>
		</a>
	{/each}
</div>
