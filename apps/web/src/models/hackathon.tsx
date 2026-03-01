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

export interface SearchResponse {
  data: HackathonProps[];
  metadata: {
    totalRecords: number;
    currentPage: number;
    limit: number;
    totalPages: number;
  };
}