import { useState, useEffect, useCallback } from 'react';
import { HackathonProps } from '@/models/hackathon';
import { useAuth } from '@/context/AuthContext';

const STORAGE_KEY = 'hack-bookmarks';
const API_BASE_URL = '/api';

// Global state to sync across all hook instances
let globalBookmarks: HackathonProps[] = [];
let listeners: Set<(bookmarks: HackathonProps[]) => void> = new Set();

// Notify all listeners of state change
const notifyListeners = () => {
    listeners.forEach(listener => listener([...globalBookmarks]));
};

const setGlobalBookmarks = (bookmarks: HackathonProps[]) => {
    globalBookmarks = bookmarks;
    notifyListeners();
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

type BookmarkApiResponse = Omit<HackathonProps, 'isBookmarked' | 'prize_usd'> & {
    prize_usd?: number | null;
    isBookmarked?: boolean;
};

const normalizeBookmark = (hackathon: BookmarkApiResponse): HackathonProps => ({
    ...hackathon,
    prize_usd: hackathon.prize_usd ?? 0,
    isBookmarked: true,
});

export function useBookmarks() {
    const [bookmarks, setBookmarks] = useState<HackathonProps[]>(globalBookmarks);
    const [mounted, setMounted] = useState(false);
    const { user, isLoading } = useAuth();

    useEffect(() => {
        setMounted(true);

        // Subscribe to global state changes
        const handleChange = (newBookmarks: HackathonProps[]) => {
            setBookmarks(newBookmarks);
        };
        listeners.add(handleChange);

        return () => {
            listeners.delete(handleChange);
        };
    }, []);

    useEffect(() => {
        if (!mounted || isLoading) return;

        if (user) {
            const fetchBookmarks = async () => {
                try {
                    const response = await fetch(`${API_BASE_URL}/v1/bookmarks`, {
                        method: 'GET',
                        credentials: 'include',
                    });
                    if (!response.ok) {
                        const errText = await response.text();
                        console.error('Bookmark fetch failed:', response.status, errText);
                        throw new Error(`Failed to fetch bookmarks: ${response.status}`);
                    }
                    const data: BookmarkApiResponse[] = await response.json();
                    setGlobalBookmarks((data ?? []).map(normalizeBookmark));
                } catch (error) {
                    console.error('Failed to fetch bookmarks:', error);
                }
            };

            fetchBookmarks();
        } else {
            const stored = getStoredBookmarks();
            setGlobalBookmarks(stored);
        }
    }, [mounted, isLoading, user]);

    useEffect(() => {
        if (!mounted || user) return;

        const handleStorage = (e: StorageEvent) => {
            if (e.key === STORAGE_KEY) {
                const newBookmarks = getStoredBookmarks();
                setGlobalBookmarks(newBookmarks);
            }
        };
        window.addEventListener('storage', handleStorage);

        return () => {
            window.removeEventListener('storage', handleStorage);
        };
    }, [mounted, user]);

    const toggleBookmark = useCallback(async (hackathon: HackathonProps) => {
        if (isLoading) return;

        const alreadyBookmarked = globalBookmarks.some((b) => b.id === hackathon.id);

        if (user) {
            try {
                const endpoint = alreadyBookmarked
                    ? `${API_BASE_URL}/v1/bookmarks?hackathon_id=${encodeURIComponent(hackathon.id)}`
                    : `${API_BASE_URL}/v1/bookmarks`;
                const response = await fetch(endpoint, {
                    method: alreadyBookmarked ? 'DELETE' : 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    credentials: 'include',
                    body: alreadyBookmarked ? undefined : JSON.stringify({ hackathon_id: hackathon.id }),
                });

                if (!response.ok) {
                    const errText = await response.text();
                    console.error('Bookmark update failed:', response.status, errText);
                    throw new Error(`Failed to update bookmark: ${response.status}`);
                }

                const nextBookmarks = alreadyBookmarked
                    ? globalBookmarks.filter((b) => b.id !== hackathon.id)
                    : [...globalBookmarks, { ...hackathon, isBookmarked: true }];

                setGlobalBookmarks(nextBookmarks);
            } catch (error) {
                console.error('Failed to update bookmark:', error);
            }
            return;
        }

        // Guest flow: read current state from localStorage to ensure latest
        const currentBookmarks = getStoredBookmarks();
        const guestAlreadyBookmarked = currentBookmarks.some((b) => b.id === hackathon.id);
        const updatedBookmarks = guestAlreadyBookmarked
            ? currentBookmarks.filter((b) => b.id !== hackathon.id)
            : [...currentBookmarks, hackathon];

        localStorage.setItem(STORAGE_KEY, JSON.stringify(updatedBookmarks));
        setGlobalBookmarks(updatedBookmarks);
    }, [isLoading, user]);

    const isBookmarked = useCallback((id: string) => {
        return globalBookmarks.some((b) => b.id === id);
    }, []);

    return { bookmarks, toggleBookmark, isBookmarked, mounted };
}
