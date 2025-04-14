import axios from 'axios';
import { PUBLIC_API_URL } from "$env/static/public";
import { checkSession } from './service';
import { userStore } from '$lib/store';

export async function getActions() {
	try {
		const response = await axios.get(`${PUBLIC_API_URL}/actions`, {
			headers: {
				Authorization: `Bearer ${userStore.current.token}`,
			},
		});

		return [response.data, null];
	} catch (error) {
		checkSession(error.response);
		console.error('Error fetching actions:', error);

		return [null, error];
	}
}

export async function invokeAction(path, args) {
	const stringArgs = args.map(arg => {
		if (typeof arg === 'boolean') {
			return arg ? 'true' : 'false';
		}

		if (typeof arg === 'object' && arg !== null) {
			try {
				return JSON.stringify(arg);
			} catch (e) {
				console.warn("Could not JSON stringify argument:", arg, e);
				return String(arg);
			}
		}
		return String(arg);
	});

	try {
		const response = await axios.post(`${PUBLIC_API_URL}/actions/invoke`, { path, args: stringArgs }, {
			headers: {
				Authorization: `Bearer ${userStore.current.token}`,
			},
		});

		return [response.data, null];
	} catch (error) {
		checkSession(error.response);
		console.error(`Error invoking action '${path}':`, error);

		return [null, error];
	}
} 