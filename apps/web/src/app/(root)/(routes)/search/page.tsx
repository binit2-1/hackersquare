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
import { SearchResponse } from "@/models/hackathon";

const fetchHackathons = async (
  queryString: string,
): Promise<SearchResponse> => {
  const apiBaseUrl =
    process.env.NEXT_PUBLIC_API_URL || "https://hackersquare-api.up.railway.app";
  const url = queryString
    ? `${apiBaseUrl}/v1/search?${queryString}`
    : `${apiBaseUrl}/v1/search`;
  const response = await fetch(url, {
    cache: "no-store",
  });
  if (!response.ok) {
    throw new Error("Failed to fetch hackathons");
  }
  const data = await response.json();
  return data;
};

const SearchPage = async ({
  searchParams,
}: {
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>;
}) => {
  const resolvedSearchParams = await searchParams;
  const params = new URLSearchParams();
  for (const [key, value] of Object.entries(resolvedSearchParams)) {
    if (value) params.append(key, String(value));
  }

  const { data: hackathons, metadata } = await fetchHackathons(
    params.toString(),
  );

  const currentPage = metadata.currentPage;
  const totalPages = metadata.totalPages;

  // Helper function to generate pagination URLs
  const buildPageUrl = (page: number) => {
    const newParams = new URLSearchParams(params);
    newParams.set("page", page.toString());
    return `/search?${newParams.toString()}`;
  };

  // Generate page numbers to display
  const getPageNumbers = () => {
    const pages: (number | "ellipsis")[] = [];
    const maxVisible = 7;

    if (totalPages <= maxVisible) {
      // Show all pages
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i);
      }
    } else {
      // Always show first page
      pages.push(1);

      if (currentPage > 3) {
        pages.push("ellipsis");
      }

      // Show pages around current page
      const start = Math.max(2, currentPage - 1);
      const end = Math.min(totalPages - 1, currentPage + 1);

      for (let i = start; i <= end; i++) {
        pages.push(i);
      }

      if (currentPage < totalPages - 2) {
        pages.push("ellipsis");
      }

      // Always show last page
      pages.push(totalPages);
    }

    return pages;
  };

  return (
    <div className=" min-h-screen max-w-196.5">
      <div className="mx-auto px-4 py-8 max-w-7xl">
        <p className="mb-4 text-sm text-gray-500">
          Showing {hackathons.length} of {metadata.totalRecords} hackathons
        </p>

        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          {hackathons.map((hackathon) => (
            <HackathonCard key={hackathon.id} hackathon={hackathon} />
          ))}
        </div>
      </div>

      {totalPages > 1 && (
        <Pagination className="mb-8">
          <PaginationContent>
            <PaginationItem>
              <PaginationPrevious
                href={currentPage > 1 ? buildPageUrl(currentPage - 1) : "#"}
                className={
                  currentPage === 1 ? "pointer-events-none opacity-50" : ""
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
                    isActive={currentPage === page}
                  >
                    {page}
                  </PaginationLink>
                </PaginationItem>
              )
            )}

            <PaginationItem>
              <PaginationNext
                href={
                  currentPage < totalPages
                    ? buildPageUrl(currentPage + 1)
                    : "#"
                }
                className={
                  currentPage === totalPages
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

export default SearchPage;
