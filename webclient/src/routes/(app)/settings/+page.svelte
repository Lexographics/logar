<script>
  import BaseView from "$lib/widgets/BaseView.svelte";
  import userService from "$lib/service/userService";
  import { settingsStore, userStore } from "$lib/store";
  import LL, { setLocale } from "../../../i18n/i18n-svelte";
  import { setMomentLocale } from "$lib/moment";
  
  let selectedLanguage = $state(settingsStore.current.selectedLanguage);
  async function languageChanged() {
    settingsStore.current.selectedLanguage = selectedLanguage;
    setLocale(settingsStore.current.selectedLanguage || "en");
    await setMomentLocale(settingsStore.current.selectedLanguage);
  }

  let displayName = $state(userStore.current.user?.display_name || "");

  function changePassword() {
    alert('Password change not implemented yet.');
  }

  function deleteAccount() {
    if (confirm('Are you sure you want to delete your account? This action cannot be undone.')) {
      alert('Account deletion not implemented yet.');
    }
  }

  async function saveDisplayName() {
    const [user, err] = await userService.updateUser(displayName);
    if (err) {
      showToast(err.message);
    }
    userStore.current.user = user;
  }
</script>

<BaseView>
  
  <div class="settings-container">
    <h2>{$LL.settings.title()}</h2>
    
    <div class="settings-section card">
      <h3>{$LL.settings.theme.title()}</h3>
      <div class="setting-item">
        <label for="theme-select">{$LL.settings.theme.select_theme()}</label>
        <select id="theme-select" bind:value={settingsStore.current.currentTheme}>
          <option value="light">{$LL.settings.theme.options.light()}</option>
          <option value="dark">{$LL.settings.theme.options.dark()}</option>
          <option value="system">{$LL.settings.theme.options.system()}</option>
        </select>
      </div>
      <div class="setting-info">
        <p>{$LL.settings.theme.current_theme(settingsStore.current.currentTheme)}</p>
        <p class="setting-description">{$LL.settings.theme.theme_description()}</p>
      </div>
    </div>

    <div class="settings-section card">
      <h3>{$LL.settings.language.title()}</h3>
      <div class="setting-item">
        <label for="language-select">{$LL.settings.language.language.title()}</label>
        <select id="language-select" bind:value={selectedLanguage} onchange={languageChanged}>
          <option value="en">English (US)</option>
          <option value="zh">中文</option>
          <option value="ru">Русский</option>
          <option value="tr">Türkçe</option>
          <option value="kk">Қазақша</option>
          <option value="az">Azərbaycanca</option>
        </select>
      </div>
    </div>

    <div class="settings-section card" id="profile">
      <h3>{$LL.settings.profile.title()}</h3>
      <div class="setting-item">
        <span>{$LL.settings.profile.username()}</span>
        <span>@{userStore.current.user?.username}</span>
      </div>
      <div class="setting-item">
        <label for="display-name">{$LL.settings.profile.display_name()}</label>
        <input type="text" id="display-name" bind:value={displayName} placeholder="Your display name">
      </div>
      <div class="setting-item">
        <span></span>
        <button onclick={saveDisplayName}>{$LL.settings.profile.save()}</button>
      </div>
    </div>

    <div class="settings-section card">
      <h3>{$LL.settings.account.title()}</h3>
      <div class="setting-item">
        <span>{$LL.settings.account.change_password()}</span>
        <button onclick={changePassword}>{$LL.settings.account.change_password()}</button>
      </div>
      <div class="setting-item">
        <span>{$LL.settings.account.delete_account()}</span>
        <button class="danger-button" onclick={deleteAccount}>{$LL.settings.account.delete_account()}</button>
      </div>
    </div>

  </div>
</BaseView>


<style>
  .settings-container {
    max-width: 800px;
    margin: 1rem auto;
    padding: 1rem;
    display: grid;
    gap: 1.5rem;
  }

  h2 {
    text-align: center;
    margin-bottom: 2rem;
    color: var(--text-color);
  }

  .settings-section.card {
    background-color: var(--card-background);
    border-radius: 8px;
    padding: 1.5rem;
    box-shadow: 0 2px 4px var(--shadow-color);
    border: none;
    margin-bottom: 0;
  }

  h3 {
    margin-top: 0;
    margin-bottom: 1.5rem;
    font-size: 1.3em;
    color: var(--text-color);
    border-bottom: 1px solid var(--border-color);
    padding-bottom: 0.75rem;
  }

  .setting-item {
    margin-bottom: 1rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
  }

  .setting-item:last-child {
      margin-bottom: 0;
  }

  label,
  .setting-item > span {
    color: var(--text-secondary-color);
    font-weight: 500;
  }

  select,
  button,
  input[type="text"] {
    padding: 0.6rem 0.8rem;
    border: 1px solid var(--input-border);
    border-radius: 5px;
    font-size: 0.95em;
    background-color: var(--input-background);
    color: var(--input-text);
    line-height: 1.4;
  }

  select {
    padding-right: 2rem;
    background-image: url('data:image/svg+xml;charset=US-ASCII,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%22292.4%22%20height%3D%22292.4%22%3E%3Cpath%20fill%3D%22%23007AFF%22%20d%3D%22M287%2069.4a17.6%2017.6%200%200%200-13-5.4H18.4c-5%200-9.3%201.8-12.9%205.4A17.6%2017.6%200%200%200%200%2082.2c0%205%201.8%209.3%205.4%2012.9l128%20127.9c3.6%203.6%207.8%205.4%2012.8%205.4s9.2-1.8%2012.8-5.4L287%2095c3.5-3.5%205.4-7.8%205.4-12.8%200-5-1.9-9.2-5.5-12.8z%22%2F%3E%3C%2Fsvg%3E');
    background-repeat: no-repeat;
    background-position: right .7em top 50%;
    background-size: .65em auto;
    -webkit-appearance: none;
    -moz-appearance: none;
    appearance: none;
    text-align: left;
  }
  
  option {
    background-color: var(--input-background-opaque);
  }

  select:focus,
  button:focus,
  input:focus {
    outline: none;
    border-color: var(--primary-color);
    box-shadow: 0 0 0 2px rgba(52, 152, 219, 0.25);
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

  /*
  .toggle-switch {
    position: relative;
    display: inline-block;
    width: 45px;
    height: 24px;
  }

  .toggle-switch input {
    opacity: 0;
    width: 0;
    height: 0;
  }
  .slider {
    position: absolute;
    cursor: pointer;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: #ccc;
    transition: .3s;
    border-radius: 24px;
  }

  .slider:before {
    position: absolute;
    content: "";
    height: 18px;
    width: 18px;
    left: 3px;
    bottom: 3px;
    background-color: white;
    transition: .3s;
    border-radius: 50%;
  }

  input:checked + .slider {
    background-color: var(--primary-color);
  }

  input:focus + .slider {
    box-shadow: 0 0 1px var(--primary-color);
  }

  input:checked + .slider:before {
    transform: translateX(21px);
  }
  */

  .settings-section.card {
    position: relative;
    overflow: hidden;
    transition: transform 0.2s, box-shadow 0.2s;
  }

  .setting-info {
    margin-top: 1rem;
    font-size: 0.9em;
    color: var(--text-secondary-color);
  }

  .setting-description {
    margin-top: 0.5rem;
    font-style: italic;
    opacity: 0.8;
  }
</style>