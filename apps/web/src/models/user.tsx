export type User = {
    id: string;
    name: string;
    email: string;
    headline: string;
    location: string;
    github_handle: string;
    website_url: string;
    linkedin_url: string;
    twitter_url: string;
    profileReadme?: string;
};

export type AuthContextType = {
    user: User | null;
    isLoading: boolean;
    logout: () => void;
    refreshUser: () => Promise<void>;
}
