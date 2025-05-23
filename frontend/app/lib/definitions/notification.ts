export interface Notification {
    id: string;
    offer_id: string;
    title: string;
    description: string;
    date: string;
    is_read: boolean;
}

export interface NotificationSearchParams {
    order_key?: string;
    pagination?: {
        page: number;
        page_size: number;
    };
}