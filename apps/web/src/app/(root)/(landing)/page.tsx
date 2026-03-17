"use client";

import { useState } from "react";
import { ArrowRightIcon, ChevronDownIcon } from "lucide-react";
import { AnimatePresence, motion } from "framer-motion";
import { useRouter } from "next/navigation";

import { HackathonCard } from "@/components/hackathon-card";
import { SearchBar } from "@/components/search-bar";
import {
  CraftButton,
  CraftButtonIcon,
  CraftButtonLabel,
} from "@/components/ui/craft-button";
import type { HackathonProps } from "@/models/hackathon";

type FeedPhase = "search" | "transitioning" | "recommendations";

const LandingPage = () => {
  const router = useRouter();
  const [phase, setPhase] = useState<FeedPhase>("search");
  const [isLoadingRecommendations, setIsLoadingRecommendations] = useState(false);
  const [recommendationsError, setRecommendationsError] = useState<string | null>(
    null,
  );
  const [recommendations, setRecommendations] = useState<HackathonProps[]>([]);

  const handleSearchSubmit = (value?: string) => {
    const params = new URLSearchParams();

    if (typeof value === "string") {
      const query = value.trim();
      if (query) params.set("q", query);
    }

    params.delete("page");
    router.push(`/search?${params.toString()}`);
  };

  const fetchRecommendations = async () => {
    const response = await fetch("/api/v1/recommendations", {
      method: "GET",
      credentials: "include",
      cache: "no-store",
    });

    if (response.status === 401) {
      throw new Error("Sign in to see your recommendations.");
    }

    if (!response.ok) {
      throw new Error("Unable to load recommendations right now.");
    }

    const payload = (await response.json()) as { data?: HackathonProps[] };
    return payload.data ?? [];
  };

  const handleRecommendationsClick = async () => {
    if (isLoadingRecommendations || phase === "transitioning") return;

    setPhase("recommendations");
    setRecommendationsError(null);
    setIsLoadingRecommendations(true);

    try {
      const recData = await fetchRecommendations();
      setRecommendations(recData);
    } catch (error) {
      setRecommendations([]);
      if (error instanceof Error) {
        setRecommendationsError(error.message);
      } else {
        setRecommendationsError("Unable to load recommendations right now.");
      }
    } finally {
      setIsLoadingRecommendations(false);
    }
  };

  const handleBackToSearch = () => {
    if (phase !== "recommendations") return;
    setPhase("search");
  };

  return (
    <div className="relative mx-auto min-h-screen w-full max-w-196.5 bg-background text-foreground">
      <AnimatePresence mode="wait" initial={false}>
        {phase !== "recommendations" ? (
          <motion.section
            key="landing-search"
            initial={{ opacity: 0, y: 36 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -190 }}
            transition={{ duration: 0.75, ease: [0.22, 1, 0.36, 1] }}
            className="absolute inset-0 flex items-center justify-center"
          >
            <div className="flex w-full max-w-2xl flex-col items-center gap-6 px-4">
              <SearchBar
                placeholders={[
                  "Search for hackathons...",
                  "Try 'Hackathons near me'",
                  "Hackathons in ...",
                ]}
                interval={2500}
                onChange={() => {}}
                icon={null}
                onSubmit={handleSearchSubmit}
              />

              <CraftButton
                type="button"
                onClick={handleRecommendationsClick}
                disabled={isLoadingRecommendations}
                className="min-w-52"
              >
                <CraftButtonLabel>Recommendations</CraftButtonLabel>
                <CraftButtonIcon>
                  <ArrowRightIcon className="size-3 stroke-2 transition-transform duration-500 group-hover:rotate-90" />
                </CraftButtonIcon>
              </CraftButton>
            </div>
          </motion.section>
        ) : (
          <motion.section
            key="landing-recommendations"
            initial={{ opacity: 0, y: 110 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: 120 }}
            transition={{ duration: 0.75, ease: [0.22, 1, 0.36, 1] }}
            className="min-h-screen w-full"
          >
            <div className="mx-auto max-w-7xl px-4 py-8">
              <button
                type="button"
                onClick={handleBackToSearch}
                className="mb-4 inline-flex items-center gap-1.5 text-sm text-muted-foreground transition-colors hover:text-foreground"
              >
                <span>Recommendations</span>
                <ChevronDownIcon className="size-4" />
              </button>

              {isLoadingRecommendations && (
                <p className="mb-4 text-sm text-gray-500">Loading recommendations...</p>
              )}

              {recommendationsError && (
                <p className="mb-4 text-sm text-gray-500">{recommendationsError}</p>
              )}

              {!isLoadingRecommendations && !recommendationsError && recommendations.length === 0 && (
                <p className="text-sm text-muted-foreground">
                  No recommendations yet. Add more profile details and tech tags to get better matches.
                </p>
              )}

              {!recommendationsError && (
                <motion.div
                  initial="hidden"
                  animate="visible"
                  variants={{
                    hidden: {},
                    visible: {
                      transition: {
                        staggerChildren: 0.1,
                        delayChildren: 0.05,
                      },
                    },
                  }}
                  className="grid grid-cols-1 gap-4 sm:grid-cols-2"
                >
                  {recommendations.map((hackathon, index) => (
                    <motion.div
                      key={hackathon.id}
                      variants={{
                        hidden: { opacity: 0, y: 42 },
                        visible: { opacity: 1, y: 0 },
                      }}
                      transition={{ duration: 0.55, ease: [0.22, 1, 0.36, 1], delay: index * 0.02 }}
                    >
                      <HackathonCard hackathon={hackathon} />
                    </motion.div>
                  ))}
                </motion.div>
              )}
            </div>
          </motion.section>
        )}
      </AnimatePresence>
    </div>
  );
};

export default LandingPage;
