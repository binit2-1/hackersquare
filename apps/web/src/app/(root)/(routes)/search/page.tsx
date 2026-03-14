import { HackathonCard } from "@/components/hackathon-card";
import { NearMeEmptyState } from "@/components/near-me-empty-state";
import { headers } from "next/headers";
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
  origin: string,
  endpoint: string,
  queryString: string,
): Promise<SearchResponse> => {
  const url = queryString
    ? `${origin}${endpoint}?${queryString}`
    : `${origin}${endpoint}`;
  const response = await fetch(url, {
    cache: "no-store",
  });
  if (!response.ok) {
    throw new Error(`Failed to fetch hackathons (${response.status})`);
  }
  const data = await response.json();
  return data;
};

const buildEmptyResponse = (page: number, limit: number): SearchResponse => ({
  data: [],
  metadata: {
    totalRecords: 0,
    currentPage: page,
    limit,
    totalPages: 0,
  },
});

const normalizeGeoHeader = (value: string | null): string => {
  if (!value) return "";
  let decoded = value;
  try {
    decoded = decodeURIComponent(value);
  } catch {
    decoded = value;
  }
  decoded = decoded.trim();
  if (!decoded || decoded.toLowerCase() === "unknown") return "";
  return decoded;
};

const normalizeCountryHeader = (value: string | null): string => {
  const normalized = normalizeGeoHeader(value);
  if (!normalized) return "";

  // Vercel often sends ISO country codes (e.g. IN, US). Expand to full names for SQL matching.
  if (/^[A-Za-z]{2}$/.test(normalized)) {
    try {
      const regionName = new Intl.DisplayNames(["en"], { type: "region" }).of(
        normalized.toUpperCase(),
      );
      return regionName || "";
    } catch {
      return "";
    }
  }

  return normalized;
};

const SearchPage = async ({
  searchParams,
}: {
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>;
}) => {
  const requestHeaders = await headers();
  const forwardedHost = requestHeaders.get("x-forwarded-host");
  const host = forwardedHost || requestHeaders.get("host") || "localhost:3000";
  const proto =
    requestHeaders.get("x-forwarded-proto") ||
    (host.includes("localhost") ? "http" : "https");
  const origin = `${proto}://${host}`;

  const city =
    normalizeGeoHeader(requestHeaders.get("x-vercel-ip-city")) ||
    (host.includes("localhost")
      ? (process.env.NEARBY_FALLBACK_CITY || "")
      : "");
  const country =
    normalizeCountryHeader(requestHeaders.get("x-vercel-ip-country")) ||
    (host.includes("localhost")
      ? (process.env.NEARBY_FALLBACK_COUNTRY || "")
      : "");

  const resolvedSearchParams = await searchParams;
  const params = new URLSearchParams();
  for (const [key, value] of Object.entries(resolvedSearchParams)) {
    if (value) params.append(key, String(value));
  }

  const queryText = (resolvedSearchParams.q || "").toString().trim();
  const nearMeRequested = /\bnear me\b/i.test(queryText);
  const nearMeGeoUnavailable = nearMeRequested && !city && !country;
  const page = Number(params.get("page") || "1") || 1;
  const limit = Number(params.get("limit") || "20") || 20;

  let endpoint = "/api/v1/search";
  let responsePayload: SearchResponse;

  if (nearMeRequested) {
    endpoint = "/api/v1/hackathons/nearby";
    params.delete("q");
    if (city) params.set("city", city);
    if (country) params.set("country", country);

    if (!city && !country) {
      responsePayload = buildEmptyResponse(page, limit);
    } else {
      try {
        responsePayload = await fetchHackathons(origin, endpoint, params.toString());
      } catch {
        const fallbackParams = new URLSearchParams(params);
        fallbackParams.delete("city");
        fallbackParams.delete("country");

        // Nearby endpoint fallback should stay geo-scoped, never global.
        fallbackParams.set("q", [city, country].filter(Boolean).join(" "));

        endpoint = "/api/v1/search";
        responsePayload = await fetchHackathons(
          origin,
          endpoint,
          fallbackParams.toString(),
        );
      }
    }
  } else {
    responsePayload = await fetchHackathons(origin, endpoint, params.toString());
  }

  const { data: hackathons, metadata } = responsePayload;

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
        {nearMeGeoUnavailable ? (
          <NearMeEmptyState />
        ) : (
          <>
            <p className="mb-4 text-sm text-gray-500">
              Showing {hackathons.length} of {metadata.totalRecords} hackathons
            </p>

            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              {hackathons.map((hackathon) => (
                <HackathonCard key={hackathon.id} hackathon={hackathon} />
              ))}
            </div>
          </>
        )}
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
