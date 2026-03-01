import { HackathonCard } from "@/components/hackathon-card";
import { HackathonProps, SearchResponse } from "@/models/hackathon";

const fetchHackathons = async (
  queryString: string,
): Promise<SearchResponse> => {
  const response = await fetch("http://localhost:8080/v1/search", {
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
  searchParams: { [key: string]: string | string[] | undefined };
}) => {
  const params = new URLSearchParams();
  for (const [key, value] of Object.entries(searchParams)) {
    if (value) params.append(key, String(value));
  }

  const { data: hackathons, metadata } = await fetchHackathons(
    params.toString(),
  );

  return (
    <div className="min-h-screen bg-background">
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
    </div>
  );
};

export default SearchPage;
