import { createContext, useContext, useEffect, useMemo, useState } from 'react';

const ThemeContext = createContext(null);
const storageKey = 'inbox-inspector-theme';

export function ThemeProvider({ children }) {
  const [preference, setPreferenceState] = useState(() => localStorage.getItem(storageKey) || 'auto');
  const [resolvedTheme, setResolvedTheme] = useState(resolveTheme(preference));

  useEffect(() => {
    function refreshTheme() {
      const next = resolveTheme(preference);
      setResolvedTheme(next);
      document.documentElement.dataset.theme = next;
    }
    refreshTheme();
    const timer = window.setInterval(refreshTheme, 60 * 1000);
    return () => window.clearInterval(timer);
  }, [preference]);

  function setPreference(nextPreference) {
    setPreferenceState(nextPreference);
    localStorage.setItem(storageKey, nextPreference);
  }

  const value = useMemo(() => ({ preference, resolvedTheme, setPreference }), [preference, resolvedTheme]);
  return <ThemeContext.Provider value={value}>{children}</ThemeContext.Provider>;
}

export function useTheme() {
  return useContext(ThemeContext);
}

function resolveTheme(preference) {
  if (preference === 'dark' || preference === 'light') {
    return preference;
  }
  const hour = new Date().getHours();
  return hour >= 21 || hour < 7 ? 'dark' : 'light';
}

