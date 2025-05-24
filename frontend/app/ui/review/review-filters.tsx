'use client';

import {
  ArrowDownIcon,
  ArrowUpIcon,
  XMarkIcon,
} from '@heroicons/react/24/outline';
import clsx from 'clsx';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { useCallback, useEffect, useState } from 'react';
import { ReviewToggle } from '../account/reviews/created-received-toggle';

const orderOptions = [
  { label: 'Date', value: 'date' },
  { label: 'Rating', value: 'rating' },
];

const ratingOptions = [1, 2, 3, 4, 5];

interface FilterState {
  orderKey: string;
  isOrderDesc: boolean;
  selectedRatings: number[];
}

export function ReviewFilterBox({
  includeReviewToggle = false,
}: {
  includeReviewToggle: boolean;
}) {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  const orderKey = searchParams.get('orderKey') || 'date';
  const isOrderDesc = searchParams.get('isOrderDesc') !== 'false';
  const ratingsParam = searchParams.get('ratings');
  const selectedRatings = ratingsParam
    ? ratingsParam.split(',').map(Number)
    : [];

  const [localFilters, setLocalFilters] = useState<FilterState>({
    orderKey,
    isOrderDesc,
    selectedRatings,
  });

  const [hasChanges, setHasChanges] = useState(false);

  useEffect(() => {
    const urlOrderKey = searchParams.get('orderKey') || 'date';
    const urlIsOrderDesc = searchParams.get('isOrderDesc') !== 'false';
    const urlRatingsParam = searchParams.get('ratings');
    const urlSelectedRatings = urlRatingsParam
      ? urlRatingsParam.split(',').map(Number)
      : [];

    const currentRatings = [...localFilters.selectedRatings].sort().join(',');
    const urlRatings = [...urlSelectedRatings].sort().join(',');

    const filtersMatch =
      urlOrderKey === localFilters.orderKey &&
      urlIsOrderDesc === localFilters.isOrderDesc &&
      currentRatings === urlRatings;

    if (!filtersMatch) {
      setLocalFilters({
        orderKey: urlOrderKey,
        isOrderDesc: urlIsOrderDesc,
        selectedRatings: urlSelectedRatings,
      });

      setHasChanges(false);
    }
  }, [searchParams]);

  const applyFilters = useCallback(() => {
    const params = new URLSearchParams(searchParams.toString());

    params.set('orderKey', localFilters.orderKey);
    params.set('isOrderDesc', localFilters.isOrderDesc.toString());

    if (localFilters.selectedRatings.length > 0) {
      params.set('ratings', localFilters.selectedRatings.join(','));
    } else {
      params.delete('ratings');
    }

    params.set('page', '1');
    router.push(`${pathname}?${params.toString()}`);
    setHasChanges(false);
  }, [localFilters, router, pathname, searchParams]);

  const updateLocalFilters = useCallback((newFilters: FilterState) => {
    setLocalFilters(newFilters);
    setHasChanges(true);
  }, []);

  const toggleRating = useCallback(
    (rating: number): void => {
      const newRatings = [...localFilters.selectedRatings];
      const index = newRatings.indexOf(rating);

      if (index >= 0) {
        newRatings.splice(index, 1);
      } else {
        newRatings.push(rating);
      }

      updateLocalFilters({
        ...localFilters,
        selectedRatings: newRatings,
      });
    },
    [localFilters, updateLocalFilters],
  );

  const toggleSortDirection = useCallback((): void => {
    updateLocalFilters({
      ...localFilters,
      isOrderDesc: !localFilters.isOrderDesc,
    });
  }, [localFilters, updateLocalFilters]);

  const changeOrderKey = useCallback(
    (key: string): void => {
      updateLocalFilters({
        ...localFilters,
        orderKey: key,
      });
    },
    [localFilters, updateLocalFilters],
  );

  const clearRatings = useCallback((): void => {
    updateLocalFilters({
      ...localFilters,
      selectedRatings: [],
    });
  }, [localFilters, updateLocalFilters]);

  const clearFilters = useCallback((): void => {
    const defaultFilters = {
      orderKey: 'date',
      isOrderDesc: true,
      selectedRatings: [],
    };

    setLocalFilters(defaultFilters);

    const currentOrderKey = searchParams.get('orderKey') || 'date';
    const currentIsOrderDesc = searchParams.get('isOrderDesc') !== 'false';
    const currentRatings = searchParams.get('ratings')
      ? searchParams.get('ratings')!.split(',').map(Number)
      : [];

    const filtersMatch =
      defaultFilters.orderKey === currentOrderKey &&
      defaultFilters.isOrderDesc === currentIsOrderDesc &&
      defaultFilters.selectedRatings.length === currentRatings.length &&
      defaultFilters.selectedRatings.every((r) => currentRatings.includes(r));

    setHasChanges(!filtersMatch);
  }, [searchParams, updateLocalFilters]);

  return (
    <div className='flex h-full flex-col rounded-lg border-2 border-gray-300 bg-white p-4'>
      <div className='flex flex-col space-y-4'>
        <div className='mb-4'>
          {includeReviewToggle ? (
            <div className='mb-4'>
              <h3 className='mb-2 text-sm font-medium text-gray-700'>
                Reviews
              </h3>

              <ReviewToggle />
            </div>
          ) : (
            <></>
          )}
          <h3 className='mb-2 text-sm font-medium text-gray-700'>Sort by</h3>
          <div className='flex items-center space-x-2'>
            <div className='inline-flex rounded-md shadow-sm'>
              {orderOptions.map((option) => (
                <button
                  key={option.value}
                  onClick={() => changeOrderKey(option.value)}
                  className={clsx(
                    'border px-4 py-2 text-sm font-medium',
                    option.value === localFilters.orderKey
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
            <button
              onClick={toggleSortDirection}
              className='rounded-md border border-gray-300 bg-white p-2 hover:bg-gray-50'
            >
              {localFilters.isOrderDesc ? (
                <ArrowDownIcon className='h-5 w-5 text-gray-500' />
              ) : (
                <ArrowUpIcon className='h-5 w-5 text-gray-500' />
              )}
            </button>
          </div>
        </div>

        <div>
          <div className='flex items-center justify-between'>
            <h3 className='mb-2 text-sm font-medium text-gray-700'>
              Filter by rating
            </h3>
            {localFilters.selectedRatings.length > 0 && (
              <button
                className='mb-2 text-xs text-gray-500 hover:text-gray-700'
                onClick={clearRatings}
              >
                Clear
              </button>
            )}
          </div>
          <div className='flex space-x-2'>
            {ratingOptions.map((rating) => (
              <button
                key={rating}
                onClick={() => toggleRating(rating)}
                className={clsx(
                  'flex h-8 w-8 items-center justify-center rounded-md border text-sm font-medium',
                  localFilters.selectedRatings.includes(rating)
                    ? 'border-blue-500 bg-blue-500 text-white'
                    : 'border-gray-300 bg-white text-gray-700 hover:bg-gray-50',
                )}
              >
                {rating}
              </button>
            ))}
          </div>
        </div>
      </div>

      <div className='mt-4 flex w-full items-center justify-between'>
        <div>
          {(localFilters.selectedRatings.length > 0 ||
            localFilters.orderKey !== 'date' ||
            !localFilters.isOrderDesc) && (
            <button
              className='flex items-center text-sm text-gray-500 hover:text-gray-700'
              onClick={clearFilters}
            >
              <XMarkIcon className='mr-1 h-4 w-4' />
              Clear all filters
            </button>
          )}
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
    </div>
  );
}
