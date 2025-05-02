<script>
  import { goto } from '$app/navigation';
  import { userStore } from '$lib/store';
  import { onMount } from 'svelte';
  import { getBasePath } from '$lib/utils';

  onMount(() => {
    if (userStore.current.token) {
      goto(`${getBasePath()}/dashboard`);
    } else {
      goto(`${getBasePath()}/login`);
    }
  });
</script>
