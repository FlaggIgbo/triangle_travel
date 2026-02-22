<script lang="ts">
  import { onMount } from 'svelte';
  import { getToken, setToken, setPhone, clearToken, isLoggedIn, DevPhone, DevOTP, isDev } from '$lib/auth';

  interface Flight {
    id: number;
    airline: string;
    flight_number: string;
    from_iata: string;
    to_iata: string;
    departure_date: string;
    departure_time?: string;
    confirmation?: string;
  }

  let loggedIn = $state(false);
  let phone = $state('');
  let otp = $state('');
  let step = $state<'phone' | 'otp'>('phone');
  let authError = $state('');
  let authLoading = $state(false);

  let flights = $state<Flight[]>([]);
  let showAddForm = $state(false);
  let form = $state({
    airline: '',
    flight_number: '',
    from_iata: '',
    to_iata: '',
    departure_date: '',
    departure_time: '',
    confirmation: ''
  });
  let formError = $state('');
  let formLoading = $state(false);

  function headers(): HeadersInit {
    const token = getToken();
    return {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {})
    };
  }

  onMount(() => {
    loggedIn = isLoggedIn();
    if (loggedIn) loadFlights();
  });

  async function requestOtp() {
    const p = phone.replace(/\D/g, '');
    if (p.length !== 11 || !p.startsWith('1')) {
      authError = 'Enter a valid US phone number (e.g. 5551234567)';
      return;
    }
    authLoading = true;
    authError = '';
    try {
      const res = await fetch('/api/auth/send-otp', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ phone: `+1${p}` })
      });
      const data = await res.json().catch(() => ({}));
      if (!res.ok) throw new Error(data.error || 'Failed to send OTP');
      step = 'otp';
      otp = '';
    } catch (e) {
      authError = e instanceof Error ? e.message : 'Failed to send OTP';
    } finally {
      authLoading = false;
    }
  }

  async function verifyOtp() {
    if (otp.length !== 6) {
      authError = 'Enter the 6-digit code';
      return;
    }
    authLoading = true;
    authError = '';
    try {
      const p = phone.replace(/\D/g, '');
      const res = await fetch('/api/auth/verify-otp', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ phone: `+1${p}`, code: otp })
      });
      const data = await res.json().catch(() => ({}));
      if (!res.ok) throw new Error(data.error || 'Invalid code');
      setToken(data.token);
      setPhone(`+1${p}`);
      loggedIn = true;
      step = 'phone';
      loadFlights();
    } catch (e) {
      authError = e instanceof Error ? e.message : 'Invalid code';
    } finally {
      authLoading = false;
    }
  }

  function logout() {
    clearToken();
    loggedIn = false;
    flights = [];
    showLogin = true;
  }

  async function loadFlights() {
    const token = getToken();
    if (!token) return;
    try {
      const res = await fetch('/api/flights', { headers: headers() });
      if (res.ok) flights = await res.json();
    } catch {
      flights = [];
    }
  }

  async function addFlight() {
    if (!form.airline || !form.flight_number || !form.from_iata || !form.to_iata || !form.departure_date) {
      formError = 'Fill in required fields';
      return;
    }
    formLoading = true;
    formError = '';
    try {
      const res = await fetch('/api/flights', {
        method: 'POST',
        headers: headers(),
        body: JSON.stringify(form)
      });
      const data = await res.json().catch(() => ({}));
      if (!res.ok) throw new Error(data.error || 'Failed to add flight');
      flights = [...flights, data];
      showAddForm = false;
      form = { airline: '', flight_number: '', from_iata: '', to_iata: '', departure_date: '', departure_time: '', confirmation: '' };
    } catch (e) {
      formError = e instanceof Error ? e.message : 'Failed to add flight';
    } finally {
      formLoading = false;
    }
  }

  async function deleteFlight(id: number) {
    try {
      const res = await fetch(`/api/flights/${id}`, { method: 'DELETE', headers: headers() });
      if (res.ok) flights = flights.filter((f) => f.id !== id);
    } catch {
      /* ignore */
    }
  }
</script>

<div class="container">
  {#if !loggedIn}
    <div class="auth-gate">
      <h1>My Flights</h1>
      <p class="tagline">Log in to add and view your booked flights</p>

      <div class="auth-card">
          {#if step === 'phone'}
            <h2>Sign in with US phone</h2>
            {#if isDev}
              <p class="dev-hint">Dev mode: use {DevPhone} and OTP {DevOTP}</p>
            {/if}
            <p class="hint">We'll send a one-time code to your number</p>
            <input
              type="tel"
              bind:value={phone}
              placeholder={isDev ? DevPhone : '(555) 123-4567'}
              maxlength="14"
            />
            <button onclick={requestOtp} disabled={authLoading}>
              {authLoading ? 'Sending…' : 'Send code'}
            </button>
          {:else}
            <h2>Enter code</h2>
            {#if isDev}
              <p class="dev-hint">Dev: use {DevOTP}</p>
            {/if}
            <p class="hint">Check your phone for the 6-digit code</p>
            <input
              type="text"
              bind:value={otp}
              placeholder={isDev ? DevOTP : '000000'}
              maxlength="6"
              pattern="[0-9]*"
              inputmode="numeric"
            />
            <button onclick={verifyOtp} disabled={authLoading}>
              {authLoading ? 'Verifying…' : 'Verify'}
            </button>
            <button class="secondary" onclick={() => { step = 'phone'; authError = ''; }}>
              Change number
            </button>
          {/if}
          {#if authError}
            <p class="error">{authError}</p>
          {/if}
        </div>
    </div>
  {:else}
    <header class="header">
      <div>
        <h1>My Flights</h1>
        <p class="tagline">Your booked flights</p>
      </div>
      <button class="logout" onclick={logout}>Log out</button>
    </header>

    <button class="add-btn" onclick={() => (showAddForm = !showAddForm)}>
      {showAddForm ? 'Cancel' : '+ Add flight'}
    </button>

    {#if showAddForm}
      <form class="flight-form" onsubmit={(e) => { e.preventDefault(); addFlight(); }}>
        <h3>Add flight</h3>
        <div class="grid">
          <label>
            <span>Airline *</span>
            <input type="text" bind:value={form.airline} placeholder="e.g. United" />
          </label>
          <label>
            <span>Flight number *</span>
            <input type="text" bind:value={form.flight_number} placeholder="e.g. UA123" />
          </label>
          <label>
            <span>From (IATA) *</span>
            <input type="text" bind:value={form.from_iata} placeholder="SFO" maxlength="4" />
          </label>
          <label>
            <span>To (IATA) *</span>
            <input type="text" bind:value={form.to_iata} placeholder="JFK" maxlength="4" />
          </label>
          <label>
            <span>Date *</span>
            <input type="date" bind:value={form.departure_date} />
          </label>
          <label>
            <span>Time</span>
            <input type="time" bind:value={form.departure_time} />
          </label>
          <label>
            <span>Confirmation</span>
            <input type="text" bind:value={form.confirmation} placeholder="Booking ref" />
          </label>
        </div>
        {#if formError}<p class="error">{formError}</p>{/if}
        <button type="submit" disabled={formLoading}>{formLoading ? 'Adding…' : 'Add flight'}</button>
      </form>
    {/if}

    <section class="flights-list">
      {#if flights.length === 0}
        <p class="empty">No flights yet. Add one above.</p>
      {:else}
        {#each flights as f}
          <div class="flight-card">
            <div class="flight-main">
              <strong>{f.airline} {f.flight_number}</strong>
              <span>{f.from_iata} → {f.to_iata}</span>
              <span>{f.departure_date}{f.departure_time ? ` ${f.departure_time}` : ''}</span>
              {#if f.confirmation}<span class="conf">Ref: {f.confirmation}</span>{/if}
            </div>
            <button class="delete" onclick={() => deleteFlight(f.id)}>Remove</button>
          </div>
        {/each}
      {/if}
    </section>
  {/if}
</div>

<style>
  .container { max-width: 640px; margin: 0 auto; padding: 2rem 1.5rem; }
  .auth-gate { text-align: center; padding: 3rem 1rem; }
  .auth-gate h1 { font-size: 1.5rem; margin: 0 0 0.5rem 0; }
  .tagline { color: var(--muted); font-size: 0.9rem; margin: 0 0 1.5rem 0; }
  .auth-card {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 12px;
    padding: 2rem;
    max-width: 320px;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  .auth-card h2 { font-size: 1.1rem; margin: 0; }
  .dev-hint { font-size: 0.85rem; color: var(--accent); margin: 0 0 0.25rem 0; }
  .hint { font-size: 0.85rem; color: var(--muted); margin: 0; }
  .auth-card input {
    padding: 0.6rem 1rem;
    border: 1px solid var(--border);
    border-radius: 6px;
    background: var(--bg);
    color: var(--text);
    font-size: 1rem;
  }
  .auth-card input:focus { outline: none; border-color: var(--accent); }
  .auth-card button {
    padding: 0.6rem 1.2rem;
    background: var(--accent);
    color: var(--bg);
    border: none;
    border-radius: 6px;
    font-weight: 600;
    cursor: pointer;
  }
  .auth-card button.secondary { background: transparent; color: var(--muted); }
  .auth-card button.secondary:hover { color: var(--text); }
  .error { color: #f85149; font-size: 0.9rem; margin: 0; }
  .header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 1.5rem; }
  .header h1 { font-size: 1.5rem; margin: 0 0 0.25rem 0; }
  .logout { padding: 0.4rem 0.8rem; background: transparent; color: var(--muted); border: 1px solid var(--border); border-radius: 6px; cursor: pointer; font-size: 0.9rem; }
  .logout:hover { color: var(--text); }
  .add-btn {
    padding: 0.6rem 1.2rem;
    background: var(--accent);
    color: var(--bg);
    border: none;
    border-radius: 6px;
    font-weight: 600;
    cursor: pointer;
    margin-bottom: 1rem;
  }
  .add-btn:hover { background: var(--accent-hover); }
  .flight-form {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 12px;
    padding: 1.5rem;
    margin-bottom: 1.5rem;
  }
  .flight-form h3 { margin: 0 0 1rem 0; font-size: 1rem; }
  .grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(140px, 1fr)); gap: 1rem; margin-bottom: 1rem; }
  label { display: flex; flex-direction: column; gap: 0.35rem; }
  label span { font-size: 0.8rem; color: var(--muted); }
  .flight-form input {
    padding: 0.5rem 0.75rem;
    border: 1px solid var(--border);
    border-radius: 6px;
    background: var(--bg);
    color: var(--text);
  }
  .flight-form input:focus { outline: none; border-color: var(--accent); }
  .flight-form button { padding: 0.5rem 1rem; background: var(--accent); color: var(--bg); border: none; border-radius: 6px; font-weight: 600; cursor: pointer; }
  .flights-list { display: flex; flex-direction: column; gap: 0.75rem; }
  .empty { color: var(--muted); font-style: italic; padding: 2rem; }
  .flight-card {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem;
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 8px;
  }
  .flight-main { display: flex; flex-wrap: wrap; gap: 0.5rem 1rem; align-items: center; }
  .flight-main strong { margin-right: 0.5rem; }
  .conf { font-size: 0.85rem; color: var(--muted); }
  .flight-card .delete { padding: 0.35rem 0.7rem; background: transparent; color: #f85149; border: 1px solid #f85149; border-radius: 4px; cursor: pointer; font-size: 0.85rem; }
  .flight-card .delete:hover { background: rgba(248, 81, 73, 0.1); }
</style>
