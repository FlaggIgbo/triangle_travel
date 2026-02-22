<script lang="ts">
  import '../app.css';
  import { page } from '$app/stores';

  const tabs = [
    { href: '/triangle-travel', label: 'Triangle Travel' },
    { href: '/chat', label: 'AI Chat' },
    { href: '/my-flights', label: 'My Flights' }
  ];
</script>

<div class="app">
  <header class="nav-header">
    <a href="/" class="logo">Travel App</a>
    <nav class="tabs">
      {#each tabs as tab}
        <a
          href={tab.href}
          class="tab"
          class:active={$page.url.pathname === tab.href || ($page.url.pathname === '/' && tab.href === '/triangle-travel')}
        >
          {tab.label}
        </a>
      {/each}
    </nav>
  </header>
  <main class="content">
    <slot />
  </main>
</div>

<style>
  .app {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
  }

  .nav-header {
    display: flex;
    align-items: center;
    gap: 2rem;
    padding: 1rem 1.5rem;
    background: var(--surface);
    border-bottom: 1px solid var(--border);
  }

  .logo {
    font-weight: 700;
    font-size: 1.25rem;
    color: var(--text);
    text-decoration: none;
  }

  .logo:hover {
    color: var(--accent);
  }

  .tabs {
    display: flex;
    gap: 0.5rem;
  }

  .tab {
    padding: 0.5rem 1rem;
    border-radius: 6px;
    color: var(--muted);
    text-decoration: none;
    font-size: 0.95rem;
  }

  .tab:hover {
    color: var(--text);
    background: rgba(255, 255, 255, 0.05);
  }

  .tab.active {
    color: var(--accent);
    background: rgba(88, 166, 255, 0.15);
  }

  .content {
    flex: 1;
  }
</style>
