import { useState, useEffect } from 'react';
import { HackathonProps } from '@/temp/hackathons';

export function useBookmarks() {
    const [bookmarks, setBookmarks] = useState<HackathonProps[]>([]);
    const [mounted, setMounted] = useState(false);

    useEffect(() => {
        setMounted(true);
        const storedBookmarks = localStorage.getItem('hack-bookmarks');
        if (storedBookmarks) {
            try {
                setBookmarks(JSON.parse(storedBookmarks));
            } catch (error) {
                if (error instanceof Error) {
                    console.error('Failed to parse bookmarks from localStorage:', error.message);
                } else {
                    console.error('An unknown error occurred while parsing bookmarks from localStorage');
                }
            }
        }
    }, []);

    const toggleBookmark = (hackathon: HackathonProps) => {
        const alreadyBookmarked = bookmarks.some((b) => b.title === hackathon.title);
        let updatedBookmarks;
        if (alreadyBookmarked) {
            updatedBookmarks = bookmarks.filter((b) => b.title != hackathon.title);
        } else {
            updatedBookmarks = [...bookmarks, hackathon];
        }
        setBookmarks(updatedBookmarks);
        localStorage.setItem('hack-bookmarks', JSON.stringify(updatedBookmarks));
    };
    const isBookmarked = (title: string) => {
        return bookmarks.some((b) => b.title === title);
    };

    return { bookmarks, toggleBookmark, isBookmarked, mounted };

}