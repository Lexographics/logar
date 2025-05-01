<script lang="ts">
  let modal: HTMLDivElement | null = $state(null);

  let { title, children = null, onClose = null }: { title: string, children?: any, onClose?: () => void } = $props();

  export function openModal() {
    modal.classList.add('visible');
  }
  export function closeModal() {
    modal.classList.remove('visible');
    onClose?.();
  }
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<div bind:this={modal} class="modal" role="dialog" tabindex="-1" onmousedown={(event) => {
  if (event.target === modal) { closeModal(); } 
}}>
  <div class="modal-content">
    <div style="display: flex; align-items: center; justify-content: space-between;">
      <h2>{title}</h2>
      <button class="button" onclick={closeModal} aria-label="Close">
        <i class="fa-solid fa-xmark fs-4"></i>
      </button>
    </div>
    <hr>
    {@render children?.()}
  </div>
</div>

<style>
  .modal > * {
    pointer-events: none;
  }

  :global(.modal.visible) > * {
    pointer-events: initial;
  }

  .modal {
    /* display: none; */
    display: block;
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.0);
    pointer-events: none;
    transition: 0.3s;
    z-index: 1000;
  }

  :global(.modal.visible) {
    /* display: block; */
    background-color: rgba(0, 0, 0, 0.5);
    pointer-events: initial;
  }

  .modal-content {
    background-color: #fff;
    padding: 20px;
    border-radius: 5px;
    position: relative;
    max-width: 60vw;
    transition: 0.3s;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%) scale(0.8);
    opacity: 0;
  }

  :global(.modal.visible .modal-content) {
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%) scale(1);
    opacity: 1;
  }

  .button {
    background-color: transparent;
    border: none;
    cursor: pointer;
    padding: 0;
    margin: 0;
    font-size: 1rem;
    
  }
</style>