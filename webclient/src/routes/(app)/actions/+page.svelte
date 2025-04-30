<script>
	import BaseView from '$lib/BaseView.svelte';
	import ActionNode from '$lib/ActionNode.svelte';
	import { getActions } from '$lib/service/action';
	import { onMount } from 'svelte';
	import LL from '../../../i18n/i18n-svelte';

	let actionsTree = {};
	let isLoading = true;
	let error = null;

	onMount(async () => {
		try {
			const [actionDetails, err] = await getActions();
			if (err) {
				throw err;
			}

			actionsTree = actionDetails ? buildTree(actionDetails) : {};

		} catch (err) {
			error = err.message || 'Failed to load actions';
		} finally {
			isLoading = false;
		}
	});

	function buildTree(actionDetails) {
		const tree = {
			children: {},
			childrenOrder: []
		};

		actionDetails.forEach(detail => {
			const parts = detail.path.split('/');
			let currentLevel = tree;

			parts.forEach((part, index) => {
				if (index === parts.length - 1) {
					currentLevel.children[part] = {
						_isAction: true,
						details: {
							args: detail.args,
							description: detail.description
						}
					};
					if (!currentLevel.childrenOrder.includes(part)) {
						currentLevel.childrenOrder.push(part);
					}
				} else {
					if (!currentLevel.children[part]) {
						currentLevel.children[part] = {
							children: {},
							childrenOrder: []
						};
						currentLevel.childrenOrder.push(part);
					}
					currentLevel = currentLevel.children[part];
				}
			});
		});

		return tree;
	}
</script>

<BaseView>
	<div class="actions-page">
		<h1 class="page-title">{$LL.remote_actions.title()}</h1>

		{#if isLoading}
			<div class="status-message">{$LL.remote_actions.loading()}</div>
		{:else if error}
			<div class="status-message error">
				<strong>{$LL.remote_actions.error()}:</strong> {error}
			</div>
		{:else if Object.keys(actionsTree.children || {}).length === 0}
			<div class="status-message">{$LL.remote_actions.no_actions()}</div>
		{:else}
			{#each actionsTree.childrenOrder as name}
				<ActionNode {name} node={actionsTree.children[name]} currentPath="" isOpen={actionsTree.childrenOrder.length <= 2} />
			{/each}
		{/if}
	</div>
</BaseView>

<style>
	.actions-page {
		margin: 3rem 0;
		padding: 0 1rem;
	}

	.page-title {
		text-align: center;
		margin-bottom: 1rem;
		font-size: 1.8rem;
		font-weight: 600;
	}

	.status-message {
		text-align: center;
		padding: 1rem 0;
		color: var(--text-secondary-color);
		font-size: 0.95rem;
	}
	
	.status-message.error {
		color: var(--error-color);
		font-weight: 500;
	}
	
	.status-message.error strong {
		font-weight: 600;
	}

</style>