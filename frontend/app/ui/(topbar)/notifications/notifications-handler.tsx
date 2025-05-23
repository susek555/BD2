import NotificationsButton from './notifications-button'
import { fetchNotificationsTopbar } from '@/app/lib/data/notifications/data'

export default async function NotificationsHandler() {
    // fetch notifications from the server
    const { newNotifications, notifications } = await fetchNotificationsTopbar()

    return (
        <>
            <NotificationsButton newNotifications={newNotifications} notifications={notifications}/>
        </>
    )
}