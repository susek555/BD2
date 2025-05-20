'use client';
import { RatingPercentages } from '@/app/lib/definitions/reviews';

interface RatingDistributionProps {
  averageRating: number;
  distribution: RatingPercentages;
}

export function RatingDistribution({
  averageRating,
  distribution,
}: RatingDistributionProps) {
  return (
    <div className='rounded-lg bg-white p-6 text-gray-800 shadow'>
      <div className='mb-4 flex items-center justify-between'>
        <div className='text-6xl font-light text-gray-700'>
          {averageRating.toFixed(1)}
        </div>
      </div>

      <div className='mb-4 space-y-2'>
        {['5', '4', '3', '2', '1'].map((rating) => {
          const percentage = distribution
            ? distribution[rating as keyof RatingPercentages]
            : 0;

          return (
            <div key={rating} className='flex items-center'>
              <div className='mr-2 w-6 text-right font-medium'>{rating}</div>
              <div className='relative h-5 flex-1 overflow-hidden rounded-sm bg-gray-200'>
                <div
                  className='h-full rounded-sm bg-yellow-500'
                  style={{ width: `${percentage}%` }}
                />
                <div className='absolute inset-y-0 right-2 flex items-center text-xs font-medium'>
                  {percentage}%
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
