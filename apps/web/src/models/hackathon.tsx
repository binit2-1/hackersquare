export interface HackathonProps {
  id: string;
  title: string;
  host: string;
  location: string;
  prize: string;
  startDate: string;
  endDate: string;
  applyUrl: string;
  tags: string[];
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