import axios, { type AxiosRequestConfig, type AxiosResponse } from 'axios';

// const API_URL = `https://epx5-api-service-666023067909652511.rcf2.deploys.app`;
const API_URL = `http://localhost:8080`;
const responseBody = (res: AxiosResponse) => res.data;

const api = axios.create({
	baseURL: API_URL
});

const token: null | string = null;

const authHeader = (config: AxiosRequestConfig) => {
	if (token) {
		config.headers = {
			Authorization: `${token}`
		};
	}
	return config;
};
api.interceptors.request.use(authHeader);
api.interceptors.response.use(
	(e) => e,
	(error) => {
		return Promise.reject(error);
	}
);

export const FileManagerService = {
	list: (prefix?: string) => api.get(`/api/file?prefix=${prefix || ''}`).then(responseBody),
	createFolder: (name: string) =>
		api
			.post(`/api/file/folder`, {
				name
			})
			.then(responseBody),
	upload: (path: string, file: File) => {
		const d = new FormData();
		d.append('path', path);
		d.append('file', file);
		return api.post(`/api/file/upload`, d).then(responseBody);
	}
};
