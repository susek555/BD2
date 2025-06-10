'use client';

import { XMarkIcon } from '@heroicons/react/20/solid';
import * as Dialog from '@radix-ui/react-dialog';
import { VisuallyHidden } from '@radix-ui/react-visually-hidden';
import { useState } from 'react';
import NotificationsTable from '@/app/ui/(notifications-table)/notifications-table';
import { Notification } from '@/app/lib/definitions/notification';

interface NotificationsDialogProps {
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
  notifications?: Notification[];
  changeNumberOfUnread: (count: number) => void;
  updateNotificationSeenStatus: (id: string, seen: boolean) => void;
}

export default function NotificationsModal({
  open,
  onOpenChange,
  notifications,
  changeNumberOfUnread,
  updateNotificationSeenStatus
}: NotificationsDialogProps) {
  const [dialogOpen, setDialogOpen] = useState(open || false);

  const isControlled = open !== undefined && onOpenChange !== undefined;
  const isOpen = isControlled ? open : dialogOpen;
  const setIsOpen = isControlled ? onOpenChange : setDialogOpen;

  return (
    <Dialog.Root open={isOpen} onOpenChange={setIsOpen}>
      <Dialog.Portal>
        <VisuallyHidden>
          <Dialog.Title>Latest notifications</Dialog.Title>
        </VisuallyHidden>
        <Dialog.Overlay className='animate-fade-in fixed inset-0 z-40 bg-black/50' />
        <Dialog.Content
          className='animate-slide-up fixed top-1/2 left-1/2 z-50 max-h-[90vh] w-[90vw] max-w-5/12 -translate-x-1/2 -translate-y-1/2 overflow-auto rounded-lg bg-white shadow-lg'
          onClick={(e) => {
            e.stopPropagation();
          }}
        >
          <div className='relative'>
            <div>
              <p className='p-4 text-lg font-bold'>Notifications</p>
              {notifications && notifications.length > 0 ? (
                <NotificationsTable notifications={notifications} changeNumberOfUnread={changeNumberOfUnread} updateNotificationSeenStatus={updateNotificationSeenStatus}/>
              ) : (
                <div className="p-4 text-center text-gray-500">
                  You don&apos;t have any notifications
                </div>
              )}
              {notifications && notifications.length > 0 && (
              <div className="p-4 border-t border-gray-200">
                <a
                  href="/account/notifications"
                  className="block text-center text-blue-600 hover:text-blue-800 font-medium"
                >
                  See more notifications
                </a>
              </div>
              )}
            </div>

            <Dialog.Close asChild>
              <button
                className='absolute top-4 right-4 rounded-full p-1 hover:bg-gray-100'
                aria-label='Close'
              >
                <XMarkIcon className='h-6 w-6' />
              </button>
            </Dialog.Close>
          </div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
