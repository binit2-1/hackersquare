"use client";

import { HackathonCard } from "@/components/hackathon-card";
import { useBookmarks } from "@/hooks/useBookmarks";

const BookmarksPage = () => {
  const { bookmarks, mounted } = useBookmarks();
  if (!mounted) {
    return null;
  }
  return (
    <div className="min-h-screen bg-background">
      <div className="max-w-196.5 mx-auto px-4 py-8">
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          {bookmarks.map((hackathon) => (
            <HackathonCard key={hackathon.id} hackathon={hackathon} />
          ))}
        </div>
      </div>
    </div>
  );
};

export default BookmarksPage;
