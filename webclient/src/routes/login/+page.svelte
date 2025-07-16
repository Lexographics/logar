<script>
  import { goto } from '$app/navigation';
  import { getBasePath } from '$lib/utils';
  import userService from '$lib/service/userService';
  import { userStore } from '$lib/store';
  import { onMount } from 'svelte';
  import LL from '../../i18n/i18n-svelte';
  import { showToast } from '$lib/toast';
  import { StatusCode } from '$lib/types/response';

  let username = $state("");
  let password = $state("");

  async function handleSubmit(e) {
    e.preventDefault();
    
    const [data, error] = await userService.login(username, password);
    if (error) {
      if(error.response?.data?.status_code === StatusCode.InvalidCredentials) {
        showToast("Invalid username or password");
      } else {
        showToast(error.message);
      }
    } else {
      userStore.current = data;
      goto(`${getBasePath()}/dashboard`);
    }
  }

  let mouseX = $state(0);
  let mouseY = $state(0);
  let windowWidth = $state(0);
  let windowHeight = $state(0);

  let stars = $state([]);

  const PROXIMITY_RADIUS = 500;

  onMount(() => {
    if (userStore.current?.token) {
      goto(`${getBasePath()}/dashboard`);
    }
  });

  onMount(() => {
    windowWidth = window.innerWidth;
    windowHeight = window.innerHeight;

    const numStars = 100;
    const newStars = [];
    for (let i = 0; i < numStars; i++) {
      newStars.push({
        id: i,
        topPercent: Math.random() * 100,
        leftPercent: Math.random() * 100,
        width: Math.random() * 2 + 1,
        height: Math.random() * 2 + 1,
      });
    }
    stars = newStars;

    const handleMouseMove = (event) => {
      mouseX = event.clientX;
      mouseY = event.clientY;
    };

    const handleResize = () => {
      windowWidth = window.innerWidth;
      windowHeight = window.innerHeight;
    };

    window.addEventListener('mousemove', handleMouseMove);
    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('mousemove', handleMouseMove);
      window.removeEventListener('resize', handleResize);
    };
  });

  function calculateStarStyle(star) {
    const starX = (star.leftPercent / 100) * windowWidth;
    const starY = (star.topPercent / 100) * windowHeight;

    const dx = starX - mouseX;
    const dy = starY - mouseY;
    const distance = Math.sqrt(dx * dx + dy * dy);

    let scale = 1;

    if (distance < PROXIMITY_RADIUS) {
      const proximityFactor = 1 - (distance / PROXIMITY_RADIUS);
      scale = 1 + proximityFactor * 1.5;
    }

    return `
      top: ${star.topPercent}%;
      left: ${star.leftPercent}%;
      width: ${star.width}px;
      height: ${star.height}px;
      transform: scale(${scale});
    `;
  }
</script>

<div class="auth-container">
  {#each stars as star (star.id)}
    <div class="star" style={calculateStarStyle(star)}></div>
  {/each}

  <div class="auth-card">
    <h1>{$LL.login.title()}</h1>
    <form onsubmit={handleSubmit}>
      <div class="form-group">
        <label for="username">{$LL.login.username()}</label>
        <input type="text" id="username" bind:value={username} required />
      </div>
      <div class="form-group">
        <label for="password">{$LL.login.password()}</label>
        <input type="password" id="password" bind:value={password} required />
      </div>
      <button type="submit" class="submit-button">{$LL.login.login()}</button>
    </form>
  </div>
</div>

<style>
  .auth-container {
    position: relative;
    overflow: hidden;
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100dvh;
    max-height: 100dvh;
  }

  :global([data-theme="dark"]) .auth-container {
    background-color: #0a0a1f;
  }
  
  @media (prefers-color-scheme: dark) {
    :global([data-theme="system"]) .auth-container {
      background-color: #0a0a1f;
    }
  }

  :global([data-theme="light"]) .auth-container {
    background-color: #f0f5ff;
  }

  @media (prefers-color-scheme: light) {
    :global([data-theme="system"]) .auth-container {
      background-color: #f0f5ff;
    }
  }

  .star {
    position: absolute;
    background-color: white;
    border-radius: 50%;
    will-change: transform, opacity, box-shadow;
    transition: transform 0.15s linear, opacity 0.15s linear, box-shadow 0.15s linear;
    opacity: 0.7;
    box-shadow: 0 0 var(--animation-shadow, 3px) white;
  }

  .auth-card {
    background-color: var(--card-background);
    padding: 2rem 3rem;
    border-radius: 2rem;
    box-shadow: 0 0 15px var(--shadow-color);
    width: 100%;
    max-width: 400px;
    text-align: center;
    position: relative;
    z-index: 10;
  }

  h1 {
    margin-bottom: 1.5rem;
    color: var(--text-color);
  }

  .form-group {
    margin-bottom: 1rem;
    text-align: left;
  }

  label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
    color: var(--text-secondary-color);
  }

  input[type="text"],
  input[type="password"] {
    width: 100%;
    padding: 0.75rem;
    font-size: 1rem;
    background-color: var(--input-background);
    border: 1px solid var(--input-border);
    color: var(--input-text);
    border-radius: 4px;
    box-sizing: border-box;
    transition: border-color 0.2s, box-shadow 0.2s, background-color 0.2s;
  }

  input[type="text"]:focus,
  input[type="password"]:focus {
    outline: none;
    border-color: var(--primary-color);
    background-color: var(--input-background);
    box-shadow: 0 0 5px var(--shadow-color);
  }

  input::placeholder {
    color: var(--text-secondary-color);
    opacity: 0.7;
  }

  input:-webkit-autofill,
  input:-webkit-autofill:hover,
  input:-webkit-autofill:focus,
  input:-webkit-autofill:active {
    -webkit-text-fill-color: var(--input-text) !important;
    transition: background-color 5000s ease-in-out 0s;
    caret-color: var(--input-text);
  }

  .submit-button {
    width: 100%;
    padding: 0.75rem;
    font-size: 1.1rem;
    font-weight: 600;
    color: white;
    background-color: var(--primary-color);
    border: none;
    border-radius: 4rem;
    cursor: pointer;
    transition: background-color 0.2s, transform 0.2s, box-shadow 0.2s;
    margin-top: 1rem;
  }

  .submit-button:hover, .submit-button:focus {
    background-color: var(--primary-hover-color);
    transform: scale(1.03);
    box-shadow: 0 4px 15px var(--shadow-color);
  }

  .submit-button:focus {
     outline: 2px solid var(--primary-hover-color);
     outline-offset: 2px;
  }
</style>