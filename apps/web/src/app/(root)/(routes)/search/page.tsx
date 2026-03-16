import { HackathonCard } from "@/components/hackathon-card";
import { NearMeEmptyState } from "@/components/near-me-empty-state";
import { SearchAIOverview } from "@/components/search-ai-overview";
import { SearchPagination } from "@/components/search-pagination";
import { headers } from "next/headers";
import { SearchResponse } from "@/models/hackathon";

const fetchHackathons = async (
  origin: string,
  queryString: string,
): Promise<SearchResponse> => {
  const url = queryString
    ? `${origin}/api/v1/search?${queryString}`
    : `${origin}/api/v1/search`;

  const response = await fetch(url, {
    cache: "no-store",
  });

  if (!response.ok) {
    throw new Error(`Failed to fetch hackathons (${response.status})`);
  }

  return response.json();
};

const normalizeGeoHeader = (value: string | null): string => {
  if (!value) return "";

  let decoded = value;
  try {
    decoded = decodeURIComponent(value);
  } catch {
    decoded = value;
  }

  const normalized = decoded.trim();
  if (!normalized || normalized.toLowerCase() === "unknown") return "";
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
    (host.includes("localhost") ? (process.env.NEARBY_FALLBACK_CITY || "") : "");

  const resolvedSearchParams = await searchParams;
  const params = new URLSearchParams();

  for (const [key, value] of Object.entries(resolvedSearchParams)) {
    if (value) params.append(key, String(value));
  }

  const queryText = (resolvedSearchParams.q || "").toString().trim();
  const nearMeRequested = /\bnear\s*me\b|\bnearby\b/i.test(queryText);

  if (nearMeRequested && !city) {
    return (
      <div className=" min-h-screen max-w-196.5">
        <div className="mx-auto px-4 py-8 max-w-7xl">
          <NearMeEmptyState />
        </div>
      </div>
    );
  }

  if (nearMeRequested && city) {
    params.delete("clientCity");
    params.append("clientCity", city);
  }

  const responsePayload = await fetchHackathons(origin, params.toString());

  const { data: hackathons, metadata } = responsePayload;
  const currentPage = metadata.currentPage;
  const totalPages = metadata.totalPages;

  return (
    <div className=" min-h-screen max-w-196.5">
      <div className="mx-auto px-4 py-8 max-w-7xl">
        {queryText && <SearchAIOverview query={queryText} />}
        <p className="mb-4 text-sm text-gray-500">
          Showing {hackathons.length} of {metadata.totalRecords} hackathons
        </p>

        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          {hackathons.map((hackathon) => (
            <HackathonCard key={hackathon.id} hackathon={hackathon} />
          ))}
        </div>
      </div>

      <SearchPagination currentPage={currentPage} totalPages={totalPages} />
    </div>
  );
};

export default SearchPage;
