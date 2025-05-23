import { fetchNotifications } from "@/app/lib/data/notifications/data";
import { NotificationSearchParams } from "@/app/lib/definitions/notification";
import NotificationsTable from "@/app/ui/(notifications-table)/notifications-table";
import Pagination from "@/app/ui/(offers-table)/pagination";
import NotificationsFilterBox from "@/app/ui/account/notifications/notifications-filter-box";
import { Suspense } from "react";

export default async function Page({
    searchParams,
}: {
    searchParams?: Promise<{
        orderKey?: 'newest' | 'unread';
        page?: string;
    }>
}) {
  const params = (await searchParams) || {};
  const orderKey = params.orderKey || 'unread';
  const page = Number(params.page || 1);

  const notificationsSearchParams: NotificationSearchParams = {
    order_key: orderKey,
    pagination: {
      page: page,
      page_size: 8,
    },
  };

  const totalPages = 10; // TODO get page count from API
  const totalNotifications = 100; // TODO get total notifications from API
  const { newNotifications, notifications } = await fetchNotifications();

  return (
      <div className='flex flex-grow flex-row'>
        <div className='h-full w-80 flex-none py-4'>
          <NotificationsFilterBox />
        </div>

        <div className='flex flex-1 flex-col'>
          <div className='flex-1, p-6 px-12'>
            <NotificationsTable notifications={notifications} />
          </div>

          <div className='mt-5 flex w-full justify-center pb-4'>
            <Suspense>
              <Pagination totalPages={totalPages} />
            </Suspense>
          </div>
        </div>
      </div>
    );
}