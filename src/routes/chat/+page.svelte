<script lang="ts">
  interface Message {
    role: 'user' | 'assistant';
    content: string;
  }

  let messages = $state<Message[]>([]);
  let input = $state('');
  let loading = $state(false);

  async function send() {
    const text = input.trim();
    if (!text || loading) return;

    messages = [...messages, { role: 'user', content: text }];
    input = '';
    loading = true;

    try {
      const res = await fetch('/api/chat', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ message: text })
      });
      const data = await res.json().catch(() => ({}));
      const reply = data.reply ?? 'Sorry, I couldn\'t process that. The AI chat is not yet configured.';
      messages = [...messages, { role: 'assistant', content: reply }];
    } catch {
      messages = [...messages, { role: 'assistant', content: 'Something went wrong. Please try again.' }];
    } finally {
      loading = false;
    }
  }
</script>

<div class="chat-container">
  <header class="header">
    <h1>AI Chat</h1>
    <p class="tagline">Ask travel-related questions</p>
  </header>

  <div class="messages">
    {#if messages.length === 0}
      <p class="placeholder">Ask me anything about travel—flights, destinations, tips, or trip planning.</p>
    {:else}
      {#each messages as msg}
        <div class="message" class:user={msg.role === 'user'} class:assistant={msg.role === 'assistant'}>
          <span class="role">{msg.role === 'user' ? 'You' : 'Assistant'}</span>
          <p>{msg.content}</p>
        </div>
      {/each}
    {/if}
  </div>

  <form class="input-form" onsubmit={(e) => { e.preventDefault(); send(); }}>
    <input
      type="text"
      bind:value={input}
      placeholder="Ask a question..."
      disabled={loading}
    />
    <button type="submit" disabled={loading || !input.trim()}>
      {loading ? '…' : 'Send'}
    </button>
  </form>
</div>

<style>
  .chat-container {
    max-width: 640px;
    margin: 0 auto;
    padding: 2rem 1.5rem;
    display: flex;
    flex-direction: column;
    min-height: calc(100vh - 120px);
  }
  .header { margin-bottom: 1.5rem; }
  .header h1 { font-size: 1.5rem; font-weight: 600; margin: 0 0 0.25rem 0; }
  .tagline { color: var(--muted); font-size: 0.9rem; margin: 0; }
  .messages {
    flex: 1;
    overflow-y: auto;
    margin-bottom: 1rem;
  }
  .placeholder {
    color: var(--muted);
    font-style: italic;
    padding: 2rem;
    text-align: center;
  }
  .message {
    padding: 1rem;
    margin-bottom: 0.75rem;
    border-radius: 8px;
    border: 1px solid var(--border);
  }
  .message.user { background: rgba(88, 166, 255, 0.1); border-color: var(--accent); }
  .message.assistant { background: var(--surface); }
  .role {
    font-size: 0.75rem;
    color: var(--muted);
    display: block;
    margin-bottom: 0.35rem;
  }
  .message p { margin: 0; font-size: 0.95rem; }
  .input-form {
    display: flex;
    gap: 0.5rem;
    padding: 0.5rem 0;
  }
  .input-form input {
    flex: 1;
    padding: 0.6rem 1rem;
    border: 1px solid var(--border);
    border-radius: 6px;
    background: var(--surface);
    color: var(--text);
    font-size: 0.95rem;
  }
  .input-form input:focus { outline: none; border-color: var(--accent); }
  .input-form button {
    padding: 0.6rem 1.2rem;
    background: var(--accent);
    color: var(--bg);
    border: none;
    border-radius: 6px;
    font-weight: 600;
    cursor: pointer;
  }
  .input-form button:hover:not(:disabled) { background: var(--accent-hover); }
  .input-form button:disabled { opacity: 0.6; cursor: not-allowed; }
</style>
