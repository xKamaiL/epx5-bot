<script>
	import { FileManagerService } from '../../services';
    import {fileList} from "../../store.ts";

	FileManagerService.list()
		.then((files) => {
			fileList.set(files);
		})
		.catch((err) => {
			fileList.set([]);
		});
	export let data = [];
	export let path = '/';
</script>

<h1>Folder {path}</h1>
{JSON.stringify(data)}
