<script>
	import FileDrop from 'filedrop-svelte';
	import fileSize from 'filesize';
	import Button from '../components/common/Button.svelte';

	let inputRef;
	let files;
</script>

<div class="py-10 text-center md:w-1/2 mx-auto">
	<input type="file" bind:this={inputRef} class="hidden" />
	<FileDrop
		on:filedrop={(e) => {
			files = e.detail.files;
		}}
	>
		<div
			class="bg-gray-100 transition-all hover:bg-gray-200 rounded-lg shadow cursor-pointer px-6 py-20"
		>
			<i class="fa fa-folder-plus text-5xl mb-4" />
			<h3 class="text-3xl  text-center mb-2">1. Drag and drop your pdf file or folder</h3>
		</div>
	</FileDrop>
	{#if files}
		<div class="my-5">
			<ul>
				{#each files.accepted as file}
					<li>
						<span class="text-primary-500 hover:text-primary-600 transition-all font-bold"
							>{file.name}</span
						>
						- {fileSize(file.size)}
					</li>
				{/each}
			</ul>
		</div>
	{/if}
	<div class="mt-10 mb-5">
		<h4 class="text-3xl   text-center">2. Select a <b class="text-gray-700">#</b> text-channel</h4>
	</div>
	<div>
		<ul>
			<li class="bg-gray-100 py-2 px-1 rounded transition-all cursor-pointer">
				<b class="text-lg">#</b>
				general chat
			</li>
		</ul>
	</div>
</div>
