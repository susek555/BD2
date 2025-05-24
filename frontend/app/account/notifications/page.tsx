import { fetchNotifications } from "@/app/lib/data/notifications/data";
import { NotificationSearchParams } from "@/app/lib/definitions/notification";
import NotificationsTable from "@/app/ui/(notifications-table)/notifications-table";
import FoundInfo from "@/app/ui/(common)/found-info";
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

  const {totalPages, totalNotifications, notifications} = await fetchNotifications(notificationsSearchParams);

  return (
    <>
      <div className="flex flex-grow flex-row">
        <div className="mx-10 md:mx-45"/>
        <FoundInfo title={"Total notifications"} totalOffers={totalNotifications} />
      </div>
      <div className="my-4" />
      <div className='flex flex-grow flex-row'>
        <div className='h-full w-80 flex-none'>
          <NotificationsFilterBox />
        </div>

        <div className='flex flex-1 flex-col'>
          <div className='flex-1 px-12'>
            <NotificationsTable notifications={notifications} />
          </div>

          <div className='mt-5 flex w-full justify-center pb-4'>
            <Suspense>
              <Pagination totalPages={totalPages} />
            </Suspense>
          </div>
        </div>
      </div>
    </>
    );
}