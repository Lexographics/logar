<script>
  import BaseView from "$lib/widgets/BaseView.svelte";

</script>

<BaseView>
  <div class="page">
    <h2 class="title">Feature Flags</h2>
    <p>Feature flags are not yet implemented.</p>
  </div>
</BaseView>

<style>
  .page {
    padding: 24px;
    display: flex;
    flex-direction: column;
    color: var(--text-color);
  }

  .title {
    font-size: 1.5rem;
    font-weight: 600;
    color: var(--text-color);
  }
</style>