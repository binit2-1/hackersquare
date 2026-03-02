export interface HackathonProps {
  id: string;
  title: string;
  host: string;
  location: string;
  prize_usd: number;
  start_date: string;
  end_date: string;
  apply_url: string;
  tags?: string[];
  isBookmarked: boolean;
}

export interface SearchResponse {
  data: HackathonProps[];
  metadata: {
    totalRecords: number;
    currentPage: number;
    limit: number;
    totalPages: number;
  };
}