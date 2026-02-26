<script lang="ts">
  import { page } from '$app/stores';

  const status = $derived($page.status ?? 500);
  const message = $derived($page.error?.message ?? 'Something went wrong');
</script>

<div class="error-page">
  <h1>{status}</h1>
  <p class="status-label">
    {#if status === 404}
      Page not found
    {:else if status === 500}
      Server error
    {:else}
      {message}
    {/if}
  </p>
  {#if status === 404}
    <p class="hint">The page you're looking for doesn't exist or has been moved.</p>
  {:else}
    <p class="hint">We're sorry, something went wrong on our end. Please try again later.</p>
  {/if}
  <a href="/" class="back-link">← Back to home</a>
</div>

<style>
  .error-page {
    max-width: 480px;
    margin: 4rem auto;
    padding: 2rem;
    text-align: center;
  }
  h1 {
    font-size: 4rem;
    font-weight: 700;
    margin: 0 0 0.5rem 0;
    color: var(--muted);
  }
  .status-label {
    font-size: 1.25rem;
    font-weight: 600;
    margin: 0 0 0.5rem 0;
  }
  .hint {
    color: var(--muted);
    font-size: 0.95rem;
    margin: 0 0 1.5rem 0;
  }
  .back-link {
    display: inline-block;
    padding: 0.6rem 1.2rem;
    background: var(--accent);
    color: var(--bg);
    border-radius: 6px;
    text-decoration: none;
    font-weight: 600;
  }
  .back-link:hover {
    background: var(--accent-hover);
  }
</style>
