<script>
  let isDropdownOpen = false;
  let userProfile = {
    name: "",
    avatar: "",
  };

  function toggleDropdown() {
    isDropdownOpen = !isDropdownOpen;
  }

  function handleLogout() {
    console.log("Logging out...");
  }

  function handleClickOutside(event) {
    if (isDropdownOpen && !event.target?.closest(".profile-section")) {
      isDropdownOpen = false;
    }
  }
</script>

<svelte:window on:click={handleClickOutside} />

<nav class="topbar">
  <div class="profile-section">
    <button class="profile-button" on:click={toggleDropdown}>
      <span class="username">{userProfile.name}</span>
      <img src={userProfile.avatar} alt="Profile" class="avatar" />
    </button>

    <div class="dropdown-menu" class:active={isDropdownOpen}>
      <button class="dropdown-item"> Profile </button>
      <button class="dropdown-item"> My Account </button>
      <button class="dropdown-item" on:click={handleLogout}> Log out </button>
    </div>
  </div>
</nav>

<style>
  .topbar {
    width: 100%;
    height: 60px;
    background-color: #ffffff;
    border-bottom: 1px solid #e5e5e5;
    display: flex;
    justify-content: flex-end;
    align-items: center;
    padding: 0 20px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
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
    background-color: #f5f5f5;
  }

  .avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    object-fit: cover;
  }

  .username {
    font-size: 14px;
    color: #333;
  }

  .dropdown-menu {
    position: absolute;
    top: 100%;
    right: 0;
    background-color: white;
    border-radius: 4px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
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
    color: #333;
  }

  .dropdown-item:hover {
    background-color: #f5f5f5;
  }
</style>
