<script>
  // Removed unused import: import { time } from 'svelte/internal';

  let {arg, loading, value = $bindable()} = $props();

	const units = $state({
    nsec: 1,
    µsec: 1000,
    msec: 1_000_000,
    sec: 1_000_000_000,
    min: 60 * 1_000_000_000,
    hour: 3600 * 1_000_000_000,
    day: 86400 * 1_000_000_000,
    week: 7 * 86400 * 1_000_000_000,
    month: 30 * 86400 * 1_000_000_000,
    year: 365 * 86400 * 1_000_000_000,
  });
	const unitKeys = Object.keys(units);

	function composeDurationNs(num, unit) {
    const multiplier = units[unit];
    if (typeof num !== 'number' || isNaN(num) || !multiplier) {
      return 0;
    }
		return Math.round(num * multiplier);
	}

	let currentNum = $state(0);
	let currentUnit = $state('sec');
	let isDuration = $derived(arg.kind === 'duration');

	$effect(() => {
		if (isDuration) {
			const newValue = composeDurationNs(currentNum, currentUnit);
			if (newValue !== value) {
				value = newValue;
			}
		}
	});
</script>

<label class="arg-label">
  <div>
    {#if arg.name}
      <span class="arg-type">{arg.name}</span>
      <br>
    {/if}
    <span class="arg-type">{arg.type}</span>
  </div>
  {#if arg.kind === "text"}
    <input type="text" bind:value={value} class="input input-sm arg-input" disabled={loading} />
  {:else if arg.kind === "int" || arg.kind === "int8" || arg.kind === "int16" || arg.kind === "int32" || arg.kind === "int64" || arg.kind === "uint" || arg.kind === "uint8" || arg.kind === "uint16" || arg.kind === "uint32" || arg.kind === "uint64"}
    <input type="number" bind:value={value} class="input input-sm arg-input" disabled={loading} step="1" />
  {:else if arg.kind === "float" || arg.kind === "float32" || arg.kind === "float64"}
    <input type="number" bind:value={value} class="input input-sm arg-input" disabled={loading} step="any" />
  {:else if arg.kind === "bool"}
    <input type="checkbox" bind:checked={value} class="checkbox arg-input" disabled={loading} />
  {:else if arg.kind === "time"}
    <input type="datetime-local" bind:value={value} class="input input-sm arg-input" disabled={loading} />
  {:else if arg.kind === "duration"}
    <div class="duration-input-group">
      <input type="number" bind:value={currentNum} class="input input-sm arg-input duration-number" disabled={loading} step="any" />
      <select bind:value={currentUnit} class="select select-sm arg-input duration-unit" disabled={loading}>
        {#each unitKeys as unit}
          <option value={unit}>{unit}</option>
        {/each}
      </select>
    </div>
  {:else}
    <input type="text" bind:value={value} class="input input-sm arg-input" disabled={loading} placeholder="Unknown kind" />
  {/if}
</label>

<style>
  .arg-label {
		display: flex;
		align-items: center;
		gap: 0.35rem;
		background-color: var(--input-background);
		padding: 0.2rem 0.4rem;
		border-radius: 4px;
		border: 1px solid var(--border-color);
		font-size: 0.85rem;
    flex-wrap: nowrap;
	}
	.arg-type {
		font-size: 0.75rem;
		color: var(--text-secondary-color);
    white-space: nowrap;
	}
  .arg-type:nth-child(3) {
    font-size: 0.5rem;
  }

	.arg-input {
		flex-shrink: 0;
	}

	.arg-label .input.arg-input,
  .arg-label .select.arg-input {
		min-width: 50px;
		max-width: 180px;
		flex-grow: 1;
		background-color: var(--input-background);
    border: 1px solid var(--input-border);
    border-radius: 5px;
    font-size: 0.90em;
    color: var(--input-text);
    line-height: 1.4;
    padding: 0.4rem 0.6rem;
	}

  input[type="text"].arg-input,
  input[type="number"].arg-input,
  input[type="datetime-local"].arg-input {
     padding: 0.4rem 0.6rem;
  }

  .arg-label .select.arg-input {
    padding-right: 1.5rem;
    min-width: 60px;
    flex-grow: 0;
    flex-basis: auto;
  }


	.arg-label > .checkbox.arg-input {
		margin: 0;
    width: 1rem;
    height: 1rem;
	}

  .duration-input-group {
    display: flex;
    align-items: stretch;
    gap: 0.2rem;
    flex-grow: 1;
    flex-shrink: 1;
    min-width: 120px;
  }

  .duration-input-group .duration-number {
    flex-grow: 1;
    min-width: 40px;
  }

  .duration-input-group .duration-unit {
    flex-grow: 0;
  }

  .arg-input:disabled {
    cursor: not-allowed;
    opacity: 0.6;
  }
</style>