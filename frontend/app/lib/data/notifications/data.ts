import { Notification } from "@/app/lib/definitions/notification"

export async function fetchNotifications() {
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
            is_read: false,
        },
        {
            id: "2",
            offer_id: "2",
            title: "Offer accepted",
            description: "Your offer has been accepted",
            date: "2023-10-02",
            is_read: false,
        },
        {
            id: "3",
            offer_id: "3",
            title: "Offer declined",
            description: "Your offer has been declined",
            date: "2023-10-03",
            is_read: true,
        },
    ]

    return {
        newNotifications,
        notifications,
    }
}