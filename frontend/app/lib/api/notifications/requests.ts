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