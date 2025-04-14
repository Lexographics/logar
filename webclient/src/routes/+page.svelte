<script>
  import { goto } from '$app/navigation';
  import { userStore } from '$lib/store';
  import { onMount } from 'svelte';
  import { base } from '$app/paths';

  onMount(() => {
    if (userStore.current.token) {
      goto(`${base}/dashboard`);
    } else {
      goto(`${base}/login`);
    }
  });
</script>
