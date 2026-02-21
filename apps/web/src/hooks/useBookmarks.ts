import { useState, useEffect, useCallback } from 'react';
import { HackathonProps } from '@/temp/hackathons';

const STORAGE_KEY = 'hack-bookmarks';

// Global state to sync across all hook instances
let globalBookmarks: HackathonProps[] = [];
let listeners: Set<(bookmarks: HackathonProps[]) => void> = new Set();

// Notify all listeners of state change
const notifyListeners = () => {
    listeners.forEach(listener => listener([...globalBookmarks]));
};

// Get bookmarks from localStorage (client-side only)
const getStoredBookmarks = (): HackathonProps[] => {
    if (typeof window === 'undefined') return [];
    const stored = localStorage.getItem(STORAGE_KEY);
    if (stored) {
        try {
            return JSON.parse(stored);
        } catch (error) {
            console.error('Failed to parse bookmarks:', error);
        }
    }
    return [];
};

export function useBookmarks() {
    const [bookmarks, setBookmarks] = useState<HackathonProps[]>(globalBookmarks);
    const [mounted, setMounted] = useState(false);

    useEffect(() => {
        setMounted(true);
        // Initialize from localStorage on mount
        const stored = getStoredBookmarks();
        if (stored.length > 0) {
            globalBookmarks = stored;
            setBookmarks(stored);
        }

        // Subscribe to global state changes
        const handleChange = (newBookmarks: HackathonProps[]) => {
            setBookmarks(newBookmarks);
        };
        listeners.add(handleChange);

        // Listen for storage events from other tabs
        const handleStorage = (e: StorageEvent) => {
            if (e.key === STORAGE_KEY) {
                const newBookmarks = getStoredBookmarks();
                globalBookmarks = newBookmarks;
                notifyListeners();
            }
        };
        window.addEventListener('storage', handleStorage);

        return () => {
            listeners.delete(handleChange);
            window.removeEventListener('storage', handleStorage);
        };
    }, []);

    const toggleBookmark = useCallback((hackathon: HackathonProps) => {
        // Read current state from localStorage to ensure we have latest
        const currentBookmarks = getStoredBookmarks();

        const alreadyBookmarked = currentBookmarks.some((b) => b.id === hackathon.id);
        let updatedBookmarks;
        if (alreadyBookmarked) {
            updatedBookmarks = currentBookmarks.filter((b) => b.id !== hackathon.id);
        } else {
            updatedBookmarks = [...currentBookmarks, hackathon];
        }

        // Update global state
        globalBookmarks = updatedBookmarks;

        // Update localStorage
        localStorage.setItem(STORAGE_KEY, JSON.stringify(updatedBookmarks));

        // Notify all hook instances
        notifyListeners();
    }, []);

    const isBookmarked = useCallback((id: string) => {
        // Check global state first, then fall back to localStorage
        if (globalBookmarks.length > 0) {
            return globalBookmarks.some((b) => b.id === id);
        }
        const stored = getStoredBookmarks();
        return stored.some((b) => b.id === id);
    }, []);

    return { bookmarks, toggleBookmark, isBookmarked, mounted };
}
