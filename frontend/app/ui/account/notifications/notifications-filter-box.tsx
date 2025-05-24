'use client'

import clsx from "clsx"
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { useCallback, useEffect, useState } from 'react';
import MarkAllAsReadButton from "./mark-all-as-read-button";


const orderOptions = [
    { label: 'Unread', value: 'unread' },
    { label: 'Newest', value: 'newest' },
]

export default function NotificationsFilterBox() {
    const router = useRouter();
    const pathname = usePathname();
    const searchParams = useSearchParams();

    const orderKey = searchParams.get('orderKey') || 'unread';

    const [localOrderKey, setLocalOrderKey] = useState(orderKey);
    const [hasChanges, setHasChanges] = useState(false);

    const setOrderKey = useCallback((key: string) => {
        setLocalOrderKey(key);
        setHasChanges(true);
    }, []);

    useEffect(() => {
        const urlOrderKey = searchParams.get('orderKey') || 'unread';
        if (urlOrderKey !== localOrderKey) {
            setLocalOrderKey(urlOrderKey);
            setHasChanges(false);
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [searchParams]);

    const applyFilters  = useCallback(() => {
        const params = new URLSearchParams(searchParams.toString());

        if (localOrderKey === 'unread') {
            params.delete('orderKey');
        } else {
            params.set('orderKey', localOrderKey);
        }

        params.delete('page');
        router.push(`${pathname}?${params.toString()}`);
        setHasChanges(false);
    }, [localOrderKey, searchParams, pathname, router]);


    return (
    <div className='flex h-full flex-col rounded-lg border-2 border-gray-300 bg-white p-4'>
      <div className='flex flex-col space-y-4'>
        <MarkAllAsReadButton />
        <div className='mb-4'>
          <h3 className='mb-2 text-sm font-medium text-gray-700'>Sort by</h3>
          <div className='flex items-center space-x-2'>
            <div className='inline-flex rounded-md shadow-sm'>
              {orderOptions.map((option) => (
                <button
                  key={option.value}
                  onClick={() => setOrderKey(option.value)}
                  className={clsx(
                    'border px-4 py-2 text-sm font-medium',
                    option.value === localOrderKey
                      ? 'z-10 border-blue-500 bg-blue-50 text-blue-700'
                      : 'border-gray-300 bg-white text-gray-700 hover:bg-gray-50',
                    orderOptions[0].value === option.value && 'rounded-l-md',
                    orderOptions[orderOptions.length - 1].value ===
                      option.value && 'rounded-r-md',
                  )}
                >
                  {option.label}
                </button>
              ))}
            </div>
          </div>
        </div>
      </div>

        {hasChanges ? (
          <div className='ml-auto'>
            <button
              onClick={applyFilters}
              className={`rounded-md px-4 py-2 text-sm font-medium ${
                hasChanges
                  ? 'bg-blue-500 text-white hover:bg-blue-600'
                  : 'cursor-not-allowed bg-gray-200 text-gray-500'
              }`}
            >
              Apply Filters
            </button>
          </div>
        ) : (
          <></>
        )}
    </div>
  );
}