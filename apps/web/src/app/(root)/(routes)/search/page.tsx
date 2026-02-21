
import { HackathonCard } from "@/components/hackathon-card";
import { HackathonProps } from "@/temp/hackathons";

const fetchHackathons = async() : Promise<HackathonProps[]> => {
  const response = await fetch("http://localhost:8080/api/hackathons", {
    cache: "no-store",
  });
  if (!response.ok) {
    throw new Error("Failed to fetch hackathons");
  }
  const data = await response.json();
  return data;
}


const SearchPage = async () => {

  const hackathons = await fetchHackathons();

  
  return (
    <div className="min-h-screen bg-background">
      <div className="max-w-196.5 mx-auto px-4 py-8">
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
