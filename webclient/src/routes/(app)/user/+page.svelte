<script>
  import { goto } from "$app/navigation";
  import BaseView from "$lib/BaseView.svelte";
  import { getActiveSessions, logout, revokeSession } from "$lib/service/user";
  import { userStore } from "$lib/store";
  import moment from "moment";
  import { onMount } from "svelte";

  let myActiveSessions = $state([]);

  async function onRevokeSession(session) {
    if (session.is_current) {
      const yes = confirm("Are you sure you want to revoke your own session?");
      if (!yes) {
        return;
      }
      await logout();
      userStore.current = null;
      goto("/login");
      return;
    }
    
    const yes = confirm("Are you sure you want to revoke this session?");
    if (!yes) {
      return;
    }

    const error = await revokeSession(session.token);
    if (error) {
      console.error(error);
    }
    myActiveSessions = myActiveSessions.filter(s => s.token !== session.token);
  }

  onMount(async () => {
    const [data, error] = await getActiveSessions();
    if (error) {
      console.error(error);
    }
    myActiveSessions = data;
  });
</script>

<BaseView>
  <div class="content">
    <h2 class="title">User / Sessions</h2>
    
    
    <h3 style="">My Active Sessions</h3>
    <ul class="item-list">
      {#each myActiveSessions as session}
      <li class="card" style="list-style: none; display: flex; justify-content: space-between; align-items: center;">
        <div>
          <span><span class="bold">Device:</span> {session.device} <span class="bold">{session.is_current ? '(You)' : ''}</span></span> <br>
          <span><span class="bold">Last Activity:</span> { moment(session.last_activity).fromNow() }</span> <br>
          <span><span class="bold">Created At:</span> { moment(session.created_at).fromNow() }</span>
        </div>
        <button onclick={() => onRevokeSession(session)} class="danger-button"> Revoke </button>
      </li>
      {/each}
    </ul>
  </div>
</BaseView>

<style>
  .content {
    padding: 0 2rem;
  }
  .title {
    margin-bottom: 20px;
    margin-top: 20px;
    text-align: center;
  }
  .card {
    padding: 20px 40px;
    border-radius: 10px;
    background-color: var(--background-color);
    box-shadow: rgba(0, 0, 0, 0.35) 0px 5px 15px;
  }
  .bold {
    font-weight: bold;
  }
  .item-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    margin-top: 1.5rem;
  }
  
  button:focus {
    outline: none;
    border-color: var(--primary-color);
    box-shadow: 0 0 0 2px rgba(52, 152, 219, 0.25);
  }
  
  button {
    padding: 0.6rem 0.8rem;
    border-radius: 5px;
    font-size: 0.95em;
    line-height: 1.4;
    display: inline-block;
    cursor: pointer;
    background-color: var(--primary-color);
    color: white;
    border: none;
    min-width: 80px;
    transition: background-color 0.2s, transform 0.1s;
    font-weight: 500;
  }

  button:hover {
    background-color: var(--primary-hover-color);
  }

  button:active {
    transform: translateY(1px);
  }

  .danger-button {
    background-color: var(--error-color);
  }

  .danger-button:hover {
    background-color: #c0392b;
  }
</style>