"use client";

import { SearchBar } from "@/components/search-bar";
import { useRouter } from "next/navigation";

const page = () => {
  const router = useRouter();
  const handleSubmit = (value?: string | React.SyntheticEvent) => {
    console.log("Submitting search...", value);
    router.push("/search");
  };
  return (
    <div className="w-screen h-screen flex justify-center items-center bg-background text-foreground">
      <SearchBar
        placeholders={[
          "Search for hackathons...",
          "Try 'Hackathons near me'",
          "Hackathons in ...",
        ]}
        interval={2500}
        onSubmit={handleSubmit}
      />
    </div>
  );
};

export default page;
