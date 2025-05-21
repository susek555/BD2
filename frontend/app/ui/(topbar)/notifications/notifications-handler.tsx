import NotificationsButton from './notifications-button'

export default async function NotificationsHandler() {
    // fetch notifications from the server
    const newNotifications = 7

    return (
        <>
            <NotificationsButton newNotifications={newNotifications} />
        </>
    )
}