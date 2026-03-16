"use client";

import * as React from "react";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";

type SearchPaginationProps = {
  currentPage: number;
  totalPages: number;
};

export function SearchPagination({ currentPage, totalPages }: SearchPaginationProps) {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  const buildPageUrl = React.useCallback(
    (page: number) => {
      const currentParams = new URLSearchParams(searchParams.toString());
      currentParams.set("page", page.toString());
      return `${pathname}?${currentParams.toString()}`;
    },
    [pathname, searchParams],
  );

  const handlePageClick = React.useCallback(
    (event: React.MouseEvent<HTMLAnchorElement>, page: number) => {
      event.preventDefault();
      const currentParams = new URLSearchParams(searchParams.toString());
      currentParams.set("page", page.toString());
      router.push(`${pathname}?${currentParams.toString()}`);
    },
    [pathname, router, searchParams],
  );

  const getPageNumbers = React.useCallback(() => {
    const pages: (number | "ellipsis")[] = [];
    const maxVisible = 7;

    if (totalPages <= maxVisible) {
      for (let i = 1; i <= totalPages; i++) pages.push(i);
      return pages;
    }

    pages.push(1);

    if (currentPage > 3) pages.push("ellipsis");

    const start = Math.max(2, currentPage - 1);
    const end = Math.min(totalPages - 1, currentPage + 1);
    for (let i = start; i <= end; i++) pages.push(i);

    if (currentPage < totalPages - 2) pages.push("ellipsis");

    pages.push(totalPages);
    return pages;
  }, [currentPage, totalPages]);

  if (totalPages <= 1) return null;

  return (
    <Pagination className="mb-8">
      <PaginationContent>
        <PaginationItem>
          <PaginationPrevious
            href={currentPage > 1 ? buildPageUrl(currentPage - 1) : "#"}
            onClick={(event) => {
              if (currentPage <= 1) {
                event.preventDefault();
                return;
              }
              handlePageClick(event, currentPage - 1);
            }}
            className={currentPage === 1 ? "pointer-events-none opacity-50" : ""}
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
                isActive={currentPage === page}
                onClick={(event) => handlePageClick(event, page)}
              >
                {page}
              </PaginationLink>
            </PaginationItem>
          ),
        )}

        <PaginationItem>
          <PaginationNext
            href={currentPage < totalPages ? buildPageUrl(currentPage + 1) : "#"}
            onClick={(event) => {
              if (currentPage >= totalPages) {
                event.preventDefault();
                return;
              }
              handlePageClick(event, currentPage + 1);
            }}
            className={currentPage === totalPages ? "pointer-events-none opacity-50" : ""}
          />
        </PaginationItem>
      </PaginationContent>
    </Pagination>
  );
}
