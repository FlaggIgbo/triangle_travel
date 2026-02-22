const TOKEN_KEY = 'travel_app_token';
const PHONE_KEY = 'travel_app_phone';

// Dev/sandbox: use sessionStorage (clears on tab close). Prod: use localStorage
// Dev = Vite dev server (import.meta.env.DEV) or localhost (built app served locally)
// Lazy: never access storage at module load (SSR has no window/localStorage)
function getStorage(): Storage | null {
  if (typeof window === 'undefined') return null;
  if (import.meta.env.DEV) return sessionStorage;
  if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') return sessionStorage;
  return localStorage;
}

export const DevPhone = '+15550000000';
export const DevOTP = '123456';
export const isDev =
  typeof window !== 'undefined' &&
  (import.meta.env?.DEV === true || window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1');

export function getToken(): string | null {
  const s = getStorage();
  return s ? s.getItem(TOKEN_KEY) : null;
}

export function setToken(token: string): void {
  const s = getStorage();
  if (s) s.setItem(TOKEN_KEY, token);
}

export function clearToken(): void {
  const s = getStorage();
  if (s) {
    s.removeItem(TOKEN_KEY);
    s.removeItem(PHONE_KEY);
  }
}

export function setPhone(phone: string): void {
  const s = getStorage();
  if (s) s.setItem(PHONE_KEY, phone);
}

export function getPhone(): string | null {
  const s = getStorage();
  return s ? s.getItem(PHONE_KEY) : null;
}

export function isLoggedIn(): boolean {
  return !!getToken();
}
