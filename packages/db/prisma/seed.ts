import "dotenv/config";
import { PrismaClient } from "./generated/prisma/client";
import { PrismaPg } from "@prisma/adapter-pg";
import pg from "pg";

const connectionString = process.env.DATABASE_URL;
const pool = new pg.Pool({ connectionString });
const adapter = new PrismaPg(pool);
const prisma = new PrismaClient({ adapter });

async function main() {
  console.log("🌱 Seeding database with 15 hackathons...");

  const hackathonsData = [
    {
      title: "Global AI Hackathon",
      host: "TechCorp",
      location: "Online",
      prize: "$10,000",
      startDate: new Date("2026-04-10"),
      endDate: new Date("2026-04-13"),
      applyUrl: "https://example.com/apply",
      tags: ["AI", "Web3"],
    },
    {
      title: "Ethereum Denver 2026",
      host: "ETH Global",
      location: "Denver, CO",
      prize: "$50,000",
      startDate: new Date("2026-03-01"),
      endDate: new Date("2026-03-05"),
      applyUrl: "https://ethdenver.com",
      tags: ["Blockchain", "Smart Contracts", "DeFi"],
    },
    {
      title: "Web3 Hackathon NYC",
      host: "NYC Web3 Foundation",
      location: "New York, NY",
      prize: "$25,000",
      startDate: new Date("2026-06-15"),
      endDate: new Date("2026-06-20"),
      applyUrl: "https://nycweb3hackathon.com",
      tags: ["Web3", "NFTs", "Metaverse"],
    },
    {
      title: "AI for Social Good Hackathon",
      host: "Global Impact Lab",
      location: "Online",
      prize: "$15,000",
      startDate: new Date("2026-09-10"),
      endDate: new Date("2026-09-12"),
      applyUrl: "https://globalimpactlab.org/ai-hackathon",
      tags: ["AI", "Social Impact", "Sustainability"],
    },
    {
      title: "DeFi Summer Hackathon",
      host: "DeFi Alliance",
      location: "Miami, FL",
      prize: "$30,000",
      startDate: new Date("2026-07-15"),
      endDate: new Date("2026-07-20"),
      applyUrl: "https://defialliance.com/summer-hackathon",
      tags: ["DeFi", "Blockchain", "Smart Contracts"],
    },
    {
      title: "Next.js UI Challenge",
      host: "Vercel Enthusiasts",
      location: "Online",
      prize: "$5,000",
      startDate: new Date("2026-05-01"),
      endDate: new Date("2026-05-03"),
      applyUrl: "https://nextjs-challenge.dev",
      tags: ["Frontend", "React", "Next.js"],
    },
    {
      title: "Go Cloud-Native Jam",
      host: "Gopher Foundation",
      location: "San Francisco, CA",
      prize: "$20,000",
      startDate: new Date("2026-08-12"),
      endDate: new Date("2026-08-14"),
      applyUrl: "https://gopherjam.org",
      tags: ["Golang", "Backend", "Microservices"],
    },
    {
      title: "DevOps Days Hackathon",
      host: "Cloud Native Computing Foundation",
      location: "Seattle, WA",
      prize: "$12,000",
      startDate: new Date("2026-10-05"),
      endDate: new Date("2026-10-07"),
      applyUrl: "https://devopsdayshack.com",
      tags: ["DevOps", "Docker", "Kubernetes"],
    },
    {
      title: "Tokyo Anime-Tech Hack",
      host: "Otaku Devs",
      location: "Tokyo, Japan",
      prize: "¥2,000,000",
      startDate: new Date("2026-11-20"),
      endDate: new Date("2026-11-22"),
      applyUrl: "https://anime-tech.jp",
      tags: ["Animation", "Media", "Japanese Culture"],
    },
    {
      title: "Bare Metal Assembly Hackathon",
      host: "Low Level Devs",
      location: "Online",
      prize: "$8,000",
      startDate: new Date("2026-12-01"),
      endDate: new Date("2026-12-03"),
      applyUrl: "https://baremetalhack.org",
      tags: ["Assembly", "C", "Systems Programming"],
    },
    {
      title: "Open Source Contributors Jam",
      host: "Linux Foundation",
      location: "Berlin, Germany",
      prize: "$15,000",
      startDate: new Date("2026-09-25"),
      endDate: new Date("2026-09-27"),
      applyUrl: "https://osjam.org",
      tags: ["Open Source", "Linux", "Community"],
    },
    {
      title: "Global Game Jam 2026",
      host: "IGDA",
      location: "Online",
      prize: "$0 (Glory)",
      startDate: new Date("2026-01-29"),
      endDate: new Date("2026-01-31"),
      applyUrl: "https://globalgamejam.org",
      tags: ["GameDev", "Unity", "Godot"],
    },
    {
      title: "DefCon Capture the Flag",
      host: "DEF CON",
      location: "Las Vegas, NV",
      prize: "$50,000",
      startDate: new Date("2026-08-06"),
      endDate: new Date("2026-08-09"),
      applyUrl: "https://defcon.org/ctf",
      tags: ["Cybersecurity", "Hacking", "InfoSec"],
    },
    {
      title: "Swift vs Kotlin Clash",
      host: "Mobile Dev Weekly",
      location: "London, UK",
      prize: "$10,000",
      startDate: new Date("2026-04-22"),
      endDate: new Date("2026-04-24"),
      applyUrl: "https://mobileclash.dev",
      tags: ["Mobile", "iOS", "Android"],
    },
    {
      title: "Space Apps Challenge",
      host: "NASA",
      location: "Global / Online",
      prize: "NASA Visit",
      startDate: new Date("2026-10-10"),
      endDate: new Date("2026-10-12"),
      applyUrl: "https://spaceapps.nasa.gov",
      tags: ["Space", "Data", "Science"],
    }
  ];

  for (const hackathon of hackathonsData) {
    await prisma.hackathon.create({
      data: hackathon,
    });
  }

  console.log(`✅ Successfully created 15 hackathons in the database!`);
}

main()
  .catch((e) => {
    console.error(e);
    process.exit(1);
  })
  .finally(async () => {
    await prisma.$disconnect();
    await pool.end();
  });