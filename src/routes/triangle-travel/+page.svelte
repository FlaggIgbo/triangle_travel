<script lang="ts">
  interface TriangleResult {
    driveThenFly: Record<string, number>;
    flyThenFly: Record<string, number>;
    avgPrice: number;
  }

  let start = $state('JFK');
  let end = $state('SFO');
  let startDate = $state('2025-03-01');
  let endDate = $state('2025-03-15');
  let alliance = $state('None');
  let airline = $state('');
  let cabin = $state('economy');
  let loading = $state(false);
  let error = $state<string | null>(null);
  let result = $state<TriangleResult | null>(null);

  const alliances = [
    { value: 'None', label: 'Any' },
    { value: 'ONE_WORLD', label: 'Oneworld' },
    { value: 'SKY_TEAM', label: 'SkyTeam' },
    { value: 'STAR_ALLIANCE', label: 'Star Alliance' },
    { value: 'ALL', label: 'All Alliances' }
  ];

  const airlines = [
    { value: '', label: 'Any airline' },
    { value: 'UA', label: 'United' },
    { value: 'AA', label: 'American' },
    { value: 'DL', label: 'Delta' },
    { value: 'B6', label: 'JetBlue' },
    { value: 'WN', label: 'Southwest' },
    { value: 'AS', label: 'Alaska' },
    { value: 'LH', label: 'Lufthansa' },
    { value: 'BA', label: 'British Airways' }
  ];

  function kayakMultiCity(from: string, to: string, via: string, dep: string, ret: string): string {
    let url = `https://www.kayak.com/flights/${from}-${to}/${dep}/${via}-${from}/${ret}?sort=bestflight_a`;
    if (alliance === 'ONE_WORLD' || alliance === 'SKY_TEAM' || alliance === 'STAR_ALLIANCE') {
      url += `&fs=alliance=${alliance}`;
    } else if (alliance === 'ALL') {
      url += '&fs=alliance=ONE_WORLD,SKY_TEAM,STAR_ALLIANCE';
    }
    if (airline) {
      url += `&fs=airline=${airline}`;
    }
    return url;
  }

  async function search() {
    loading = true;
    error = null;
    result = null;
    try {
      const res = await fetch('/api/search', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ start, end, startDate, endDate, cabin, alliance })
      });
      if (!res.ok) {
        const err = await res.json().catch(() => ({}));
        throw new Error(err.error || res.statusText);
      }
      result = await res.json();
    } catch (e) {
      error = e instanceof Error ? e.message : 'Search failed';
    } finally {
      loading = false;
    }
  }
</script>

<div class="container">
  <header class="header">
    <h1>Triangle Travel</h1>
    <p class="tagline">Find a third city to explore on your round-trip itinerary</p>
  </header>

  <form class="search-form" onsubmit={(e) => { e.preventDefault(); search(); }}>
    <div class="grid">
      <label>
        <span>Start airport</span>
        <input type="text" bind:value={start} placeholder="e.g. JFK" maxlength="4" />
      </label>
      <label>
        <span>End airport</span>
        <input type="text" bind:value={end} placeholder="e.g. SFO" maxlength="4" />
      </label>
      <label>
        <span>Outbound date</span>
        <input type="date" bind:value={startDate} />
      </label>
      <label>
        <span>Return date</span>
        <input type="date" bind:value={endDate} />
      </label>
      <label>
        <span>Alliance</span>
        <select bind:value={alliance}>
          {#each alliances as a}
            <option value={a.value}>{a.label}</option>
          {/each}
        </select>
      </label>
      <label>
        <span>Airline</span>
        <select bind:value={airline}>
          {#each airlines as a}
            <option value={a.value}>{a.label}</option>
          {/each}
        </select>
      </label>
    </div>
    <button type="submit" disabled={loading}>
      {loading ? 'Searching…' : 'Find triangle options'}
    </button>
  </form>

  {#if error}
    <div class="error">{error}</div>
  {/if}

  {#if result}
    <section class="results">
      <h2>Places you can drive or take a train to, then fly</h2>
      <p class="hint">Within 55–300 miles of your destination. Check Kayak for prices.</p>
      {#if Object.keys(result.driveThenFly).length > 0}
        <ul class="card-list">
          {#each Object.entries(result.driveThenFly) as [iata, dist]}
            <li class="card">
              <strong>{iata}</strong>
              <span>{dist.toFixed(1)} mi</span>
              <a href={kayakMultiCity(start, end, iata, startDate, endDate)} target="_blank" rel="noopener noreferrer">
                Check Kayak →
              </a>
            </li>
          {/each}
        </ul>
      {:else}
        <p class="empty">No nearby airports in database for this city.</p>
      {/if}

      <h2>Places you can fly to, then fly out of</h2>
      <p class="hint">Direct routes from your destination. Add a stopover and check prices on Kayak.</p>
      {#if Object.keys(result.flyThenFly).length > 0}
        <ul class="card-list">
          {#each Object.keys(result.flyThenFly) as iata}
            <li class="card">
              <strong>{iata}</strong>
              <a href={kayakMultiCity(start, end, iata, startDate, endDate)} target="_blank" rel="noopener noreferrer">
                Check Kayak →
              </a>
            </li>
          {/each}
        </ul>
      {:else}
        <p class="empty">No direct routes in database for this city/alliance.</p>
      {/if}
    </section>
  {/if}
</div>

<style>
  .container {
    max-width: 720px;
    margin: 0 auto;
    padding: 2rem 1.5rem;
  }
  .header { margin-bottom: 2rem; }
  .header h1 { font-size: 1.75rem; font-weight: 600; margin: 0 0 0.25rem 0; }
  .tagline { color: var(--muted); font-size: 0.95rem; margin: 0; }
  .search-form {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 12px;
    padding: 1.5rem;
    margin-bottom: 1.5rem;
  }
  .grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
    gap: 1rem;
    margin-bottom: 1rem;
  }
  label { display: flex; flex-direction: column; gap: 0.35rem; }
  label span { font-size: 0.8rem; color: var(--muted); }
  input, select {
    padding: 0.5rem 0.75rem;
    border: 1px solid var(--border);
    border-radius: 6px;
    background: var(--bg);
    color: var(--text);
    font-size: 0.95rem;
  }
  input:focus, select:focus { outline: none; border-color: var(--accent); }
  button {
    padding: 0.6rem 1.2rem;
    background: var(--accent);
    color: var(--bg);
    border: none;
    border-radius: 6px;
    font-weight: 600;
    cursor: pointer;
    font-size: 0.95rem;
  }
  button:hover:not(:disabled) { background: var(--accent-hover); }
  button:disabled { opacity: 0.6; cursor: not-allowed; }
  .error {
    color: #f85149;
    padding: 1rem;
    background: rgba(248, 81, 73, 0.1);
    border-radius: 8px;
    margin-bottom: 1rem;
  }
  .results h2 { font-size: 1.1rem; margin: 1.5rem 0 0.5rem 0; }
  .hint { font-size: 0.85rem; color: var(--muted); margin: 0 0 0.75rem 0; }
  .card-list {
    list-style: none;
    padding: 0;
    margin: 0 0 1.5rem 0;
    display: flex;
    flex-wrap: wrap;
    gap: 0.75rem;
  }
  .card {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.6rem 1rem;
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 8px;
    font-size: 0.9rem;
  }
  .card strong { min-width: 2.5rem; }
  .card a { color: var(--accent); text-decoration: none; font-size: 0.85rem; }
  .card a:hover { text-decoration: underline; }
  .empty { color: var(--muted); font-style: italic; margin: 0 0 1.5rem 0; }
</style>
