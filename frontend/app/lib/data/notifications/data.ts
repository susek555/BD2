import { Notification, NotificationSearchParams } from "@/app/lib/definitions/notification"

export async function fetchNotificationsTopbar() {
    // TODO connect API

    await new Promise(resolve => setTimeout(resolve, 1000));

    const newNotifications = 7
    const notifications: Notification[] = [
        {
            id: "1",
            offer_id: "1",
            title: "New offer",
            description: "You have a new offer",
            date: "2023-10-01",
            seen: false,
        },
        {
            id: "2",
            offer_id: "2",
            title: "Offer accepted",
            description: "Your offer has been accepted",
            date: "2023-10-02",
            seen: false,
        },
        {
            id: "3",
            offer_id: "3",
            title: "Offer declined",
            description: "Your offer has been declined",
            date: "2023-10-03",
            seen: true,
        },
        {
            id: "4",
            offer_id: "4",
            title: "New message",
            description: "You have a new message",
            date: "2023-10-04",
            seen: false,
        }
    ]

    return {
        newNotifications,
        notifications,
    }
}

export async function fetchNotifications(params: NotificationSearchParams)  {
    // TODO connect API

    await new Promise(resolve => setTimeout(resolve, 1000));

    const notifications: Notification[] = [
        {
            id: "1",
            offer_id: "1",
            title: "New offer",
            description: "You have a new offer",
            date: "2023-10-01",
            seen: false,
        },
        {
            id: "2",
            offer_id: "2",
            title: "Offer accepted",
            description: "Your offer has been accepted",
            date: "2023-10-02",
            seen: false,
        },
        {
            id: "3",
            offer_id: "3",
            title: "Offer declined",
            description: "Your offer has been declined",
            date: "2023-10-03",
            seen: true,
        },
        {
            id: "4",
            offer_id: "4",
            title: "New message",
            description: "You have a new message",
            date: "2023-10-04",
            seen: false,
        }
    ]

    return {
        totalPages: 10,
        totalNotifications: 100,
        notifications,
    }
}