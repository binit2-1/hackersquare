"use client";

import { Suspense, useMemo } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { HackathonCard } from "@/components/hackathon-card";
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { useBookmarks } from "@/hooks/useBookmarks";

const ITEMS_PER_PAGE = 10;

const BookmarksPageContent = () => {
  const { bookmarks, mounted } = useBookmarks();
  const searchParams = useSearchParams();
  const router = useRouter();

  const currentPage = Number(searchParams.get("page") || "1");
  const totalPages = Math.max(1, Math.ceil(bookmarks.length / ITEMS_PER_PAGE));
  const safePage = Math.min(Math.max(1, currentPage), totalPages);

  const paginatedBookmarks = useMemo(() => {
    const start = (safePage - 1) * ITEMS_PER_PAGE;
    return bookmarks.slice(start, start + ITEMS_PER_PAGE);
  }, [bookmarks, safePage]);

  const buildPageUrl = (page: number) => `/bookmarks?page=${page}`;

  const getPageNumbers = () => {
    const pages: (number | "ellipsis")[] = [];
    const maxVisible = 7;

    if (totalPages <= maxVisible) {
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i);
      }
    } else {
      pages.push(1);

      if (safePage > 3) {
        pages.push("ellipsis");
      }

      const start = Math.max(2, safePage - 1);
      const end = Math.min(totalPages - 1, safePage + 1);

      for (let i = start; i <= end; i++) {
        pages.push(i);
      }

      if (safePage < totalPages - 2) {
        pages.push("ellipsis");
      }

      pages.push(totalPages);
    }

    return pages;
  };

  if (!mounted) {
    return null;
  }

  return (
    <div className="min-h-screen bg-background max-w-196.5">
      <div className="mx-auto px-4 py-8 max-w-7xl">
        <p className="mb-4 text-sm text-gray-500">
          Showing {paginatedBookmarks.length} of {bookmarks.length} bookmarks
        </p>

        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          {paginatedBookmarks.map((hackathon) => (
            <HackathonCard key={hackathon.id} hackathon={hackathon} />
          ))}
        </div>
      </div>

      {totalPages > 1 && (
        <Pagination className="mb-8">
          <PaginationContent>
            <PaginationItem>
              <PaginationPrevious
                href={safePage > 1 ? buildPageUrl(safePage - 1) : "#"}
                className={
                  safePage === 1 ? "pointer-events-none opacity-50" : ""
                }
              />
            </PaginationItem>

            {getPageNumbers().map((page, index) =>
              page === "ellipsis" ? (
                <PaginationItem key={`ellipsis-${index}`}>
                  <PaginationEllipsis />
                </PaginationItem>
              ) : (
                <PaginationItem key={page}>
                  <PaginationLink
                    href={buildPageUrl(page)}
                    isActive={safePage === page}
                  >
                    {page}
                  </PaginationLink>
                </PaginationItem>
              ),
            )}

            <PaginationItem>
              <PaginationNext
                href={
                  safePage < totalPages ? buildPageUrl(safePage + 1) : "#"
                }
                className={
                  safePage === totalPages
                    ? "pointer-events-none opacity-50"
                    : ""
                }
              />
            </PaginationItem>
          </PaginationContent>
        </Pagination>
      )}
    </div>
  );
};

export default function BookmarksPage() {
  return (
    <Suspense fallback={null}>
      <BookmarksPageContent />
    </Suspense>
  );
}
