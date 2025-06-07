export interface Notification {
    id: string;
    offer_id: string;
    title: string;
    description: string;
    date: string;
    seen: boolean;
}

export interface NotificationSearchParams {
    order_key?: "seen" | "created_at";
    pagination?: {
        page: number;
        page_size: number;
    };
}
