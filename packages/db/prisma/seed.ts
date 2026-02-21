import "dotenv/config";
import { PrismaClient } from "./generated/prisma/client";
import { PrismaPg } from "@prisma/adapter-pg";
import pg from "pg";


const connectionString = process.env.DATABASE_URL;
const pool = new pg.Pool({ connectionString });
const adapter = new PrismaPg(pool);
const prisma = new PrismaClient({ adapter });

async function main() {
  console.log("🌱 Seeding database...");

  const hackathon1 = await prisma.hackathon.create({
    data: {
      title: "Global AI Hackathon",
      host: "TechCorp",
      location: "Online",
      prize: "$10,000",
      startDate: new Date(),
      endDate: new Date(new Date().setDate(new Date().getDate() + 3)),
      applyUrl: "https://example.com/apply",
      tags: ["AI", "Web3"],
    },
  });

  const hackathon2 = await prisma.hackathon.create({
    data: {
      title: "Ethereum Denver 2026",
      host: "ETH Global",
      location: "Denver, CO",
      prize: "$50,000",
      startDate: new Date("2026-03-01"),
      endDate: new Date("2026-03-05"),
      applyUrl: "https://ethdenver.com",
      tags: ["Blockchain", "Smart Contracts", "DeFi"],
    },
  });

  console.log(`✅ Created Hackathons: ${hackathon1.title} & ${hackathon2.title}`);
}

main()
  .catch((e) => {
    console.error(e);
    process.exit(1);
  })
  .finally(async () => {
    await prisma.$disconnect();
    await pool.end(); // Must close the pg pool explicitly in v7
  });