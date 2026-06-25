"use client";

import * as React from "react";

const uuidStorageKey = "jolly:user_uuid";
const tokenStorageKey = "jolly:token";
const sessionEventName = "jolly:session:changed";

export function getStoredUserUUID(): string | null {
  if (typeof window === "undefined") {
    return null;
  }
  const value = window.localStorage.getItem(uuidStorageKey);
  return value && value.trim() !== "" ? value : null;
}

export function getStoredToken(): string | null {
  if (typeof window === "undefined") {
    return null;
  }
  const value = window.localStorage.getItem(tokenStorageKey);
  return value && value.trim() !== "" ? value : null;
}

export function storeSession(userUUID: string, token: string) {
  if (typeof window === "undefined") {
    return;
  }
  window.localStorage.setItem(uuidStorageKey, userUUID);
  window.localStorage.setItem(tokenStorageKey, token);
  window.dispatchEvent(new Event(sessionEventName));
}

export function clearSession() {
  if (typeof window === "undefined") {
    return;
  }
  window.localStorage.removeItem(uuidStorageKey);
  window.localStorage.removeItem(tokenStorageKey);
  window.dispatchEvent(new Event(sessionEventName));
}

export function useUserSession() {
  const userUUID = React.useSyncExternalStore(
    (onStoreChange) => {
      if (typeof window === "undefined") {
        return () => {};
      }

      const onStorage = (e: StorageEvent) => {
        if (e.key === uuidStorageKey || e.key === tokenStorageKey) {
          onStoreChange();
        }
      };
      const onSessionChanged = () => onStoreChange();

      window.addEventListener("storage", onStorage);
      window.addEventListener(sessionEventName, onSessionChanged);

      return () => {
        window.removeEventListener("storage", onStorage);
        window.removeEventListener(sessionEventName, onSessionChanged);
      };
    },
    () => getStoredUserUUID(),
    () => null,
  );

  const token = React.useSyncExternalStore(
    (onStoreChange) => {
      if (typeof window === "undefined") {
        return () => {};
      }

      const onStorage = (e: StorageEvent) => {
        if (e.key === uuidStorageKey || e.key === tokenStorageKey) {
          onStoreChange();
        }
      };
      const onSessionChanged = () => onStoreChange();

      window.addEventListener("storage", onStorage);
      window.addEventListener(sessionEventName, onSessionChanged);

      return () => {
        window.removeEventListener("storage", onStorage);
        window.removeEventListener(sessionEventName, onSessionChanged);
      };
    },
    () => getStoredToken(),
    () => null,
  );

  const signin = React.useCallback((uuid: string, token: string) => {
    storeSession(uuid, token);
  }, []);

  const signout = React.useCallback(() => {
    clearSession();
  }, []);

  return { userUUID, token, signin, signout };
}
