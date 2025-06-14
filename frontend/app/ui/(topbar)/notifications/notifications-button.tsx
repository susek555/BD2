'use client'

import { BellAlertIcon } from "@heroicons/react/20/solid"
import { BaseAccountButton } from "@/app/ui/(topbar)/base-account-buttons/base-account-button"
import NotificationModal from "@/app/ui/(topbar)/notifications/notifications-modal"
import { useState } from "react"
import { Notification } from "@/app/lib/definitions/notification"
import useNotificationsSocket from "./notifications-socket"
import { useEffect } from "react"

export default function NotificationsButton() {
    const [newNotifications, setNewNotifications] = useState(0);
    const [notifications, setNotificationsData] = useState<Notification[]>([]);

    function changeNumberOfUnread(delta: number) {
        setNewNotifications(prev => prev + delta);
    }

    function updateNotificationSeenStatus(id: string, seen: boolean) {
        setNotificationsData(prev => {
            return prev.map(notification => {
                if (notification.id === id) {
                    return { ...notification, seen };
                }
                return notification;
            })
        })
    };


    const { messages } = useNotificationsSocket()

    // Update notifications when new messages are received
    useEffect(() => {
        if (messages.length > 0) {
        try {
            const latest = messages[messages.length - 1]
            const newData = JSON.parse(latest)
            setNotificationsData(newData.notifications || [])
            setNewNotifications(newData.unseen_notifs_count || 0)
        } catch (e) {
            console.error('Błąd parsowania wiadomości:', e)
        }
        }
    }, [messages])

    const [isDialogOpen, setIsDialogOpen] = useState(false);

    function handleNotificationsButtonClick() {
        setIsDialogOpen(true);
    }

    return (
        <>
            <BaseAccountButton
                className="w-12 flex items-center justify-center relative"
                onClick={() => handleNotificationsButtonClick()}
            >
                <BellAlertIcon className="w-5 text-gray-50" />
                {newNotifications > 0 && (
                    <div className="absolute -top-1 -right-1 bg-red-500 text-white rounded-full w-5 h-5 flex items-center justify-center text-xs font-bold">
                        {newNotifications > 99 ? '99+' : newNotifications}
                    </div>
                )}
            </BaseAccountButton>
            <NotificationModal
                open={isDialogOpen}
                onOpenChange={setIsDialogOpen}
                notifications={notifications}
                changeNumberOfUnread={changeNumberOfUnread}
                updateNotificationSeenStatus={updateNotificationSeenStatus}
            />
        </>
    )
}