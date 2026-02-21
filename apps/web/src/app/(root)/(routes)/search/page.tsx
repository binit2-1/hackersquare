import { HackathonCard } from "@/components/hackathon-card";
import { hackathons } from "@/temp/hackathons";

const SearchPage = () => {
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
