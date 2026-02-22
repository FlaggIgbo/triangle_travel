const TOKEN_KEY = 'travel_app_token';
const PHONE_KEY = 'travel_app_phone';

// Dev/sandbox: use sessionStorage (clears on tab close). Prod: use localStorage
// Dev = Vite dev server (import.meta.env.DEV) or localhost (built app served locally)
function getStorage(): Storage {
  if (typeof window === 'undefined') return localStorage;
  if (import.meta.env.DEV) return sessionStorage;
  if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') return sessionStorage;
  return localStorage;
}
const storage = getStorage();

export const DevPhone = '+15550000000';
export const DevOTP = '123456';
export const isDev =
  typeof window !== 'undefined' &&
  (import.meta.env?.DEV === true || window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1');

export function getToken(): string | null {
  if (typeof window === 'undefined') return null;
  return storage.getItem(TOKEN_KEY);
}

export function setToken(token: string): void {
  storage.setItem(TOKEN_KEY, token);
}

export function clearToken(): void {
  storage.removeItem(TOKEN_KEY);
  storage.removeItem(PHONE_KEY);
}

export function setPhone(phone: string): void {
  storage.setItem(PHONE_KEY, phone);
}

export function getPhone(): string | null {
  return storage.getItem(PHONE_KEY);
}

export function isLoggedIn(): boolean {
  return !!getToken();
}
