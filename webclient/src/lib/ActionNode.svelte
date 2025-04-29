<script>
	let { name = '', node = {}, currentPath = '' } = $props();
	import { invokeAction } from '$lib/service/action';
	import Self from './ActionNode.svelte';

	let result = $state(null);
	let error = $state(null);
	let loading = $state(false);

	let isAction = $state(false);
	$effect(() => {
		isAction = node && node._isAction === true;
	});

	let actionDetails = $state({ args: [], description: '' });
	$effect(() => {
		actionDetails = isAction ? node.details : { args: [], description: '' };
	});

	let nodePath = $state('');
	$effect(() => {
		nodePath = currentPath ? `${currentPath}/${name}` : name;
	});

	let inputValues = $state({});

	function getInputType(goType) {
		switch (goType) {
			case 'int': case 'int8': case 'int16': case 'int32': case 'int64':
			case 'uint': case 'uint8': case 'uint16': case 'uint32': case 'uint64':
			case 'float32': case 'float64': return 'number';
			case 'bool': return 'checkbox';
			case 'string': default: return 'text';
		}
	}

	function parseArgument(value, goType) {
		switch (goType) {
			case 'int': case 'int8': case 'int16': case 'int32': case 'int64':
			case 'uint': case 'uint8': case 'uint16': case 'uint32': case 'uint64':
				const intVal = parseInt(value, 10);
				return isNaN(intVal) ? 0 : intVal;
			case 'float32': case 'float64':
				const floatVal = parseFloat(value);
				return isNaN(floatVal) ? 0.0 : floatVal;
			case 'bool': return !!value;
			case 'string': default: return String(value);
		}
	}

	async function handleInvoke() {
		if (!isAction) return;
		loading = true;
		result = null;
		error = null;
		const argsToSend = actionDetails.args.map((argType, i) => {
			const rawValue = inputValues[i];
			return parseArgument(rawValue, argType);
		});
		try {
			const [res, invokeError] = await invokeAction(nodePath, argsToSend);
			if (invokeError) throw invokeError;
			if (res && res.error) error = res.error;
			else if (res) result = res.result;
			else error = 'Received unexpected null response';
		} catch (err) {
			error = err.response?.data || err.message || `Failed to invoke action '${nodePath}'`;
		} finally {
			loading = false;
		}
	}
</script>

{#if isAction}
	<!-- Action -->
	<div class="action-node">
		<div class="action-header">
			<div class="action-icon">♦</div>
			<span class="action-name">{name}</span>
		</div>
		{#if actionDetails.description}
			<p class="action-description">{actionDetails.description}</p>
		{/if}
		<div class="action-body">
			<div class="input-row">
				{#if actionDetails.args && actionDetails.args.length > 0}
					<div class="args-container">
						{#each actionDetails.args as argType, i}
							<label class="arg-label">
								<span class="arg-type">{argType}</span>
								{#if getInputType(argType) === 'checkbox'}
									<input type="checkbox" bind:checked={inputValues[i]} class="checkbox arg-input" disabled={loading} />
								{:else}
									<input
										type={getInputType(argType)}
										bind:value={inputValues[i]}
										placeholder={`Arg ${i + 1}`}
										class="input input-sm arg-input"
										disabled={loading}
										step={getInputType(argType) === 'number' ? 'any' : undefined}
									/>
								{/if}
							</label>
						{/each}
					</div>
				{:else}
					<p class="no-args-message">(No arguments)</p>
				{/if}

				<button onclick={handleInvoke} disabled={loading} class="btn btn-sm invoke-button variant-ghost-primary">
					{#if loading}...
					{:else}Invoke{/if}
				</button>
			</div>

			{#if error}
				<div class="output-box error-display">
					<pre>{error}</pre>
				</div>
			{/if}
			{#if result !== null}
				<div class="output-box result-box">
					<pre>{JSON.stringify(result, null, 2)}</pre>
				</div>
			{/if}
		</div>
	</div>
{:else}
	<!-- Directory -->
	<details class="directory-node">
		<summary class="directory-summary list-item">
			<span class="directory-icon"></span>
			<span class="directory-name">{name}</span>
		</summary>
		<div class="directory-content">
			{#each node.childrenOrder as childName}
				<Self name={childName} node={node.children[childName]} currentPath={nodePath} />
			{/each}
		</div>
	</details>
{/if}

<style>

	.action-node,
	.directory-node {
		padding: 0.6rem 0;
		border-top: 1px solid var(--border-color);
	}

	:global(.actions-page > .action-node:first-child),
	:global(.actions-page > .directory-node:first-child) {
		border-top: none;
		padding-top: 0.1rem; 
	}

	.action-header {
		display: flex;
		align-items: baseline;
		gap: 0.4rem;
		margin-bottom: 0.4rem;
	}
	.action-icon {
		font-size: 0.8rem;
		color: var(--text-secondary-color);
		line-height: 1.2;
	}
	.action-name {
		font-weight: 500;
		font-size: 1rem;
		color: var(--text-color);
	}
	.action-body {
		padding-left: 1.2rem;
		display: flex;
		flex-direction: column;
		gap: 0.6rem;
	}

	.input-row {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		gap: 0.75rem;
	}

	.args-container {
		display: flex;
		flex-wrap: wrap;
		gap: 0.75rem;
		flex-grow: 1;
	}
	.arg-label {
		display: flex;
		align-items: center;
		gap: 0.35rem;
		background-color: var(--input-background);
		padding: 0.2rem 0.4rem;
		border-radius: 4px;
		border: 1px solid var(--border-color);
		font-size: 0.85rem;
	}
	.arg-type {
		font-size: 0.75rem;
		color: var(--text-secondary-color);
	}
	.arg-input {
		flex-shrink: 0;
	}

	.arg-label > .input.arg-input {
		min-width: 100px;
		max-width: 180px;
		flex-grow: 1;
		background-color: var(--input-background);
	}
	.arg-label > .checkbox.arg-input {
		margin: 0;
	}

	.no-args-message {
		font-size: 0.8rem;
		color: var(--text-secondary-color);
		font-style: italic;
		flex-grow: 1;
	}

	.invoke-button {
		flex-shrink: 0;
	}

	.output-box {
		padding: 0.4rem 0.6rem;
		border-radius: 4px;
		font-size: 0.8rem;
		border: 1px solid;
	}
	.output-box pre {
		white-space: pre-wrap;
		word-break: break-all;
		font-family: monospace;
	}
	.result-box {
		background-color: var(--card-background);
		border-color: var(--border-color);
		color: var(--text-color);
	}
	.error-display {
		border-color: var(--error-color);
		background-color: var(--row-error-bg);
		color: var(--error-color);
	}

	.directory-node > summary {
		list-style: none;
	}
	.directory-node > summary::-webkit-details-marker {
		display: none;
	}
	.directory-summary {
		padding: 0.1rem 0;
		cursor: pointer;
		display: flex;
		align-items: center;
		gap: 0.4rem;
	}
	.directory-summary:hover .directory-name {
		color: var(--primary-color);
	}

	.directory-icon {
		display: inline-block;
		width: 1em;
		height: 1em;
		font-size: 0.7rem;
		line-height: 1em;
		text-align: center;
		transition: transform 0.2s ease-in-out;
		color: var(--text-secondary-color);
	}
	.directory-icon::before {
		content: '➤';
	}
	.directory-node[open] > .directory-summary .directory-icon {
		transform: rotate(90deg);
	}

	.directory-name {
		font-weight: 500;
		font-size: 1.4rem;
		transition: color 0.15s ease-in-out;
		color: var(--text-color);
		user-select: none;
	}
	.directory-content {
		margin-left: 0.5rem; 
		padding-left: 1.5rem;
		border-left: 1px solid var(--border-color);
	}

	.action-description {
		font-size: 0.85rem;
		color: var(--text-secondary-color);
		margin-left: 1.2rem;
		margin-top: -0.2rem;
		margin-bottom: 0.4rem;
	}

	button {
    display: inline-block;
    cursor: pointer;
    background-color: var(--primary-color);
    color: white;
    border: none;
    min-width: 80px;
    transition: background-color 0.2s, transform 0.1s;
    font-weight: 500;
    padding: 0.6rem 0.8rem;
  }

  button:hover {
    background-color: var(--primary-hover-color);
  }

  button:active {
    transform: translateY(1px);
  }

  input[type="text"],
  input[type="number"] {
    padding: 0.6rem 0.8rem;
    border: 1px solid var(--input-border);
    border-radius: 5px;
    font-size: 0.95em;
    background-color: var(--input-background);
    color: var(--input-text);
    line-height: 1.4;
  }

</style> 