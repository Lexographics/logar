<script>
  import { goto } from "$app/navigation";
  import BaseView from "$lib/widgets/BaseView.svelte";
  import userService from "$lib/service/userService";
  import { userStore } from "$lib/store";
  import Modal from "$lib/widgets/Modal.svelte";
  import moment from "moment";
  import { onMount } from "svelte";
  import LL from "../../../i18n/i18n-svelte";

  let myActiveSessions = $state([]);
  
  let createUserModal = $state(null);
  let newUser = $state({
    username: "",
    display_name: "",
    password: "",
    is_admin: false,
  });
  async function onCreateUser() {
    const [data, error] = await userService.createUser(newUser.username, newUser.password, newUser.display_name, newUser.is_admin);
    if (error) {
      console.error(error);
    }
    users = [...users, data];
    createUserModal?.closeModal();
  }

  async function onRevokeSession(session) {
    if (session.is_current) {
      const yes = confirm("Are you sure you want to revoke your own session?");
      if (!yes) {
        return;
      }
      await userService.logout();
      userStore.current = null;
      goto("/login");
      return;
    }
    
    const yes = confirm("Are you sure you want to revoke this session?");
    if (!yes) {
      return;
    }

    const error = await userService.revokeSession(session.token);
    if (error) {
      console.error(error);
    }
    myActiveSessions = myActiveSessions.filter(s => s.token !== session.token);
  }

  onMount(async () => {
    const [data, error] = await userService.getActiveSessions();
    if (error) {
      console.error(error);
    }
    myActiveSessions = data;

    const [usersData, err] = await userService.getAllUsers();
    if (err) {
      console.error(err);
    }
    users = usersData;
    loaded = true;
  });

  let loaded = $state(false);
  let users = $state([]);
</script>

<BaseView loaded={loaded}>
  <div class="content">
    <h2 class="title">{$LL.user_sessions.title()}</h2>
    
    
    <h3 style="">{$LL.user_sessions.my_active_sessions.title()}</h3>
    <ul class="item-list">
      {#each myActiveSessions as session}
      <li class="card" style="list-style: none; display: flex; justify-content: space-between; align-items: center;">
        <div>
          <span><span class="bold">{$LL.user_sessions.my_active_sessions.device()}:</span> {session.device} <span class="bold">{session.is_current ? $LL.user_sessions.my_active_sessions.you() : ''}</span></span> <br>
          <span><span class="bold">{$LL.user_sessions.my_active_sessions.last_activity()}:</span> { moment(session.last_activity).fromNow() }</span> <br>
          <span><span class="bold">{$LL.user_sessions.my_active_sessions.created_at()}:</span> { moment(session.created_at).fromNow() }</span>
        </div>
        <div>
          <button onclick={() => onRevokeSession(session)} class="danger-button"> {$LL.user_sessions.my_active_sessions.revoke()} </button>
        </div>
      </li>
      {/each}
    </ul>

    <h3 style="margin-top: 2rem;">{$LL.user_sessions.all_users.title()}</h3>
    <ul class="item-list">
      {#each users as user}
      <li class="card" style="list-style: none; display: flex; justify-content: space-between; align-items: center;">
        <div style="display: flex; align-items: center; gap: 1rem;">
          <img src="https://api.dicebear.com/9.x/thumbs/svg?seed={user.username}" alt="avatar" class="avatar" />

          <div>
            <span>
              <span>{user.display_name} (@{user.username}) {user.is_admin ? '(' + $LL.user_sessions.all_users.admin() + ')' : ''}</span>
            </span>
            <br>
            {#if moment(user.last_activity).isAfter(moment().subtract(10, 'second'))}
            <span style="color: var(--success-color);" class="blink">
              <div style="background-color: var(--success-color); padding: 0.3rem; margin: 0.1rem; border-radius: 50%; display: inline-block;"></div>
              {$LL.user_sessions.all_users.online()}
            </span>
            {:else}
              <span style="color: var(--error-color);">{$LL.user_sessions.all_users.last_seen()}: {moment(user.last_activity).fromNow()}</span>
            {/if}
          </div>
        </div>
      </li>
      {/each}
    </ul>
    
    {#if userStore.current?.user?.is_admin}
      <button style="margin-top: 1rem;" onclick={() => createUserModal?.openModal()} class="">Create User</button>
    {/if}

    <div style="margin-top: 2rem;"></div>
  </div>

</BaseView>

<Modal bind:this={createUserModal} title="Create User" onClose={() => {
  newUser = {
    username: "",
    password: "",
    is_admin: false,
  };
}}>
  <div style="display: flex; flex-direction: column; gap: 1rem; width: 100%; margin-top: 1rem;">
    <input type="text" bind:value={newUser.username} placeholder="Username" />
    <input type="text" bind:value={newUser.display_name} placeholder="Display Name" />
    <input type="text" bind:value={newUser.password} placeholder="Password" />
    <label>
      <input type="checkbox" bind:checked={newUser.is_admin} />
      <span>Is Admin</span>
    </label>
    <button onclick={onCreateUser}>Create User</button>
  </div>
</Modal>

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
    margin-top: 0.5rem;
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

  .avatar {
    width: 2rem;
    height: 2rem;
    border-radius: 50%;
  }

  input {
    padding: 0.6rem 0.8rem;
    border-radius: 5px;
    font-size: 0.95em;
    line-height: 1.4;
    border: 1px solid var(--border-color);
  }

  .blink {
    animation: blink 1s  infinite;
  }

  @keyframes blink {
    50% {
      opacity: 0.5;
    }
    100% {
      opacity: 1;
    }
  }
</style>