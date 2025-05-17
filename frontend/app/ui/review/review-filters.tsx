'use client';

import {
  ArrowDownIcon,
  ArrowUpIcon,
  XMarkIcon,
} from '@heroicons/react/24/outline';
import clsx from 'clsx';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { ReviewToggle } from '../account/reviews/created-received-toggle';

const orderOptions = [
  { label: 'Date', value: 'date' },
  { label: 'Rating', value: 'rating' },
];

const ratingOptions = [1, 2, 3, 4, 5];

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

  const updateFilters = (updates: Record<string, string | null>) => {
    const params = new URLSearchParams(searchParams);

    Object.entries(updates).forEach(([key, value]) => {
      if (value === null) {
        params.delete(key);
      } else {
        params.set(key, value);
      }
    });

    if (!('page' in updates)) {
      params.set('page', '1');
    }

    router.push(`${pathname}?${params.toString()}`);
  };

  const toggleRating = (rating: number) => {
    const newRatings = [...selectedRatings];
    const index = newRatings.indexOf(rating);

    if (index >= 0) {
      newRatings.splice(index, 1);
    } else {
      newRatings.push(rating);
    }

    updateFilters({
      ratings: newRatings.length > 0 ? newRatings.join(',') : null,
    });
  };

  const toggleSortDirection = () => {
    updateFilters({ isOrderDesc: (!isOrderDesc).toString() });
  };

  const changeOrderKey = (key: string) => {
    updateFilters({ orderKey: key });
  };

  const clearFilters = () => {
    updateFilters({
      ratings: null,
      orderKey: 'date',
      isOrderDesc: 'true',
    });
  };

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
                    option.value === orderKey
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
              {isOrderDesc ? (
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
            {selectedRatings.length > 0 && (
              <button
                className='mb-2 text-xs text-gray-500 hover:text-gray-700'
                onClick={() => updateFilters({ ratings: null })}
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
                  selectedRatings.includes(rating)
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

      <div className='mt-4 hidden justify-end md:flex'>
        {(selectedRatings.length > 0 ||
          orderKey !== 'date' ||
          !isOrderDesc) && (
          <button
            className='flex items-center text-sm text-gray-500 hover:text-gray-700'
            onClick={clearFilters}
          >
            <XMarkIcon className='mr-1 h-4 w-4' />
            Clear all filters
          </button>
        )}
      </div>
    </div>
  );
}
