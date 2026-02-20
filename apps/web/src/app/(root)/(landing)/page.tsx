"use client";

import { SearchBar } from "@/components/search-bar";

const page = () => {
  return (
    <div className="w-screen h-screen flex justify-center items-center bg-background text-foreground">
      <SearchBar
        placeholders={[
          "Search for hackathons...",
          "Try 'Hackathons near me'",
          "Hackathons in ...",
        ]}
        interval={2500}
        onSubmit={(value: any) => console.log(value)}
      />
    </div>
  );
};

export default page;
