export type User = {
    id: string;
    name: string;
    email: string;
};

export type AuthContextType = {
    user: User | null;
    isLoading: boolean;
    logout: () => void;
    refreshUser: () => Promise<void>;
}
