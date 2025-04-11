import { tick } from 'svelte';

export interface StorageLike {
	getItem(key: string): string | null;
	setItem(key: string, value: string): void;
}

export class WebStorage<T> {
	#key: string;
	#version = $state(0);
	#listeners = 0;
	#value: T | undefined;
	#storage: StorageLike;

	#handler = (e: StorageEvent) => {
		if (e.storageArea !== this.#storage) return;
		if (e.key !== this.#key) return;

		this.#version += 1;
	};

	constructor(key: string, initial?: T, storage?: StorageLike) {
		this.#key = key;
		this.#value = initial;
		this.#storage = storage || (typeof localStorage !== 'undefined' ? localStorage : undefined);

		if (this.#storage) {
			if (this.#storage.getItem(key) === null && initial !== undefined) {
				this.#storage.setItem(key, JSON.stringify(initial));
			}
		}
	}

	get current() {
		this.#version;

		const root =
			this.#storage
				? JSON.parse(this.#storage.getItem(this.#key) as any)
				: this.#value;

		const proxies = new WeakMap();

		const proxy = (value: unknown) => {
			if (typeof value !== 'object' || value === null) {
				return value;
			}

			let p = proxies.get(value);

			if (!p) {
				p = new Proxy(value, {
					get: (target, property) => {
						this.#version;
						return proxy(Reflect.get(target, property));
					},
					set: (target, property, value) => {
						this.#version += 1;
						Reflect.set(target, property, value);

						if (this.#storage) {
							this.#storage.setItem(this.#key, JSON.stringify(root));
						}

						return true;
					}
				});

				proxies.set(value, p);
			}

			return p;
		};

		if ($effect.tracking()) {
			$effect(() => {
				if (this.#listeners === 0 && typeof window !== 'undefined') {
					window.addEventListener('storage', this.#handler);
				}

				this.#listeners += 1;

				return () => {
					tick().then(() => {
						this.#listeners -= 1;
						if (this.#listeners === 0 && typeof window !== 'undefined') {
							window.removeEventListener('storage', this.#handler);
						}
					});
				};
			});
		}

		return proxy(root);
	}

	set current(value) {
		if (this.#storage) {
			this.#storage.setItem(this.#key, JSON.stringify(value));
		}

		this.#version += 1;
	}
}

export class LocalStorage<T> extends WebStorage<T> {
	constructor(key: string, initial?: T) {
		super(key, initial, typeof localStorage !== 'undefined' ? localStorage : undefined);
	}
}

export class SessionStorage<T> extends WebStorage<T> {
	constructor(key: string, initial?: T) {
		super(key, initial, typeof sessionStorage !== 'undefined' ? sessionStorage : undefined);
	}
}