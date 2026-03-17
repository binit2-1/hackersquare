import { headers } from "next/headers";
import iso3166 from "iso-3166-2";

import LandingClient from "./landing-client";
import type { HackathonProps } from "@/models/hackathon";

const LandingPage = async () => {
  const requestHeaders = await headers();

  // 1. Grab Vercel headers with safe fallbacks for local dev.
  const city = requestHeaders.get("x-vercel-ip-city") || process.env.NEARBY_FALLBACK_CITY || "Bengaluru";
  const countryCode = requestHeaders.get("x-vercel-ip-country") || "IN";
  const regionCode = requestHeaders.get("x-vercel-ip-country-region") || "KA";

  // 2. Convert raw codes to full names for backend ranking.
  const state = iso3166.subdivision(countryCode, regionCode)?.name || "Karnataka";
  const country =
    iso3166.country(countryCode)?.name ||
    (countryCode === "IN" ? "India" : (countryCode === "US" ? "United States" : countryCode));

  // 3. Build query params expected by the Go backend.
  const params = new URLSearchParams({
    clientCity: city,
    clientState: state,
    clientCountry: country,
  });

  // 4. Fetch directly from Go API.
  const baseApiUrl = (process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080").replace(/\/+$/, "");
  const recommendationsApiUrl = `${baseApiUrl}/v1/recommendations?${params.toString()}`;

  let initialRecommendations: HackathonProps[] = [];
  let initialRecommendationsError: string | null = null;

  try {
    const response = await fetch(recommendationsApiUrl, {
      cache: "no-store",
      headers: {
        cookie: requestHeaders.get("cookie") || "",
      },
    });

    if (response.ok) {
      const payload = (await response.json()) as { data?: HackathonProps[] };
      initialRecommendations = payload.data ?? [];
    } else if (response.status === 401) {
      initialRecommendationsError = "Sign in to see your recommendations.";
    } else {
      initialRecommendationsError = "Unable to load recommendations right now.";
    }
  } catch {
    initialRecommendationsError = "Unable to load recommendations right now.";
  }

  return (
    <LandingClient
      recommendationsApiUrl={recommendationsApiUrl}
      initialRecommendations={initialRecommendations}
      initialRecommendationsError={initialRecommendationsError}
    />
  );
};

export default LandingPage;
