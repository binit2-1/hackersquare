export interface HackathonProps {
  id: string;
  title: string;
  host: string;
  startDate: string;
  endDate: string;
  location: string;
  prize: string;
  tags: string[];
  isBookmarked: boolean;
  applyUrl: string;
}

export const hackathons: HackathonProps[] = [
  {
    id: "1",
    title: "AI Build Week 2025",
    host: "Vercel",
    startDate: "2025-03-10",
    endDate: "2025-03-17",
    location: "San Francisco, CA",
    prize: "$25,000",
    tags: ["AI/ML", "Next.js", "TypeScript"],
    isBookmarked: false,
    applyUrl: "#",
  },
  {
    id: "2",
    title: "Web3 Frontier Hackathon",
    host: "Ethereum Foundation",
    startDate: "2025-04-01",
    endDate: "2025-04-03",
    location: "Online",
    prize: "$50,000",
    tags: ["Web3", "Solidity", "DeFi", "Blockchain"],
    isBookmarked: true,
    applyUrl: "#",
  },
  {
    id: "3",
    title: "Climate Tech Challenge",
    host: "Google.org",
    startDate: "2025-03-22",
    endDate: "2025-03-24",
    location: "New York, NY",
    prize: "$15,000",
    tags: ["Python", "Data Science", "IoT"],
    isBookmarked: false,
    applyUrl: "#",
  },
  {
    id: "4",
    title: "Open Source Sprint",
    host: "GitHub",
    startDate: "2025-05-05",
    endDate: "2025-05-07",
    location: "Online",
    prize: "$10,000",
    tags: ["Open Source", "React", "DevTools"],
    isBookmarked: false,
    applyUrl: "#",
  },
  {
    id: "5",
    title: "HealthTech Innovation Cup",
    host: "Mayo Clinic",
    startDate: "2025-04-14",
    endDate: "2025-04-16",
    location: "Austin, TX",
    prize: "$30,000",
    tags: ["HealthTech", "Python", "AI/ML", "IoT"],
    isBookmarked: true,
    applyUrl: "#",
  },
  {
    id: "6",
    title: "Mobile First Hackathon",
    host: "Meta",
    startDate: "2025-06-01",
    endDate: "2025-06-03",
    location: "Seattle, WA",
    prize: "$20,000",
    tags: ["React Native", "Mobile", "AR/VR"],
    isBookmarked: false,
    applyUrl: "#",
  },
  {
    id: "7",
    title: "FinTech Disrupt 2025",
    host: "Stripe",
    startDate: "2025-05-20",
    endDate: "2025-05-22",
    location: "Online",
    prize: "$40,000",
    tags: ["FinTech", "APIs", "TypeScript", "Payments"],
    isBookmarked: false,
    applyUrl: "#",
  },
  {
    id: "8",
    title: "Game Dev Jam",
    host: "Unity Technologies",
    startDate: "2025-07-11",
    endDate: "2025-07-13",
    location: "Los Angeles, CA",
    prize: "$12,000",
    tags: ["Game Dev", "Unity", "C#"],
    isBookmarked: false,
    applyUrl: "#",
  },
];
