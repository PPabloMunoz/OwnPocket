import { create } from "zustand";
import { persist } from "zustand/middleware";

export const useThemeStore = create<{
  theme: string;
  setTheme: (theme: string) => void;
  toggleTheme: () => void;
}>()(
  persist(
    (set, get) => ({
      theme: "dark",
      setTheme: (theme) => set({ theme }),
      toggleTheme: () =>
        set({ theme: get().theme === "dark" ? "light" : "dark" }),
    }),
    { name: "ownpocket-theme" },
  ),
);
