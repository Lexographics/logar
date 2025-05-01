<script>
  import { goto } from "$app/navigation";
  import { base } from "$app/paths";
  import LL from "../../i18n/i18n-svelte";
  import userService from "../service/userService";
  import { userStore } from "../store";

  let isDropdownOpen = $state(false);

  function toggleDropdown() {
    isDropdownOpen = !isDropdownOpen;
  }

  async function handleLogout() {
    await userService.logout();
    userStore.current = {
      token: null,
    };
    goto(`${base}/login`);
  }

  function handleClickOutside(event) {
    if (isDropdownOpen && !event.target?.closest(".profile-section")) {
      isDropdownOpen = false;
    }
  }
</script>

<svelte:window onclick={handleClickOutside} />

<nav class="topbar">
  <div class="profile-section">
    <button class="profile-button" onclick={toggleDropdown}>
      <div style="display: flex; flex-direction: column; align-items: center;">
        <span class="username">{userStore.current.user?.display_name}</span>
        <span class="username" style="font-size: 0.75rem;">@{userStore.current.user?.username}</span>
      </div>
      
      <img src={"https://api.dicebear.com/9.x/thumbs/svg?seed=" + userStore.current.user?.username} alt="avatar" class="avatar" />
    </button>

    <div class="dropdown-menu" class:active={isDropdownOpen}>
      <button class="dropdown-item" onclick={() => { goto(`${base}/settings#profile`) }}> {$LL.topbar.my_account()} </button>
      <button class="dropdown-item" onclick={handleLogout}> {$LL.topbar.sign_out()} </button>
    </div>
  </div>
</nav>

<style>
  .topbar {
    height: 60px;
    background-color: var(--header-background);
    border-bottom: 1px solid var(--border-color);
    display: flex;
    justify-content: flex-end;
    align-items: center;
    padding: 0 20px;
    box-shadow: 0 2px 4px var(--shadow-color);
  }

  .profile-section {
    position: relative;
  }

  .profile-button {
    display: flex;
    align-items: center;
    gap: 10px;
    background: none;
    border: none;
    cursor: pointer;
    padding: 8px;
    border-radius: 4px;
  }

  .profile-button:hover {
    background-color: var(--input-background);
  }

  .avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    object-fit: cover;
  }

  .username {
    font-size: 14px;
    color: var(--text-color);
  }

  .dropdown-menu {
    position: absolute;
    top: 100%;
    right: 0;
    background-color: var(--card-background);
    border-radius: 4px;
    box-shadow: 0 2px 10px var(--shadow-color);
    min-width: 150px;
    margin-top: 8px;
    opacity: 0;
    visibility: hidden;
    transform: translateY(-10px);
    transition: all 0.2s ease-in-out;
  }

  .dropdown-menu.active {
    opacity: 1;
    visibility: visible;
    transform: translateY(0);
  }

  .dropdown-item {
    width: 100%;
    padding: 12px 16px;
    border: none;
    background: none;
    text-align: left;
    cursor: pointer;
    color: var(--text-color);
  }

  .dropdown-item:hover {
    background-color: var(--input-background);
  }
</style>
