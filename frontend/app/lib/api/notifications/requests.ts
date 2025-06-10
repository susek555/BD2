import { NextResponse } from 'next/server';
import { getServerSession } from 'next-auth';
import { authConfig } from '@/app/lib/authConfig';
import { fetchWithRefresh } from '@/app/lib/api/fetchWithRefresh';
import { API_URL } from '@/app/lib/constants';

export const markNotificationAs = async (id: string, seen: boolean) => {
    let response;
    if (seen) {
        response = await fetch(`/api/notifications/${id}/seen`, {
            method: 'PUT',
        });
    } else {
        response = await fetch(`/api/notifications/${id}/unseen`, {
            method: 'PUT',
        });
    }

    if (!response.ok) {
        throw new Error('Failed to update listing');
    }

    return response.json();
};

export const markAllNotificationsAsSeen = async () => {
    const response = await fetch('/api/notifications/seen', {
        method: 'PUT',
    });

    if (!response.ok) {
        throw new Error('Failed to update notifications');
    }

    return response.json();
};

export const getNotifications = async (params: {
    pagination?: {
        page?: number;
        page_size?: number;
    };
    order_key?: 'seen' | 'created_at';
    is_order_desc?: boolean;
}) => {
    const session = await getServerSession(authConfig);
    if (!session) {
        return NextResponse.json({ error: "Not authenticated" }, { status: 401 });
    }

    if (params.order_key ===  'seen') {
        params.is_order_desc = false;
    }

    const response = await fetchWithRefresh(`${API_URL}/notification/filter`, {
        method: "POST",
        body: JSON.stringify(params), // Sending params directly instead of wrapping in {params}
        headers: {
            'Content-Type': 'application/json',
        }
    });

    if (!response.ok) {
        return NextResponse.json({ error: "Failed to fetch notifications" }, { status: response.status });
    }

    const data = await response.json();
    console.log("Response from notifications API:", data);
    return data;
}